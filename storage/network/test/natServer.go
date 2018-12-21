package main

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/nat"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/gogo/protobuf/proto"
	"net"
	"time"
)

const NatServerTestPort = 8001

type NatServer struct {
	server   *net.UDPConn
	natCache map[string]*nat.HostBehindNat
}

func NewServer() *NatServer {

	s, err := net.ListenUDP("udp4", &net.UDPAddr{
		Port: NatServerTestPort,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(s.LocalAddr().String())

	server := &NatServer{
		server:   s,
		natCache: make(map[string]*nat.HostBehindNat),
	}

	return server
}

func (s *NatServer) Processing() {

	fmt.Println("start to run......")

	for {
		data := make([]byte, 2048)
		n, peerAddr, err := s.server.ReadFromUDP(data)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		request := &net_pb.NatMsg{}
		if err := proto.Unmarshal(data[:n], request); err != nil {
			fmt.Println("can't parse the nat request", err)
			continue
		}

		fmt.Println("\nNat request:->", request)

		switch request.Typ {
		case nbsnet.NatKeepAlive:
			ack := request.KeepAlive
			nodeId := ack.NodeId
			item, ok := s.natCache[nodeId]
			if !ok {
				item := &nat.HostBehindNat{
					UpdateTIme: time.Now(),
					PubAddr:    peerAddr,
					PriAddr:    request.KeepAlive.PriAddr,
				}
				s.natCache[nodeId] = item
			} else {
				item.UpdateTIme = time.Now()
			}

			ack.PubAddr = peerAddr.String()
			data, _ := proto.Marshal(request)
			if _, err := s.server.WriteToUDP(data, peerAddr); err != nil {
				fmt.Println("keep alive ack err:->", err)
			}
		case nbsnet.NatDigApply:
			app := request.DigApply
			item, ok := s.natCache[app.TargetId]
			if !ok {
				panic(err)
			}
			app.Public = peerAddr.String()
			data, _ := proto.Marshal(request)
			if _, err := s.server.WriteToUDP(data, item.PubAddr); err != nil {
				panic(err)
			}
		case nbsnet.NatDigConfirm:
			ack := request.DigConfirm
			item, ok := s.natCache[ack.TargetId]
			if !ok {
				panic(ok)
			}
			ack.Public = peerAddr.String()
			data, _ := proto.Marshal(request)
			if _, err := s.server.WriteToUDP(data, item.PubAddr); err != nil {
				panic(err)
			}
		default:
			fmt.Println("unknown msg type")
		}
	}
}

func (s *NatServer) answerKA(peerAddr *net.UDPAddr, request *net_pb.BootReg) error {

	peerAddrStr := peerAddr.String()
	host, port, err := nbsnet.SplitHostPort(peerAddrStr)
	response := &net_pb.NatMsg{
		Typ: nbsnet.NatBootAnswer,
		BootAnswer: &net_pb.BootAnswer{
			NatType:    net_pb.NatType_BehindNat,
			PublicIp:   host,
			PublicPort: port,
		},
	}

	responseData, err := proto.Marshal(response)
	if err != nil {
		fmt.Println("failed to marshal nat response data", err)
		return err
	}

	if _, err := s.server.WriteTo(responseData, peerAddr); err != nil {
		return fmt.Errorf(err.Error())
	}

	item := &nat.HostBehindNat{
		UpdateTIme: time.Now(),
		PubAddr:    peerAddr,
		PriAddr:    nbsnet.JoinHostPort(request.PrivateIp, request.PrivatePort),
	}
	s.natCache[request.NodeId] = item
	fmt.Println("online:->", request.NodeId)

	return nil
}
