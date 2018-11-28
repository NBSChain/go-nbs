package memership

import (
	"crypto/rand"
	"github.com/NBSChain/go-nbs/thirdParty/gossip/pb"
	"github.com/golang/protobuf/proto"
	"math/big"
	"time"
)

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

func (node *MemManager) keepAlive() {

	keepAlive := &pb.Gossip{
		MessageType: pb.MsgType_heartBeat,
		HeartBeat: &pb.HeartBeat{
			Type:    pb.MsgType_heartBeat,
			Sender:  node.nodeID,
			SeqNo:   time.Now().Unix(),
			Payload: nil,
		},
	}

	data, _ := proto.Marshal(keepAlive)
	now := time.Now()

	for nodeId, item := range node.partialView {

		if now.Sub(item.updateTime) < MemberShipKeepAlive {
			continue
		}

		if _, err := item.conn.Write(data); err != nil {
			logger.Warning("node in partial view is expired:->", nodeId, err)
			delete(node.partialView, nodeId) //TODO::make sure the timeout logic
			item.conn.Close()
			continue
		}

		item.updateTime = now
	}
}

func (node *MemManager) keepAliveWithData(Typ pb.MsgType, payLoad []byte, nodeId string) {

	keepAlive := &pb.Gossip{
		MessageType: pb.MsgType_heartBeat,
		HeartBeat: &pb.HeartBeat{
			Type:    Typ,
			Sender:  node.nodeID,
			SeqNo:   time.Now().Unix(),
			Payload: payLoad,
		},
	}

	data, _ := proto.Marshal(keepAlive)

	if nodeId != "" {
		item, ok := node.partialView[nodeId]
		if !ok {
			logger.Error("can't find the target peer node.")
			return
		}

		if _, err := item.conn.Write(data); err != nil {
			logger.Warning("node in partial view is expired:->", nodeId, err)
			delete(node.partialView, nodeId) //TODO::make sure the timeout logic
			item.conn.Close()
			return
		}

		item.updateTime = time.Now()

		return
	}

	//TIPS::broadcast
	for nodeId, item := range node.partialView {
		if _, err := item.conn.Write(data); err != nil {
			logger.Warning("node in partial view is expired:->", nodeId, err)
			delete(node.partialView, nodeId) //TODO::make sure the timeout logic
			item.conn.Close()
			continue
		}

		item.updateTime = time.Now()
	}

}
