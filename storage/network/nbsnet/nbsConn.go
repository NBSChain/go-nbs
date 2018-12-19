package nbsnet

import (
	"context"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/NBSChain/go-nbs/utils/crypto"
	"github.com/gogo/protobuf/proto"
	"net"
	"time"
)

type ConnType int8

var logger = utils.GetLogInstance()

const (
	NatHoleKATime          = time.Second * 20
	_             ConnType = iota
	CTypeNormal
	CTypeNatSimplex
	CTypeNatDuplex
	CTypeNatListen
)

type NbsUdpConn struct {
	SessionID string
	CType     ConnType
	ctx       context.Context
	close     context.CancelFunc
	RealConn  *net.UDPConn
	LocAddr   *NbsUdpAddr
}

func NewNbsConn(c *net.UDPConn, sessionID string, cType ConnType, natAddr *NbsUdpAddr) *NbsUdpConn {
	ctx, cancel := context.WithCancel(context.Background())

	conn := &NbsUdpConn{
		ctx:       ctx,
		close:     cancel,
		RealConn:  c,
		CType:     cType,
		SessionID: sessionID,
		LocAddr:   natAddr,
	}

	if cType == CTypeNatSimplex ||
		cType == CTypeNatDuplex {
		go conn.KeepHoleOpened()
	}

	return conn
}

func (conn *NbsUdpConn) KeepHoleOpened() {
	for {
		now := time.Now().String()
		msg := &net_pb.ConnKA{
			KA: crypto.MD5SS(now),
		}
		data, _ := proto.Marshal(msg)

		select {
		case <-time.After(NatHoleKATime):
			if _, err := conn.Write(data); err != nil {
				logger.Warning("the keep alive for hole msg err:->", err)
				return
			}
		case <-conn.ctx.Done():
			logger.Debug("bye")
		}
	}
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
	conn.close()
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
	return conn.RealConn.Write(b)
}

func (conn *NbsUdpConn) Receive(b []byte) (int, error) {
	return conn.RealConn.Read(b)
}

func (conn *NbsUdpConn) ReceiveFromUDP(b []byte) (int, *net.UDPAddr, error) {
	return conn.RealConn.ReadFromUDP(b)
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
