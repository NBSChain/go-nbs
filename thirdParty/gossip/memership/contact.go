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
			Seq:       request.Seq,
			TTL:       int32(counter),
			ApplierID: request.NodeId,
			Applier:   request.Addr,
		},
	}

	node.indirectTheSubRequest(req)
}

func (node *MemManager) actAsContact(request *pb.InitSub, applierAddr *nbsnet.NbsUdpAddr) {

	count := len(node.partialView)
	if count == 0 {
		node.acceptSub(request, applierAddr)
		return
	}

	for _, item := range node.partialView {
		node.forwardSub(item, request, applierAddr)
	}

	for i := 0; i < utils.AdditionalCopies; i++ {
		item := node.choseRandomInPartialView()
		node.forwardSub(item, request, applierAddr)
	}
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

func (node *MemManager) acceptSub(sub *pb.InitSub, addr *nbsnet.NbsUdpAddr) {

	_, ok := node.partialView[sub.NodeId]
	if ok {
		item := node.choseRandomInPartialView()
		node.forwardSub(item, sub, addr)
		return
	}

	item := &peerNodeItem{
		nodeId:      sub.NodeId,
		addr:        addr,
		probability: 1, //TODO::
	}
	node.partialView[sub.NodeId] = item

	//TODO:: ? need to update probability?
	node.updateProbability(node.partialView)

	node.notifySubscriber(sub, addr)
}

func (node *MemManager) forwardSub(item *peerNodeItem, sub *pb.InitSub, addr *nbsnet.NbsUdpAddr) {
	//TODO::
}

func (node *MemManager) notifySubscriber(sub *pb.InitSub, addr *nbsnet.NbsUdpAddr) {

	port := utils.GetConfig().GossipCtrlPort
	conn, err := network.GetInstance().Connect(nil, addr, port)
	if err != nil {
		logger.Error("the contact failed to notify the subscriber:", err)
		return
	}
	defer conn.Close()

	msg := &pb.Gossip{
		MessageType: pb.MsgType_reqContractAck,
		ContactRes: &pb.ReqContactACK{
			Seq:        sub.Seq,
			SupplierID: node.nodeID,
			Supplier:   nbsnet.ConvertToGossipAddr(conn.LocalAddr()),
		},
	}

	msgData, err := proto.Marshal(msg)
	if err != nil {
		logger.Error("failed to marshal the contact init msg:", err)
		return
	}

	conn.Send(msgData)
}
