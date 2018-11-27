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

	response := &net_pb.NatRequest{
		MsgType: net_pb.NatMsgType_Connect,
		ConnReq: connReq,
	}

	toItemData, err := proto.Marshal(response)
	if err != nil {
		return nil, err
	}

	if _, err := tunnel.kaConn.Write(toItemData); err != nil {
		return nil, err
	}

	connChan := &ConnTask{
		Err: make(chan error),
	}

	tunnel.inviteTask[connId] = connChan

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

		conn.Write(data)

		select {
		case err := <-task.Err:
			if err == nil {
				if task.UdpConn != nil { //private network succes
					return task.UdpConn, nil
				} else { //this conn works
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

func (tunnel *KATunnel) answerInvite(invite *net_pb.ReverseInvite) {

	myPort := strconv.Itoa(int(invite.ToPort))

	conn, err := shareport.DialUDP("udp4", "0.0.0.0:"+myPort,
		invite.PubIp+":"+invite.FromPort)
	if err != nil {
		logger.Errorf("failed to dial up peer to answer inviter:", err)
		return
	}
	defer conn.Close()

	req := &net_pb.NatRequest{
		MsgType: net_pb.NatMsgType_ReverseDigACK,
		InviteAck: &net_pb.ReverseInviteAck{
			SessionId: invite.SessionId,
		},
	}

	data, _ := proto.Marshal(req)
	if _, err := conn.Write(data); err != nil {
		logger.Errorf("failed to write answer to inviter:", err)
		return
	}

	logger.Debug("Step4: answer the invite:->", conn.LocalAddr().String(), invite, req)
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

		tunnel.serverHub.WriteTo(data, remNatAddr)

		select {
		case err := <-task.digResult:
			logger.Errorf("failed to dig out", err)
			return
		default:
			time.Sleep(time.Millisecond * 500)
		}
	}

	logger.Errorf("dig out time out")
	delete(tunnel.workLoad, sessionId)
}

func (tunnel *KATunnel) digSuccess(msg *net_pb.HoleDig) {
	sessionId := msg.SessionId

	if cTask, ok := tunnel.inviteTask[sessionId]; ok {
		cTask.Err <- nil
		delete(tunnel.inviteTask, sessionId)
	}

	if pTask, ok := tunnel.workLoad[sessionId]; ok {
		pTask.digResult <- nil
		go pTask.relayData()
	}
}
