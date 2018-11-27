package nat

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"time"
)

func (tunnel *KATunnel) readKeepAlive() {
	logger.Info("windows-> no need to get data from keep alive connection")
	tunnel.server()
}

func (tunnel *KATunnel) listening() {

	tunnel.server()
}

func (tunnel *KATunnel) server() {

	for {
		responseData := make([]byte, utils.NormalReadBuffer)
		hasRead, peerAddr, err := tunnel.serverHub.ReadFrom(responseData)
		if err != nil {
			logger.Warning("receiving port:", err, peerAddr)
			continue
		}

		response := &net_pb.NatResponse{}
		if err := proto.Unmarshal(responseData[:hasRead], response); err != nil {
			fmt.Println("failed to unmarshal nat response data", err)
			continue
		}

		switch response.MsgType {
		case net_pb.NatMsgType_KeepAlive:
			tunnel.updateTime = time.Now()
			tunnel.natAddr.NatIp = response.KeepAlive.PubIP
			tunnel.natAddr.NatPort = response.KeepAlive.PubPort
		case net_pb.NatMsgType_ReverseDig:
			tunnel.answerInvite(response.Invite)
		case net_pb.NatMsgType_Connect:
			tunnel.digOut(response.ConnRes)
		}
	}
}
