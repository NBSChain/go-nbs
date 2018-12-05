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
func (node *MemManager) choseRandomInPartialView() *peerNodeItem {
	count := len(node.partialView)
	j := 0
	random, _ := rand.Int(rand.Reader, big.NewInt(int64(count)))

	for _, item := range node.partialView {
		if j == int(random.Int64()) {
			return item
		} else {
			j++
		}
	}
	return nil
}

func (node *MemManager) PushOut(isHeartBeat bool, nodeId string, payLoad []byte) error {

	keepAlive := &pb.Gossip{
		MsgType: nbsnet.GspHeartBeat,
		HeartBeat: &pb.HeartBeat{
			Sender:  node.nodeID,
			SeqNo:   time.Now().Unix(),
			Payload: payLoad,
		},
	}

	data, _ := proto.Marshal(keepAlive)
	if nodeId != "" {
		item, ok := node.partialView[nodeId]
		if !ok {
			return PartialViewItemNotFound
		}

		return node.pushCtrlChan(item, data)
	}

	now := time.Now()
	for _, item := range node.partialView {

		if isHeartBeat && now.Sub(item.updateTime) < MemShipHeartBeat {
			continue
		}

		if err := node.pushCtrlChan(item, data); err != nil {
			logger.Warning(err)
			continue
		}
	}

	logger.Debug(" broad cast payload through control channel:->", len(node.partialView))

	return nil
}

func (node *MemManager) pushCtrlChan(item *peerNodeItem, data []byte) error {

	if _, err := item.ctrlConn.Write(data); err != nil {
		node.removePartialItem(item)
		return fmt.Errorf("push msg err:->nodeId:%s,err:%s", item.nodeId, err.Error())
	}

	item.updateTime = time.Now()
	return nil
}

func (node *MemManager) removePartialItem(item *peerNodeItem) {
	delete(node.partialView, item.nodeId) //TODO::make sure the timeout logic
	if err := item.ctrlConn.Close(); err != nil {
		logger.Warning(err)
	}
	logger.Warning("remove node from partial view:->", item.nodeId)
}
