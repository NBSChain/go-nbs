// +build linux

package shareport

import (
	"net"
	"os"
	"strconv"
	"syscall"
)

const filePrefix = "linux_sharePort."

func newShareSocket() (int, error) {

	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	if err != nil {
		return 0, os.NewSyscallError("Socket", err)
	}

	if err := syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
		return 0, os.NewSyscallError("setSocketOption", err)
	}

	if err := syscall.SetNonblock(fd, true); err != nil {
		return 0, os.NewSyscallError("SetNonblock", err)
	}

	return fd, nil
}

func fdToPacketConn(fd int) (net.PacketConn, error) {

	file := os.NewFile(uintptr(fd), filePrefix+strconv.Itoa(os.Getpid()))

	packetConn, err := net.FilePacketConn(file)
	if err != nil {
		syscall.Close(fd)
		return nil, err
	}

	if err := file.Close(); err != nil {
		syscall.Close(fd)
		packetConn.Close()
		return nil, err
	}

	return packetConn, nil
}

func listenUDP(address *syscall.SockaddrInet4) (net.PacketConn, error) {

	fd, err := newShareSocket()
	if err != nil {
		return nil, err
	}

	if err := syscall.Bind(fd, addr); err != nil {
		return nil, os.NewSyscallError("Bind", err)
	}

	return fdToPacketConn(fd)
}

func fdToConn(fd int) (net.Conn, error) {

	file := os.NewFile(uintptr(fd), filePrefix+strconv.Itoa(os.Getpid()))
	conn, err := net.FileConn(file)
	if err != nil {
		file.Close()
		return nil, err
	}

	if err = file.Close(); err != nil {
		conn.Close()
		return nil, err
	}

	return conn, nil
}

func dial(localAddr, remoteAddr *syscall.SockaddrInet4) (net.Conn, error) {

	fd, err := newShareSocket()
	if err != nil {
		return nil, err
	}

	if err := syscall.Bind(fd, address); err != nil {
		return nil, os.NewSyscallError("Bind", err)
	}

	if localAddr.Port == AddrInet4AnyPort {
		addr, err := getLocalAddr(fd)
		if err != nil {
			return nil, err
		}
		localAddr.Port = addr.Port
	}

	if err := syscall.Connect(fd, remoteAddr); err != nil {
		syscall.Close(fd)
		return nil, err
	}

	if localAddr.Addr == AddrInet4AnyIp {
		addr, err := getLocalAddr(fd)
		if err != nil {
			return nil, err
		}
		localAddr.Addr = addr.Addr
	}

	return fdToConn(fd)
}
