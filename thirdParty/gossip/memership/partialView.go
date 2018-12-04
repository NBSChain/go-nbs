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

func (node *MemManager) sendHeartBeat() {

	keepAlive := &pb.Gossip{
		MsgType: nbsnet.GspHeartBeat,
		HeartBeat: &pb.HeartBeat{
			Sender:  node.nodeID,
			SeqNo:   time.Now().Unix(),
			Payload: nil,
		},
	}

	data, _ := proto.Marshal(keepAlive)
	now := time.Now()

	for nodeId, item := range node.partialView {

		if now.Sub(item.updateTime) < MemShipHeartBeat {
			continue
		}

		if _, err := item.ctrlConn.Write(data); err != nil {
			logger.Warning("node in partial view is expired:->", nodeId, err)
			delete(node.partialView, nodeId) //TODO::make sure the timeout logic
			if err := item.ctrlConn.Close(); err != nil {
				logger.Warning(err)
			}
			continue
		}

		item.updateTime = now

		logger.Debug("send empty heart beat :->", nodeId, item.ctrlConn.String())
	}
}

func (node *MemManager) sendHBWithPayLoad(nodeId string, payLoad []byte) error {

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
			logger.Error("can't find the target peer node.")
			return fmt.Errorf("can't find the target peer node")
		}

		if _, err := item.ctrlConn.Write(data); err != nil {
			logger.Warning("node in partial view is expired:->", nodeId, err)
			delete(node.partialView, nodeId) //TODO::make sure the timeout logic
			if err := item.ctrlConn.Close(); err != nil {
				logger.Warning(err)
			}
			return fmt.Errorf("can't find the target peer node:->nodeId:%s,err:%s", nodeId, err.Error())
		}

		item.updateTime = time.Now()

		logger.Debug("payload heart beat :->", nodeId, item.addr)
		return nil
	}

	//TIPS::broadcast
	for nodeId, item := range node.partialView {
		if _, err := item.ctrlConn.Write(data); err != nil {
			logger.Warning("node in partial view is expired:->", nodeId, err)
			delete(node.partialView, nodeId) //TODO::make sure the timeout logic
			if err := item.ctrlConn.Close(); err != nil {
				logger.Warning(err)
			}
			continue
		}

		item.updateTime = time.Now()
	}
	logger.Debug(" broad cast heart beat with payload :->", nodeId, len(node.partialView))
	return nil
}
