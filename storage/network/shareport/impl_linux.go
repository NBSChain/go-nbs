// +build linux

package shareport

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"syscall"
)

func dial(network, localAddr, remoteAddr string) (*net.UDPConn, error) {

	var localUdpAddr *net.UDPAddr

	if localAddr != "" {

		host, port, err := net.SplitHostPort(localAddr)
		if err != nil {
			return nil, err
		}

		intPort, err := strconv.Atoi(port)
		if err != nil {
			return nil, err
		}

		localUdpAddr = &net.UDPAddr{
			Port: intPort,
			IP:   net.ParseIP(host),
		}
	}

	d := &net.Dialer{
		LocalAddr: localUdpAddr,
		Control:   sharePort,
	}

	c, err := d.Dial(network, remoteAddr)
	if err != nil {
		return nil, err
	}

	conn, ok := c.(*net.UDPConn)
	if !ok {
		return nil, fmt.Errorf("not a udp connection")
	}

	return conn, nil
}

func sharePort(network, address string, rawConn syscall.RawConn) (err error) {

	fn := func(s uintptr) {
		if err = syscall.SetsockoptInt(int(s), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
			return
		}
		if err = syscall.SetNonblock(int(s), true); err != nil {
			return
		}
	}

	if err != nil {
		return err
	}
	if err = rawConn.Control(fn); err != nil {
		return err
	}

	return nil
}

func listenUDP(network, address string) (*net.UDPConn, error) {

	lc := &net.ListenConfig{
		Control: sharePort,
	}

	c, err := lc.ListenPacket(context.Background(), network, address)
	if err != nil {
		return nil, err
	}

	conn, ok := c.(*net.UDPConn)
	if !ok {
		return nil, ENotSupported
	}

	return conn, nil
}
