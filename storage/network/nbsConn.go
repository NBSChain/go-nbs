package network

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/nat"
	"net"
	"time"
)

type NbsUdpConn struct {
	connId    string
	cType     nat.ConnType
	proxyAddr *net.UDPAddr
	realConn  *net.UDPConn
	isClosed  bool
	//parent    *ConnManager
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
	//conn.parent.Close(conn.connId)
	return conn.realConn.Close()
}

func (conn *NbsUdpConn) ReadFromUDP(b []byte) (int, *net.UDPAddr, error) {
	return conn.realConn.ReadFromUDP(b)
}

func (conn *NbsUdpConn) WriteToUDP(b []byte, addr *net.UDPAddr) (int, error) {
	return conn.realConn.WriteToUDP(b, addr)
}

//------------------nat conn---------------
func (conn *NbsUdpConn) Send(b []byte) (int, error) {
	if conn.cType == nat.ConnTypeNat {
		return conn.realConn.Write(b)
	} else if conn.cType == nat.ConnTypeNatInverse {
		return conn.realConn.Write(b)
	}
	return 0, fmt.Errorf("unkown nat connection type")
}

func (conn *NbsUdpConn) Receive(b []byte) (int, error) {

	if conn.cType == nat.ConnTypeNat {
		return conn.realConn.Read(b)
	} else if conn.cType == nat.ConnTypeNatInverse {
		n, err := conn.realConn.Read(b)
		if err != nil {
			return 0, err
		}

		return conn.realConn.WriteToUDP(b[:n], conn.proxyAddr)
	}
	return 0, fmt.Errorf("unkown nat connection type")
}
