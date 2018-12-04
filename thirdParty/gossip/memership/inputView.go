package memership

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/thirdParty/gossip/pb"
	"github.com/gogo/protobuf/proto"
	"net"
	"time"
)

func (node *MemManager) heartBeat(task *innerTask) error {
	beat := task.msg.HeartBeat
	addr := task.addr
	item, ok := node.inputView[beat.Sender]
	if !ok {
		err := fmt.Errorf("no such input view item:->%s", beat.Sender)
		logger.Warning(err)
		return err
	}

	item.updateTime = time.Now()

	payLoad := beat.Payload
	if payLoad == nil {
		return nil
	}

	msg := &pb.Gossip{}
	if err := proto.Unmarshal(payLoad, msg); err != nil {
		logger.Warning("heart beat payload err:->", err)
		return err
	}
	return node.ctrlMsg(msg, addr)
}

func (node *MemManager) ctrlMsg(msg *pb.Gossip, addr *net.UDPAddr) error {

	logger.Debug("heart beat with payload:", msg)

	switch msg.MsgType {
	case nbsnet.GspRegContact:
		req := msg.ContactReq
		sub := &newSub{
			nodeId: req.ApplierID,
			seq:    req.Seq,
			addr:   req.Applier,
		}

		node.indirectTheSubRequest(sub, int(req.TTL))
	}
	return nil
}
