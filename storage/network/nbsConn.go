package network

import (
	"net"
	"time"
)

type ConnType int8

const (
	_ ConnType = iota
	ConnType_Normal
	ConnType_Nat
	ConnType_NatInverse
)

type NbsUdpConn struct {
	connId   string
	cType    ConnType
	realConn *net.UDPConn
	isClosed bool
	parent   *ConnManager
}

func (conn *NbsUdpConn) SetDeadline(t time.Time) {
	conn.realConn.SetDeadline(t)
}

func (conn *NbsUdpConn) Write(d []byte) (int, error) {
	return conn.realConn.Write(d)
}

func (conn *NbsUdpConn) Read(b []byte) (int, error) {
	return conn.realConn.Read(b)
}

func (conn *NbsUdpConn) Close() error {
	conn.isClosed = true
	conn.parent.Close(conn.connId)
	return conn.realConn.Close()
}

func (conn *NbsUdpConn) ReadFromUDP(b []byte) (int, *net.UDPAddr, error) {
	return conn.realConn.ReadFromUDP(b)
}

func (conn *NbsUdpConn) WriteToUDP(b []byte, addr *net.UDPAddr) (int, error) {
	return conn.realConn.WriteToUDP(b, addr)
}

func (conn *NbsUdpConn) Send([]byte) (int, error) {
	return 0, nil
}
