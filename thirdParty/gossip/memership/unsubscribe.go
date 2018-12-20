package memership

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/thirdParty/gossip/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"time"
)

func (node *MemManager) DestroyNode() error {
	lenIn := len(node.InputView)
	lenOut := len(node.PartialView)

	tempOut := make([]string, len(node.PartialView))
	for nodeId := range node.PartialView {
		tempOut = append(tempOut, nodeId)
	}

	tempIn := make([]string, len(node.InputView))
	for nodeId := range node.InputView {
		tempIn = append(tempIn, nodeId)
	}

	node.replaceMeByMyOutView(lenIn, lenOut, tempOut, tempIn)
	node.removeMeFromOutView(lenIn, tempIn)
	node.removeMeFromInView()

	node.close()

	if err := node.serviceConn.Close(); err != nil {
		logger.Warning("gossip offline err:->", err)
		return err
	}

	close(node.taskQueue)

	return nil
}

func (node *MemManager) replaceMeByMyOutView(lenIn, lenOut int, tempOut, tempIn []string) {

	for i := lenIn - utils.AdditionalCopies - 1 - 1; i >= 0 && lenOut > 0; i-- {
		j := i % lenOut
		outId := tempOut[j]
		outItem := node.PartialView[outId]

		msg := &pb.Gossip{
			MsgType: nbsnet.GspReplaceArc,
			FromId:  node.nodeID,
			ArcReplace: &pb.ArcReplace{
				ToId: outId,
				Addr: nbsnet.ConvertToGossipAddr(outItem.outAddr, outItem.nodeId),
			},
		}

		inId := tempIn[i]
		inItem := node.InputView[inId]

		data, _ := proto.Marshal(msg)
		if _, err := node.serviceConn.WriteToUDP(data, inItem.inAddr); err != nil {
			logger.Warning("")
			continue
		}
	}
}

func (node *MemManager) replaceForUnsubPeer(task *gossipTask) error {
	replace := task.msg.ArcReplace
	item, ok := node.PartialView[task.msg.FromId]
	if !ok {
		return ItemNotFound
	}

	logger.Debug("remove old unsub node and replace it with new one:->")
	node.removeFromView(item, node.PartialView)

	if _, ok := node.PartialView[replace.ToId]; ok {
		return fmt.Errorf("no need to make a new item, I have got it")
	}

	expT := time.Now().Add(DefaultSubExpire)

	item, err := node.newOutViewNode(replace.Addr, expT)
	if err != nil {
		logger.Warning("new node err:->", err)
		return err
	}
	netId, nbsAddr := network.GetInstance().GetNatAddr()
	msg := &pb.Gossip{
		MsgType: nbsnet.GspReplaceAck,
		ReplaceAck: &pb.Subscribe{
			SeqNo:  1,
			Expire: expT.Unix(),
			NodeId: netId,
			Addr:   nbsnet.ConvertToGossipAddr(nbsAddr, netId),
		},
	}

	logger.Debug("replace node cause'of unsub:->", task.msg.FromId, replace.ToId)

	return node.send(item, msg)
}

func (node *MemManager) acceptAsReplacedPeer(task *gossipTask) error {

	ack := task.msg.ReplaceAck

	_, ok := node.InputView[ack.NodeId]
	if !ok {
		return fmt.Errorf("no need to replace, I have got it")
	}

	node.newInViewNode(ack.NodeId, task.addr)

	logger.Debug("get new input item cause'of some unsub:->", ack.NodeId)

	return nil
}

func (node *MemManager) removeMeFromOutView(lenIn int, tempIn []string) {

	msg := &pb.Gossip{
		MsgType: nbsnet.GspRemoveIVArc,
		FromId:  node.nodeID,
	}
	data, _ := proto.Marshal(msg)

	for i := lenIn - utils.AdditionalCopies - 1; i < lenIn && i >= 0; i++ {

		inId := tempIn[i]
		inItem := node.InputView[inId]

		if _, err := node.serviceConn.WriteToUDP(data, inItem.inAddr); err != nil {
			logger.Warning("")
			continue
		}
	}
}

func (node *MemManager) removeMeFromInView() {
	msg := &pb.Gossip{
		MsgType: nbsnet.GspRemoveOVAcr,
		FromId:  node.nodeID,
	}

	for _, item := range node.PartialView {
		if err := node.send(item, msg); err != nil {
			logger.Warning("notify remove me err:->", err)
		}
	}
}

func (node *MemManager) removeUnsubPeerFromOut(task *gossipTask) error {

	nodeId := task.msg.FromId
	item, ok := node.PartialView[nodeId]
	if !ok {
		return ItemNotFound
	}
	logger.Debug("remove unsub node form output view:->")
	node.removeFromView(item, node.PartialView)

	return nil
}

func (node *MemManager) removeUnsubPeerFromIn(task *gossipTask) error {

	nodeId := task.msg.FromId
	item, ok := node.InputView[nodeId]
	if !ok {
		return ItemNotFound
	}
	logger.Debug("remove unsub node form input view:->")
	node.removeFromView(item, node.InputView)

	return nil
}
