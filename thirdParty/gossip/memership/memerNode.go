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
	MemShipHeartBeat  = time.Second * 100
	IsolatedTime      = time.Second * 3
	MaxInnerTaskSize  = 1 << 10
	MaxForwardTimes   = 10
	DefaultSubExpire  = time.Hour
	SubscribeTimeOut  = time.Second * 2
	MSGTrashCollect   = time.Minute * 10
	MaxItemPerRound   = 1 << 10
	ProbUpdateInter   = 10
)

var (
	ItemNotFound = fmt.Errorf("no such peer node in my view")
)

type msgTask struct {
	msg  *pb.Gossip
	addr *net.UDPAddr
}
type innerTask struct {
	params interface{}
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
	isBootNode  bool
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
		taskQueue:   make(chan *gossipTask, MaxInnerTaskSize),
		InputView:   make(map[string]*ViewNode),
		PartialView: make(map[string]*ViewNode),
		taskRouter:  make(map[int]worker),
		msgCounter:  make(map[string]*msgCounter),
	}

	node.taskRouter[int(nbsnet.GspSub)] = node.firstInitSub
	node.taskRouter[int(nbsnet.GspVoteContact)] = node.getVoteApply
	node.taskRouter[int(nbsnet.GspVoteResult)] = node.subToContract
	node.taskRouter[int(nbsnet.GspIntroduce)] = node.getForwardSub
	node.taskRouter[int(nbsnet.GspWelcome)] = node.subAccepted
	node.taskRouter[int(nbsnet.GspVoteResAck)] = node.voteAck
	node.taskRouter[SendHeartBeat] = node.sendHeartBeat
	node.taskRouter[int(nbsnet.GspHeartBeat)] = node.noop
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
	return node
}

func (node *MemManager) InitNode() error {

	if err := node.initMsgService(); err != nil {
		return err
	}

	go node.receivingCmd()

	go node.msgProcessor()

	if err := node.RegisterMySelf(); err != nil {
		logger.Warning(err)
		return err
	}

	go node.timer()

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

		n, peerAddr, err := node.serviceConn.ReadFromUDP(buffer)
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
			msgTask: msgTask{
				msg:  message,
				addr: peerAddr,
			},
		}
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
				logger.Debug("this is not my duty to process:->", task.taskType)
				continue
			}

			if err := handler(task); err != nil {
				logger.Error("gossip run loop err:->", err, task)
				continue
			}

			if task.msg != nil {
				node.freshInputView(task.msg.FromId)
			}

		case <-node.ctx.Done():
			logger.Info("gossip offline")
			return
		}
	}
}

func (node *MemManager) timer() {
	var isolateCheck, heartBeat, msgCollect time.Duration
	for {
		select {
		case <-time.After(time.Second):
			isolateCheck += time.Second
			heartBeat += time.Second
			msgCollect += time.Second

			if isolateCheck >= IsolatedTime {
				node.taskQueue <- &gossipTask{
					taskType: CheckItemInView,
				}
				isolateCheck = 0
			}
			if heartBeat >= MemShipHeartBeat {
				node.taskQueue <- &gossipTask{
					taskType: SendHeartBeat,
				}
				heartBeat = 0
			}
			if msgCollect >= MSGTrashCollect { //TODO::need a test.
				node.taskQueue <- &gossipTask{
					taskType: MsgCounterCollect,
				}
				msgCollect = 0
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
		if now.Sub(item.updateTime) > (MemShipHeartBeat*2 + IsolatedTime) {
			logger.Debug("more than isolate check:->")
			node.removeFromView(item, node.InputView)
		}
	}

	if len(node.InputView) == 0 && !node.isBootNode {
		return node.reSubscribe()
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
		FromId:  node.nodeID,
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
		if err := node.sendData(item, data); err != nil {
			logger.Warning("send data failed:->", err)
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
		item := node.randomSelectItem()
		logger.Debug("hey, don't introduce me to myself, forward:->", item.nodeId)
		return node.send(item, task.msg)
	}

	if _, ok := node.PartialView[subId]; ok {
		item := node.randomSelectItem()
		if subId == item.nodeId {
			logger.Debug("I think you have been introduced by yourself.")
			return nil
		}
		logger.Debug("I have got you, so forward to next node:->", item.nodeId)
		return node.send(item, task.msg)
	}

	prob := float64(1) / float64(1+len(node.PartialView))
	random, _ := rand.Int(rand.Reader, big.NewInt(100))
	logger.Debug("get introduced req:->", random, prob*100)

	if random.Int64() > int64(prob*100) {
		item := node.randomSelectItem()
		logger.Debug("no lucky, forward you, sorry:->", item.nodeId)
		return node.send(item, task.msg)
	}

	logger.Debug("yeah, lucky enough, accept this introduce, I am always your backup")
	return node.asSubAdapter(task.msg.Subscribe)
}

func (node *MemManager) noop(task *gossipTask) error {
	logger.Debug("noop:->", task.taskType, task.msg)
	return nil
}
