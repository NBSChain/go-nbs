package memership

import (
	"github.com/NBSChain/go-nbs/storage/network"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/thirdParty/gossip/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"net"
)

type TaskType int

const (
	ProxyInitSubRequest TaskType = iota + 1
)

type peerNodeItem struct {
	nodeId      string
	probability float64
	addr        *nbsnet.NbsUdpAddr
}

type innerTask struct {
	tType TaskType
	param []interface{}
}

type MemManager struct {
	nodeID      string
	serviceConn *nbsnet.NbsUdpConn
	inputView   map[string]*peerNodeItem
	partialView map[string]*peerNodeItem
	taskSignal  chan innerTask
}

var (
	logger = utils.GetLogInstance()
)

func NewMemberNode(peerId string) *MemManager {

	node := &MemManager{
		nodeID:      peerId,
		taskSignal:  make(chan innerTask),
		inputView:   make(map[string]*peerNodeItem),
		partialView: make(map[string]*peerNodeItem),
	}

	return node
}

func (node *MemManager) InitNode() error {

	if err := node.initMsgService(); err != nil {
		return err
	}

	go node.receivingCmd()

	go node.taskDispatcher()

	if err := node.registerMySelf(); err != nil {
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

		n, peerAddr, err := node.serviceConn.ReadFromUDP(buffer)
		if err != nil {
			logger.Warning("reading contract application err:", err)
			continue
		}

		message := &pb.Gossip{}
		if err := proto.Unmarshal(buffer[:n], message); err != nil {
			logger.Warning("this is not a gossip message:->", buffer)
			continue
		}

		logger.Debug("gossip server:->", message, peerAddr)

		switch message.MessageType {
		case pb.MsgType_init:
			node.confirmAndPrepare(message.InitMsg, peerAddr)
		case pb.MsgType_reqContractAck:
			node.subToContract(message.ContactRes, peerAddr)
		default:
			continue
		}
	}
}

func (node *MemManager) taskWorker(task innerTask) {

	switch task.tType {
	case ProxyInitSubRequest:
		node.findProperContactNode(task.param[0].(*pb.InitSub), task.param[1].(*nbsnet.NbsUdpAddr))
	}
}

func (node *MemManager) taskDispatcher() {

	for {
		select {
		case task := <-node.taskSignal:
			node.taskWorker(task)
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
