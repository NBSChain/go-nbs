package network

import (
	"github.com/NBSChain/go-nbs/storage/network/nat"
	"net"
	"time"
)

type NbsUdpConn struct {
	connId    string
	c         *net.UDPConn
	kaChannel *nat.KAChannel
}

func (conn *NbsUdpConn) SetDeadline(t time.Time) {
	conn.c.SetDeadline(t)
}

func (conn *NbsUdpConn) Write(d []byte) (int, error) {
	return conn.c.Write(d)
}

func (conn *NbsUdpConn) Read(b []byte) (int, error) {
	return conn.c.Read(b)
}

func (conn *NbsUdpConn) Close() error {
	return conn.c.Close()
}
