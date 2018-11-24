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
		}
	}
}

func (tunnel *KATunnel) process(buffer []byte, peerAddr *net.UDPAddr) error {

	request := &net_pb.NatRequest{}
	proto.Unmarshal(buffer, request)

	switch request.MsgType {

	case net_pb.NatMsgType_Connect:
		go tunnel.natHoleStep4Answer(request.ConnReq)
	case net_pb.NatMsgType_DigOut, net_pb.NatMsgType_DigIn:
		if err := tunnel.DigSuccess(request.HoleMsg); err != nil {
			logger.Error(err)
		}
	case net_pb.NatMsgType_ReverseDig:
		tunnel.answerInvite(request.Invite)
	case net_pb.NatMsgType_ReverseDigACK:
		tunnel.setupReverseChan(request.Invite, peerAddr)
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

//TODO::
func (tunnel *KATunnel) connManage() {
}

//TODO::
func (tunnel *KATunnel) restoreNatChannel() {
}
