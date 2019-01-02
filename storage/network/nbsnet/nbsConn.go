package nbsnet

import (
	"context"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/gogo/protobuf/proto"
	"net"
	"sync"
	"time"
)

type ConnType int8

var logger = utils.GetLogInstance()

const (
	NatHoleKATime          = time.Second * 18
	_             ConnType = iota
	CTypeNormal
	CTypeNatSimplex
	CTypeNatDuplex
	CTypeNatListen
)

type NbsUdpConn struct {
	CType    ConnType
	ctx      context.Context
	close    context.CancelFunc
	RealConn *net.UDPConn
	sync.Mutex
	updateTime time.Time
}

func NewNbsConn(c *net.UDPConn, cType ConnType) *NbsUdpConn {
	ctx, cancel := context.WithCancel(context.Background())

	conn := &NbsUdpConn{
		ctx:        ctx,
		close:      cancel,
		RealConn:   c,
		CType:      cType,
		updateTime: time.Now(),
	}

	if cType == CTypeNatSimplex ||
		cType == CTypeNatDuplex {
		go conn.KeepHoleOpened()
	}

	return conn
}

func (conn *NbsUdpConn) KeepHoleOpened() {
	logger.Debug("setup keep live routine", conn.String())

	for {
		select {
		case <-time.After(NatHoleKATime):
			if err := conn.keepAlive(); err != nil {
				return
			}
		case <-conn.ctx.Done():
			logger.Debug("hole closed, bye")
			return
		}
	}
}

//TODO::need a response ?ka-ack?
func (conn *NbsUdpConn) keepAlive() error {

	now := time.Now()
	msg := &net_pb.NatMsg{
		Typ: NatBlankKA,
	}
	data, _ := proto.Marshal(msg)

	conn.Lock()
	defer conn.Unlock()
	if now.Sub(conn.updateTime) < NatHoleKATime {
		return nil
	}
	if _, err := conn.RealConn.Write(data); err != nil {
		logger.Warning("the keep alive for hole msg err:->", err)
		return err
	}

	conn.updateTime = now
	logger.Debug("try to keep hole opened:->", conn.String())
	return nil
}

func (conn *NbsUdpConn) natMsgFilter(buf []byte, peerAddr *net.UDPAddr) (bool, error) {
	if conn.CType != CTypeNatListen {
		return false, nil
	}

	b := make([]byte, len(buf))
	copy(b, buf)

	msg := net_pb.NatMsg{}
	if err := proto.Unmarshal(b, &msg); err != nil {
		return false, nil
	}
	if msg.Typ < NatMsgBase || msg.Typ > NatEnd {
		return false, nil
	}

	logger.Debug("this is a inner msg:->", msg)
	var data []byte = nil
	switch msg.Typ {
	case NatFindPubIpSyn:
		data, _ = proto.Marshal(&net_pb.NatMsg{
			Typ: NatFindPubIpACK,
		})
	}
	if data == nil {
		return true, nil
	}

	if peerAddr != nil {
		if _, err := conn.RealConn.WriteToUDP(data, peerAddr); err != nil {
			return true, err
		}
		return true, nil
	}
	if _, err := conn.RealConn.Write(data); err != nil {
		return true, err
	}

	return true, nil
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
	conn.Lock()
	conn.updateTime = time.Now()
	conn.Unlock()
	return conn.RealConn.Write(d)
}

func (conn *NbsUdpConn) Read(b []byte) (int, error) {
	conn.Lock()
	conn.updateTime = time.Now()
	conn.Unlock()

reading:
	n, err := conn.RealConn.Read(b)
	if err != nil {
		return 0, err
	}

	isInnerMsg, err := conn.natMsgFilter(b[:n], nil)
	if err != nil {
		return 0, err
	}

	if isInnerMsg {
		goto reading
	}

	return n, err
}

func (conn *NbsUdpConn) Close() error {
	conn.close()
	return conn.RealConn.Close()
}

func (conn *NbsUdpConn) ReadFromUDP(b []byte) (int, *net.UDPAddr, error) {
	conn.Lock()
	conn.updateTime = time.Now()
	conn.Unlock()
reading:
	n, addr, err := conn.RealConn.ReadFromUDP(b)
	if err != nil {
		return 0, nil, err
	}

	isInnerMsg, err := conn.natMsgFilter(b[:n], addr)
	if err != nil {
		return 0, nil, err
	}

	if isInnerMsg {
		goto reading
	}

	return n, addr, err
}

func (conn *NbsUdpConn) WriteToUDP(b []byte, addr *net.UDPAddr) (int, error) {
	conn.Lock()
	conn.updateTime = time.Now()
	conn.Unlock()
	return conn.RealConn.WriteToUDP(b, addr)
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
