//+build !windows

package nat

import (
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/storage/network/shareport"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"net"
	"strconv"
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
			//TODO::recovery and broadcast the new information
			logger.Warning("reading keep alive message failed:", err)
			continue
		}

		if err := tunnel.process(buffer[:n]); err != nil {
			continue
		}
	}
}

func (tunnel *KATunnel) process(buffer []byte) error {

	response := &net_pb.NatMsg{}
	if err := proto.Unmarshal(buffer, response); err != nil {
		logger.Warning("keep alive response Unmarshal failed:", err)
		return err
	}

	logger.Debug("keep alive:->", response)

	switch response.T {
	case nbsnet.NatKeepAlive:
		tunnel.refreshNatInfo(response.V)
	case nbsnet.NatReversDig:
		tunnel.answerInvite(response.V)
	case nbsnet.NatConnect:
		tunnel.digOut(response.V)
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

		request := &net_pb.NatMsg{}
		if err := proto.Unmarshal(buffer[:n], request); err != nil {
			logger.Warning("parse message failed", err)
			continue
		}

		logger.Debug("listen connection:", request)

		switch request.T {
		default:
			logger.Warning("unknown msg for linux/unix/bsd systems.")
		}
	}
}

//TIPS::unix/bsd/linux need to read the response from same connection
func (tunnel *KATunnel) DigInPubNet(lAddr, rAddr *nbsnet.NbsUdpAddr, task *ConnTask, sessionID string) {

	DigMsg := &net_pb.HoleDig{
		SessionId:   sessionID,
		NetworkType: FromPubNet,
	}
	digData, _ := proto.Marshal(DigMsg)
	holeMsg := &net_pb.NatMsg{
		T: nbsnet.NatDigIn,
		L: int32(len(digData)),
		V: digData,
	}
	data, _ := proto.Marshal(holeMsg)

	port := strconv.Itoa(int(rAddr.NatPort))
	pubAddr := net.JoinHostPort(rAddr.NatIp, port)
	conn, err := shareport.DialUDP("udp4", tunnel.sharedAddr, pubAddr)
	if err != nil {
		logger.Warning("dig hole in pub network failed", err)
		task.err <- err
		return
	}

	go tunnel.waitDigResponse(task, conn)

	logger.Info("Step 4:-> I start to dig in:->")
	tunnel.digDig(data, conn, task)
}
