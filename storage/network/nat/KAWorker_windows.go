package nat

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/storage/network/shareport"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"net"
	"strconv"
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

		response := &net_pb.NatResponse{}
		if err := proto.Unmarshal(responseData[:hasRead], response); err != nil {
			fmt.Println("failed to unmarshal nat response data", err)
			continue
		}

		logger.Debug("server hub connection :", response, peerAddr)

		switch response.MsgType {
		case net_pb.NatMsgType_KeepAlive:
			tunnel.refreshNatInfo(response.KeepAlive)
		case net_pb.NatMsgType_ReverseDig:
			tunnel.answerInvite(response.Invite)
		case net_pb.NatMsgType_Connect:
			tunnel.digOut(response.ConnRes)
		case net_pb.NatMsgType_DigIn, net_pb.NatMsgType_DigOut:
			tunnel.digSuccess(response.HoleMsg, peerAddr)
		case net_pb.NatMsgType_DigSuccess:
			logger.Debug("dig success:->", peerAddr)
		case net_pb.NatMsgType_BoardOnHole:
			tunnel.unpackMsg(response.HolePayLoad)
		}
	}
}

func (tunnel *KATunnel) DigInPubNet(lAddr, rAddr *nbsnet.NbsUdpAddr, task *ConnTask, sessionID string) (*net.UDPConn, error) {

	holeMsg := &net_pb.NatRequest{
		MsgType: net_pb.NatMsgType_DigIn,
		DigMsg: &net_pb.HoleDig{
			SessionId:   sessionID,
			NetworkType: FromPubNet,
		},
	}
	data, _ := proto.Marshal(holeMsg)

	port := strconv.Itoa(int(rAddr.NatPort))
	pubAddr := net.JoinHostPort(rAddr.NatIp, port)
	conn, err := shareport.DialUDP("udp4", tunnel.sharedAddr, pubAddr)
	if err != nil {
		logger.Warning("dig hole in pub network failed", err)
		return nil, err
	}

	return nil, fmt.Errorf("time out")
}

func (tunnel *KATunnel) digSuccess(msg *net_pb.HoleDig, peerAddr *net.UDPAddr) {

	res := &net_pb.NatResponse{
		MsgType: net_pb.NatMsgType_DigSuccess,
		HoleMsg: msg,
	}

	data, _ := proto.Marshal(res)

	if _, err := tunnel.serverHub.WriteTo(data, peerAddr); err != nil {
		logger.Warning("failed to response the dig confirm.")
		return
	}
	sid := msg.SessionId
	logger.Info("Step 5:-> dig success:->", sid)

	if pTask, ok := tunnel.workLoad[sid]; ok {
		pTask.err <- nil
		go pTask.relayData()
	} else {
		tunnel.workLoad[sid] = &ProxyTask{
			sessionID: sid,
			toAddr:    peerAddr,
		}
	}
}
