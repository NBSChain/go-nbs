package main

import (
	"context"
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/storage/network/shareport"
	"github.com/golang/protobuf/proto"
	"net"
	"syscall"
	"time"
)

type SimplePeer struct {
}

func NewSimplePeer() *SimplePeer {
	return &SimplePeer{}
}

func sharePort(network, address string, rawConn syscall.RawConn) error {
	fn := func(s uintptr) {
		if err := syscall.SetsockoptInt(syscall.Handle(s), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
			panic(err)
		}
		if err := syscall.SetNonblock(syscall.Handle(s), true); err != nil {
			panic(err)
		}
	}
	if err := rawConn.Control(fn); err != nil {
		panic(err)
	}

	return nil
}

func (peer *SimplePeer) probe() error {

	conn, err := shareport.DialUDP("udp4", "192.168.30.12:7001", "192.168.103.155:8001")
	if err != nil {
		panic(err)
	}

	lc := &net.ListenConfig{
		Control: sharePort,
	}

	_, err = lc.ListenPacket(context.Background(), "udp4", "192.168.30.12:7001")
	if err != nil {
		panic(err)
	}

	host, port, _ := net.SplitHostPort(conn.LocalAddr().String())

	request := &nat_pb.NatRequest{
		MsgType: nat_pb.NatMsgType_BootStrapReg,
		BootRegReq: &nat_pb.BootNatRegReq{
			NodeId:      "this is a probe message.",
			PrivateIp:   host,
			PrivatePort: port,
		},
	}

	requestData, err := proto.Marshal(request)
	if err != nil {
		fmt.Println("failed to marshal nat request", err)
	}

	fmt.Println("simple peer probe......")

	go read(conn)

	for {

		time.Sleep(time.Second * 3)

		conn.SetWriteDeadline(time.Now().Add(time.Second * 2))
		if no, err := conn.Write(requestData); err != nil || no == 0 {
			fmt.Println("failed to send nat request to natServer ", err, no)
			continue
		}

		fmt.Println("-request->done.")
	}

	return nil
}

func read(conn *net.UDPConn) {

	time.Sleep(time.Second * 2)

	for {
		conn.SetReadDeadline(time.Now().Add(time.Second * 5))
		responseData := make([]byte, 2048)
		hasRead, addr, err := conn.ReadFrom(responseData)
		if err != nil {
			fmt.Println("failed to read nat response from natServer", err)
			continue
		}

		response := &nat_pb.Response{}
		if err := proto.Unmarshal(responseData[:hasRead], response); err != nil {
			fmt.Println("failed to unmarshal nat response data", err)
			continue
		}

		fmt.Println("=======response=>", response, addr)
	}
}
