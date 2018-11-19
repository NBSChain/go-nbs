package network

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/nat"
	"github.com/NBSChain/go-nbs/storage/network/pb"
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
		return err
	}

	addr.PeerId = peerId
	network.addresses = addr

	if addr.CanBeService {
		return nil
	}

	if err := network.natManager.NewKAChannel(); err != nil {
		return err
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

func (network *nbsNetwork) GetAddress() *net_pb.NbsAddress {
	return network.addresses
}

func (network *nbsNetwork) DialUDP(nt string, localAddr, remoteAddr *net.UDPAddr) (*NbsUdpConn, error) {

	c, err := net.DialUDP(nt, localAddr, remoteAddr)
	if err != nil {
		return nil, err
	}

	conn := &NbsUdpConn{
		realConn: c,
		cType:    nat.ConnTypeNormal,
		//parent:   network.connManager,
		connId: localAddr.String() + ConnectionSeparator + remoteAddr.String(),
	}

	//network.connManager.put(conn)

	return conn, nil
}

func (network *nbsNetwork) ListenUDP(nt string, lAddr *net.UDPAddr) (*NbsUdpConn, error) {

	c, err := net.ListenUDP(nt, lAddr)
	if err != nil {
		return nil, err
	}

	conn := &NbsUdpConn{
		realConn: c,
		cType:    nat.ConnTypeNormal,
		//parent:   network.connManager,
		connId: lAddr.String(),
	}

	//network.connManager.put(conn)

	return conn, nil
}

func (network *nbsNetwork) Connect(fromId, toId, toPubIp string, toPort int) (*NbsUdpConn, error) {

	var conn *NbsUdpConn
	var connId = fromId + ConnectionSeparator + toId

	if toPubIp == "" { //TIPS::it means the target is behind a nat server.

		natTunnel := network.natManager.NatKATun

		connTask := natTunnel.MakeANatConn(fromId, toId, connId, toPort)

		var c *net.UDPConn
		select {
		case c = <-connTask.ConnCh:
			if c == nil {
				return nil, connTask.Err
			}
			logger.Debug("nat connection success.")
		case <-time.After(time.Second * 5):
			return nil, fmt.Errorf("time out")
		}
		conn = &NbsUdpConn{
			realConn:  c,
			cType:     connTask.CType,
			proxyAddr: connTask.ProxyAddr,
			connId:    connId,
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
			cType:    nat.ConnTypeNormal,
			connId:   connId,
		}
	}

	return conn, nil
}
