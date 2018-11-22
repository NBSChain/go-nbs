package nbsnet

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/gogo/protobuf/proto"
	"net"
	"strconv"
	"time"
)

type ConnType int8

const (
	_ ConnType = iota
	CTypeNormal
	CTypeNat
)

type ConnTask struct {
	SessionId string
	PubErr    chan error
	PriErr    chan error
	PubConn   *net.UDPConn
	PriConn   *net.UDPConn
}

type NbsUdpConn struct {
	ConnId   string
	CType    ConnType
	RealConn *net.UDPConn
	IsClosed bool
	LocAddr  *NbsUdpAddr
}

/************************************************************************
*
*			normal function
*
*************************************************************************/
func (conn *NbsUdpConn) SetDeadline(t time.Time) {
	conn.RealConn.SetDeadline(t)
}

func (conn *NbsUdpConn) Write(d []byte) (int, error) { //mark::::

	data := conn.packAddr(d)
	//TODO:: judge the data length value returned by raw udp socket.
	return conn.RealConn.Write(data)
}

func (conn *NbsUdpConn) Read(b []byte) (int, error) { //mark::::
	n, err := conn.RealConn.Read(b)

	data, _ := conn.unpackAddr(b[:n], nil)
	b = make([]byte, len(data))
	copy(b, data)
	return n, err
}

func (conn *NbsUdpConn) Close() error {
	conn.IsClosed = true
	//conn.parent.Close(conn.connId)
	return conn.RealConn.Close()
}

func (conn *NbsUdpConn) ReadFromUDP(b []byte) (int, *NbsUdpAddr, error) { //mark::::
	n, peerAddr, err := conn.RealConn.ReadFromUDP(b)
	data, pAddr := conn.unpackAddr(b[:n], peerAddr)
	b = make([]byte, len(data))
	copy(b, data)
	return n, pAddr, err
}

func (conn *NbsUdpConn) WriteToUDP(b []byte, addr *NbsUdpAddr) (int, error) { //mark::::

	data := conn.packAddr(b)
	peerAddr := &net.UDPAddr{
		IP:   net.ParseIP(addr.PubIp),
		Port: int(addr.PubPort),
	}
	//TODO:: judge the data length value returned by raw udp socket.
	return conn.RealConn.WriteToUDP(data, peerAddr)
}

func (conn *NbsUdpConn) LocalAddr() *NbsUdpAddr {
	return conn.LocAddr
}

func SplitHostPort(addr string) (string, int32, error) {
	host, port, err := net.SplitHostPort(addr)
	intPort, _ := strconv.Atoi(port)

	return host, int32(intPort), err
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
