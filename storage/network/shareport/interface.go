package shareport

import "net"

func Dial(network, laddr, raddr string) (net.Conn, error) {
	return nil, nil
}

func ListenTCP(network, address string) (net.Listener, error) {
	return nil, nil
}

func ListenUDP(network, address string) (net.PacketConn, error) {
	return nil, nil
}
