package main

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/golang/protobuf/proto"
	"net"
	"time"
)

type SimplePeer struct {
}

func NewSimplePeer() *SimplePeer {
	return &SimplePeer{}
}

func (peer *SimplePeer) probe() error {

	conn, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   net.ParseIP("192.168.103.101"),
		Port: NatServerTestPort,
	})

	if err != nil {
		return err
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

	for {

		time.Sleep(time.Second * 2)

		conn.SetDeadline(time.Now().Add(time.Second * 5))
		if no, err := conn.Write(requestData); err != nil || no == 0 {
			fmt.Println("failed to send nat request to natServer ", err, no)
			continue
		}

		responseData := make([]byte, 2048)
		hasRead, err := conn.Read(responseData)
		if err != nil {
			fmt.Println("failed to read nat response from natServer", err)
			continue
		}

		response := &nat_pb.Response{}
		if err := proto.Unmarshal(responseData[:hasRead], response); err != nil {
			fmt.Println("failed to unmarshal nat response data", err)
			continue
		}

		fmt.Println(response)

	}

	return nil
}
