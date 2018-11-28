package memership

import (
	"github.com/NBSChain/go-nbs/thirdParty/gossip/pb"
	"github.com/gogo/protobuf/proto"
	"net"
	"time"
)

func (node *MemManager) handleKeepAlive(beat *pb.HeartBeat, addr *net.UDPAddr) {

	logger.Debug("gossip heart beat :", beat, addr)

	item, ok := node.inputView[beat.Sender]
	if !ok {
		logger.Warning("no such input view item:->", beat.Sender)
		return
	}

	item.updateTime = time.Now()

	payLoad := beat.Payload
	if payLoad == nil {
		return
	}

	msg := &pb.Gossip{}
	if err := proto.Unmarshal(payLoad, msg); err != nil {
		logger.Warning("keep alive payload err:->", err)
		return
	}

	node.ctrlMsg(msg, addr)
}

func (node *MemManager) ctrlMsg(msg *pb.Gossip, addr *net.UDPAddr) {
	switch msg.MessageType {
	case pb.MsgType_reqContract:
		req := msg.ContactReq
		sub := &newSub{
			nodeId: req.ApplierID,
			seq:    req.Seq,
			addr:   req.Applier,
		}

		node.indirectTheSubRequest(sub, int(req.TTL))
	}
}
