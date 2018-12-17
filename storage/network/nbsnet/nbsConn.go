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

func (conn *NbsUdpConn) ReceiveFromUDP(b []byte) (int, *net.UDPAddr, error) {

	switch conn.CType {
	case CTypeNormal:
		return conn.RealConn.ReadFromUDP(b)
	case CTypeNatSimplex:
		return conn.RealConn.ReadFromUDP(b)
	case CTypeNatDuplex:
		return conn.RealConn.ReadFromUDP(b)
	case CTypeNatListen: //TODO::I don't think this is the final resolution.
	GOON:
		n, peerAddr, err := conn.RealConn.ReadFromUDP(b)
		if err != nil {
			return 0, nil, err
		}
		msg := &net_pb.NatMsg{}
		if err2 := proto.Unmarshal(b[:n], msg); err2 != nil {
			return n, peerAddr, nil
		}

		if msg.Typ != NatPriDigSyn {
			return n, peerAddr, nil
		}
		res := &net_pb.NatMsg{
			Typ: NatPriDigAck,
			Seq: msg.Seq + 1,
		}
		data, _ := proto.Marshal(res)
		logger.Info("Step 1-6:->answer dig in private:->", res)
		if _, err := conn.RealConn.WriteToUDP(data, peerAddr); err != nil {
			logger.Warning("answer NatPriDigAck err:->", err)
		}
		goto GOON

	default:
		return 0, nil, fmt.Errorf("unkown nat connection type")
	}
}

/************************************************************************
*
*			private functions
*
*************************************************************************/
func (conn *NbsUdpConn) String() string {
	return "[" + conn.RealConn.LocalAddr().String() + "]-->[" +
		conn.RealConn.RemoteAddr().String() + "]"
}

func ConnString(conn net.Conn) string {

	return "[" + conn.LocalAddr().String() + "]-->[" +
		conn.RemoteAddr().String() + "]"
}
