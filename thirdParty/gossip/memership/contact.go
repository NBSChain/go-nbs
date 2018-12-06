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

	if len(node.partialView) == 0 {
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
	for _, item := range node.partialView {
		if err := item.sendData(data); err != nil {
			logger.Error("forward sub as contact err :->", err)
			continue
		}
		forwardTime++
	}

	for i := 0; i < utils.AdditionalCopies; i++ {
		item := node.choseRandomInPartialView()

		if item == nil {
			continue
		}

		if err := item.sendData(data); err != nil {
			logger.Error("forward extra C sub as contact err :->", err)
			continue
		}
		forwardTime++
	}

	return forwardTime
}

func (node *MemManager) publishVoteResult(sub *pb.Subscribe) error {

	item, err := newOutViewNode(sub, node.partialView)
	if err != nil {
		logger.Error("create view node err:->", err)
		return err
	}

	msg := &pb.Gossip{
		MsgType: nbsnet.GspVoteResult,
		VoteResult: &pb.Subscribe{
			Duration: sub.Duration,
			SeqNo:    sub.SeqNo + 1,
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

	newInViewNode(sub.Addr.NetworkId, addr, node.inputView)

	return nil
}

func (node *MemManager) asContactServer(sub *pb.Subscribe) error {

	node.broadCastSub(sub)

	if err := node.publishVoteResult(sub); err != nil {
		return err
	}

	return nil
}

func (node *MemManager) asContactProxy(sub *pb.Subscribe, counter int) error {

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

func (node *MemManager) getVoteApply(task *msgTask) error {
	req := task.msg.VoteContact
	return node.asContactProxy(req.Subscribe, int(req.TTL))
}

func (node *MemManager) sendVoteApply(pb *pb.Gossip) error {
	data, err := proto.Marshal(pb)
	if err != nil {
		return err
	}

	var forwardTime int
	for _, item := range node.partialView {

		pro, _ := rand.Int(rand.Reader, big.NewInt(100))
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

	item, ok := node.partialView[nodeId]
	if ok {
		if item := node.choseRandomInPartialView(); item != nil {
			msg := &pb.Gossip{
				MsgType:   nbsnet.GspIntroduce,
				Subscribe: sub,
			}
			msg.MsgId = crypto.MD5SS(msg.String())
			return item.send(msg)
		}
	}

	item, err := newOutViewNode(sub, node.partialView)
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
