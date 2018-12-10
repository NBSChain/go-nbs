package memership

import (
	"crypto/rand"
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/thirdParty/gossip/pb"
	"github.com/golang/protobuf/proto"
	"math/big"
	"time"
)

//TODO:: make sure this random is ok
func (node *MemManager) choseRandomInPartialView() *viewNode {
	count := len(node.partialView)
	j := 0
	random, _ := rand.Int(rand.Reader, big.NewInt(int64(count)))
	logger.Debug("chose random in partialView :->", random)
	for _, item := range node.partialView {
		if j == int(random.Int64()) {
			return item
		} else {
			j++
		}
	}
	return nil
}

func (node *MemManager) removeFromView(item *viewNode, views map[string]*viewNode) {

	delete(views, item.nodeId)

	if item.outConn != nil {
		if err := item.outConn.Close(); err != nil {
			logger.Warning(err)
		}
	}
	logger.Warning("remove node from view:->", item.nodeId)
}

func (node *MemManager) getHeartBeat(task *msgTask) error {
	beat := task.msg.HeartBeat
	item, ok := node.inputView[beat.FromID]
	if !ok {
		err := fmt.Errorf("no such input view item:->%s", beat.FromID)
		logger.Warning(err)
		return err
	}
	item.updateTime = time.Now()
	return nil
}

func (node *MemManager) updateMyInProb(task *msgTask) error {

	wei := task.msg.OVWeight
	item, ok := node.inputView[wei.NodeId]
	if !ok {
		return ItemNotFound
	}

	item.probability = wei.Weight
	return nil
}

func (node *MemManager) updateMyOutProb(task *msgTask) error {
	wei := task.msg.IVWeight
	item, ok := node.partialView[wei.NodeId]
	if !ok {
		return ItemNotFound
	}

	item.probability = wei.Weight

	return nil
}

func (node *MemManager) normalizeWeight(views map[string]*viewNode) {
	var summerOut float64
	for _, item := range views {
		summerOut += item.probability
	}

	for _, item := range views {
		item.probability = item.probability / summerOut
	}
}

func (node *MemManager) updateProbability(task *msgTask) error {

	node.normalizeWeight(node.partialView)

	for _, item := range node.partialView {
		msg := &pb.Gossip{
			MsgType: nbsnet.GspUpdateOVWei,
			OVWeight: &pb.WeightUpdate{
				NodeId: node.nodeID,
				Weight: item.probability,
			},
		}

		if err := item.send(msg); err != nil {
			logger.Warning("send weight update to partial view item err:->", err, item.nodeId)
		}
	}

	node.normalizeWeight(node.inputView)

	for _, item := range node.inputView {
		msg := &pb.Gossip{
			MsgType: nbsnet.GspUpdateIVWei,
			IVWeight: &pb.WeightUpdate{
				NodeId: node.nodeID,
				Weight: item.probability,
			},
		}

		data, _ := proto.Marshal(msg)
		if _, err := node.serviceConn.WriteToUDP(data, item.inAddr); err != nil {
			logger.Warning("send weight update to input view item err:->", err, item.nodeId)
		}
	}

	node.subNo = 0
	return nil
}
