package main

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/storage/network/shareport"
	"github.com/golang/protobuf/proto"
	"net"
	"os"
	"time"
)

type NatPeer struct {
	peerID        string
	keepAliveConn *net.UDPConn
	holePunchConn *net.UDPConn
	privateIP     string
	privatePort   string
	receivingHub  *net.UDPConn
	isApplier     bool
}

func NewPeer() *NatPeer {

	l, err := shareport.ListenUDP("udp4", "0.0.0.0:7001")
	if err != nil {
		fmt.Println("********************dial share port failed:-> ", err)
		panic(err)
	}

	c1, err := shareport.DialUDP("udp4", "0.0.0.0:7001", "52.8.190.235:8001")
	if err != nil {
		panic(err)
	}

	dialHost := c1.LocalAddr().String()
	host, port, _ := net.SplitHostPort(dialHost)

	client := &NatPeer{
		keepAliveConn: c1,
		peerID:        os.Args[2],
		privateIP:     host,
		privatePort:   port,
		receivingHub:  l,
	}

	fmt.Println("dialed----1--->", c1.LocalAddr().String(), c1.RemoteAddr())

	return client
}

func (peer *NatPeer) runLoop() {

	request := &net_pb.NatRequest{
		MsgType: net_pb.NatMsgType_BootStrapReg,
		BootRegReq: &net_pb.BootNatRegReq{
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

	go peer.readingHub()

	for {
		if no, err := peer.keepAliveConn.Write(requestData); err != nil || no == 0 {
			fmt.Println("---gun1---->", err, no)
		}
		time.Sleep(time.Second * 10)
	}
}

func (peer *NatPeer) readingHub() {

	for {
		responseData := make([]byte, 2048)
		hasRead, peerAddr, err := peer.receivingHub.ReadFrom(responseData)
		if err != nil {
			fmt.Println("failed to read nat response from natServer", err)
			continue
		}

		response := &net_pb.Response{}
		if err := proto.Unmarshal(responseData[:hasRead], response); err != nil {
			fmt.Println("failed to unmarshal nat response data", err)
			continue
		}

		switch response.MsgType {
		case net_pb.NatMsgType_BootStrapReg:
			fmt.Println("---reading hub--->", peerAddr, response)
		case net_pb.NatMsgType_Connect:
			fmt.Println("*******hole punch message--->", peerAddr, response)
			go peer.connectToPeers(response.ConnRes)
		case net_pb.NatMsgType_Ping:
			fmt.Println("$$$$$$$$$$$punching message---->", peerAddr, response)
		}
	}
}

func (peer *NatPeer) punchAHole(targetId string) {

	peer.isApplier = true

	inviteRequest := &net_pb.NatConReq{
		FromPeerId: peer.peerID,
		ToPeerId:   targetId,
	}

	request := &net_pb.NatRequest{
		MsgType: net_pb.NatMsgType_Connect,
		ConnReq: inviteRequest,
	}

	requestData, err := proto.Marshal(request)
	if err != nil {
		panic(err)
	}

	if _, err := peer.keepAliveConn.Write(requestData); err != nil {
		panic(err)
	}
}

func (peer *NatPeer) connectToPeers(response *net_pb.NatConRes) {

	holeMsg := &net_pb.Response{
		MsgType: net_pb.NatMsgType_Ping,
		Pong: &net_pb.NatPing{
			Ping: peer.peerID,
		},
	}

	requestData, err := proto.Marshal(holeMsg)
	if err != nil {
		fmt.Println("********************err hole message:->", err)
		return
	}

	c2, err := shareport.DialUDP("udp4", "0.0.0.0:7001", response.PublicIp+":"+response.PublicPort)
	if err != nil {
		fmt.Println("---gun2---dial failed->", err)
		return
	}
	peer.holePunchConn = c2
	fmt.Println("dialed----2--->", c2.LocalAddr().String(), c2.RemoteAddr())

	for {
		if no, err := peer.holePunchConn.Write(requestData); err != nil || no == 0 {
			fmt.Println("---gun2---->", err, no)
		}

		time.Sleep(10 * time.Second)
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

	} else if os.Args[1] == "-p" {
		probe := NewSimplePeer()
		probe.probe()
	} else {
		fmt.Println("unknown action")
	}
}
