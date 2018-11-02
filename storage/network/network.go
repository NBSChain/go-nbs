package network

type HostOption func() error
type SetupOption func() error

type Network interface {
	StartUp(id string, options ...SetupOption) error

	NewHost(options ...HostOption) Host

	ListenAddrString(address string) HostOption
}
