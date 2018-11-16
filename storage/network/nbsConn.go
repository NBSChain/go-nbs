package network

import (
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/gogo/protobuf/proto"
	"net"
	"time"
)

type NbsUdpConn struct {
	addr net_pb.NbsAddress
	c    *net.UDPConn
}

func (conn *NbsUdpConn) SetDeadline(t time.Time) {
	conn.c.SetDeadline(t)
}

func (conn *NbsUdpConn) Write(d []byte) (int, error) {

	msg := &net_pb.NbsMessage{
		Data:    d,
		NbsAddr: &conn.addr,
	}
	nbsData, err := proto.Marshal(msg)
	if err != nil {
		return 0, err
	}

	return conn.c.Write(nbsData)
}

func (conn *NbsUdpConn) Read(b []byte) (int, error) {
	buffer := make([]byte, NormalReadBuffer)
	n, err := conn.c.Read(buffer)
	if err != nil {
		return 0, err
	}

	msg := &net_pb.NbsMessage{}

	if err := proto.Unmarshal(buffer[:n], msg); err != nil {
		return 0, err
	}

	copy(b, msg.Data)

	GetInstance().StorePeerInfo(msg.NbsAddr)

	return len(msg.Data), nil
}

func (conn *NbsUdpConn) Close() error {
	return conn.c.Close()
}
