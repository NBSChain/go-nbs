package nat

import (
	"fmt"
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
		}
	}
}

func (tunnel *KATunnel) DigInPubNet(lAddr, rAddr *nbsnet.NbsUdpAddr, task *ConnTask, sessionID string) (*net.UDPConn, error) {

	holeMsg := &net_pb.NatRequest{
		MsgType: net_pb.NatMsgType_DigIn,
		HoleMsg: &net_pb.HoleDig{
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

	for i := 0; i < TryDigHoleTimes; i++ {

		logger.Info("Step 4:-> I start to dig in:->", i, pubAddr, tunnel.sharedAddr)

		if _, err := conn.Write(data); err != nil {
			logger.Error(err)
		}
		conn.SetReadDeadline(time.Now().Add(time.Second))
		buffer := make([]byte, utils.NormalReadBuffer)
		if _, err := conn.Read(buffer); err != nil {
			logger.Warning("dig read failed:->", err)
		} else {
			msg := &net_pb.NatResponse{}
			proto.Unmarshal(buffer, msg)
			logger.Info("read dig res:->", msg)
		}

		select {
		case err := <-task.err:
			if err == nil {
				if task.udpConn != nil { //private network succes
					logger.Info("Step 6-1:-> create connection from task.udpConn:->")
					return task.udpConn, nil
				} else { //this conn works
					logger.Info("Step 6-2:-> create connection from this conn:->")
					return conn, nil
				}
			} else {
				return nil, err
			}
		default:
			logger.Debug("retry again......")
		}
	}

	return nil, fmt.Errorf("time out")
}
