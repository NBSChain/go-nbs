package shareport

import (
	"errors"
	"net"
	"syscall"
)

/*********************************************************************************
*
*We only support IPV4 UDP right now.
*
*********************************************************************************/

var ENotSupported = errors.New("we only support udp4 right now")

var AddrInet4AnyIp = [4]byte{0, 0, 0, 0}

func UDPAddrToSockAddr(network, addr string) (*syscall.SockaddrInet4, error) {

	address, err := net.ResolveUDPAddr(network, addr)
	if err != nil {
		return nil, err
	}

	socket := &syscall.SockaddrInet4{
		Port: address.Port,
	}

	if ipv4 := address.IP.To4(); ipv4 != nil {
		socket.Addr = [4]byte{ipv4[0], ipv4[1], ipv4[2], ipv4[3]}
	}

	return socket, nil
}

func DialUDP(network, localAddr, remoteAddr string) (conn *net.UDPConn, err error) {

	return dial(network, localAddr, remoteAddr)
}

func ListenUDP(network, addr string) (*net.UDPConn, error) {
	return listenUDP(network, addr)
}
