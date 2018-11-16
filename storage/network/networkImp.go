package network

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/nat"
	"net"
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

func (network *nbsNetwork) DialUDP(nt string, localAddr, remoteAddr *net.UDPAddr) (*NbsUdpConn, error) {

	c, err := net.DialUDP(nt, localAddr, remoteAddr)
	if err != nil {
		return nil, err
	}

	conn := &NbsUdpConn{
		realConn: c,
		parent:   network.connManager,
		connId:   localAddr.String() + "-" + remoteAddr.String(),
	}

	network.connManager.put(conn)

	return conn, nil
}
