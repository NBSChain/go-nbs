package memership

import (
	"crypto/rand"
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/thirdParty/gossip/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/NBSChain/go-nbs/utils/crypto"
	"github.com/golang/protobuf/proto"
	"math/big"
	"net"
)

func (node *MemManager) broadCastSub(sub *pb.Subscribe) int {

	if len(node.PartialView) == 0 {
		logger.Info("no partial view node to broadcast ")
		return 0
	}
	sub.SeqNo++
	msg := &pb.Gossip{
		MsgType:   nbsnet.GspIntroduce,
		Subscribe: sub,
	}
	msg.MsgId = crypto.MD5SS(msg.String())
	data, _ := proto.Marshal(msg)

	forwardTime := 0
	logger.Debug("broad cast sub to all partial views:->", len(node.PartialView))
	for _, item := range node.PartialView {
		if err := item.sendData(data); err != nil {
			logger.Error("forward sub as contact err :->", err)
			continue
		}

		forwardTime++
	}

	for i := 0; i < utils.AdditionalCopies; i++ {
		item := node.choseRandomInPartialView()

		logger.Debug("random chose target:->", item.nodeId)

		if err := item.sendData(data); err != nil {
			logger.Error("forward extra C sub as contact err :->", err)
			continue
		}
		forwardTime++
	}

	return forwardTime
}

func (node *MemManager) publishVoteResult(sub *pb.Subscribe) error {

	item, err := node.newOutViewNode(sub.Addr, sub.Duration)
	if err != nil {
		logger.Error("create view node err:->", err)
		return err
	}
	sub.SeqNo++
	msg := &pb.Gossip{
		MsgType: nbsnet.GspVoteResult,
		VoteResult: &pb.Subscribe{
			Duration: sub.Duration,
			SeqNo:    sub.SeqNo,
			Addr:     nbsnet.ConvertToGossipAddr(item.outConn.LocAddr, node.nodeID),
		},
	}

	if err := item.send(msg); err != nil {
		logger.Error("send contact vote result err:->", err)
		return err
	}

	addr, ok := item.outConn.RealConn.RemoteAddr().(*net.UDPAddr)
	if !ok {
		return fmt.Errorf("get connection remote addr failed")
	}

	node.newInViewNode(sub.Addr.NetworkId, addr)

	return nil
}

func (node *MemManager) asContactServer(sub *pb.Subscribe) error {

	logger.Debug("ok I am your contact server")

	node.broadCastSub(sub)

	if err := node.publishVoteResult(sub); err != nil {
		return err
	}

	return nil
}

func (node *MemManager) asContactProxy(sub *pb.Subscribe, counter int) error {

	if node.subNo++; node.subNo >= ProbUpdateInter {
		node.taskQueue <- &gossipTask{
			taskType: UpdateProbability,
		}
	}

	if counter == 0 {
		return node.asContactServer(sub)
	}

	sub.SeqNo++
	req := &pb.Gossip{
		MsgType: nbsnet.GspVoteContact,
		VoteContact: &pb.VoteContact{
			TTL:       int32(counter) - 1,
			Subscribe: sub,
		},
	}

	if err := node.sendVoteApply(req); err != nil {
		logger.Warning(err)
		return node.asContactServer(sub)
	}

	return nil
}

func (node *MemManager) getVoteApply(task *gossipTask) error {
	req := task.msg.VoteContact
	return node.asContactProxy(req.Subscribe, int(req.TTL))
}

func (node *MemManager) sendVoteApply(pb *pb.Gossip) error {
	data, err := proto.Marshal(pb)
	if err != nil {
		return err
	}
	node.normalizeWeight(node.PartialView)

	var forwardTime int
	for _, item := range node.PartialView {

		pro, _ := rand.Int(rand.Reader, big.NewInt(100))
		logger.Debug("vote apply pro and itemPro:->", pro, item.probability*100)
		if pro.Int64() > int64(item.probability*100) {
			continue
		}

		if err := item.sendData(data); err != nil {
			continue
		}
		forwardTime++
	}

	if forwardTime == 0 {
		return fmt.Errorf("no contact node vote")
	}

	return nil
}

func (node *MemManager) asSubAdapter(sub *pb.Subscribe) error {

	logger.Debug("accept the subscriber:->", sub)
	nodeId := sub.Addr.NetworkId

	_, ok := node.PartialView[nodeId]
	if ok {
		return fmt.Errorf("duplicate accept subscribe=%s request:->", nodeId)
	}

	if node.nodeID == nodeId {
		return fmt.Errorf("hey it's yourself")
	}

	item, err := node.newOutViewNode(sub.Addr, sub.Duration)
	if err != nil {
		return err
	}

	msg := &pb.Gossip{
		MsgType: nbsnet.GspWelcome,
		SubConfirm: &pb.SynAck{
			FromId: node.nodeID,
		},
	}

	return item.send(msg)
}

func (node *MemManager) voteAck(task *gossipTask) error {
	nodeId := task.msg.VoteAck.FromId
	if _, ok := node.InputView[nodeId]; !ok {
		logger.Warning("no such node in input view:->", task.msg)
	}
	return nil
}
