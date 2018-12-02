package nat

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"net"
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
		case nbsnet.NatDigIn, utils.NatDigOut:
			tunnel.digSuccess(response.PayLoad, peerAddr)
		case nbsnet.NatDigSuccess:
			logger.Debug("dig success:->", peerAddr)
		}
	}
}

func (tunnel *KATunnel) DigInPubNet(lAddr, rAddr *nbsnet.NbsUdpAddr, task *ConnTask, sessionID string) (*net.UDPConn, error) {

	//holeMsg := &net_pb.NatManage{
	//	MsgType: utils.NatDigIn,
	//	DigMsg: &net_pb.HoleDig{
	//		SessionId:   sessionID,
	//		NetworkType: FromPubNet,
	//	},
	//}
	//data, _ := proto.Marshal(holeMsg)
	//
	//port := strconv.Itoa(int(rAddr.NatPort))
	//pubAddr := net.JoinHostPort(rAddr.NatIp, port)
	//conn, err := shareport.DialUDP("udp4", tunnel.sharedAddr, pubAddr)
	//if err != nil {
	//	logger.Warning("dig hole in pub network failed", err)
	//	return nil, err
	//}

	return nil, fmt.Errorf("time out")
}

func (tunnel *KATunnel) digSuccess(data []byte, peerAddr *net.UDPAddr) {
	msg := &net_pb.HoleDig{}
	if err := proto.Unmarshal(data, msg); err != nil {
		logger.Warning("dig success unmarshal err:->", err)
		return
	}

	res := &net_pb.NatMsg{
		Typ:     nbsnet.NatDigSuccess,
		Len:     int32(len(data)),
		PayLoad: data,
	}

	resData, _ := proto.Marshal(res)

	if _, err := tunnel.serverHub.WriteTo(resData, peerAddr); err != nil {
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
