package network

import "github.com/NBSChain/go-nbs/thirdParty/idService"

type HostOption func() error
type SetupOption func() error

type Network interface {
	StartUp(id *idService.Identity, options ...SetupOption) error

	NewHost(options ...HostOption) Host

	ListenAddrString(address string) HostOption
}
