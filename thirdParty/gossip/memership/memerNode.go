package memership

import (
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
	MemShipHeartBeat  = time.Second * 11 //TODO::?? heart beat time interval.
	MaxInnerTaskSize  = 1 << 10
	MaxForwardTimes   = 10
	DefaultSubExpire  = time.Hour
	SubscribeTimeOut  = time.Second * 2
	IsolatedTime      = MemShipHeartBeat * 3
	MSGTrashCollect   = time.Minute * 10
	MaxItemPerRound   = 1 << 10
)

var (
	HandlerNotFound = fmt.Errorf("no such gossip task handler")
)

type msgTask struct {
	isInner  bool
	taskType int
	msg      *pb.Gossip
	addr     *net.UDPAddr
}

type worker func(*msgTask) error

type msgCounter struct {
	counter int
	time    time.Time
}

type MemManager struct {
	nodeID      string
	updateTime  time.Time
	taskQueue   chan *msgTask
	serviceConn *nbsnet.NbsUdpConn
	inputView   map[string]*viewNode
	partialView map[string]*viewNode
	taskRouter  map[int]worker
	msgCounter  map[string]*msgCounter
}

var (
	logger = utils.GetLogInstance()
)

func NewMemberNode(peerId string) *MemManager {

	node := &MemManager{
		nodeID:      peerId,
		taskQueue:   make(chan *msgTask, MaxInnerTaskSize),
		inputView:   make(map[string]*viewNode),
		partialView: make(map[string]*viewNode),
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
			continue
		}

		message := &pb.Gossip{}
		if err := proto.Unmarshal(buffer[:n], message); err != nil {
			logger.Warning("this is not a gossip message:->", buffer)
			continue
		}

		logger.Debug("gossip server:->", peerAddr, message)

		node.taskQueue <- &msgTask{
			msg:  message,
			addr: peerAddr,
		}
	}
}

func (node *MemManager) msgProcessor() {

	for {
		select {
		case task := <-node.taskQueue:
			var handler worker
			var ok bool
			if task.isInner {
				handler, ok = node.taskRouter[task.taskType]
			} else {

				msgType := int(task.msg.MsgType)
				handler, ok = node.taskRouter[msgType]
			}
			if !ok {
				logger.Error("gossip msg handler err:->", HandlerNotFound)
			}
			if err := handler(task); err != nil {
				logger.Error("gossip run loop err:->", err)
			}
		}
	}
}

func (node *MemManager) timer() {
	for {
		select {
		case <-time.After(MemShipHeartBeat):
			node.taskQueue <- &msgTask{
				isInner:  true,
				taskType: SendHeartBeat,
			}

		case <-time.After(MSGTrashCollect):
			node.taskQueue <- &msgTask{
				isInner:  true,
				taskType: MsgCounterCollect,
			}
		}
	}
}

func (node *MemManager) msgCounterClean(task *msgTask) error {
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

func (node *MemManager) sendHeartBeat(task *msgTask) error {

	now := time.Now()
	if now.Sub(node.updateTime) > IsolatedTime {
		node.Resub()
	}

	msg := &pb.Gossip{
		MsgType: nbsnet.GspHeartBeat,
		HeartBeat: &pb.HeartBeat{
			FromID: node.nodeID,
		},
	}

	data, _ := proto.Marshal(msg)
	for _, item := range node.partialView {
		if !item.needUpdate() {
			continue
		}
		if err := item.sendData(data); err != nil {
			node.removeFromView(item, node.partialView)
		}

		if now.After(item.expiredTime) {
			node.removeFromView(item, node.partialView)
			node.unsubItem(item) //TODO::???
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

	if c.counter++; c.counter > MaxForwardTimes {
		return fmt.Errorf("msg(%s)forward too many times:->", msgId)
	}

	return nil
}

func (node *MemManager) getForwardSub(task *msgTask) error {

	if err := node.msgCache(task.msg.MsgId); err != nil {
		return err
	}

	if node.nodeID == task.msg.Subscribe.Addr.NetworkId {
		return fmt.Errorf("hey, it's me, no need to introduce me to myself")
	}

	req := task.msg.Subscribe
	prob := float64(1) / float64(1+len(node.partialView))
	random, _ := rand.Int(rand.Reader, big.NewInt(100))

	//TODO:: make sure this probability is fine.
	if random.Int64() < int64(prob*100) {
		return node.asSubAdapter(req)
	}

	if item := node.choseRandomInPartialView(task.msg.Subscribe.Addr.NetworkId); item != nil {
		return item.send(task.msg)
	}
	return node.asSubAdapter(req)
}
