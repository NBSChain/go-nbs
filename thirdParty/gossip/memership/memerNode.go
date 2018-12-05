package memership

import (
	"crypto/rand"
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/thirdParty/gossip/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"math/big"
	"net"
	"time"
)

const (
	MemShipHeartBeat = time.Second * 10 //TODO::?? heart beat time interval.
	MaxInnerTaskSize = 1 << 10
	MaxForwardTimes  = 10
)

var (
	HandlerNotFound         = fmt.Errorf("no such gossip task handler")
	PartialViewItemNotFound = fmt.Errorf("no such item in partial view")
)

type subOnline struct {
	nodeId string
	seq    int64
	addr   *pb.BasicHost
}

type peerNodeItem struct {
	nodeId      string
	encounterNo int
	probability float64
	addr        *nbsnet.NbsUdpAddr
	updateTime  time.Time
	ctrlConn    *nbsnet.NbsUdpConn
}

type innerTask struct {
	msg  *pb.Gossip
	addr *net.UDPAddr
}

type worker func(*innerTask) error

type MemManager struct {
	nodeID      string
	serviceConn *nbsnet.NbsUdpConn
	inputView   map[string]*peerNodeItem
	partialView map[string]*peerNodeItem
	taskQueue   chan *innerTask
	taskRouter  map[net_pb.MsgType]worker
}

var (
	logger = utils.GetLogInstance()
)

func NewMemberNode(peerId string) *MemManager {

	node := &MemManager{
		nodeID:      peerId,
		taskQueue:   make(chan *innerTask, MaxInnerTaskSize),
		inputView:   make(map[string]*peerNodeItem),
		partialView: make(map[string]*peerNodeItem),
		taskRouter:  make(map[net_pb.MsgType]worker),
	}

	node.taskRouter[nbsnet.GspInitSub] = node.firstInitSub
	node.taskRouter[nbsnet.GspProxySub] = node.proxySubReq
	node.taskRouter[nbsnet.GspContactAck] = node.subToContract
	node.taskRouter[nbsnet.GspHeartBeat] = node.getHeartBeat
	node.taskRouter[nbsnet.GspForwardSub] = node.getForwardedRequest

	return node
}

func (node *MemManager) InitNode() error {

	if err := node.initMsgService(); err != nil {
		return err
	}

	go node.receivingCmd()

	go node.RunLoop()

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

		node.taskQueue <- &innerTask{
			msg:  message,
			addr: peerAddr,
		}
	}
}

func (node *MemManager) RunLoop() {

	for {
		select {
		case task := <-node.taskQueue:
			msgType := task.msg.MsgType
			handler, ok := node.taskRouter[msgType]
			if !ok {
				logger.Error("gossip msg handler err:->", HandlerNotFound)
			}
			if err := handler(task); err != nil {
				logger.Error("gossip run loop err:->", err)
			}

		case <-time.After(MemShipHeartBeat):
			node.PushOut(true, "", nil)
		}
	}
}

func (node *MemManager) updateProbability(view map[string]*peerNodeItem) {

	var summerOut float64
	for _, item := range view {
		summerOut += item.probability
	}

	for _, item := range view {
		item.probability = item.probability / summerOut
	}
}

func (node *MemManager) getForwardedRequest(task *innerTask) error {

	prob := float64(1) / float64(1+len(node.partialView))

	random, _ := rand.Int(rand.Reader, big.NewInt(100))

	//TODO:: make sure this probability is fine.
	if random.Int64() > int64(prob*100) {
		return nil
	}

	req := task.msg.TransSubReq
	sub := &subOnline{
		nodeId: req.ApplierID,
		seq:    req.SeqNo,
		addr:   req.Addr,
	}

	item, ok := node.partialView[req.ApplierID]
	if ok {

		return node.forwardSub(item, sub)
	}

	return node.acceptSubAsNode(sub)
}

func (node *MemManager) acceptSubAsNode(sub *subOnline) error {

	return nil
}
