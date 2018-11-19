package network

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/nat"
	"net"
	"time"
)

const (
	ConnectionSeparator = "-"
)

func (network *nbsNetwork) StartUp(peerId string) error {

	network.netWorkId = peerId

	network.natManager = nat.NewNatManager(network.netWorkId)

	addr, err := network.natManager.FindWhoAmI()
	if err != nil {
		logger.Warning("boot strap err:", err)
		return err
	}
	addr.PeerId = peerId
	network.addresses = addr

	network.connManager = newConnManager()

	if !network.addresses.CanBeService {
		ka, err := network.natManager.NewKAChannel()
		if err != nil {
			logger.Warning("failed to create nat server ka channel.")
			return err
		}

		if err := ka.InitNatChannel(); err != nil {
			logger.Warning("create NAT keep alive tunnel failed", err)
			return err
		}
		network.connManager.natKATun = ka
	}

	return nil
}

func (network *nbsNetwork) GetNatInfo() string {
	if network.natManager == nil {
		return "nat manager isn't initialized."
	}

	status := fmt.Sprintf("\n=========================================================================\n"+
		"\tnetworkId:\t%s\n"+
		"\tCanBeService:\t%v\n"+
		"\tpublicIP:\t%s\n"+
		"\tprivateIP:\t%s\n"+
		"=========================================================================",
		network.netWorkId,
		network.addresses.CanBeService,
		network.addresses.PublicIp,
		network.addresses.PrivateIp)

	return status
}

func (network *nbsNetwork) GetPublicIp() string {
	if network.addresses.CanBeService {
		return network.addresses.PublicIp
	}

	return ""
}

func (network *nbsNetwork) DialUDP(nt string, localAddr, remoteAddr *net.UDPAddr) (*NbsUdpConn, error) {

	c, err := net.DialUDP(nt, localAddr, remoteAddr)
	if err != nil {
		return nil, err
	}

	conn := &NbsUdpConn{
		realConn: c,
		cType:    ConnType_Normal,
		parent:   network.connManager,
		connId:   localAddr.String() + ConnectionSeparator + remoteAddr.String(),
	}

	network.connManager.put(conn)

	return conn, nil
}

func (network *nbsNetwork) ListenUDP(nt string, lAddr *net.UDPAddr) (*NbsUdpConn, error) {

	c, err := net.ListenUDP(nt, lAddr)
	if err != nil {
		return nil, err
	}

	conn := &NbsUdpConn{
		realConn: c,
		cType:    ConnType_Normal,
		parent:   network.connManager,
		connId:   lAddr.String(),
	}

	network.connManager.put(conn)

	return conn, nil
}

func (network *nbsNetwork) Connect(fromId, toId, toPubIp string, toPort int) (*NbsUdpConn, error) {

	var conn *NbsUdpConn
	var connId = fromId + ConnectionSeparator + toId

	if toPubIp == "" { //TIPS::it means the target is behind a nat server.

		natTunnel := network.connManager.natKATun

		connChan, err := natTunnel.MakeANatConn(fromId, toId, connId, toPort)
		if err != nil {
			return nil, err
		}
		var c *net.UDPConn
		select {
		case c = <-connChan:
			logger.Debug("nat connection success.")
		case <-time.After(time.Second * 5):
			return nil, fmt.Errorf("time out")
		}
		conn = &NbsUdpConn{
			realConn: c,
			cType:    ConnType_Nat,
			connId:   connId,
		}

	} else {

		c, err := net.DialUDP("udp4", nil, &net.UDPAddr{
			Port: toPort,
			IP:   net.ParseIP(toPubIp),
		})
		if err != nil {
			return nil, err
		}

		conn = &NbsUdpConn{
			realConn: c,
			cType:    ConnType_Normal,
			connId:   connId,
		}
	}

	conn.parent = network.connManager
	network.connManager.put(conn)

	return conn, nil
}
