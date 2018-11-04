package main

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/nat"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/golang/protobuf/proto"
	"net"
	"os"
	"time"
)

type NatPeer struct {
	peerID    string
	conn      *net.UDPConn
	privateIP string
}

const NatHoleListeningPort = 7001

func NewPeer() *NatPeer {

	c, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   net.ParseIP("52.8.190.235"),
		Port: NatServerTestPort,
	})

	if err != nil {
		panic(err)
	}

	ips := nat.ExternalIP()

	client := &NatPeer{
		conn:      c,
		peerID:    os.Args[1],
		privateIP: ips[0],
	}

	return client
}

func (peer *NatPeer) KeepAlive() {

	request := &nat_pb.NatRequest{
		NodeId:      peer.peerID,
		PrivateIp:   peer.privateIP,
		PrivatePort: NatHoleListeningPort,
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

		response := &nat_pb.NatResponse{}
		if err := proto.Unmarshal(responseData[:hasRead], response); err != nil {
			fmt.Println("failed to unmarshal nat response data", err)
			continue
		}

		fmt.Println("get response data from nat natServer:", response)

		time.Sleep(20 * time.Second)
	}
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

		client.KeepAlive()
	} else if os.Args[1] == "-s" {

		server := NewServer()
		server.Processing()
	} else {
		fmt.Println("unknown action")
	}
}
