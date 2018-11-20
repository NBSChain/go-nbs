package network

import (
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"net"
)

type Network interface {
	StartUp(id string) error

	GetNatInfo() string

	GetAddress() nbsnet.NbsUdpAddr

	DialUDP(network string, localAddr, remoteAddr *net.UDPAddr) (*nbsnet.NbsUdpConn, error)

	ListenUDP(network string, lisAddr *net.UDPAddr) (*nbsnet.NbsUdpConn, error)

	Connect(lAddr, rAddr *nbsnet.NbsUdpAddr) (*nbsnet.NbsUdpConn, error)
}
