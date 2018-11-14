package network

import (
	"github.com/NBSChain/go-nbs/storage/network/nat"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"net"
)

type NbsAddress struct {
	PublicAddr *net.UDPAddr
	PrivateIp  string
	NatType    nat_pb.NatType
}

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

	if err := network.natManager.FindWhoAmI(); err != nil {
		logger.Warning("boot strap err:", err)
	}
	return nil
}

func (network *nbsNetwork) GetNatInfo() string {
	if network.natManager == nil {
		return "nat manager isn't initialized."
	}

	return network.natManager.GetStatus()
}

func (network *nbsNetwork) LocalAddrInfo() NbsAddress {

	addr := NbsAddress{
		PublicAddr: network.natManager.PublicAddress,
		PrivateIp:  network.natManager.PrivateIP,
		NatType:    network.natManager.NatType,
	}

	return addr
}
