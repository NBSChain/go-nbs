package network

import (
	"net"
)

type Network interface {
	StartUp(id string) error

	GetNatInfo() string

	DialUDP(network string, localAddr, remoteAddr *net.UDPAddr) (*NbsUdpConn, error)
}
