package network

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/storage/network/shareport"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"net"
	"strconv"
	"time"
)

const (
	HolePunchTimeOut = 6 * time.Second
)

type connTask struct {
	err         error
	locPort     string
	udpConn     chan *net.UDPConn
	portCapConn *net.UDPConn
}

//TODO:: data race
func (task *connTask) Close() {
	if task.portCapConn != nil {
		_ = task.portCapConn.Close()
	}
}

/************************************************************************
*
*			private functions
*
*************************************************************************/

func (network *nbsNetwork) makeDirectConn(lAddr, rAddr *nbsnet.NbsUdpAddr, toPort int) (*nbsnet.NbsUdpConn, error) {

	if rAddr == nil {
		return nil, fmt.Errorf("remote address can't be nil")
	}

	var sessionID = lAddr.NetworkId + ConnectionSeparator + rAddr.NetworkId

	rUdpAddr := &net.UDPAddr{
		Port: toPort,
		IP:   net.ParseIP(rAddr.PubIp),
	}

	lUdpAddr := &net.UDPAddr{
		IP: net.ParseIP(lAddr.PriIp),
	}

	c, err := net.DialUDP("udp4", lUdpAddr, rUdpAddr)
	if err != nil {
		return nil, err
	}

	logger.Debug("Step6:make direct connection:->", c.LocalAddr().String(), c.RemoteAddr().String())
	natAddr := network.natClient.NatAddr
	host, _, _ := nbsnet.SplitHostPort(c.LocalAddr().String())
	conn := &nbsnet.NbsUdpConn{
		RealConn:  c,
		CType:     nbsnet.CTypeNormal,
		SessionID: sessionID,
		LocAddr: &nbsnet.NbsUdpAddr{
			NetworkId: network.networkId,
			CanServe:  natAddr.CanServe,
			NatServer: natAddr.NatServer,
			NatIp:     natAddr.NatIp,
			PubIp:     natAddr.PubIp,
			NatPort:   natAddr.NatPort,
			PriIp:     host,
		},
	}
	return conn, nil
}

func (network *nbsNetwork) invitePeerBehindNat(lAddr, rAddr *nbsnet.NbsUdpAddr,
	connId string, toPort int) (*net.UDPConn, error) {

	conn, err := shareport.DialUDP("udp4", "", rAddr.NatServer)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localHost := conn.LocalAddr().String()
	_, fromPort, _ := net.SplitHostPort(localHost)

	req := &net_pb.NatMsg{
		Typ: nbsnet.NatReversInvite,
		ReverseInvite: &net_pb.ReverseInvite{
			SessionId: connId,
			PubIp:     lAddr.PubIp,
			ToPort:    int32(toPort),
			PeerId:    rAddr.NetworkId,
			FromPort:  fromPort,
		},
	}
	reqData, _ := proto.Marshal(req)
	if _, err := conn.Write(reqData); err != nil {
		return nil, err
	}

	connChan := &connTask{
		udpConn: make(chan *net.UDPConn),
	}

	logger.Debug("Step1: notify applier's nat server:", req)

	go network.waitInviteAnswer(localHost, connId, connChan)

	select {
	case conn := <-connChan.udpConn:
		if conn != nil {
			return conn, nil
		} else {
			return nil, connChan.err
		}
	case <-time.After(HolePunchTimeOut):
		return nil, fmt.Errorf("time out")
	}
}

func (network *nbsNetwork) punchANatHole(lAddr, rAddr *nbsnet.NbsUdpAddr,
	connId string, toPort int) (*net.UDPConn, nbsnet.ConnType, error) {

	priConnTask := &connTask{
		udpConn: make(chan *net.UDPConn),
	}
	defer priConnTask.Close()

	go network.directDialInPriNet(lAddr, rAddr, priConnTask, toPort, connId)

	pubConnTask := &connTask{
		udpConn: make(chan *net.UDPConn),
	}
	defer pubConnTask.Close()
	go network.noticeNatAndWait(lAddr, rAddr, connId, toPort, pubConnTask)

	var pubFail, priFail bool
	for i := 2; i > 0; i-- {
		select {
		case conn := <-priConnTask.udpConn:
			logger.Debug("hole punch step1-8 dig direct in private network finished:->", priConnTask.err)
			if conn != nil {
				return conn, nbsnet.CTypeNormal, nil
			} else {
				priFail = true
				if pubFail {
					return nil, 0, priConnTask.err
				}
			}
		case conn := <-pubConnTask.udpConn:
			logger.Debug("hole punch step2-x dig in public network finished:->", pubConnTask.err)
			if conn != nil {
				return conn, nbsnet.CTypeNatDuplex, nil
			} else {
				pubFail = true
				if priFail {
					return nil, 0, pubConnTask.err
				}
			}
		case <-time.After(HolePunchTimeOut / 2):
			return nil, 0, fmt.Errorf("time out")
		}
	}
	return nil, 0, fmt.Errorf("time out")
}

func (network *nbsNetwork) waitInviteAnswer(host, sessionID string, task *connTask) {

	lisConn, err := shareport.ListenUDP("udp4", host)
	if err != nil {
		task.err = err
		task.udpConn <- nil
		return
	}
	defer lisConn.Close()

	logger.Debug("Step2: wait the answer:", host)

	buffer := make([]byte, utils.NormalReadBuffer)
	n, peerAddr, err := lisConn.ReadFromUDP(buffer)
	if err != nil {
		task.err = err
		task.udpConn <- nil
		return
	}

	res := &net_pb.NatMsg{}
	if err := proto.Unmarshal(buffer[:n], res); err != nil {
		task.err = err
		task.udpConn <- nil
		return
	}

	if res.Typ != nbsnet.NatReversInviteAck {
		task.err = fmt.Errorf("didn't get the answer")
		task.udpConn <- nil
		return
	}

	conn, err := shareport.DialUDP("udp4", host, peerAddr.String())
	if err != nil {
		task.err = err
		task.udpConn <- nil
		return
	}
	task.udpConn <- conn
	task.err = err

	logger.Debug("Step5: get answer and make a connection:->", host, peerAddr.String())
}

func (network *nbsNetwork) directDialInPriNet(lAddr, rAddr *nbsnet.NbsUdpAddr, task *connTask, toPort int, sessionID string) {

	conn, err := net.DialUDP("udp4", &net.UDPAddr{
		IP: net.ParseIP(lAddr.PriIp),
	}, &net.UDPAddr{
		IP:   net.ParseIP(rAddr.PriIp),
		Port: toPort,
	})
	if err != nil {
		logger.Warning("Step 1-1:can't dial by private network.", err)
		task.err = err
		task.udpConn <- nil
		return
	}
	conStr := "[" + conn.LocalAddr().String() + "]-->[" + conn.RemoteAddr().String() + "]"

	logger.Debug("hole punch step1-2 start in private network:->", conStr)

	holeMsg := &net_pb.NatMsg{
		Typ: nbsnet.NatPriDigSyn,
		Seq: time.Now().Unix(),
	}
	data, _ := proto.Marshal(holeMsg)
	if _, err := conn.Write(data); err != nil {
		logger.Error("Step 1-3:private network dig dig failed:->", err)
		task.err = err
		task.udpConn <- nil
		return
	}

	if err := conn.SetReadDeadline(time.Now().Add(HolePunchTimeOut / 2)); err != nil {
		task.err = err
		task.udpConn <- nil
		return
	}

	buffer := make([]byte, utils.NormalReadBuffer)
	n, err := conn.Read(buffer)
	if err != nil {
		logger.Error("Step 1-5:private network reading dig result err:->", err)
		task.err = err
		task.udpConn <- nil
		return
	}
	resMsg := &net_pb.NatMsg{}
	if err := proto.Unmarshal(buffer[:n], resMsg); err != nil {
		logger.Info("Step 1-4:->dig in private network Unmarshal err:->", err)
		task.err = err
		task.udpConn <- nil
		return
	}

	if resMsg.Typ != nbsnet.NatPriDigAck || resMsg.Seq != holeMsg.Seq+1 {
		task.err = fmt.Errorf("wrong ack package")
		task.udpConn <- nil
		return
	}

	if err := conn.SetDeadline(time.Time{}); err != nil {
		task.err = err
		task.udpConn <- nil
	}

	task.err = nil
	task.udpConn <- conn
}

func (network *nbsNetwork) answerInvite(params interface{}) error {
	invite, ok := params.(*net_pb.ReverseInvite)
	if !ok {
		return CmdTaskErr
	}
	myPort := strconv.Itoa(int(invite.ToPort))

	conn, err := shareport.DialUDP("udp4", "0.0.0.0:"+myPort,
		invite.PubIp+":"+invite.FromPort)
	if err != nil {
		logger.Errorf("failed to dial up peer to answer inviter:", err)
		return err
	}
	defer conn.Close()

	req := &net_pb.NatMsg{
		Typ: nbsnet.NatReversInviteAck,
		ReverseInviteAck: &net_pb.ReverseInviteAck{
			SessionId: invite.SessionId,
		},
	}

	reqData, _ := proto.Marshal(req)
	logger.Debug("Step4: answer the invite:->", conn.LocalAddr().String(), invite.SessionId)

	for i := 0; i < 3; i++ {
		if _, err := conn.Write(reqData); err != nil {
			logger.Errorf("failed to write answer to inviter:", err)
			return err
		}
	}
	return nil
}
