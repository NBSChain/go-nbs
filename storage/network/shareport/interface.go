package shareport

import (
	"net"
	"syscall"
)

/*********************************************************************************
*
*We only support IPV4 UDP right now.
*
*********************************************************************************/

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

func DialUDP(network, laddr, raddr string) (conn net.Conn, err error) {

	localAddress, err := UDPAddrToSockAddr(network, laddr)
	if err != nil {
		return nil, err
	}
	remoteAddress, err := UDPAddrToSockAddr(network, raddr)
	if err != nil {
		return nil, err
	}

	return dial(localAddress, remoteAddress)
}

func ListenUDP(network, addr string) (net.PacketConn, error) {

	address, err := UDPAddrToSockAddr(network, addr)
	if err != nil {
		return nil, err
	}
	return listenUDP(address)
}
