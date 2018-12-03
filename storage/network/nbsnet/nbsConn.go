package nbsnet

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/gogo/protobuf/proto"
	"net"
	"time"
)

type ConnType int8

var logger = utils.GetLogInstance()

const (
	_ ConnType = iota
	CTypeNormal
	CTypeNatSimplex
	CTypeNatDuplex
	CTypeNatListen
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
func (conn *NbsUdpConn) SetDeadline(t time.Time) error {
	return conn.RealConn.SetDeadline(t)
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
//TODO::
func (conn *NbsUdpConn) Send(b []byte) (int, error) {
	switch conn.CType {
	case CTypeNormal:
		return conn.RealConn.Write(b)
	case CTypeNatSimplex:
		return conn.RealConn.Write(b)

	case CTypeNatDuplex:
		return conn.RealConn.Write(b)
	default:
		return 0, fmt.Errorf("unkown nat connection type")
	}
}

//TODO::
func (conn *NbsUdpConn) Receive(b []byte) (int, error) {

	switch conn.CType {
	case CTypeNormal:
		return conn.RealConn.Read(b)
	case CTypeNatSimplex:
		return conn.RealConn.Read(b)
	default:
		return 0, fmt.Errorf("unkown nat connection type")
	}
}

/************************************************************************
*
*			private functions
*
*************************************************************************/
//TODO::I don't think this is the final resolution.
func (conn *NbsUdpConn) preHandleMsg(msg *net_pb.NatMsg, addr *net.UDPAddr) bool {

	if msg.Typ == NatPriDigSyn {
		res := &net_pb.NatMsg{
			Typ: NatPriDigAck,
			Seq: msg.Seq + 1,
		}
		b, _ := proto.Marshal(res)
		if _, err := conn.RealConn.WriteTo(b, addr); err != nil {
			logger.Warning("write back nat message err:->", err)
		}
		return true
	}

	return false
}

func (conn *NbsUdpConn) String() string {
	return "[" + conn.RealConn.LocalAddr().String() + "]-->[" +
		conn.RealConn.RemoteAddr().String() + "]"
}
