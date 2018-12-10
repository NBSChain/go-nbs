package memership

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/thirdParty/gossip/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
)

func (node *MemManager) DestroyNode() error {
	node.close()

	if err := node.serviceConn.Close(); err != nil {
		logger.Warning("gossip offline err:->", err)
		return err
	}

	lenIn := len(node.inputView)
	lenOut := len(node.partialView)

	tempOut := make([]string, len(node.partialView))
	for nodeId := range node.partialView {
		tempOut = append(tempOut, nodeId)
	}

	tempIn := make([]string, len(node.inputView))
	for nodeId := range node.inputView {
		tempIn = append(tempIn, nodeId)
	}

	node.replaceMeByMyOutView(lenIn, lenOut, tempOut, tempIn)
	node.removeMeFromOutView(lenIn, tempIn)
	node.removeMeFromInView()

	node.inputView = make(map[string]*viewNode)
	node.partialView = make(map[string]*viewNode)
	node.msgCounter = make(map[string]*msgCounter)

	close(node.taskQueue)

	return nil
}

func (node *MemManager) replaceMeByMyOutView(lenIn, lenOut int, tempOut, tempIn []string) {

	for i := lenIn - utils.AdditionalCopies - 1 - 1; i >= 0 && lenOut > 0; i-- {
		j := i % lenOut
		outId := tempOut[j]
		outItem := node.partialView[outId]

		msg := &pb.Gossip{
			MsgType: nbsnet.GspReplaceArc,
			ArcReplace: &pb.ArcReplace{
				FromId: node.nodeID,
				ToId:   outId,
				Addr:   nbsnet.ConvertToGossipAddr(outItem.outAddr, outItem.nodeId),
			},
		}

		inId := tempIn[i]
		inItem := node.inputView[inId]

		data, _ := proto.Marshal(msg)
		if _, err := node.serviceConn.WriteToUDP(data, inItem.inAddr); err != nil {
			logger.Warning("")
			continue
		}
	}
}

func (node *MemManager) replaceForUnsubPeer(task *msgTask) error {
	replace := task.msg.ArcReplace

	item, ok := node.partialView[replace.FromId]
	if !ok {
		return ItemNotFound
	}

	node.removeFromView(item, node.partialView)

	if _, ok := node.partialView[replace.ToId]; ok {
		return fmt.Errorf("no need to make a new item, I have got it")
	}

	item, err := node.newOutViewNode(replace.Addr, int64(DefaultSubExpire))
	if err != nil {
		logger.Warning("new node err:->", err)
		return err
	}

	msg := &pb.Gossip{
		MsgType: nbsnet.GspReplaceAck,
		ReplaceAck: &pb.Subscribe{
			SeqNo:    1,
			Duration: int64(DefaultSubExpire),
			Addr:     nbsnet.ConvertToGossipAddr(item.outConn.LocAddr, node.nodeID),
		},
	}

	logger.Debug("replace node cause'of unsub:->", replace.FromId, replace.ToId)

	return item.send(msg)
}

func (node *MemManager) acceptAsReplacedPeer(task *msgTask) error {

	ack := task.msg.ReplaceAck
	nodeId := ack.Addr.NetworkId

	_, ok := node.inputView[nodeId]
	if !ok {
		return fmt.Errorf("no need to replace, I have got it")
	}

	node.newInViewNode(nodeId, task.addr)

	logger.Debug("get new input item cause'of some unsub:->", nodeId)

	return nil
}

func (node *MemManager) removeMeFromOutView(lenIn int, tempIn []string) {

	msg := &pb.Gossip{
		MsgType: nbsnet.GspRemoveIVArc,
		ArcDrop: &pb.ArcDrop{
			NodeId: node.nodeID,
		},
	}
	data, _ := proto.Marshal(msg)

	for i := lenIn - utils.AdditionalCopies - 1; i < lenIn && i >= 0; i++ {

		inId := tempIn[i]
		inItem := node.inputView[inId]

		if _, err := node.serviceConn.WriteToUDP(data, inItem.inAddr); err != nil {
			logger.Warning("")
			continue
		}
	}
}

func (node *MemManager) removeMeFromInView() {
	msg := &pb.Gossip{
		MsgType: nbsnet.GspRemoveOVAcr,
		ArcDrop: &pb.ArcDrop{
			NodeId: node.nodeID,
		},
	}

	for _, item := range node.partialView {
		if err := item.send(msg); err != nil {
			logger.Warning("notify remove me err:->", err)
		}
	}
}

func (node *MemManager) removeUnsubPeerFromOut(task *msgTask) error {

	rm := task.msg.ArcDrop
	nodeId := rm.NodeId
	item, ok := node.partialView[nodeId]
	if !ok {
		return ItemNotFound
	}

	node.removeFromView(item, node.partialView)

	return nil
}

func (node *MemManager) removeUnsubPeerFromIn(task *msgTask) error {

	rm := task.msg.ArcDrop
	nodeId := rm.NodeId
	item, ok := node.inputView[nodeId]
	if !ok {
		return ItemNotFound
	}

	node.removeFromView(item, node.inputView)

	return nil
}
