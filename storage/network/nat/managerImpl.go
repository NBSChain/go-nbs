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

const (
	MsgPoolSize = 1 << 10
)

func NewNatManager(networkId string) *Manager {

	denat.GetDeNatSerIns().Setup(networkId)

	natObj := &Manager{
		networkId:   networkId,
		canServe:    make(chan bool),
		cache:       make(map[string]*HostBehindNat),
		task:        make(chan *MsgTask, MsgPoolSize),
		msgHandlers: make(map[net_pb.MsgType]taskProcess),
	}

	natObj.initService()

	go natObj.TaskReceiver()
	go natObj.RunLoop()

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
		digTask:    make(map[string]*ConnTask),
	}

	go tunnel.runLoop()

	go tunnel.listening()

	go tunnel.readKeepAlive()

	go tunnel.connManage()

	nat.NatKATun = tunnel
	select {
	case <-tunnel.natChanged:
	case <-time.After(time.Second * 2): //TODO::
	}

	return nil
}

func (nat *Manager) WaitNatConfirm() chan bool {
	return nat.canServe
}

func (nat *Manager) PunchANatHole(lAddr, rAddr *nbsnet.NbsUdpAddr,
	connId string, toPort int) (*net.UDPConn, nbsnet.ConnType, error) {

	priConnTask := &ConnTask{
		udpConn: make(chan *net.UDPConn),
	}
	defer priConnTask.Close()

	go nat.NatKATun.directDialInPriNet(lAddr, rAddr, priConnTask, toPort, connId)

	pubConnTask := &ConnTask{
		udpConn: make(chan *net.UDPConn),
	}
	defer pubConnTask.Close()
	go nat.NatKATun.DigHoeInPubNet(lAddr, rAddr, connId, toPort, pubConnTask)

	var pubFail, priFail bool
	for i := 2; i > 0; i-- {
		select {
		case conn := <-priConnTask.udpConn:
			logger.Debug("hole punch step1-2 dig direct in private network finished:->", priConnTask.err)
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

func (nat *Manager) InvitePeerBehindNat(lAddr, rAddr *nbsnet.NbsUdpAddr,
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

	connChan := &ConnTask{
		udpConn: make(chan *net.UDPConn),
	}

	logger.Debug("Step1: notify applier's nat server:", req)

	go nat.waitInviteAnswer(localHost, connId, connChan)

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

func (nat *Manager) waitInviteAnswer(host, sessionID string, task *ConnTask) {

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
