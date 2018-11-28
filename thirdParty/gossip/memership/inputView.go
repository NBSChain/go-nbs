package memership

import (
	"github.com/NBSChain/go-nbs/thirdParty/gossip/pb"
	"net"
	"time"
)

func (node *MemManager) pareKeepAlive(beat *pb.HeartBeat, addr *net.UDPAddr) {

	logger.Debug("gossip heart beat :", beat, addr)

	item, ok := node.inputView[beat.Sender]
	if !ok {
		logger.Warning("no such input view item:->", beat.Sender)
		return
	}

	item.updateTime = time.Now()

	switch beat.Type {
	case pb.MsgType_heartBeat:
	default:
		node.ctrlMsg(beat, addr)
	}
}

func (node *MemManager) ctrlMsg(beat *pb.HeartBeat, addr *net.UDPAddr) {

}
