package shareport

import (
	"errors"
	"net"
)

/*********************************************************************************
*
*We only support IPV4 UDP right now.
*
*********************************************************************************/

var ENotSupported = errors.New("we only support udp4 right now")

var AddrInet4AnyIp = [4]byte{0, 0, 0, 0}

func DialUDP(network, localAddr, remoteAddr string) (conn *net.UDPConn, err error) {

	return dial(network, localAddr, remoteAddr)
}

func ListenUDP(network, addr string) (*net.UDPConn, error) {
	return listenUDP(network, addr)
}
