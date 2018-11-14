package memership

import (
	"github.com/NBSChain/go-nbs/storage/network"
	"github.com/NBSChain/go-nbs/thirdParty/gossip/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"net"
)

type MemberManager interface {
}

type memberNode struct {
	contractNode *contractNode
	serviceConn  *net.UDPConn
}

func NewMemberNode() MemberManager {

	node := &memberNode{
		contractNode: newContractNode(),
	}

	if err := node.startService(); err != nil {
		panic(err)
	}

	go node.registerToNetwork()

	return node
}

func (node *memberNode) startService() error {

	conn, err := net.ListenUDP("udp4", &net.UDPAddr{
		Port: utils.GetConfig().GossipContractServicePort,
	})

	if err != nil {
		logger.Error("can't start contract service:", err)
		return err
	}

	node.serviceConn = conn

	go node.runLoop()

	return nil
}

func (node *memberNode) registerToNetwork() {

}

func (node *memberNode) runLoop() {

	for {
		buffer := make([]byte, network.NormalReadBuffer)

		n, peerAddr, err := node.serviceConn.ReadFrom(buffer)
		if err != nil {
			logger.Warning("reading contract application err:", err)
			continue
		}

		if n >= network.NormalReadBuffer {
			//TODO:: check what we can to support this situation.
		}

		logger.Debug("receive contract apply:", peerAddr)

		message := &pb.Gossip{}
		if err := proto.Unmarshal(buffer[:n], message); err != nil {
			logger.Warning("this is not a gossip message:->", buffer)
			continue
		}

		switch message.MsgType {
		case pb.Type_init:
			node.contractNode.initSub(message.InitMsg)
		default:
			continue
		}
	}
}
