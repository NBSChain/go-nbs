package nbsnet

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/gogo/protobuf/proto"
	"net"
	"time"
)

type ConnType int8

const (
	_ ConnType = iota
	CTypeNormal
	CTypeNat
)

type NbsUdpConn struct {
	SessionID string
	CType     ConnType
	RealConn  *net.UDPConn
	IsClosed  bool
	LocAddr   *NbsUdpAddr
}

/************************************************************************
*
*			normal function
*
*************************************************************************/
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
	return conn.RealConn.Close()
}

func (conn *NbsUdpConn) ReadFromUDP(b []byte) (int, *net.UDPAddr, error) {
	return conn.RealConn.ReadFromUDP(b)
}

func (conn *NbsUdpConn) WriteToUDP(b []byte, addr *net.UDPAddr) (int, error) {
	return conn.RealConn.WriteToUDP(b, addr)
}

func (conn *NbsUdpConn) LocalAddr() *NbsUdpAddr {
	return conn.LocAddr
}

/************************************************************************
*
*			nat connection
*
*************************************************************************/
func (conn *NbsUdpConn) Send(b []byte) (int, error) {

	switch conn.CType {
	case CTypeNormal:
		return conn.RealConn.Write(b)
	case CTypeNat:
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
	default:
		return 0, fmt.Errorf("unkown nat connection type")
	}

	return 0, nil
}

/************************************************************************
*
*			private functions
*
*************************************************************************/
func (conn *NbsUdpConn) packAddr(d []byte) []byte {
	lAddr := conn.LocAddr
	msg := &net_pb.NbsNetMsg{
		RawData: d,
		FromAddr: &net_pb.NbsAddr{
			NetworkId: lAddr.NetworkId,
			CanServer: lAddr.CanServe,
			PriIp:     lAddr.PriIp,
			PriPort:   lAddr.PriPort,
			NatIP:     lAddr.NatIp,
			NatPort:   lAddr.NatPort,
		},
	}
	data, _ := proto.Marshal(msg)
	return data
}

func (conn *NbsUdpConn) unpackAddr(d []byte, pAddr net.Addr) ([]byte, *NbsUdpAddr) {

	msg := &net_pb.NbsNetMsg{}

	proto.Unmarshal(d, msg)
	if pAddr == nil {
		pAddr = conn.RealConn.RemoteAddr()
	}

	host, port, _ := SplitHostPort(pAddr.String())
	addr := msg.FromAddr
	peerAddr := &NbsUdpAddr{
		NetworkId: addr.NetworkId,
		CanServe:  addr.CanServer,
		PubIp:     host,
		PubPort:   port,
		PriIp:     addr.PriIp,
		PriPort:   addr.PriPort,
		NatIp:     addr.NatIP,
		NatPort:   addr.NatPort,
	}

	return msg.RawData, peerAddr
}
func (conn *NbsUdpConn) String() string {
	return "[" + conn.RealConn.LocalAddr().String() + "]-->[" +
		conn.RealConn.RemoteAddr().String() + "]"
}
