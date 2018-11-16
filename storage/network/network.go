package network

import (
	"net"
)

type Network interface {
	StartUp(id string) error

	GetNatInfo() string

	GetPublicIp() string

	DialUDP(network string, localAddr, remoteAddr *net.UDPAddr) (*NbsUdpConn, error)

	ListenUDP(network string, lisAddr *net.UDPAddr) (*NbsUdpConn, error)

	Connect(fromId, toId, toIP string, toPort int) (*NbsUdpConn, error)
}
