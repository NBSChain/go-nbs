package memership

import (
	"crypto/rand"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/thirdParty/gossip/pb"
	"github.com/golang/protobuf/proto"
	"math/big"
)

//TODO:: make sure this random is ok
func (node *MemManager) randomSelectItem() *ViewNode {
	count := len(node.PartialView)
	j := 0
	random, _ := rand.Int(rand.Reader, big.NewInt(int64(count)))
	logger.Debug("chose random in PartialView :->", random, count)
	for _, item := range node.PartialView {
		if j == int(random.Int64()) {
			return item
		} else {
			j++
		}
	}
	return nil
}

func (node *MemManager) removeFromView(item *ViewNode, views map[string]*ViewNode) {
	if item.outConn != nil {
		item.outConn.Close()
	}
	delete(views, item.nodeId)
	logger.Warning("remove node from view:->", item.nodeId)
}

func (node *MemManager) updateMyInProb(task *gossipTask) error {

	wei := task.msg.OVWeight
	item, ok := node.InputView[wei.NodeId]
	if !ok {
		return ItemNotFound
	}

	item.probability = wei.Weight
	return nil
}

func (node *MemManager) updateMyOutProb(task *gossipTask) error {
	wei := task.msg.IVWeight
	item, ok := node.PartialView[wei.NodeId]
	if !ok {
		return ItemNotFound
	}

	item.probability = wei.Weight

	return nil
}

func (node *MemManager) normalizeWeight(views map[string]*ViewNode) {
	var summerOut float64
	for _, item := range views {
		summerOut += item.probability
	}

	for _, item := range views {
		item.probability = item.probability / summerOut
	}
}

func (node *MemManager) updateProbability(task *gossipTask) error {

	node.normalizeWeight(node.PartialView)

	for _, item := range node.PartialView {
		msg := &pb.Gossip{
			MsgType: nbsnet.GspUpdateOVWei,
			OVWeight: &pb.WeightUpdate{
				NodeId: node.nodeID,
				Weight: item.probability,
			},
		}

		if err := node.send(item, msg); err != nil {
			logger.Warning("send weight update to partial view item err:->", err, item.nodeId)
		}
	}

	node.normalizeWeight(node.InputView)

	for _, item := range node.InputView {
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

func (node *MemManager) meanProb(views map[string]*ViewNode) float64 {

	pLen := len(views)
	if pLen == 0 {
		return 1.0
	}

	var summerOut float64
	for _, item := range views {
		summerOut += item.probability
	}

	return summerOut / float64(pLen)
}

//TIPS void to use lock
func (node *MemManager) viewNodeError(task *gossipTask) error {
	item := task.params.(*ViewNode)
	node.removeFromView(item, node.InputView)
	node.removeFromView(item, node.PartialView)
	return nil
}
