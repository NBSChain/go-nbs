package network

import (
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"net"
)

type HostOption func() error
type SetupOption func() error

const (
	NormalReadBuffer = 1 << 11
)

type Network interface {
	StartUp(id string, options ...SetupOption) error

	GetNatInfo() string

	DialUDP(network string, localAddr, remoteAddr *net.UDPAddr) (*NbsUdpConn, error)

	StorePeerInfo(addr *net_pb.NbsAddress)
}
