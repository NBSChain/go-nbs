package main

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/gogo/protobuf/proto"
	"net"
)

const NatServerTestPort = 8001

type NatServer struct {
	server *net.UDPConn
}

func NewServer() *NatServer {

	s, err := net.ListenUDP("udp4", &net.UDPAddr{
		Port: NatServerTestPort,
	})
	if err != nil {
		panic(err)
	}
	server := &NatServer{
		server: s,
	}

	return server
}

func (s *NatServer) Processing() {
	fmt.Println("start to run......")
	for {
		data := make([]byte, 2048)

		n, peerAddr, err := s.server.ReadFromUDP(data)
		if err != nil {
			fmt.Errorf(err.Error())
			continue
		}

		fmt.Println("receiver ka from :", peerAddr)

		request := &nat_pb.NatRequest{}
		if err := proto.Unmarshal(data[:n], request); err != nil {
			fmt.Println("can't parse the nat request", err, peerAddr)
			continue
		}

		fmt.Println("get nat request from client:", request)

		response := &nat_pb.NatResponse{}
		if peerAddr.IP.Equal(net.ParseIP(request.PrivateIp)) {
			response.IsAfterNat = false
		} else {
			response.IsAfterNat = true
			response.PublicIp = peerAddr.IP.String()
			response.PublicPort = int32(peerAddr.Port)
			response.Zone = peerAddr.Zone
		}

		responseData, err := proto.Marshal(response)
		if err != nil {
			fmt.Println("failed to marshal nat response data", err)
			continue
		}

		if _, err := s.server.WriteToUDP(responseData, peerAddr); err != nil {
			fmt.Errorf(err.Error())
		}
	}
}
