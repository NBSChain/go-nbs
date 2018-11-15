package network

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/nat"
)

func (network *nbsNetwork) NewHost(options ...HostOption) Host {

	instance := &NbsHost{}

	return instance
}

func (network *nbsNetwork) ListenAddrString(address string) HostOption {

	return func() error {
		return nil
	}
}

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

	network.addresses = addr

	return nil
}

func (network *nbsNetwork) GetNatInfo() string {
	if network.natManager == nil {
		return "nat manager isn't initialized."
	}

	status := fmt.Sprintf("\n=========================================================================\n"+
		"\tnetworkId:\t%s\n"+
		"\tisInPbulic:\t%v\n"+
		"\tpublicIP:\t%s\n"+
		"\tprivateIP:\t%s\n"+
		"=========================================================================",
		network.netWorkId,
		network.addresses.IsInPub,
		network.addresses.PublicIP,
		network.addresses.PrivateIp)

	return status
}
