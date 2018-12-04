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
)

func (node *MemManager) firstSubOnline(task *innerTask) error {
	req, ok := task.param.(*pb.InitSub)
	if !ok {
		return fmt.Errorf("not enough param")
	}

	sub := &newSub{
		nodeId: req.NodeId,
		seq:    req.Seq,
		addr:   req.Addr,
	}

	counter := 2 * len(node.partialView)

	node.indirectTheSubRequest(sub, counter)
	return nil
}

func (node *MemManager) actAsContact(sub *newSub) {

	count := len(node.partialView)
	if count == 0 {
		node.acceptSub(sub)
		return
	}

	for _, item := range node.partialView {
		node.forwardSub(item, sub)
	}

	for i := 0; i < utils.AdditionalCopies; i++ {
		item := node.choseRandomInPartialView()
		node.forwardSub(item, sub)
	}
}

func (node *MemManager) indirectTheSubRequest(sub *newSub, counter int) {

	if counter == 0 {
		node.actAsContact(sub)
		return
	}

	req := &pb.Gossip{
		MsgType: nbsnet.GspRegContact,
		ReqContact: &pb.ReqContact{
			Seq:       sub.seq,
			TTL:       int32(counter) - 1,
			ApplierID: sub.nodeId,
			Applier:   sub.addr,
		},
	}

	node.updateProbability(node.partialView)

	var forwardTime int
	for _, view := range node.partialView {
		pro, _ := rand.Int(rand.Reader, big.NewInt(100))

		if pro.Int64() > int64(view.probability*100) {
			continue
		}

		if err := node.forwardContactRequest(view, req); err == nil {
			forwardTime++
		}
	}

	if forwardTime == 0 {
		node.acceptSub(sub)
	}
}

func (node *MemManager) forwardContactRequest(peerNode *peerNodeItem, gossip *pb.Gossip) error {

	data, _ := proto.Marshal(gossip)

	return node.keepAliveWithData(peerNode.nodeId, data)
}

func (node *MemManager) acceptSub(sub *newSub) {
	logger.Debug("accept the subscriber:->", sub.nodeId, sub.addr.String())

	item, ok := node.partialView[sub.nodeId]
	if ok {
		item := node.choseRandomInPartialView()
		node.forwardSub(item, sub)
		return
	}

	addr := nbsnet.ConvertFromGossipAddr(sub.addr)
	addr.NetworkId = sub.nodeId

	conn, err := node.notifySubscriber(sub.seq, addr)
	if nil != err {
		logger.Error("failed to notify the subscriber:", err)
		return
	}

	item = &peerNodeItem{
		nodeId: sub.nodeId,
		addr:   addr,
		conn:   conn,
	}

	node.partialView[sub.nodeId] = item
	item.probability = 1 / float64(len(node.partialView))

	//TODO:: ? need to update probability?
	node.updateProbability(node.partialView)
}

func (node *MemManager) forwardSub(item *peerNodeItem, sub *newSub) {
	//TODO::
}

func (node *MemManager) notifySubscriber(Seq int64, addr *nbsnet.NbsUdpAddr) (*nbsnet.NbsUdpConn, error) {

	port := utils.GetConfig().GossipCtrlPort
	conn, err := network.GetInstance().Connect(nil, addr, port)
	if err != nil {
		logger.Error("the contact failed to notify the subscriber:", err)
		return nil, err
	}

	msg := &pb.Gossip{
		MsgType: nbsnet.GspContactAck,
		ReqContactACK: &pb.ReqContactACK{
			Seq:        Seq,
			SupplierID: node.nodeID,
			Supplier:   nbsnet.ConvertToGossipAddr(conn.LocalAddr()),
		},
	}

	msgData, _ := proto.Marshal(msg)

	if _, err := conn.Send(msgData); err != nil {
		logger.Warning("failed to send data to notify subscriber.", err)
		return nil, err
	}

	return conn, nil
}
