package network

import (
	"context"
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/denat"
	"github.com/NBSChain/go-nbs/storage/network/nat"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/shareport"
	"github.com/NBSChain/go-nbs/thirdParty/account"
	"github.com/NBSChain/go-nbs/utils"
	"net"
	"sync"
)

type nbsNetwork struct {
	networkId   string
	ctx         context.Context
	CanServe    bool
	natServer   *nat.Server
	natClient   *nat.Client
	connManager *nbsnet.ConnManager
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
		ctx: context.Background(),
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

	denat.GetDeNatSerIns().Setup(peerId)

	network.networkId = peerId

	network.natServer = nat.NewNatServer(network.networkId)

	c, err := nat.NewNatClient(peerId, network.natServer.CanServe)
	if err != nil {
		return err
	}

	network.natClient = c
	network.CanServe = c.CanServer
	return nil
}

func (network *nbsNetwork) GetNatInfo() string {

	addr := network.natClient.NatAddr
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

func (network *nbsNetwork) DialUDP(nt string, localAddr, remoteAddr *net.UDPAddr) (*nbsnet.NbsUdpConn, error) {

	c, err := net.DialUDP(nt, localAddr, remoteAddr)
	if err != nil {
		return nil, err
	}

	host, _, _ := nbsnet.SplitHostPort(c.LocalAddr().String())
	conn := &nbsnet.NbsUdpConn{
		RealConn:  c,
		CType:     nbsnet.CTypeNormal,
		SessionID: c.LocalAddr().String() + ConnectionSeparator + remoteAddr.String(),
		LocAddr: &nbsnet.NbsUdpAddr{
			NetworkId: network.networkId,
			CanServe:  network.CanServe,
			NatServer: network.natAddr.NatServer,
			NatIp:     network.natAddr.NatIp,
			PubIp:     network.natAddr.PubIp,
			NatPort:   network.natAddr.NatPort,
			PriIp:     host,
		},
	}

	return conn, nil
}

func (network *nbsNetwork) ListenUDP(nt string, lAddr *net.UDPAddr) (*nbsnet.NbsUdpConn, error) {

	var realConn *net.UDPConn
	var cType nbsnet.ConnType
	if network.CanServe {
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
		cType = nbsnet.CTypeNatListen
	}

	host, _, _ := nbsnet.SplitHostPort(realConn.LocalAddr().String())
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
		},
	}

	return conn, nil
}

//TODO::bind local port and ip can't support right now.
func (network *nbsNetwork) Connect(lAddr, rAddr *nbsnet.NbsUdpAddr, toPort int) (*nbsnet.NbsUdpConn, error) {

	if lAddr == nil {
		lAddr = network.natClient.NatAddr
	}

	if rAddr.CanServe {
		return network.makeDirectConn(lAddr, rAddr, toPort)
	}

	var sessionID = lAddr.NetworkId + ConnectionSeparator + rAddr.NetworkId

	var realConn *net.UDPConn
	var connType nbsnet.ConnType
	if lAddr.CanServe {
		connType = nbsnet.CTypeNatSimplex
		c, err := network.natServer.InvitePeerBehindNat(lAddr, rAddr, sessionID, toPort)
		if err != nil {
			return nil, err
		}
		realConn = c
	} else {
		c, t, err := network.natServer.PunchANatHole(lAddr, rAddr, sessionID, toPort)
		if err != nil {
			return nil, err
		}
		connType = t
		realConn = c
	}

	host, _, _ := nbsnet.SplitHostPort(realConn.LocalAddr().String())
	conn := &nbsnet.NbsUdpConn{
		RealConn:  realConn,
		CType:     connType,
		SessionID: sessionID,
		LocAddr: &nbsnet.NbsUdpAddr{
			NetworkId: network.networkId,
			CanServe:  network.natAddr.CanServe,
			NatServer: network.natAddr.NatServer,
			NatIp:     network.natAddr.NatIp,
			PubIp:     network.natAddr.PubIp,
			NatPort:   network.natAddr.NatPort,
			PriIp:     host,
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
	host, _, _ := nbsnet.SplitHostPort(c.LocalAddr().String())
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
		},
	}
	return conn, nil
}
