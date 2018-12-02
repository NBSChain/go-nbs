package nat

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
)

func (tunnel *KATunnel) readKeepAlive() {
	logger.Info("windows-> no need to get data from keep alive connection")
}

func (tunnel *KATunnel) listening() {

	for {
		responseData := make([]byte, utils.NormalReadBuffer)
		hasRead, peerAddr, err := tunnel.serverHub.ReadFromUDP(responseData)
		if err != nil {
			logger.Warning("receiving port:", err, peerAddr)
			continue
		}

		response := &net_pb.NatMsg{}
		if err := proto.Unmarshal(responseData[:hasRead], response); err != nil {
			fmt.Println("failed to unmarshal nat response data", err)
			continue
		}

		logger.Debug("server hub connection :", response, peerAddr)

		switch response.Typ {
		case nbsnet.NatKeepAlive:
			tunnel.refreshNatInfo(response.PayLoad)
		case nbsnet.NatReversDig:
			tunnel.answerInvite(response.PayLoad)
		case nbsnet.NatConnect:
			tunnel.digOut(response.PayLoad)
		case nbsnet.NatConnectACK:
			tunnel.makeAHole(response.PayLoad)
		}
	}
}
