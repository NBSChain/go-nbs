package main

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/gogo/protobuf/proto"
	"net"
	"time"
)

const NatServerTestPort = 8001

type NatServer struct {
	server   *net.UDPConn
	natCache map[string]*NatCacheItem
}

type NatCacheItem struct {
	PeerId      string
	PublicIp    string
	PublicPort  int32
	PrivateIp   string
	PrivatePort int32
	updateTime  time.Time
}

func NewServer() *NatServer {

	s, err := net.ListenUDP("udp4", &net.UDPAddr{
		Port: NatServerTestPort,
	})

	if err != nil {
		panic(err)
	}

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

		n, peerAddr, err := s.server.ReadFromUDP(data)
		if err != nil {
			fmt.Errorf(err.Error())
			continue
		}

		fmt.Println("receiver ka from :", peerAddr)

		request := &nat_pb.Request{}
		if err := proto.Unmarshal(data[:n], request); err != nil {
			fmt.Println("can't parse the nat request", err, peerAddr)
			continue
		}

		fmt.Println("get nat request from client:", request)

		if request.ReqType == nat_pb.RequestType_KAReq {
			s.answerKA(peerAddr, request.KeepAlive)
		} else if request.ReqType == nat_pb.RequestType_inviteReq {
			s.makeAMatch(peerAddr, request.Invite)
		}
	}
}

func (s *NatServer) answerKA(peerAddr *net.UDPAddr, request *nat_pb.RegRequest) error {

	response := &nat_pb.Response{
		ResType: nat_pb.ResponseType_KARes,
	}

	resKA := &nat_pb.RegResponse{}

	if peerAddr.IP.Equal(net.ParseIP(request.PrivateIp)) {
		resKA.IsAfterNat = false
	} else {
		resKA.IsAfterNat = true
		resKA.PublicIp = peerAddr.IP.String()
		resKA.PublicPort = int32(peerAddr.Port)
		resKA.Zone = peerAddr.Zone
	}

	response.KeepAlive = resKA

	responseData, err := proto.Marshal(response)
	if err != nil {
		fmt.Println("failed to marshal nat response data", err)
		return err
	}

	if _, err := s.server.WriteToUDP(responseData, peerAddr); err != nil {
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

	s.natCache[item.PeerId] = item

	return nil
}

func (s *NatServer) makeAMatch(peerAddr *net.UDPAddr, request *nat_pb.InviteRequest) error {

	responseTo := &nat_pb.Response{
		ResType: nat_pb.ResponseType_inviteRes,
	}

	cacheItem := s.natCache[request.ToPeerId]

	toInfo := &nat_pb.InviteResponse{
		PeerId:      cacheItem.PeerId,
		PublicIp:    cacheItem.PublicIp,
		PublicPort:  cacheItem.PublicPort,
		PrivateIp:   cacheItem.PrivateIp,
		PrivatePort: cacheItem.PrivatePort,
	}
	responseTo.Invite = toInfo

	responseToData, err := proto.Marshal(responseTo)
	if err != nil {
		fmt.Println("failed to marshal target which you want to connect to", err)
		return err
	}

	if _, err := s.server.WriteToUDP(responseToData, peerAddr); err != nil {
		return fmt.Errorf(err.Error())
	}

	//<<<<<<<-------------------------------------->>>>>>>>>>>>>

	responseFrom := &nat_pb.Response{
		ResType: nat_pb.ResponseType_invitedRes,
	}

	cacheItem = s.natCache[request.FromPeerId]

	fromInfo := &nat_pb.InviteResponse{
		PeerId:      cacheItem.PeerId,
		PublicIp:    cacheItem.PublicIp,
		PublicPort:  cacheItem.PublicPort,
		PrivateIp:   cacheItem.PrivateIp,
		PrivatePort: cacheItem.PrivatePort,
	}

	responseFrom.Invite = fromInfo

	responseFromData, err := proto.Marshal(responseFrom)
	if err != nil {
		fmt.Println("failed to marshal inviter's request data", err)
		return err
	}

	if _, err := s.server.WriteToUDP(responseFromData, &net.UDPAddr{
		IP:   net.ParseIP(toInfo.PublicIp),
		Port: int(toInfo.PrivatePort),
	}); err != nil {
		return fmt.Errorf(err.Error())
	}

	return nil
}
