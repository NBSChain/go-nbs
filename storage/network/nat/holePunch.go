//+build !windows

package nat

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/denat"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/storage/network/shareport"
	"github.com/golang/protobuf/proto"
	"net"
	"strconv"
	"time"
)

//UdpConn inviter call first
func (tunnel *KATunnel) StartDigHole(lAddr, rAddr *nbsnet.NbsUdpAddr, connId string, toPort int) (*ConnTask, error) {

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

	request := &net_pb.NatRequest{
		MsgType: net_pb.NatMsgType_Connect,
		ConnReq: connReq,
	}

	toItemData, err := proto.Marshal(request)
	if err != nil {
		return nil, err
	}

	host, port, _ := nbsnet.SplitHostPort(rAddr.NatServer)

	conn, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   net.ParseIP(host),
		Port: int(port),
	})
	if err != nil {
		logger.Warning("failed to send dig request:", err)
		return nil, err
	}
	defer conn.Close()

	if _, err := conn.Write(toItemData); err != nil {
		return nil, err
	}

	connChan := &ConnTask{
		Err: make(chan error),
	}

	tunnel.inviteTask[connId] = connChan
	logger.Info("Step 1:->notify the nat server:->")
	return connChan, nil
}

func (tunnel *KATunnel) DigInPubNet(lAddr, rAddr *nbsnet.NbsUdpAddr, task *ConnTask, sessionID string) (*net.UDPConn, error) {

	holeMsg := &net_pb.NatRequest{
		MsgType: net_pb.NatMsgType_DigIn,
		HoleMsg: &net_pb.HoleDig{
			SessionId: sessionID,
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

	for i := 0; i < HolePunchingTimeOut; i++ {

		if _, err := conn.Write(data); err != nil {
			logger.Error(err)
		}

		logger.Info("Step 4:-> I start to dig in:->", i, pubAddr, tunnel.sharedAddr)
		select {
		case err := <-task.Err:
			if err == nil {
				if task.UdpConn != nil { //private network succes
					logger.Info("Step 6-1:-> create connection from task.UdpConn:->")
					return task.UdpConn, nil
				} else { //this conn works
					logger.Info("Step 6-2:-> create connection from this conn:->")
					return conn, nil
				}
			} else {
				return nil, err
			}
		default:
			time.Sleep(time.Second)
		}
	}

	return nil, fmt.Errorf("time out")
}

/************************************************************************
*
*			server side
*
*************************************************************************/
//TIPS:: the server forward the connection invite to peer
func (nat *Manager) forwardDigRequest(req *net_pb.NatConnect, peerAddr *net.UDPAddr) error {

	res := &net_pb.NatResponse{
		MsgType: net_pb.NatMsgType_Connect,
		ConnRes: req,
	}
	rawData, _ := proto.Marshal(res)

	nat.cacheLock.Lock()
	defer nat.cacheLock.Unlock()

	toItem, ok := nat.cache[req.ToAddr.NetworkId]
	if !ok {
		if err := denat.GetDeNatSerIns().ProxyConnInvite(req); err != nil {
			return err
		}
	} else {
		if _, err := nat.sysNatServer.WriteToUDP(rawData, toItem.pubAddr); err != nil {
			return err
		}
	}

	logger.Info("Step 2:-> forward to peer:->", res)

	return nil
}

func (tunnel *KATunnel) digOut(req *net_pb.NatConnect) {

	sessionId := req.SessionId

	remNatAddr := &net.UDPAddr{
		IP:   net.ParseIP(req.FromAddr.NatIP),
		Port: int(req.FromAddr.NatPort),
	}

	holeMsg := &net_pb.NatRequest{
		MsgType: net_pb.NatMsgType_DigOut,
		HoleMsg: &net_pb.HoleDig{
			SessionId: sessionId,
		},
	}
	data, _ := proto.Marshal(holeMsg)

	task := &ProxyTask{
		sessionID: sessionId,
		toAddr: &net.UDPAddr{
			IP:   net.ParseIP(tunnel.natAddr.PriIp),
			Port: int(req.ToPort),
		},
		digResult: make(chan error),
	}

	tunnel.workLoad[sessionId] = task

	for i := 0; i < HolePunchingTimeOut; i++ {

		if _, err := tunnel.serverHub.WriteTo(data, remNatAddr); err != nil {
			logger.Error(err.Error())
		}

		logger.Info("Step 3:-> peer start to dig out:->", holeMsg)

		select {
		case err := <-task.digResult:
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

func (tunnel *KATunnel) digSuccess(msg *net_pb.HoleDig) {

	sessionId := msg.SessionId

	logger.Info("Step 5:-> dig success:->", msg)

	if cTask, ok := tunnel.inviteTask[sessionId]; ok {
		cTask.Err <- nil
		delete(tunnel.inviteTask, sessionId)
	}

	if pTask, ok := tunnel.workLoad[sessionId]; ok {
		pTask.digResult <- nil
		go pTask.relayData()
	}
}
