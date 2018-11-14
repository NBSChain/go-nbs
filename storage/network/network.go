package network

import "github.com/NBSChain/go-nbs/storage/network/pb"

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

	NatType() nat_pb.NatType
}
