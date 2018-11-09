package shareport

import "net"

/*********************************************************************************
*
*We only support IPV4 UDP right now.
*
*********************************************************************************/

func ResolveAddr(network, address string) (net.Addr, error) {
	switch network {
	default:
		return nil, net.UnknownNetworkError(network)
	case "ip", "ip4", "ip6":
		return net.ResolveIPAddr(network, address)
	case "tcp", "tcp4", "tcp6":
		return net.ResolveTCPAddr(network, address)
	case "udp", "udp4", "udp6":
		return net.ResolveUDPAddr(network, address)
	case "unix", "unixgram", "unixpacket":
		return net.ResolveUnixAddr(network, address)
	}
}

func DialUDP(network, laddr, raddr string) (net.Conn, error) {
	return nil, nil
}

func ListenUDP(network, address string) (net.PacketConn, error) {
	return nil, nil
}

func DialTCP(network, laddr, raddr string) (net.Conn, error) {
	return nil, nil
}

func ListenTCP(network, address string) (net.Listener, error) {
	return nil, nil
}
