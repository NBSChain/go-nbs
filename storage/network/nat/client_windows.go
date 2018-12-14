package nat

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
)

func (tunnel *KATunnel) waitServerCmd() {

	for {
		responseData := make([]byte, utils.NormalReadBuffer)
		hasRead, peerAddr, err := tunnel.serverHub.ReadFromUDP(responseData)
		if err != nil {
			if tunnel.errNo++; tunnel.errNo > ErrNoBeforeRetry {
				logger.Warning("too many reading error:->")
				tunnel.reSetupChannel()
				return
			}
			logger.Warning("receiving port:", err, peerAddr)
			continue
		}

		response := &net_pb.NatMsg{}
		if err := proto.Unmarshal(responseData[:hasRead], response); err != nil {
			fmt.Println("failed to unmarshal nat response data", err)
			continue
		}

		logger.Debug("server hub connection :", peerAddr, response)

		switch response.Typ {
		case nbsnet.NatKeepAlive:
			tunnel.refreshNatInfo(response.KeepAlive)
		case nbsnet.NatReversInvite:
			tunnel.answerInvite(response.ReverseInvite)
		case nbsnet.NatDigApply:
			tunnel.digOut(response.DigApply)
		case nbsnet.NatDigConfirm:
			tunnel.makeAHole(response.DigConfirm)
		}
	}
}
