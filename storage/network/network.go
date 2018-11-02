package network

import "github.com/NBSChain/go-nbs/thirdParty/idService"

type HostOption func() error

type Network interface {
	NewHost(options ...HostOption) Host

	ListenAddrString(address string) HostOption

	Identity(id *idService.Identity) HostOption
}
