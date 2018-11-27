//+build !windows

package nat

import (
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"net"
	"time"
)

/************************************************************************
*
*			for linux unix darwin and so on
*
*************************************************************************/
func (tunnel *KATunnel) readKeepAlive() {

	for {
		buffer := make([]byte, utils.NormalReadBuffer)

		n, err := tunnel.kaConn.Read(buffer)
		if err != nil {
			logger.Warning("reading keep alive message failed:", err)
			continue
		}

		response := &net_pb.NatResponse{}
		if err := proto.Unmarshal(buffer[:n], response); err != nil {
			logger.Warning("keep alive response Unmarshal failed:", err)
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

func (tunnel *KATunnel) process(buffer []byte, peerAddr *net.UDPAddr) error {

	request := &net_pb.NatRequest{}
	proto.Unmarshal(buffer, request)

	switch request.MsgType {
	case net_pb.NatMsgType_DigOut, net_pb.NatMsgType_DigIn:
		tunnel.digSuccess(request.HoleMsg)
	}

	return nil
}

func (tunnel *KATunnel) listening() {

	for {
		buffer := make([]byte, utils.NormalReadBuffer)
		n, peerAddr, err := tunnel.serverHub.ReadFromUDP(buffer)
		if err != nil {
			logger.Warning("receiving port:", err, peerAddr)
			continue
		}

		if err := tunnel.process(buffer[:n], peerAddr); err != nil {
			logger.Warning("process nat response message failed")
			continue
		}
	}
}
