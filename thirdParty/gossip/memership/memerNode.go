package memership

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/thirdParty/gossip/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"math/big"
	"net"
	"time"
)

const (
	SendHeartBeat     = 1
	MsgCounterCollect = 2
	CheckItemInView   = 3
	UpdateProbability = 4
	GetInputViews     = 5
	GetOutputViews    = 6
	ClearOutputViews  = 7
	ClearInputViews   = 8

	MemShipHeartBeat = time.Second * 20 //TODO::?? heart beat time interval.
	MaxInnerTaskSize = 1 << 10
	MaxForwardTimes  = 10
	DefaultSubExpire = time.Hour
	SubscribeTimeOut = time.Second * 2
	IsolatedTime     = MemShipHeartBeat * 5
	MSGTrashCollect  = time.Minute * 10
	MaxItemPerRound  = 1 << 10
	ProbUpdateInter  = 10
)

var (
	HandlerNotFound = fmt.Errorf("no such gossip task handler")
	ItemNotFound    = fmt.Errorf("no such peer node in my view")
)

type msgTask struct {
	msg  *pb.Gossip
	addr *net.UDPAddr
}
type innerTask struct {
	params []interface{}
	result chan interface{}
}
type gossipTask struct {
	taskType int
	msgTask
	innerTask
}

type worker func(*gossipTask) error

type msgCounter struct {
	counter int
	time    time.Time
}

type MemManager struct {
	ctx         context.Context
	close       context.CancelFunc
	nodeID      string
	subNo       int
	updateTime  time.Time
	taskQueue   chan *gossipTask
	serviceConn *nbsnet.NbsUdpConn
	InputView   map[string]*ViewNode
	PartialView map[string]*ViewNode
	taskRouter  map[int]worker
	msgCounter  map[string]*msgCounter
}

var (
	logger = utils.GetLogInstance()
)

func NewMemberNode(peerId string) *MemManager {

	ctx, cal := context.WithCancel(context.Background())

	node := &MemManager{
		nodeID:      peerId,
		ctx:         ctx,
		close:       cal,
		updateTime:  time.Now(),
		taskQueue:   make(chan *gossipTask, MaxInnerTaskSize),
		InputView:   make(map[string]*ViewNode),
		PartialView: make(map[string]*ViewNode),
		taskRouter:  make(map[int]worker),
		msgCounter:  make(map[string]*msgCounter),
	}

	node.taskRouter[int(nbsnet.GspSub)] = node.firstInitSub
	node.taskRouter[int(nbsnet.GspVoteContact)] = node.getVoteApply
	node.taskRouter[int(nbsnet.GspVoteResult)] = node.subToContract
	node.taskRouter[int(nbsnet.GspHeartBeat)] = node.getHeartBeat
	node.taskRouter[int(nbsnet.GspIntroduce)] = node.getForwardSub
	node.taskRouter[int(nbsnet.GspWelcome)] = node.subAccepted
	node.taskRouter[int(nbsnet.GspVoteResAck)] = node.voteAck
	node.taskRouter[SendHeartBeat] = node.sendHeartBeat
	node.taskRouter[MsgCounterCollect] = node.msgCounterClean
	node.taskRouter[CheckItemInView] = node.checkItemInView
	node.taskRouter[int(nbsnet.GspReplaceArc)] = node.replaceForUnsubPeer
	node.taskRouter[int(nbsnet.GspReplaceAck)] = node.acceptAsReplacedPeer
	node.taskRouter[int(nbsnet.GspRemoveIVArc)] = node.removeUnsubPeerFromOut
	node.taskRouter[int(nbsnet.GspRemoveOVAcr)] = node.removeUnsubPeerFromIn
	node.taskRouter[UpdateProbability] = node.updateProbability
	node.taskRouter[int(nbsnet.GspUpdateOVWei)] = node.updateMyInProb
	node.taskRouter[int(nbsnet.GspUpdateIVWei)] = node.updateMyOutProb
	node.taskRouter[int(nbsnet.GspSubACK)] = node.reSubAckConfirm

	//TODO::refactor this debug command
	node.taskRouter[GetOutputViews] = node.outputViewTask
	node.taskRouter[GetInputViews] = node.inputViewTask
	node.taskRouter[ClearOutputViews] = node.removeOV
	node.taskRouter[ClearInputViews] = node.removeIV

	return node
}

func (node *MemManager) InitNode() error {

	if err := node.initMsgService(); err != nil {
		return err
	}

	go node.receivingCmd()

	go node.msgProcessor()

	go node.timer()

	if err := node.RegisterMySelf(); err != nil {
		logger.Warning(err)
		return err
	}

	return nil
}

func (node *MemManager) initMsgService() error {

	conn, err := network.GetInstance().ListenUDP("udp4", &net.UDPAddr{
		Port: utils.GetConfig().GossipCtrlPort,
	})

	if err != nil {
		logger.Error("can't start contract service:", err)
		return err
	}

	node.serviceConn = conn

	return nil
}

func (node *MemManager) receivingCmd() {

	for {
		buffer := make([]byte, utils.NormalReadBuffer)

		n, peerAddr, err := node.serviceConn.ReceiveFromUDP(buffer)
		if err != nil {
			logger.Warning("reading contact application err:", err)
			break
		}

		message := &pb.Gossip{}
		if err := proto.Unmarshal(buffer[:n], message); err != nil {
			logger.Warning("this is not a gossip message:->", peerAddr, n)
			continue
		}

		logger.Debug("gossip server:->", peerAddr, message)

		task := &gossipTask{
			taskType: int(message.MsgType),
		}
		task.msg = message
		task.addr = peerAddr
		node.taskQueue <- task

		select {
		case <-node.ctx.Done():
			logger.Debug("mem manager finished")
			break
		default:
		}
	}
}

func (node *MemManager) msgProcessor() {

	for {
		select {
		case task := <-node.taskQueue:
			var handler worker
			var ok bool
			handler, ok = node.taskRouter[task.taskType]
			if !ok {
				logger.Error("gossip msg handler err:->", HandlerNotFound, task.msg, task.taskType)
				continue
			}

			if err := handler(task); err != nil {
				logger.Error("gossip run loop err:->", err, task)
			}
		case <-node.ctx.Done():
			logger.Info("gossip offline")
			return
		}
	}
}

func (node *MemManager) timer() {
	for {
		select {
		case <-time.After(MemShipHeartBeat):
			node.taskQueue <- &gossipTask{
				taskType: SendHeartBeat,
			}

			node.taskQueue <- &gossipTask{
				taskType: CheckItemInView,
			}

		case <-time.After(MSGTrashCollect):
			node.taskQueue <- &gossipTask{
				taskType: MsgCounterCollect,
			}
		case <-node.ctx.Done():
			logger.Info("gossip offline")
			return
		}
	}
}

func (node *MemManager) checkItemInView(task *gossipTask) error {
	now := time.Now()

	for _, item := range node.InputView {
		if now.Sub(item.updateTime) > IsolatedTime {
			logger.Debug("more than isolate check:->")
			node.removeFromView(item, node.InputView)
		}
	}

	if len(node.InputView) == 0 && now.Sub(node.updateTime) > IsolatedTime {
		return node.Resub()
	}

	return nil
}

func (node *MemManager) msgCounterClean(task *gossipTask) error {
	no := 0
	now := time.Now()
	for id, c := range node.msgCounter {
		if now.Sub(c.time) > MSGTrashCollect {
			delete(node.msgCounter, id)
		}
		if no++; no > MaxItemPerRound {
			break
		}
	}
	return nil
}

func (node *MemManager) sendHeartBeat(task *gossipTask) error {

	now := time.Now()

	msg := &pb.Gossip{
		MsgType: nbsnet.GspHeartBeat,
		HeartBeat: &pb.HeartBeat{
			FromID: node.nodeID,
		},
	}

	data, _ := proto.Marshal(msg)
	for _, item := range node.PartialView {
		if now.After(item.expiredTime) {
			logger.Warning("subscribe expired:->", item.expiredTime, now, item.nodeId)
			node.removeFromView(item, node.PartialView)
		}

		if now.Sub(item.updateTime) < (MemShipHeartBeat / 2) {
			continue
		}
		logger.Debug("send heart beat to:->", item.nodeId, item.outConn.String())
		if err := item.sendData(data); err != nil {
			logger.Warning("send data failed:->", err)
			node.removeFromView(item, node.PartialView)
		}
	}

	return nil
}

func (node *MemManager) msgCache(msgId string) error {

	c, ok := node.msgCounter[msgId]

	if !ok {
		c = &msgCounter{
			counter: 0,
			time:    time.Now(),
		}
		node.msgCounter[msgId] = c
	}

	if c.counter++; c.counter >= MaxForwardTimes {
		return fmt.Errorf("msg(%s)forward too many times:->", msgId)
	}

	return nil
}

func (node *MemManager) getForwardSub(task *gossipTask) error {

	if err := node.msgCache(task.msg.MsgId); err != nil {
		return err
	}

	if len(node.PartialView) == 0 {
		logger.Debug("I have no friends right now, welcome you")
		return node.asSubAdapter(task.msg.Subscribe)
	}

	subId := task.msg.Subscribe.NodeId
	if subId == node.nodeID {
		item := node.choseRandomInPartialView()
		logger.Debug("hey, don't introduce me to myself, forward:->", item.nodeId)
		return item.send(task.msg)
	}

	if _, ok := node.PartialView[subId]; ok {
		item := node.choseRandomInPartialView()
		logger.Debug("I have got you, so forward to next node:->", item.nodeId)
		return item.send(task.msg)
	}

	prob := float64(1) / float64(1+len(node.PartialView))
	random, _ := rand.Int(rand.Reader, big.NewInt(100))
	logger.Debug("get introduced req:->", random, prob*100)

	if random.Int64() > int64(prob*100) {
		item := node.choseRandomInPartialView()
		logger.Debug("no lucky, forward you, sorry:->", item.nodeId)
		return item.send(task.msg)
	}

	logger.Debug("yeah, I am always your backup")
	return node.asSubAdapter(task.msg.Subscribe)
}
