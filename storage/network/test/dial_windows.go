package main

import (
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"net"
)

func (peer *NatPeer) dialMultiTarget(ack *net_pb.DigConfirm, port int32, localAddr *net.UDPAddr) (*net.UDPConn, error) {
	lis, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		return nil, err
	}
	msg := net_pb.NatMsg{
		Typ: nbsnet.NatBlankKA,
	}
	data, _ := proto.Marshal(&msg)

	for _, ips := range ack.PubIps {
		tarAddr := &net.UDPAddr{
			IP:   net.ParseIP(ips),
			Port: int(port),
		}
		if _, err := lis.WriteToUDP(data, tarAddr); err != nil {
			return nil, err
		}
	}
	buffer := make([]byte, utils.NormalReadBuffer)
	n, p, e := lis.ReadFromUDP(buffer)
	if e != nil {
		return nil, err

	}
	logger.Info("get answer:->", n, p)
	lis.Close()
	conn, err := net.DialUDP("udp", localAddr, p)
	if err != nil {
		return nil, err
	}
	return conn, err
}
