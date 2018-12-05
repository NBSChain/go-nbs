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

func (node *MemManager) proxySubReq(task *innerTask) error {
	req := task.msg.InitSub

	sub := &subOnline{
		nodeId: req.NodeId,
		seq:    req.Seq,
		addr:   req.Addr,
	}

	counter := 2 * len(node.partialView)

	return node.indirectTheSubRequest(sub, counter)
}

func (node *MemManager) actAsContact(sub *subOnline) error {

	count := len(node.partialView)
	if count == 0 {
		return node.acceptSub(sub)

	}

	forwardTime := 0
	for _, item := range node.partialView {
		if err := node.forwardSub(item, sub); err != nil {
			logger.Error("forward sub as contact err :->", err)
			continue
		}
		forwardTime++
	}

	for i := 0; i < utils.AdditionalCopies; i++ {
		item := node.choseRandomInPartialView()
		if err := node.forwardSub(item, sub); err != nil {
			logger.Error("forward extra C sub as contact err :->", err)
			continue
		}
		forwardTime++
	}
	if forwardTime == 0 {
		return fmt.Errorf("no success forward made even if the partial view is not empty:->", count)
	}

	return nil
}

func (node *MemManager) indirectTheSubRequest(sub *subOnline, counter int) error {

	if counter == 0 {
		return node.actAsContact(sub)
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
		//TODO:: make sure this probability is fine.
		if pro.Int64() > int64(view.probability*100) {
			continue
		}

		if err := node.forwardContactRequest(view, req); err == nil {
			forwardTime++
		}
	}

	if forwardTime == 0 {
		return node.acceptSub(sub)
	}
	return nil
}

func (node *MemManager) forwardContactRequest(peerNode *peerNodeItem, gossip *pb.Gossip) error {

	data, _ := proto.Marshal(gossip)

	return node.PushOut(false, peerNode.nodeId, data)
}

func (node *MemManager) acceptSub(sub *subOnline) error {
	logger.Debug("accept the subscriber:->", sub.nodeId, sub.addr.String())

	item, ok := node.partialView[sub.nodeId]
	if ok {
		item := node.choseRandomInPartialView()
		return node.forwardSub(item, sub)

	}

	addr := nbsnet.ConvertFromGossipAddr(sub.addr)
	addr.NetworkId = sub.nodeId

	conn, err := node.notifySubscriber(sub.seq, addr)
	if nil != err {
		logger.Error("failed to notify the subscriber:", err)
		return err
	}

	item = &peerNodeItem{
		nodeId:   sub.nodeId,
		addr:     addr,
		ctrlConn: conn,
	}

	node.partialView[sub.nodeId] = item
	item.probability = 1 / float64(len(node.partialView))

	//TODO:: ? need to update probability?
	node.updateProbability(node.partialView)

	return nil
}

func (node *MemManager) forwardSub(item *peerNodeItem, sub *subOnline) error {

	if item.encounterNo++; item.encounterNo > MaxForwardTimes {
		return nil
	}
	msg := &pb.Gossip{
		MsgType: nbsnet.GspForwardSub,
		TransSubReq: &pb.TransSubReq{
			From:      node.nodeID,
			ApplierID: sub.nodeId,
			SeqNo:     sub.seq,
			Addr:      sub.addr,
		},
	}

	data, _ := proto.Marshal(msg)
	if _, err := item.ctrlConn.Write(data); err != nil {
		logger.Warning("forward sub request err:->", err)
		return err
	}
	logger.Debug("decide to forward the sub request:->", msg)
	return nil
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
