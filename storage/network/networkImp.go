package network

import "github.com/NBSChain/go-nbs/thirdParty/idService"

func (network *nbsNetwork) NewHost(options ...HostOption) Host {

	instance := &NbsHost{}

	return instance
}

func (network *nbsNetwork) ListenAddrString(address string) HostOption {

	return func() error {
		return nil
	}
}

func (network *nbsNetwork) Identity(id *idService.Identity) HostOption {

	return func() error {
		return nil
	}
}
