// +build !windows

package main

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/nat"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/storage/network/shareport"
	"github.com/NBSChain/go-nbs/utils"
	"net"
	"time"
)

func (peer *NatPeer) dialMultiTarget(ack *net_pb.DigConfirm, port int32, localAddr string) (*net.UDPConn, error) {

	ch := make(chan *net.UDPConn)
	for _, ips := range ack.PubIps {
		tarAddr := nbsnet.JoinHostPort(ips, port)
		go peer.findTheRightConn(localAddr, tarAddr, ch)
	}

	select {
	case c := <-ch:
		logger.Debug("pick out the right conn and send ka:->", nbsnet.ConnString(c))
		close(ch)
		return c, nil

	case <-time.After(time.Second * 4):
		return nil, fmt.Errorf("time out")
	}
}

func (peer *NatPeer) findTheRightConn(fromAddr, toAddr string, ch chan *net.UDPConn) {
	conn, err := shareport.DialUDP("udp4", fromAddr, toAddr)
	if err != nil {
		logger.Warning("dial up a bad conn:->", toAddr)
		return
	}
	msg := net_pb.NatMsg{
		Typ: nbsnet.NatBlankKA,
		Seq: time.Now().Unix(),
	}
	data, _ := proto.Marshal(&msg)
	if _, err := conn.Write(data); err != nil {
		logger.Warning("find a bad conn:->", toAddr)
		return
	}

	logger.Debug("write to one channel:->", toAddr)

	conn.SetDeadline(time.Now().Add(time.Second * 4))
	buffer := make([]byte, utils.NormalReadBuffer)
	if _, err := conn.Read(buffer); err != nil {
		logger.Warning("this conn failed:->", nbsnet.ConnString(conn), err)
		return
	}
	logger.Debug("get response from the channel:->", toAddr)

	conn.SetDeadline(nat.NoTimeOut)
	ch <- conn
}
