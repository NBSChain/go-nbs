package network

import (
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"net"
)

type Network interface {
	StartUp(id string) error

	GetNatInfo() string

	GetAddress() *net_pb.NbsAddress

	DialUDP(network string, localAddr, remoteAddr *net.UDPAddr) (*NbsUdpConn, error)

	ListenUDP(network string, lisAddr *net.UDPAddr) (*NbsUdpConn, error)

	Connect(fromId, toId, toIP string, toPort int) (*NbsUdpConn, error)
}
