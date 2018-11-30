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

func NewNatManager(networkId string) *Manager {

	denat.GetDeNatSerIns().Setup(networkId)

	natObj := &Manager{
		networkId: networkId,
		canServe:  make(chan bool),
		cache:     make(map[string]*HostBehindNat),
	}

	natObj.startNatService()

	go natObj.natServiceListening()

	go natObj.cacheManager()

	return natObj
}

func (nat *Manager) SetUpNatChannel(netNatAddr *nbsnet.NbsUdpAddr) error {

	port := strconv.Itoa(utils.GetConfig().NatChanSerPort)
	listener, err := shareport.ListenUDP("udp4", "0.0.0.0:"+port)
	if err != nil {
		logger.Warning("create share listening udp failed.")
		return err
	}

	serverHost := netNatAddr.NatServer
	client, err := shareport.DialUDP("udp4", "0.0.0.0:"+port, serverHost)
	if err != nil {
		logger.Warning("create share port dial udp connection failed.")
		return err
	}
	netNatAddr.NatServer = serverHost

	tunnel := &KATunnel{
		natChanged: make(chan struct{}),
		networkId:  nat.networkId,
		natAddr:    netNatAddr,
		serverHub:  listener,
		kaConn:     client,
		sharedAddr: client.LocalAddr().String(),
		updateTime: time.Now(),
		workLoad:   make(map[string]*ProxyTask),
	}

	go tunnel.runLoop()

	go tunnel.listening()

	go tunnel.readKeepAlive()

	go tunnel.connManage()

	nat.NatKATun = tunnel
	select {
	case <-tunnel.natChanged:
	case <-time.After(time.Second * 2):
	}

	return nil
}

func (nat *Manager) WaitNatConfirm() chan bool {
	return nat.canServe
}

func (nat *Manager) PunchANatHole(lAddr, rAddr *nbsnet.NbsUdpAddr, connId string, toPort int) (*net.UDPConn, error) {

	if err := nat.NatKATun.StartDigHole(lAddr, rAddr, connId, toPort); err != nil {
		return nil, err
	}
	priConnTask := &ConnTask{
		err: make(chan error),
	}
	pubConnTask := &ConnTask{
		err: make(chan error),
	}
	defer close(priConnTask.err)
	defer close(pubConnTask.err)

	go nat.NatKATun.directDialInPriNet(lAddr, rAddr, priConnTask, toPort, connId)

	go nat.NatKATun.DigInPubNet(lAddr, rAddr, pubConnTask, connId)

	var pubFail, priFail bool

	for i := 2; i > 0; i-- {
		select {
		case err := <-priConnTask.err:
			if err == nil {
				return priConnTask.udpConn, nil
			} else {
				priFail = true
				if pubFail {
					return nil, err
				}
			}
			return priConnTask.udpConn, err
		case err := <-pubConnTask.err:
			if err == nil {
				return pubConnTask.udpConn, nil
			} else {
				pubFail = true
				if priFail {
					return nil, err
				}
			}
		case <-time.After(HolePunchTimeOut / 2):
			return nil, fmt.Errorf("time out")
		}
	}
	return nil, fmt.Errorf("time out")
}

func (nat *Manager) InvitePeerBehindNat(lAddr, rAddr *nbsnet.NbsUdpAddr,
	connId string, toPort int) (*net.UDPConn, error) {

	conn, err := shareport.DialUDP("udp4", "", rAddr.NatServer)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localHost := conn.LocalAddr().String()
	_, fromPort, _ := net.SplitHostPort(localHost)

	req := &net_pb.NatRequest{
		MsgType: net_pb.NatMsgType_ReverseDig,
		Invite: &net_pb.ReverseInvite{
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

	connChan := &ConnTask{
		err: make(chan error),
	}

	logger.Debug("Step1: notify applier's nat server:", req)

	go nat.waitInviteAnswer(localHost, connId, connChan)

	select {
	case err := <-connChan.err:
		if err != nil {
			return nil, err
		} else {
			return connChan.udpConn, nil
		}
	case <-time.After(HolePunchTimeOut):
		return nil, fmt.Errorf("time out")
	}
}

func (nat *Manager) waitInviteAnswer(host, sessionID string, task *ConnTask) {

	lisConn, err := shareport.ListenUDP("udp4", host)
	if err != nil {
		task.err <- err
		return
	}
	defer lisConn.Close()

	logger.Debug("Step2: wait the answer:", host)

	buffer := make([]byte, utils.NormalReadBuffer)
	n, peerAddr, err := lisConn.ReadFromUDP(buffer)
	if err != nil {
		task.err <- err
		return
	}

	res := &net_pb.NatRequest{}
	if err := proto.Unmarshal(buffer[:n], res); err != nil {
		task.err <- err
		return
	}

	if res.MsgType != net_pb.NatMsgType_ReverseDigACK ||
		res.InviteAck.SessionId != sessionID {
		task.udpConn = nil
		task.err <- fmt.Errorf("didn't get the answer")
		return
	}

	conn, err := shareport.DialUDP("udp4", host, peerAddr.String())
	if err != nil {
		task.err <- err
		return
	}
	task.udpConn = conn
	task.err <- err

	logger.Debug("Step5: get answer and make a connection:->", host, peerAddr.String())
}
