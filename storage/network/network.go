package network

type HostOption func() error
type SetupOption func() error

const (
	NormalReadBuffer = 1 << 11
)

type Network interface {
	StartUp(id string, options ...SetupOption) error

	NewHost(options ...HostOption) Host

	ListenAddrString(address string) HostOption

	GetNatInfo() string

	LocalAddrInfo() NbsAddress
}
