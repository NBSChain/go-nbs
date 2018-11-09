// +build darwin freebsd dragonfly netbsd openbsd linux

package shareport

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"syscall"
)

const filePrefix = "sharePort."

func socket(addr *syscall.SockaddrInet4) (int, error) {

	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	if err != nil {
		return 0, err
	}

	if err := syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
		return 0, os.NewSyscallError("setSocketOption", err)
	}

	if err := syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEPORT, 1); err != nil {
		return 0, os.NewSyscallError("setSocketOption", err)
	}

	if err := syscall.Bind(fd, addr); err != nil {
		return 0, fmt.Errorf("bind socket to file descrition error=%s", err.Error())
	}

	if err := syscall.SetNonblock(fd, true); err != nil {
		return 0, fmt.Errorf("set noblock flag to file descrition error=%s", err.Error())
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

	fd, err := socket(address)
	if err != nil {
		return nil, err
	}

	return fdToPacketConn(fd)
}

func dial(localAddr, remoteAddr *syscall.SockaddrInet4) (net.PacketConn, error) {
	fd, err := socket(localAddr)
	if err != nil {
		return nil, err
	}

	if err := syscall.Connect(fd, remoteAddr); err != nil {
		return nil, err
	}

	return fdToPacketConn(fd)
}
