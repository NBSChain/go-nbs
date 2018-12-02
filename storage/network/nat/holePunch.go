package nat

import "C"
import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/golang/protobuf/proto"
	"net"
	"time"
)

//udpConn inviter call first
func (tunnel *KATunnel) StartDigHole(lAddr, rAddr *nbsnet.NbsUdpAddr, connId string, toPort int) error {

	connReq := &net_pb.NatConnect{

		FromAddr: &net_pb.NbsAddr{
			NetworkId: lAddr.NetworkId,
			CanServer: lAddr.CanServe,
			PriIp:     lAddr.PriIp,
			NatIP:     lAddr.NatIp,
			NatPort:   lAddr.NatPort,
		},
		ToAddr: &net_pb.NbsAddr{
			NetworkId: rAddr.NetworkId,
			CanServer: rAddr.CanServe,
			PriIp:     rAddr.PriIp,
			NatIP:     rAddr.NatIp,
			NatPort:   rAddr.NatPort,
		},
		SessionId: connId,
		ToPort:    int32(toPort),
	}

	connData, _ := proto.Marshal(connReq)
	request := &net_pb.NatMsg{
		Typ:     nbsnet.NatConnect,
		Len:     int32(len(connData)),
		PayLoad: connData,
	}

	toItemData, err := proto.Marshal(request)
	if err != nil {
		return err
	}

	host, port, _ := nbsnet.SplitHostPort(rAddr.NatServer)

	conn, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   net.ParseIP(host),
		Port: int(port),
	})
	if err != nil {
		logger.Warning("failed to send dig request:", err)
		return err
	}
	defer conn.Close()

	if _, err := conn.Write(toItemData); err != nil {
		return err
	}
	logger.Info("Step 1:->notify the nat server:->", rAddr.NatServer)
	return nil
}

/************************************************************************
*
*			server side
*
*************************************************************************/
//TIPS:: the server forward the connection invite to peer
func (nat *Manager) forwardDigRequest(data []byte, peerAddr *net.UDPAddr) error {
	req := &net_pb.NatConnect{}
	if err := proto.Unmarshal(data, req); err != nil {
		return err
	}
	res := &net_pb.NatMsg{
		Typ:     nbsnet.NatConnect,
		Len:     int32(len(data)),
		PayLoad: data,
	}
	rawData, _ := proto.Marshal(res)

	nat.cacheLock.Lock()
	defer nat.cacheLock.Unlock()

	toItem, ok := nat.cache[req.ToAddr.NetworkId]
	if !ok {
		return fmt.Errorf("this item is no more in my cache")
	} else {
		if _, err := nat.sysNatServer.WriteToUDP(rawData, toItem.pubAddr); err != nil {
			return err
		}
	}

	logger.Info("Step 2:-> forward to peer:->", res)

	return nil
}

//TODO:: refactor
func (tunnel *KATunnel) digOut(data []byte) {

	req := &net_pb.NatConnect{}
	if err := proto.Unmarshal(data, req); err != nil {
		logger.Warning("unmarshal connect data failed:->", err)
		return
	}

	sessionId := req.SessionId
	if _, ok := tunnel.workLoad[sessionId]; ok {
		logger.Info("duplicate connect require")
		return
	}

	remNatAddr := &net.UDPAddr{
		IP:   net.ParseIP(req.FromAddr.NatIP),
		Port: int(req.FromAddr.NatPort),
	}
	DigMsg := &net_pb.HoleDig{
		SessionId:   sessionId,
		NetworkType: ToPubNet,
	}
	DigData, _ := proto.Marshal(DigMsg)
	holeMsg := &net_pb.NatMsg{
		Typ:     nbsnet.NatDigOut,
		Len:     int32(len(DigData)),
		PayLoad: DigData,
	}

	msgData, _ := proto.Marshal(holeMsg)

	task := &ProxyTask{
		sessionID: sessionId,
		toAddr: &net.UDPAddr{
			IP:   net.ParseIP(tunnel.natAddr.PriIp),
			Port: int(req.ToPort),
		},
		err: make(chan error),
	}

	tunnel.workLoad[sessionId] = task

	for i := 0; i < TryDigHoleTimes; i++ {

		if _, err := tunnel.serverHub.WriteTo(msgData, remNatAddr); err != nil {
			logger.Error(err.Error())
		}

		logger.Info("Step 3:-> peer start to dig out:->", holeMsg)

		select {
		case err := <-task.err:
			if err != nil {
				logger.Error("failed to dig out", err)
			} else {
				logger.Info("Step 7:->peer packet in :->")
			}
			return
		default:
			time.Sleep(time.Millisecond * 500)
		}
	}

	logger.Error("dig out time out")
	delete(tunnel.workLoad, sessionId)
}

func (tunnel *KATunnel) digSuccessRes(data []byte, peerAddr *net.UDPAddr) {

	res := &net_pb.NatMsg{
		Typ:     nbsnet.NatDigSuccess,
		Len:     int32(len(data)),
		PayLoad: data,
	}

	resData, _ := proto.Marshal(res)
	if _, err := tunnel.serverHub.WriteTo(resData, peerAddr); err != nil {
		logger.Warning("failed to response the dig confirm.")
	}
}
