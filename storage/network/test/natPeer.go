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
	privatePort string
	p2pConn     *net.UDPConn
}

var natServerAddr = &net.UDPAddr{
	IP:   net.ParseIP("47.52.172.234"),
	Port: NatServerTestPort,
}

func NewPeer() *NatPeer {

	//ip := nat.ExternalIP()[0]
	//c, err := reuseport.NewReusablePortPacketConn("udp4", ip+":0")
	//172.168.30.18//52.8.190.235//47.52.172.234

	c, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   net.ParseIP("47.52.172.234"),
		Port: NatServerTestPort,
	})
	if err != nil {
		panic(err)
	}

	host, port, _ := net.SplitHostPort(c.LocalAddr().String())

	client := &NatPeer{
		conn:        c,
		peerID:      os.Args[2],
		privateIP:   host,
		privatePort: port,
	}

	return client
}

func (peer *NatPeer) runLoop() {

	request := &nat_pb.NatRequest{
		MsgType: nat_pb.NatMsgType_BootStrapReg,
		BootRegReq: &nat_pb.BootNatRegReq{
			NodeId:      peer.peerID,
			PrivateIp:   peer.privateIP,
			PrivatePort: peer.privatePort,
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
		hasRead, err := peer.conn.Read(responseData)
		if err != nil {
			fmt.Println("failed to read nat response from natServer", err)
			continue
		}

		response := &nat_pb.Response{}
		if err := proto.Unmarshal(responseData[:hasRead], response); err != nil {
			fmt.Println("failed to unmarshal nat response data", err)
			continue
		}

		switch response.MsgType {
		case nat_pb.NatMsgType_BootStrapReg:
			time.Sleep(20 * time.Second)
		case nat_pb.NatMsgType_Connect:
			go peer.connectToPeers(response.ConnRes)
		}
	}
}

func (peer *NatPeer) punchAHole(targetId string) {

	time.Sleep(5 * time.Second)

	inviteRequest := &nat_pb.NatConReq{
		FromPeerId: peer.peerID,
		ToPeerId:   targetId,
	}

	request := &nat_pb.NatRequest{
		MsgType: nat_pb.NatMsgType_Connect,
		ConnReq: inviteRequest,
	}

	requestData, err := proto.Marshal(request)
	if err != nil {
		panic(err)
	}

	if _, err := peer.conn.WriteTo(requestData, natServerAddr); err != nil {
		panic(err)
	}
}

func (peer *NatPeer) connectToPeers(response *nat_pb.NatConRes) {

	holeMsg := &nat_pb.Response{
		MsgType: nat_pb.NatMsgType_Connect,
		ConnRes: response,
	}

	data, err := proto.Marshal(holeMsg)
	if err != nil {
		fmt.Println("hole message:", err)
		return
	}

	desPort, _ := strconv.Atoi(response.PublicPort)
	srcPort, _ := strconv.Atoi(peer.privatePort)

	peerAddr := &net.UDPAddr{
		IP:   net.ParseIP(response.PublicIp),
		Port: desPort,
	}

	c, err := net.DialUDP("udp4", &net.UDPAddr{
		Port: srcPort,
	}, peerAddr)

	if err != nil {
		panic(err)
	}

	peer.p2pConn = c
	for {
		var no int
		if no, err = peer.p2pConn.WriteTo(data, peerAddr); err != nil || no == 0 {
			fmt.Println("failed to make p2p connection-> ", err, no)
			time.Sleep(10 * time.Second)
			continue
		}

		fmt.Println("\n\n---write data len:->", no, peerAddr.String(), c.LocalAddr().String())

		peer.p2pReader()
	}
}

func (peer *NatPeer) p2pReader() {

	fmt.Println("-----------p2pReader----------------")

	readBuff := make([]byte, 2048)
	hasRead, peerAddr, err := peer.p2pConn.ReadFrom(readBuff)
	if err != nil {
		fmt.Println("hole message read:", err)
		return
	}

	fmt.Println(readBuff[hasRead])

	fmt.Printf("****************hole message :->(%d)from :%v->", hasRead, peerAddr)
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
