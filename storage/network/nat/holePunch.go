//+build !windows

package nat

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/denat"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/storage/network/shareport"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"net"
	"strconv"
	"time"
)

//udpConn inviter call first
func (tunnel *KATunnel) StartDigHole(lAddr, rAddr *nbsnet.NbsUdpAddr, connId string, toPort int) (chan *task, error) {

	connReq := &net_pb.NatConInvite{

		FromAddr: &net_pb.NbsAddr{
			NetworkId: lAddr.NetworkId,
			CanServer: lAddr.CanServe,
			PriIp:     lAddr.PriIp,
			PriPort:   lAddr.PriPort,
			NatIP:     lAddr.NatIp,
			NatPort:   lAddr.NatPort,
		},
		ToAddr: &net_pb.NbsAddr{
			NetworkId: rAddr.NetworkId,
			CanServer: rAddr.CanServe,
			PriIp:     rAddr.PriIp,
			PriPort:   rAddr.PriPort,
			NatIP:     rAddr.NatIp,
			NatPort:   rAddr.NatPort,
		},
		SessionId: connId,
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

	return nil, nil
}

//caller make a direct connection to peer's public address
func (tunnel *KATunnel) natHoleStep2Dig(sessionId string, rAddr *nbsnet.NbsUdpAddr) (*net.UDPConn, error) {

	holeMsg := &net_pb.NatRequest{
		MsgType: net_pb.NatMsgType_DigIn,
		HoleMsg: &net_pb.HoleDig{
			SessionId: sessionId,
		},
	}
	data, _ := proto.Marshal(holeMsg)

	connChan := make(chan *net.UDPConn)

	port := strconv.Itoa(int(rAddr.NatPort))
	pubAddr := net.JoinHostPort(rAddr.NatIp, port)
	pubConn, pubErr := shareport.DialUDP("udp4", tunnel.sharedAddr, pubAddr)

	priAddr := net.JoinHostPort(rAddr.PriIp, port)
	priConn, priErr := shareport.DialUDP("udp4", tunnel.sharedAddr, priAddr)

	if pubErr != nil && priErr != nil {
		return nil, fmt.Errorf("time out")
	}
	if pubConn != nil {
		return pubConn, nil
	}

	if priConn != nil {
		return priConn, nil
	}

	return nil, fmt.Errorf("time out")
}

func (tunnel *KATunnel) digIn(conn *net.UDPConn, digSig chan bool, data []byte) {

	for i := 0; i < HolePunchingTimeOut/2; i++ {

		if _, err := conn.Write(data); err != nil {
			logger.Info("dig in failed.")
		}

		buffer := make([]byte, utils.NormalReadBuffer)

		conn.Read(buffer)

		select {
		case success := <-digSig:
			if success {
				return
			}

		default:
			time.Sleep(time.Second)
		}
	}
}

//TIPS::get peer's addr info and make a connection.
func (tunnel *KATunnel) natHoleStep4Answer(response *net_pb.NatConInvite) {

	sessionId := response.SessionId
	digSig := make(chan bool)
	tunnel.digTask[sessionId] = digSig

	pubAddr := &net.UDPAddr{
		IP:   net.ParseIP(response.FromAddr.NatIP),
		Port: int(response.FromAddr.NatPort),
	}
	priAddr := &net.UDPAddr{
		IP:   net.ParseIP(response.FromAddr.PriIp),
		Port: int(response.FromAddr.NatPort),
	}

	holeMsg := &net_pb.NatRequest{
		MsgType: net_pb.NatMsgType_DigOut,
		HoleMsg: &net_pb.HoleDig{
			SessionId: sessionId,
		},
	}

	holeData, _ := proto.Marshal(holeMsg)

	for i := 0; i < HolePunchingTimeOut/2; i++ {

		_, errPub := tunnel.serverHub.WriteToUDP(holeData, pubAddr)
		_, errPri := tunnel.serverHub.WriteToUDP(holeData, priAddr)

		if errPri != nil && errPub != nil {
			logger.Error("failed to dig out :", errPub, errPri)
			break
		}

		select {
		case success := <-digSig:
			if success {
				break
			}
		default:
			time.Sleep(time.Second)
		}
	}

	delete(tunnel.digTask, sessionId)
}

//TIPS:: step1
func (tunnel *KATunnel) SendInvite(sessionId, pubIp string, toPort int) (chan *task, error) {
	req := &net_pb.NatRequest{
		MsgType: net_pb.NatMsgType_ReverseDig,
		Invite: &net_pb.ReverseInvite{
			SessionId: sessionId,
			PubIp:     pubIp,
			ToPort:    int32(toPort),
		},
	}

	reqData, _ := proto.Marshal(req)
	if _, err := tunnel.kaConn.Write(reqData); err != nil {
		return nil, err
	}

	connChan := make(chan *task)
	tunnel.inviteTask[sessionId] = connChan
	return connChan, nil
}

//TIPS:: step2
func (tunnel *KATunnel) answerInvite(invite *net_pb.ReverseInvite) {
	port := strconv.Itoa(int(invite.ToPort))
	natPort := strconv.Itoa(utils.GetConfig().NatChanSerPort)

	rAddr := net.JoinHostPort(invite.PubIp, natPort)
	conn, err := shareport.DialUDP("udp4", "0.0.0.0:"+port, rAddr)
	if err != nil {
		return
	}
	defer conn.Close()

	reqAck := &net_pb.NatRequest{
		MsgType: net_pb.NatMsgType_ReverseDigACK,
		InviteAck: &net_pb.ReverseInviteAck{
			SessionId: invite.SessionId,
		},
	}

	data, _ := proto.Marshal(reqAck)
	for i := 0; i < MaxDuplicateConfirm; i++ {
		conn.Write(data)
		time.Sleep(CommTTLTime)
	}
}

//TIPS:: step3
func (tunnel *KATunnel) setupReverseChan(invite *net_pb.ReverseInvite, addr *net.UDPAddr) {
	sessionId := invite.SessionId
	ch, ok := tunnel.inviteTask[sessionId]
	if !ok {
		return
	}

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		ch <- &task{
			udpConn: nil,
			err:     err,
		}
		return
	}

	ch <- &task{
		udpConn: conn,
		err:     nil,
	}

	delete(tunnel.inviteTask, sessionId)
}

/************************************************************************
*
*			server side
*
*************************************************************************/
//TIPS:: the server forward the connection invite to peer
func (nat *Manager) natHoleStep3ForwardInvite(request *net_pb.NatRequest, peerAddr *net.UDPAddr) error {

	req := request.ConnReq
	rawData, _ := proto.Marshal(request)

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

func (tunnel *KATunnel) DigSuccess(holeMsg *net_pb.HoleDig) error {

	sessionId := holeMsg.SessionId

	if digSig, ok := tunnel.digTask[sessionId]; ok {
		delete(tunnel.digTask, sessionId)
		digSig <- true
		close(digSig)
	}

	return nil
}
