package network

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/nat"
	"net"
)

func (network *nbsNetwork) StartUp(peerId string, options ...SetupOption) error {

	for _, opt := range options {

		if err := opt(); err != nil {
			logger.Warning("one network startup option applies failed", opt)
		}
	}

	network.netWorkId = peerId

	network.natManager = nat.NewNatManager(network.netWorkId)

	addr, err := network.natManager.FindWhoAmI()
	if err != nil {
		logger.Warning("boot strap err:", err)
	}
	addr.PeerId = peerId
	network.addresses = addr

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
		c:      c,
		connId: remoteAddr.String(),
	}

	if !network.addresses.CanBeService {

		kc, err := network.natManager.NewKAChannel()
		if err != nil {
			logger.Warning("failed to create nat server ka channel.")
			return nil, err
		}

		kc.InitNatChannel()

		conn.kaChannel = kc
	}

	return conn, nil
}
