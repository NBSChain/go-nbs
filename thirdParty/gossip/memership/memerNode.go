package memership

import (
	"github.com/NBSChain/go-nbs/storage/network"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/thirdParty/gossip/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"net"
)

type MemberNode struct {
	peerId       string
	isPublic     bool
	contractNode *ContractNode
	serviceConn  *net.UDPConn
}

func isInPublic() bool {

	natType := network.GetInstance().NatType()

	var canService bool
	switch natType {
	case nat_pb.NatType_UnknownRES:
		canService = false

	case nat_pb.NatType_NoNatDevice:
		canService = true

	case nat_pb.NatType_BehindNat:
		canService = false

	case nat_pb.NatType_CanBeNatServer:
		canService = true

	case nat_pb.NatType_ToBeChecked:
		canService = false
	}

	return canService
}

func NewMemberNode(peerId string) *MemberNode {

	node := &MemberNode{
		peerId:       peerId,
		isPublic:     isInPublic(),
		contractNode: newContractNode(),
	}

	if err := node.startService(); err != nil {
		panic(err)
	}

	if err := node.initSubRequest(); err != nil {
		panic(err)
	}

	return node
}

func (node *MemberNode) startService() error {

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

func (node *MemberNode) runLoop() {

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
			node.contractNode.proxyInit(message.InitMsg)
		default:
			continue
		}
	}
}
