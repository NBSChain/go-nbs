package memership

import (
	"crypto/rand"
	"github.com/NBSChain/go-nbs/storage/network"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/thirdParty/gossip/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"math/big"
)

func (node *MemManager) findProperContactNode(request *pb.InitSub, applierAddr *nbsnet.NbsUdpAddr) {

	//TODO::implement the indirect mechanism .
	counter := len(node.partialView)
	if counter == 0 {
		node.actAsContact(request, applierAddr)
		return
	}

	req := &pb.Gossip{
		MessageType: pb.MsgType_reqContract,
		ContactReq: &pb.ReqContact{
			ApplierId: request.NodeId,
			TTL:       int32(counter),
		},
	}

	node.indirectTheSubRequest(req)
}

func (node *MemManager) actAsContact(request *pb.InitSub, applierAddr *nbsnet.NbsUdpAddr) {
	count := len(node.partialView)
	if count == 0 {
		node.acceptThisSub(request, applierAddr)
		return
	}

	for _, item := range node.partialView {
		node.forwardSub(item, request, applierAddr)
	}

	//random, _ := rand.Int(rand.Reader, big.NewInt(int64(count)))

}

func (node *MemManager) indirectTheSubRequest(gossip *pb.Gossip) {

	node.updateProbability(node.partialView)

	for _, view := range node.partialView {
		pro, _ := rand.Int(rand.Reader, big.NewInt(100))

		if pro.Int64() < int64(view.probability*100) {
			continue
		}

		node.forwardContactRequest(view, gossip)
	}
}

func (node *MemManager) forwardContactRequest(peerNode *peerNodeItem, gossip *pb.Gossip) {
	//TODO:: make connection to him and send the request.
}

func (node *MemManager) acceptThisSub(sub *pb.InitSub, addr *nbsnet.NbsUdpAddr) {
	msg := &pb.Gossip{
		MessageType: pb.MsgType_reqContractAck,
		ContactRes: &pb.ReqContactACK{
			ApplierId:  sub.NodeId,
			SupplierId: node.peerId,
		},
	}

	port := utils.GetConfig().GossipCtrlPort
	conn, err := network.GetInstance().Connect(nil, addr, port)
	if err != nil {
		logger.Error("the contact failed to notify the subscriber:", err)
		return
	}
	defer conn.Close()

	msgData, err := proto.Marshal(msg)
	if err != nil {
		logger.Error("failed to marshal the contact init msg:", err)
		return
	}

	conn.Send(msgData)
}

func (node *MemManager) forwardSub(item *peerNodeItem, sub *pb.InitSub, addr *nbsnet.NbsUdpAddr) {

}
