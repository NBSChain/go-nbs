package nat

import (
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
)

/************************************************************************
*
*			for linux unix darwin and so on
*
*************************************************************************/
func (tunnel *Client) readCmd() {

	for {
		buffer := make([]byte, utils.NormalReadBuffer)

		n, peerAddr, err := tunnel.conn.ReadFromUDP(buffer)
		if err != nil {
			if tunnel.errNo++; tunnel.errNo > ErrNoBeforeRetry {
				logger.Warning("too many reading error:->")
				tunnel.reSetupChannel()
				return
			}
			logger.Warning("reading keep alive message failed:", err)
			continue
		}

		msg := &net_pb.NatMsg{}
		if err := proto.Unmarshal(buffer[:n], msg); err != nil {
			logger.Warning("keep alive msg Unmarshal failed:", err)
			continue
		}

		logger.Debug("KA tunnel receive message:->", msg, peerAddr)

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

		select {
		case <-tunnel.ctx.Done():
			logger.Info("exit reading thread cause's of context close")
			return
		default:
		}
	}
}
