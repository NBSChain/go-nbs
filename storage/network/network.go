package network

type HostOption func() error

type Network interface {
	NewHost(options ...HostOption) Host
	ListenAddrString(address string) HostOption
}
