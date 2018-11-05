package main

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-reuseport"
	"net"
	"time"
)

const NatServerTestPort = 8001

type NatServer struct {
	server   net.Listener
	natCache map[string]*NatCacheItem
}

type NatCacheItem struct {
	PeerId      string
	PublicIp    string
	PublicPort  string
	PrivateIp   string
	PrivatePort string
	updateTime  time.Time
	connection  net.Conn
}

func NewServer() *NatServer {

	localAddress := &net.UDPAddr{
		Port: NatServerTestPort,
	}

	s, err := reuseport.Listen(
		"udp4", localAddress.String())

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

		conn, err := s.server.Accept()
		if err != nil {
			fmt.Errorf(err.Error())
			continue
		}

		go s.processNewConnection(conn)
	}
}

func (s *NatServer) processNewConnection(conn net.Conn) {

	for {
		data := make([]byte, 2048)
		n, err := conn.Read(data)
		if err != nil {
			fmt.Errorf(err.Error())
			continue
		}

		request := &nat_pb.Request{}
		if err := proto.Unmarshal(data[:n], request); err != nil {
			fmt.Println("can't parse the nat request", err)
			return
		}

		fmt.Println("get nat request from client:", request)

		if request.ReqType == nat_pb.RequestType_KAReq {
			s.answerKA(conn, request.KeepAlive)
		} else if request.ReqType == nat_pb.RequestType_inviteReq {
			s.makeAMatch(conn, request.Invite)
		}
	}
}

func (s *NatServer) answerKA(conn net.Conn, request *nat_pb.RegRequest) error {

	response := &nat_pb.Response{
		ResType: nat_pb.ResponseType_KARes,
	}

	resKA := &nat_pb.RegResponse{}

	peerAddrStr := conn.RemoteAddr().String()
	host, port, err := net.SplitHostPort(peerAddrStr)

	if host == request.PrivateIp {
		resKA.IsAfterNat = false
	} else {
		resKA.IsAfterNat = true
		resKA.PublicIp = host
		resKA.PublicPort = port
	}

	response.KeepAlive = resKA

	responseData, err := proto.Marshal(response)
	if err != nil {
		fmt.Println("failed to marshal nat response data", err)
		return err
	}

	if _, err := conn.Write(responseData); err != nil {
		return fmt.Errorf(err.Error())
	}

	item := &NatCacheItem{
		PeerId:      request.NodeId,
		PublicIp:    resKA.PublicIp,
		PublicPort:  resKA.PublicPort,
		PrivateIp:   request.PrivateIp,
		PrivatePort: request.PrivatePort,
		updateTime:  time.Now(),
		connection:  conn,
	}

	s.natCache[item.PeerId] = item

	return nil
}

func (s *NatServer) makeAMatch(conn net.Conn, request *nat_pb.InviteRequest) error {

	responseTo := &nat_pb.Response{
		ResType: nat_pb.ResponseType_inviteRes,
	}

	toCachedItem := s.natCache[request.ToPeerId]
	fromCachedItem := s.natCache[request.FromPeerId]

	toInfo := &nat_pb.InviteResponse{
		PeerId:      toCachedItem.PeerId,
		PublicIp:    toCachedItem.PublicIp,
		PublicPort:  toCachedItem.PublicPort,
		PrivateIp:   toCachedItem.PrivateIp,
		PrivatePort: toCachedItem.PrivatePort,
	}
	responseTo.Invite = toInfo

	responseToData, err := proto.Marshal(responseTo)
	if err != nil {
		fmt.Println("failed to marshal target which you want to connect to", err)
		return err
	}

	if _, err := conn.Write(responseToData); err != nil {
		fmt.Println("failed to send connection request to invitor", err)
		return fmt.Errorf(err.Error())
	}

	//<<<<<<<-------------------------------------->>>>>>>>>>>>>
	responseFrom := &nat_pb.Response{
		ResType: nat_pb.ResponseType_invitedRes,
	}

	fromInfo := &nat_pb.InviteResponse{
		PeerId:      fromCachedItem.PeerId,
		PublicIp:    fromCachedItem.PublicIp,
		PublicPort:  fromCachedItem.PublicPort,
		PrivateIp:   fromCachedItem.PrivateIp,
		PrivatePort: fromCachedItem.PrivatePort,
	}

	responseFrom.Invite = fromInfo

	responseFromData, err := proto.Marshal(responseFrom)
	if err != nil {
		fmt.Println("failed to marshal inviter's request data", err)
		return err
	}

	if _, err := toCachedItem.connection.Write(responseFromData); err != nil {
		fmt.Println("failed to send connection request to target", err)
		return fmt.Errorf(err.Error())
	}

	return nil
}
