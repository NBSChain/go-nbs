package main

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/golang/protobuf/proto"
	"net"
	"os"
	"strconv"
	"time"
)

type NatPeer struct {
	peerID      string
	conn        *net.UDPConn
	privateIP   string
	privatePort int
}

func NewPeer() *NatPeer {

	c, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   net.ParseIP("52.8.190.235"),
		Port: NatServerTestPort,
	})

	if err != nil {
		panic(err)
	}

	host, port, _ := net.SplitHostPort(c.LocalAddr().String())
	portInt, _ := strconv.Atoi(port)

	client := &NatPeer{
		conn:        c,
		peerID:      os.Args[2],
		privateIP:   host,
		privatePort: portInt,
	}

	return client
}

func (peer *NatPeer) runLoop() {

	request := &nat_pb.Request{
		ReqType: nat_pb.RequestType_KAReq,
		KeepAlive: &nat_pb.RegRequest{
			NodeId:      peer.peerID,
			PrivateIp:   peer.privateIP,
			PrivatePort: int32(peer.privatePort),
		},
	}

	requestData, err := proto.Marshal(request)
	if err != nil {
		fmt.Println("failed to marshal nat request", err)
	}

	fmt.Println("start to keep alive......")

	for {
		if no, err := peer.conn.Write(requestData); err != nil || no == 0 {
			fmt.Println("failed to send nat request to natServer ", err, no)
			continue
		}

		responseData := make([]byte, 2048)
		hasRead, _, err := peer.conn.ReadFromUDP(responseData)
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

		switch response.ResType {
		case nat_pb.ResponseType_KARes:
			time.Sleep(20 * time.Second)
		case nat_pb.ResponseType_invitedRes:
			peer.connectToPeers(response.Invite)
		case nat_pb.ResponseType_inviteRes:
			peer.connectToPeers(response.Invite)
		}
	}
}

func (peer *NatPeer) punchAHole(targetId string) {
	time.Sleep(20 * time.Second)

	inviteRequest := &nat_pb.InviteRequest{
		FromPeerId: peer.peerID,
		ToPeerId:   targetId,
	}

	request := &nat_pb.Request{
		ReqType: nat_pb.RequestType_inviteReq,
		Invite:  inviteRequest,
	}

	requestData, err := proto.Marshal(request)
	if err != nil {
		panic(err)
	}

	if _, err := peer.conn.Write(requestData); err != nil {
		panic(err)
	}
}

func (peer *NatPeer) connectToPeers(response *nat_pb.InviteResponse) {

	conn, err := net.DialUDP("udp4", &net.UDPAddr{
		IP:   net.ParseIP(peer.privateIP),
		Port: peer.privatePort,
	}, &net.UDPAddr{
		IP:   net.ParseIP(response.PublicIp),
		Port: int(response.PublicPort),
	})

	if err != nil {
		fmt.Errorf("failed to send hole punch data")
		return
	}

	conn.Write([]byte("anything ok,"))
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("input run mode -c or -s")
		os.Exit(1)
	}

	if os.Args[1] == "-c" {

		if len(os.Args) < 3 {
			fmt.Println("input this peer's ID")
			os.Exit(1)
		}

		client := NewPeer()

		if len(os.Args) == 4 {

			go client.runLoop()

			client.punchAHole(os.Args[3])

			<-make(chan struct{})

		} else if len(os.Args) == 3 {

			client.runLoop()
		}

	} else if os.Args[1] == "-s" {

		server := NewServer()
		server.Processing()

	} else {
		fmt.Println("unknown action")
	}
}
