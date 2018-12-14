//+build !windows

package nat

import (
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"net"
)

/************************************************************************
*
*			for linux unix darwin and so on
*
*************************************************************************/
func (tunnel *KATunnel) readKeepAlive() {

	for {
		buffer := make([]byte, utils.NormalReadBuffer)

		n, peerAddr, err := tunnel.kaConn.ReadFromUDP(buffer)
		if err != nil {
			//TODO::recovery and broadcast the new information
			logger.Warning("reading keep alive message failed:", err)
			continue
		}

		msg := &net_pb.NatMsg{}
		if err := proto.Unmarshal(buffer[:n], msg); err != nil {
			logger.Warning("keep alive msg Unmarshal failed:", err)
			continue
		}

		logger.Debug("KA tunnel receive message:->", msg, peerAddr)

		if err := tunnel.process(msg, peerAddr); err != nil {
			continue
		}
	}
}

func (tunnel *KATunnel) process(msg *net_pb.NatMsg, peerAddr *net.UDPAddr) error {

	switch msg.Typ {
	case nbsnet.NatKeepAlive:
		tunnel.refreshNatInfo(msg.KeepAlive)
	case nbsnet.NatReversInvite:
		tunnel.answerInvite(msg.ReverseInvite)
	case nbsnet.NatDigApply:
		tunnel.digOut(msg.DigApply)
	case nbsnet.NatDigConfirm:
		tunnel.makeAHole(msg.DigConfirm)
	}

	return nil
}

func (tunnel *KATunnel) listening() {
	//logger.Debug("no need to listen for linux/unix...")
	for {
		buffer := make([]byte, utils.NormalReadBuffer)
		n, peerAddr, err := tunnel.serverHub.ReadFromUDP(buffer)
		if err != nil {
			logger.Warning("server hub receiving port:", err, peerAddr)
			continue
		}

		msg := &net_pb.NatMsg{}
		if err := proto.Unmarshal(buffer[:n], msg); err != nil {
			logger.Warning("server hub parse message failed", err)
			continue
		}
		logger.Info("server hub receive msg:->", msg)
	}
}
