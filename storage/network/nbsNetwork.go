package network

import (
	"context"
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/nat"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/storage/network/shareport"
	"github.com/NBSChain/go-nbs/thirdParty/account"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"net"
	"strconv"
	"sync"
	"time"
)

type nbsNetwork struct {
	Context    context.Context
	natManager *nat.Manager
	networkId  string
	natAddr    *nbsnet.NbsUdpAddr
}

const (
	ConnectionSeparator = "-"
)

var (
	once     sync.Once
	instance *nbsNetwork
	logger   = utils.GetLogInstance()
)

/************************************************************************
*
*			public functions
*
*************************************************************************/
func GetInstance() Network {

	once.Do(func() {
		instance = newNetwork()
	})

	return instance
}

func newNetwork() *nbsNetwork {

	peerId := account.GetAccountInstance().GetPeerID()

	network := &nbsNetwork{
		Context: context.Background(),
	}

	if peerId != "" {
		if err := network.StartUp(peerId); err != nil {
			panic(err)
		}
	} else {
		logger.Warning("no account right now, so the network is down")
	}

	return network
}

func (network *nbsNetwork) StartUp(peerId string) error {

	localPeers := nbsnet.ExternalIP()
	if len(localPeers) == 0 {
		logger.Panic("no available network")
	}

	logger.Debug("all network interfaces:->", localPeers)

	network.networkId = peerId

	network.natManager = nat.NewNatManager(network.networkId)

	success := false
	for _, serverIp := range utils.GetConfig().NatServerIP {
		//serverHost := denat.GetDeNatSerIns().GetValidServer()
		port := strconv.Itoa(utils.GetConfig().NatServerPort)
		if err := network.findWhoAmI(net.JoinHostPort(serverIp, port)); err == nil {
			success = true
			break
		}
	}
	if !success {
		return fmt.Errorf("failed to find who am I")
	}

	network.natAddr.NetworkId = peerId

	if network.natAddr.CanServe {
		return nil
	}

	if err := network.natManager.SetUpNatChannel(network.natAddr); err != nil {
		return err
	}

	return nil
}

func (network *nbsNetwork) GetNatInfo() string {
	if network.natAddr == nil {
		return "nat manager isn't initialized."
	}
	addr := network.natAddr
	status := fmt.Sprintf("\n=========================================================================\n"+
		"\tnetworkId:\t%s\n"+
		"\tcanServe:\t%v\n"+
		"\tpubIp:\t%s\n"+
		"\tpriIp:\t%s\n"+
		"=========================================================================",
		network.networkId,
		addr.CanServe,
		addr.PubIp,
		addr.PriIp)

	return status
}

func (network *nbsNetwork) GetAddress() nbsnet.NbsUdpAddr {
	return *network.natAddr
}

func (network *nbsNetwork) DialUDP(nt string, localAddr, remoteAddr *net.UDPAddr) (*nbsnet.NbsUdpConn, error) {

	c, err := net.DialUDP(nt, localAddr, remoteAddr)
	if err != nil {
		return nil, err
	}

	host, port, _ := nbsnet.SplitHostPort(c.LocalAddr().String())
	conn := &nbsnet.NbsUdpConn{
		RealConn:  c,
		CType:     nbsnet.CTypeNormal,
		SessionID: c.LocalAddr().String() + ConnectionSeparator + remoteAddr.String(),
		LocAddr: &nbsnet.NbsUdpAddr{
			NetworkId: network.networkId,
			CanServe:  network.natAddr.CanServe,
			NatServer: network.natAddr.NatServer,
			NatIp:     network.natAddr.NatIp,
			PubIp:     network.natAddr.PubIp,
			NatPort:   network.natAddr.NatPort,
			PriIp:     host,
			PriPort:   port,
		},
	}

	return conn, nil
}

func (network *nbsNetwork) ListenUDP(nt string, lAddr *net.UDPAddr) (*nbsnet.NbsUdpConn, error) {

	var realConn *net.UDPConn
	var cType nbsnet.ConnType
	if network.natAddr.CanServe {
		c, err := net.ListenUDP(nt, lAddr)
		if err != nil {
			return nil, err
		}
		realConn = c
		cType = nbsnet.CTypeNormal
	} else {
		c, err := shareport.ListenUDP(nt, lAddr.String())
		if err != nil {
			return nil, err
		}
		realConn = c
		cType = nbsnet.CTypeNat
	}

	host, port, _ := nbsnet.SplitHostPort(realConn.LocalAddr().String())
	conn := &nbsnet.NbsUdpConn{
		RealConn:  realConn,
		CType:     cType,
		SessionID: lAddr.String(),
		LocAddr: &nbsnet.NbsUdpAddr{
			NetworkId: network.networkId,
			CanServe:  network.natAddr.CanServe,
			NatServer: network.natAddr.NatServer,
			NatIp:     network.natAddr.NatIp,
			NatPort:   network.natAddr.NatPort,
			PubIp:     network.natAddr.PubIp,
			PriIp:     host,
			PriPort:   port,
		},
	}

	return conn, nil
}

func (network *nbsNetwork) Connect(lAddr, rAddr *nbsnet.NbsUdpAddr, toPort int) (*nbsnet.NbsUdpConn, error) {

	if lAddr == nil {
		lAddr = network.natAddr
	}

	if rAddr.CanServe {
		return network.makeDirectConn(lAddr, rAddr, toPort)
	}

	var sessionID = lAddr.NetworkId + ConnectionSeparator + rAddr.NetworkId

	var realConn *net.UDPConn
	if lAddr.CanServe {
		c, err := network.natManager.InvitePeerBehindNat(lAddr, rAddr, sessionID, toPort)
		if err != nil {
			return nil, err
		}
		realConn = c
	} else {
		c, err := network.natManager.PunchANatHole(lAddr, rAddr, sessionID, toPort)
		if err != nil {
			return nil, err
		}
		realConn = c
	}

	host, port, _ := nbsnet.SplitHostPort(realConn.LocalAddr().String())
	conn := &nbsnet.NbsUdpConn{
		RealConn:  realConn,
		CType:     nbsnet.CTypeNat,
		SessionID: sessionID,
		LocAddr: &nbsnet.NbsUdpAddr{
			NetworkId: network.networkId,
			CanServe:  network.natAddr.CanServe,
			NatServer: network.natAddr.NatServer,
			NatIp:     network.natAddr.NatIp,
			PubIp:     network.natAddr.PubIp,
			NatPort:   network.natAddr.NatPort,
			PriIp:     host,
			PriPort:   port,
		},
	}

	return conn, nil
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
	host, port, _ := nbsnet.SplitHostPort(c.LocalAddr().String())
	conn := &nbsnet.NbsUdpConn{
		RealConn:  c,
		CType:     nbsnet.CTypeNormal,
		SessionID: sessionID,
		LocAddr: &nbsnet.NbsUdpAddr{
			NetworkId: network.networkId,
			CanServe:  network.natAddr.CanServe,
			NatServer: network.natAddr.NatServer,
			NatIp:     network.natAddr.NatIp,
			PubIp:     network.natAddr.PubIp,
			NatPort:   network.natAddr.NatPort,
			PriIp:     host,
			PriPort:   port,
		},
	}
	return conn, nil
}

func (network *nbsNetwork) findWhoAmI(serverHost string) error {

	conn, err := network.connectToNatServer(serverHost)
	if err != nil {
		logger.Error("can't know who am I:->", err)
		return err
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(time.Second * 3))

	localHost, err := network.sendNatRequest(conn)
	if err != nil {
		logger.Error("failed to read nat response:->", err)
		return err
	}

	response, err := network.parseNatResponse(conn)
	if err != nil {
		logger.Debug("get NAT server info failed:->", err)
		return err
	}

	addr := &nbsnet.NbsUdpAddr{
		PubIp:     response.PublicIp,
		PriIp:     localHost,
		CanServe:  nbsnet.CanServe(response.NatType),
		NatServer: serverHost,
	}

	if response.NatType == net_pb.NatType_ToBeChecked {

		select {
		case c := <-network.natManager.WaitNatConfirm():
			addr.CanServe = c

		case <-time.After(nat.BootStrapTimeOut / 2):
			addr.CanServe = false
		}
	}

	network.natAddr = addr
	return nil
}

func (network *nbsNetwork) connectToNatServer(serverIP string) (*net.UDPConn, error) {

	host, port, _ := nbsnet.SplitHostPort(serverIP)
	natServerAddr := &net.UDPAddr{
		IP:   net.ParseIP(host),
		Port: int(port),
	}

	conn, err := net.DialUDP("udp4", nil, natServerAddr)
	if err != nil {
		return nil, err
	}

	err = conn.SetDeadline(time.Now().Add(nat.BootStrapTimeOut))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (network *nbsNetwork) sendNatRequest(conn *net.UDPConn) (string, error) {

	localAddr := conn.LocalAddr().String()

	host, port, err := net.SplitHostPort(localAddr)
	bootRequest := &net_pb.BootNatRegReq{
		NodeId:      network.networkId,
		PrivateIp:   host,
		PrivatePort: port,
	}

	request := &net_pb.NatRequest{
		MsgType:    net_pb.NatMsgType_BootStrapReg,
		BootRegReq: bootRequest,
	}

	requestData, err := proto.Marshal(request)
	if err != nil {
		logger.Error("failed to marshal nat request:->", err)
		return "", err
	}

	if no, err := conn.Write(requestData); err != nil || no == 0 {
		logger.Error("failed to send nat request to selfNatServer:->", err, no)
		return "", err
	}

	return host, nil
}

func (network *nbsNetwork) parseNatResponse(conn *net.UDPConn) (*net_pb.BootNatRegRes, error) {

	responseData := make([]byte, utils.NormalReadBuffer)
	hasRead, err := conn.Read(responseData)
	if err != nil {
		logger.Error("reading failed from nat server:->", err)
		return nil, err
	}

	response := &net_pb.NatResponse{}
	if err := proto.Unmarshal(responseData[:hasRead], response); err != nil {
		logger.Error("unmarshal Err:->", err)
		return nil, err
	}

	logger.Debug("response:->", response)

	return response.BootRegRes, nil
}
