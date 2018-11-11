package shareport

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"syscall"
	"time"
	"unsafe"
)

const filePrefix = "windows_sharePort."

type winShareConn struct {
	fd         syscall.Handle
	localAddr  *syscall.SockaddrInet4
	remoteAddr *syscall.SockaddrInet4
}

func (winConn *winShareConn) Read(b []byte) (n int, err error) {
	buffer := &syscall.WSABuf{
		Len: uint32(len(b)),
		Buf: &b[0],
	}
	dataReceived := uint32(0)
	flags := uint32(0)
	if err := syscall.WSARecv(winConn.fd, buffer, 1,
		&dataReceived, &flags, nil, nil); err != nil {
		return 0, err
	}

	return int(dataReceived), nil
}
func (winConn *winShareConn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {

	buffer := &syscall.WSABuf{
		Len: uint32(len(p)),
		Buf: &p[0],
	}

	dataReceived := uint32(n)
	flags := uint32(0)
	var asIp4 syscall.RawSockaddrInet4
	fromAny := (*syscall.RawSockaddrAny)(unsafe.Pointer(&asIp4))
	fromSize := int32(unsafe.Sizeof(asIp4))

	err = syscall.WSARecvFrom(winConn.fd, buffer, 1,
		&dataReceived, &flags, fromAny, &fromSize, nil, nil)
	if err != nil {
		return 0, nil, err
	}

	addr = &net.UDPAddr{
		IP:   net.IPv4(asIp4.Addr[0], asIp4.Addr[1], asIp4.Addr[2], asIp4.Addr[3]),
		Port: int(asIp4.Port),
	}

	return n, addr, nil
}

func (winConn *winShareConn) Write(b []byte) (n int, err error) {

	msgLen := uint32(len(b))
	buffer := syscall.WSABuf{
		Len: msgLen,
		Buf: &b[0],
	}

	if err != nil {
		return 0, err
	}

	err = syscall.WSASend(winConn.fd, &buffer, 1,
		&msgLen, 0, nil, nil)

	if err != nil {
		return 0, err
	}

	return int(msgLen), nil
}

func (winConn *winShareConn) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	msgLen := uint32(len(p))
	buffer := syscall.WSABuf{
		Len: msgLen,
		Buf: &p[0],
	}

	host, port, err := net.SplitHostPort(addr.String())

	portInt, err := strconv.Atoi(port)
	to := &syscall.SockaddrInet4{
		Port: portInt,
		Addr: [4]byte{host[0], host[1], host[2], host[3]},
	}

	if err != nil {
		return 0, err
	}

	err = syscall.WSASendto(winConn.fd, &buffer, 1,
		&msgLen, 0, to, nil, nil)

	if err != nil {
		return 0, err
	}

	return int(msgLen), nil
}

func (winConn *winShareConn) Close() error {

	if err := syscall.WSACleanup(); err != nil {
		return err
	}

	if err := syscall.Closesocket(winConn.fd); err != nil {
		return err
	}

	return nil
}

func (winConn *winShareConn) LocalAddr() net.Addr {
	la := winConn.localAddr.Addr
	addr := net.UDPAddr{
		IP:   net.IPv4(la[0], la[1], la[2], la[3]),
		Port: winConn.localAddr.Port,
	}
	return &addr
}
func (winConn *winShareConn) RemoteAddr() net.Addr {

	ra := winConn.remoteAddr.Addr
	addr := net.UDPAddr{
		IP:   net.IPv4(ra[0], ra[1], ra[2], ra[3]),
		Port: winConn.remoteAddr.Port,
	}
	return &addr
}

func (winConn *winShareConn) SetDeadline(t time.Time) error {
	return nil
}
func (winConn *winShareConn) SetReadDeadline(t time.Time) error {
	return nil
}
func (winConn *winShareConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func getLocalAddr(fd syscall.Handle) (*syscall.SockaddrInet4, error) {

	realLocal, err := syscall.Getsockname(fd)
	if err != nil {
		fmt.Println("get sock name failed:", err)
		return nil, err
	} else {
		switch realLocal.(type) {
		case *syscall.SockaddrInet4:
			address := realLocal.(*syscall.SockaddrInet4)
			fmt.Printf("====%v:%d====\n", address.Addr, address.Port)
			return address, nil
		default:
			return nil, fmt.Errorf("only support udp4 right now")
		}
	}
}

func socket(addr *syscall.SockaddrInet4) (syscall.Handle, error) {

	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	if err != nil {
		return 0, err
	}

	if err := syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
		return 0, os.NewSyscallError("setSocketOption", err)
	}

	if err := syscall.SetNonblock(fd, true); err != nil {
		return 0, fmt.Errorf("set noblock flag to file descrition error=%s", err.Error())
	}

	return fd, nil
}

func dial(localAddress, remoteAddress *syscall.SockaddrInet4) (net.Conn, error) {

	var wsaData syscall.WSAData
	if err := syscall.WSAStartup(makeWord(2, 2), &wsaData); err != nil {
		return nil, err
	}

	fd, err := socket(localAddress)
	if err != nil {
		return nil, err
	}

	if err := syscall.Connect(fd, remoteAddress); err != nil {
		syscall.Close(fd)
		return nil, err
	}

	if localAddress == nil {
		addr, err := getLocalAddr(fd)
		if err != nil {
			return nil, err
		}

		if err := syscall.Bind(fd, addr); err != nil {
			return nil, err
		}

		localAddress = addr
	}

	return newConn(fd, localAddress, remoteAddress), nil
}

func newConn(fd syscall.Handle, localAddr, remoteAddr *syscall.SockaddrInet4) *winShareConn {

	winConn := &winShareConn{
		fd:         fd,
		localAddr:  localAddr,
		remoteAddr: remoteAddr,
	}

	return winConn
}

func makeWord(low, high uint8) uint32 {
	var ret uint16 = uint16(high)<<8 + uint16(low)
	return uint32(ret)
}

func listenUDP(address *syscall.SockaddrInet4) (net.PacketConn, error) {

	var wsaData syscall.WSAData
	if err := syscall.WSAStartup(makeWord(2, 2), &wsaData); err != nil {
		return nil, err
	}

	fd, err := socket(address)
	if err != nil {
		return nil, err
	}

	if err := syscall.Bind(fd, address); err != nil {
		return nil, fmt.Errorf("bind socket to file descrition error=%s", err.Error())
	}

	return newConn(fd, address, nil), nil
}
