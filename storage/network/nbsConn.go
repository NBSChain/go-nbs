package network

import (
	"net"
	"time"
)

type NbsUdpConn struct {
	connId   string
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
