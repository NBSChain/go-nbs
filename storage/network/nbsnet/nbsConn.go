package nbsnet

import (
	"fmt"
	"net"
	"time"
)

type ConnType int8

const (
	_ ConnType = iota
	CTypeNormal
	CTypeNat
	CTypeNatReverseDirect
	CTypeNatReverseWithProxy
)

type ConnTask struct {
	Err       chan error
	SessionId string
	CType     ConnType
	ProxyConn *net.UDPConn
	ConnCh    *net.UDPConn
}

type NbsUdpConn struct {
	ConnId   string
	CType    ConnType
	RealConn *net.UDPConn
	IsClosed bool
}

func (conn *NbsUdpConn) SetDeadline(t time.Time) {
	conn.RealConn.SetDeadline(t)
}

func (conn *NbsUdpConn) Write(d []byte) (int, error) {
	return conn.RealConn.Write(d)
}

func (conn *NbsUdpConn) Read(b []byte) (int, error) {
	return conn.RealConn.Read(b)
}

func (conn *NbsUdpConn) Close() error {
	conn.IsClosed = true
	//conn.parent.Close(conn.connId)
	return conn.RealConn.Close()
}

func (conn *NbsUdpConn) ReadFromUDP(b []byte) (int, *net.UDPAddr, error) {
	return conn.RealConn.ReadFromUDP(b)
}

func (conn *NbsUdpConn) WriteToUDP(b []byte, addr *net.UDPAddr) (int, error) {
	return conn.RealConn.WriteToUDP(b, addr)
}

//------------------nat conn---------------
func (conn *NbsUdpConn) Send(b []byte) (int, error) {

	switch conn.CType {
	case CTypeNormal:
		return conn.RealConn.Write(b)
	case CTypeNat:
		return conn.RealConn.Write(b)
	case CTypeNatReverseDirect:
		return conn.RealConn.Write(b)
	default:
		return 0, fmt.Errorf("unkown nat connection type")
	}
}

func (conn *NbsUdpConn) Receive(b []byte) (int, error) {

	switch conn.CType {
	case CTypeNormal:
		return conn.RealConn.Read(b)
	case CTypeNat:
		return conn.RealConn.Read(b)
	case CTypeNatReverseDirect:
	default:
		return 0, fmt.Errorf("unkown nat connection type")
	}

	return 0, nil
}

type HoleConn struct {
	SessionId        string
	LocalForwardConn *net.UDPConn
	RemoteDataConn   *net.UDPConn
}
