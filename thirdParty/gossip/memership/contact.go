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
	"time"
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

		if item.nodeId == sub.NodeId {
			logger.Debug("don't introduce himself, no need.")
			continue
		}

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

	item, err := node.newOutViewNode(sub.Addr, time.Unix(sub.Expire, 0))
	if err != nil {
		logger.Error("create view node err:->", err)
		return err
	}
	sub.SeqNo++
	msg := &pb.Gossip{
		MsgType: nbsnet.GspVoteResult,
		VoteResult: &pb.Subscribe{
			Expire: sub.Expire,
			NodeId: node.nodeID,
			SeqNo:  sub.SeqNo,
			Addr:   nbsnet.ConvertToGossipAddr(item.outConn.LocAddr, node.nodeID),
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

	node.newInViewNode(sub.NodeId, addr)

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

	if item, ok := node.PartialView[sub.NodeId]; ok &&
		sub.Expire == item.expiredTime.Unix() {
		return fmt.Errorf("this sub(%s) has been accepted by me:->", item.nodeId)
	}

	if node.subNo++; node.subNo >= ProbUpdateInter {
		node.taskQueue <- &gossipTask{
			taskType: UpdateProbability,
		}
	}

	if counter == 0 {
		logger.Debug("yeah TTL is 0, I am your last station")
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

	data, _ := proto.Marshal(req)
	for i := 0; i < 3; i++ {
		if no := node.sendVoteApply(data, sub.NodeId); no > 0 {
			return nil
		}
		logger.Warning("no one want to vote, ask again:->", i)
	}

	return node.asContactServer(sub)
}

func (node *MemManager) getVoteApply(task *gossipTask) error {
	req := task.msg.VoteContact
	return node.asContactProxy(req.Subscribe, int(req.TTL))
}

func (node *MemManager) sendVoteApply(data []byte, targetId string) int {

	node.normalizeWeight(node.PartialView)

	var forwardTime int
	for _, item := range node.PartialView {

		pro, _ := rand.Int(rand.Reader, big.NewInt(100))
		logger.Debug("vote apply pro and itemPro:->", pro, item.probability*100)
		if pro.Int64() > int64(item.probability*100) {
			logger.Debug("no luck to send vote apply, try next one:->", item.nodeId)
			continue
		}
		if item.nodeId == targetId {
			logger.Debug("don't let him find himself, the life is already so hard:->", targetId)
			continue
		}

		logger.Debug("ok, vote for him please:->", item.nodeId)
		if err := item.sendData(data); err != nil {
			logger.Warning("failed to send apply for err:->", err)
			continue
		}
		forwardTime++
	}

	return forwardTime
}

func (node *MemManager) asSubAdapter(sub *pb.Subscribe) error {

	logger.Debug("accept the subscriber:->", sub)

	item, ok := node.PartialView[sub.NodeId]
	if ok && sub.Expire == item.expiredTime.Unix() {
		return fmt.Errorf("duplicate accept subscribe=%s request:->", sub.NodeId)
	}

	if node.nodeID == sub.NodeId {
		return fmt.Errorf("hey it's yourself")
	}
	item, err := node.newOutViewNode(sub.Addr, time.Unix(sub.Expire, 0))
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
