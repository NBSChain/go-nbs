package memership

import (
	"crypto/rand"
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/thirdParty/gossip/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/NBSChain/go-nbs/utils/crypto"
	"github.com/golang/protobuf/proto"
	"math/big"
	"net"
	"time"
)

func (node *MemManager) broadCastSub(sub *pb.Subscribe) {

	if len(node.PartialView) == 0 {
		logger.Info("no partial view node to broadcast ")
		return
	}
	sub.SeqNo++
	msg := &pb.Gossip{
		MsgType:   nbsnet.GspIntroduce,
		Subscribe: sub,
	}
	msg.MsgId = crypto.MD5SS(msg.String())
	data, _ := proto.Marshal(msg)

	logger.Debug("broad cast sub to all partial views:->", len(node.PartialView))
	for _, item := range node.PartialView {

		if item.nodeId == sub.NodeId {
			logger.Debug("don't introduce himself:->", item.nodeId)
			continue
		}

		if err := node.sendData(item, data); err != nil {
			logger.Error("forward sub as contact err :->", err)
			continue
		}
	}

	if sub.IsReSub {
		logger.Info("ReSub no need to make additional forward")
		return
	}

	for i := 0; i < utils.AdditionalCopies; i++ {
		item := node.randomSelectItem()

		logger.Debug("for system robust, additional forward:->", item.nodeId)

		if err := node.sendData(item, data); err != nil {
			logger.Error("forward extra C sub as contact err :->", err)
			continue
		}
	}
}

func (node *MemManager) publishVoteResult(sub *pb.Subscribe) error {

	item, err := node.newOutViewNode(sub.Addr, time.Unix(sub.Expire, 0))
	if err != nil {
		logger.Error("create view node err:->", err)
		return err
	}
	sub.SeqNo++
	netId, nbsAddr := network.GetInstance().GetNatAddr()
	msg := &pb.Gossip{
		MsgType: nbsnet.GspVoteResult,
		VoteResult: &pb.Subscribe{
			Expire: sub.Expire,
			NodeId: node.nodeID,
			SeqNo:  sub.SeqNo,
			Addr:   nbsnet.ConvertToGossipAddr(nbsAddr, netId),
		},
	}

	if err := node.send(item, msg); err != nil {
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

	logger.Debug("ok I am your contact server, receive the init subscribe")

	node.broadCastSub(sub)

	if err := node.publishVoteResult(sub); err != nil {
		return err
	}

	return nil
}

func (node *MemManager) asContactProxy(sub *pb.Subscribe, counter int) error {
	if len(node.PartialView) == 0 {
		logger.Debug("I have no output views, get you.")
		return node.asContactServer(sub)
	}

	//TODO::crash and restart again, how to process?
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
	return node.sendVoteApply(data, sub.NodeId)
}

func (node *MemManager) getVoteApply(task *gossipTask) error {
	req := task.msg.VoteContact
	return node.asContactProxy(req.Subscribe, int(req.TTL))
}

func (node *MemManager) chooseWithProb() *ViewNode {

	randValue, _ := rand.Int(rand.Reader, big.NewInt(10000))
	randPro := float64(randValue.Int64()) / 10000

	logger.Debug("vote apply pro:->", randPro)
	var start float64

	for _, item := range node.PartialView {

		start += item.probability
		logger.Debug("item with prob:->", item.nodeId, item.probability, start)

		if randPro <= start {
			logger.Debug("selected node is:->", item.nodeId)
			return item
		}
	}
	return nil
}

func (node *MemManager) sendVoteApply(data []byte, targetId string) error {

	node.normalizeWeight(node.PartialView)

	item := node.chooseWithProb()

	return node.sendData(item, data)
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
		FromId:  node.nodeID,
	}
	logger.Debug("welcome you:->", sub.NodeId)
	return node.send(item, msg)
}

func (node *MemManager) voteAck(task *gossipTask) error {
	nodeId := task.msg.FromId
	if _, ok := node.InputView[nodeId]; !ok {
		logger.Warning("no such node in input view:->", task.msg)
	}
	return nil
}
