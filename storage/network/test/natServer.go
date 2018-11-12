package main

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/gogo/protobuf/proto"
	"net"
	"strconv"
	"time"
)

const NatServerTestPort = 8001

type NatServer struct {
	server   net.PacketConn
	natCache map[string]*NatCacheItem
}

type NatCacheItem struct {
	PeerId      string
	PublicIp    string
	PublicPort  string
	PrivateIp   string
	PrivatePort string
	updateTime  time.Time
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
		natCache: make(map[string]*NatCacheItem),
	}

	return server
}

func (s *NatServer) Processing() {

	fmt.Println("start to run......")

	for {
		data := make([]byte, 2048)
		n, peerAddr, err := s.server.ReadFrom(data)
		if err != nil {
			fmt.Errorf(err.Error())
			continue
		}

		request := &nat_pb.NatRequest{}
		if err := proto.Unmarshal(data[:n], request); err != nil {
			fmt.Println("can't parse the nat request", err)
			continue
		}

		fmt.Println("\nNat request:->", request)

		if request.MsgType == nat_pb.NatMsgType_BootStrapReg {
			s.answerKA(peerAddr, request.BootRegReq)
		} else if request.MsgType == nat_pb.NatMsgType_Connect {
			s.makeAMatch(peerAddr, request.ConnReq)
		}
	}
}

func (s *NatServer) answerKA(peerAddr net.Addr, request *nat_pb.BootNatRegReq) error {

	response := &nat_pb.Response{
		MsgType: nat_pb.NatMsgType_BootStrapReg,
	}

	resKA := &nat_pb.BootNatRegRes{}

	peerAddrStr := peerAddr.String()
	host, port, err := net.SplitHostPort(peerAddrStr)

	if host == request.PrivateIp {
		resKA.NatType = nat_pb.NatType_NoNatDevice
		resKA.PublicIp = host
		resKA.PublicPort = port
	} else {
		resKA.NatType = nat_pb.NatType_BehindNat
		resKA.PublicIp = host
		resKA.PublicPort = port
	}

	response.BootRegRes = resKA

	responseData, err := proto.Marshal(response)
	if err != nil {
		fmt.Println("failed to marshal nat response data", err)
		return err
	}

	if _, err := s.server.WriteTo(responseData, peerAddr); err != nil {
		return fmt.Errorf(err.Error())
	}

	item := &NatCacheItem{
		PeerId:      request.NodeId,
		PublicIp:    resKA.PublicIp,
		PublicPort:  resKA.PublicPort,
		PrivateIp:   request.PrivateIp,
		PrivatePort: request.PrivatePort,
		updateTime:  time.Now(),
	}

	fmt.Println("Nat item cached:->", item)

	s.natCache[item.PeerId] = item

	return nil
}

func (s *NatServer) makeAMatch(peerAddr net.Addr, request *nat_pb.NatConReq) error {

	responseTo := &nat_pb.Response{
		MsgType: nat_pb.NatMsgType_Connect,
	}

	cacheItem := s.natCache[request.ToPeerId]
	cacheItemFrom := s.natCache[request.FromPeerId]
	if cacheItem == nil || cacheItemFrom == nil {
		return fmt.Errorf("some peer is not in the cache")
	}

	toInfo := &nat_pb.NatConRes{
		PeerId:      cacheItem.PeerId,
		PublicIp:    cacheItem.PublicIp,
		PublicPort:  cacheItem.PublicPort,
		PrivateIp:   cacheItem.PrivateIp,
		PrivatePort: cacheItem.PrivatePort,
	}
	responseTo.ConnRes = toInfo

	responseToData, err := proto.Marshal(responseTo)
	if err != nil {
		fmt.Println("failed to marshal target which you want to connect to", err)
		return err
	}

	if _, err := s.server.WriteTo(responseToData, peerAddr); err != nil {
		fmt.Println("failed to send connection request to invitor", err)
		return fmt.Errorf(err.Error())
	}

	//<<<<<<<-------------------------------------->>>>>>>>>>>>>
	responseFrom := &nat_pb.Response{
		MsgType: nat_pb.NatMsgType_Connect,
	}

	fromInfo := &nat_pb.NatConRes{
		PeerId:      cacheItemFrom.PeerId,
		PublicIp:    cacheItemFrom.PublicIp,
		PublicPort:  cacheItemFrom.PublicPort,
		PrivateIp:   cacheItemFrom.PrivateIp,
		PrivatePort: cacheItemFrom.PrivatePort,
	}

	responseFrom.ConnRes = fromInfo

	responseFromData, err := proto.Marshal(responseFrom)
	if err != nil {
		fmt.Println("failed to marshal inviter's request data", err)
		return err
	}

	port, _ := strconv.Atoi(toInfo.PublicPort)

	if _, err := s.server.WriteTo(responseFromData, &net.UDPAddr{
		IP:   net.ParseIP(toInfo.PublicIp),
		Port: port,
	}); err != nil {
		fmt.Println("failed to send connection request to target", err)
		return fmt.Errorf(err.Error())
	}

	return nil
}
