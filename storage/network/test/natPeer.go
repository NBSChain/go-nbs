//+build !windows

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

	c1, err := shareport.DialUDP("udp4", "0.0.0.0:7001", "192.168.103.155:8001")
	if err != nil {
		panic(err)
	}

	fmt.Println("dialed----1--->", c1.LocalAddr().String(), c1.RemoteAddr())

	dialHost := c1.LocalAddr().String()
	host, port, _ := net.SplitHostPort(dialHost)

	client := &NatPeer{
		keepAliveConn: c1,
		peerID:        os.Args[2],
		privateIP:     host,
		privatePort:   port,
		receivingHub:  l,
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

	go peer.readingHub()
	go peer.readingKA()

	for {
		if no, err := peer.keepAliveConn.Write(requestData); err != nil || no == 0 {
			fmt.Println("---gun1 write---->", err, no)
		}
		time.Sleep(time.Second * 10)
	}
}

func (peer *NatPeer) readingKA() {

	for {
		responseData := make([]byte, 2048)
		hasRead, err := peer.keepAliveConn.Read(responseData)
		if err != nil {
			fmt.Println("---keep alive read failed --->", err)
			continue
		}

		response := &nat_pb.Response{}
		if err := proto.Unmarshal(responseData[:hasRead], response); err != nil {
			fmt.Println("keep alive response data", err)
			continue
		}

		switch response.MsgType {
		case nat_pb.NatMsgType_BootStrapReg:
			fmt.Println("---keep alive reading--->", response)
		case nat_pb.NatMsgType_Connect:
			go peer.connectToPeers(response.ConnRes)
			fmt.Println("++++keep alive+++++++peer connect invite+++++++++++>", response)
		case nat_pb.NatMsgType_Ping:
			fmt.Println("\n$$$$keep alive$$$$$$$hole punching$$$$$$$$$$$$$$", response)
		}
	}
}

func (peer *NatPeer) readHoleMessage() {
	for {
		responseData := make([]byte, 2048)
		hasRead, err := peer.holePunchConn.Read(responseData)
		if err != nil {
			fmt.Println("---keep alive read failed --->", err)
			continue
		}

		response := &nat_pb.Response{}
		if err := proto.Unmarshal(responseData[:hasRead], response); err != nil {
			fmt.Println("keep alive response data", err)
			continue
		}

		fmt.Println("\n$$$$$$$$$$$hole punching$$$$$$$$$$$$$$", response)
	}
}

func (peer *NatPeer) readingHub() {

	for {
		responseData := make([]byte, 2048)
		hasRead, peerAddr, err := peer.receivingHub.ReadFrom(responseData)
		if err != nil {
			fmt.Println("---reading hub-->failed ReadFrom->", err)
			continue
		}

		response := &nat_pb.Response{}
		if err := proto.Unmarshal(responseData[:hasRead], response); err != nil {
			fmt.Println("---reading hub-->failed response data", err)
			continue
		}

		switch response.MsgType {
		case nat_pb.NatMsgType_BootStrapReg:
			fmt.Println("---reading hub--->", peerAddr, response)
		case nat_pb.NatMsgType_Connect:
			go peer.connectToPeers(response.ConnRes)
			fmt.Println("++++++reading hub+++++peer connect invite+++++++++++>", peerAddr, response)
		case nat_pb.NatMsgType_Ping:
			fmt.Println("\n$$$$$reading hub$$$$$$hole punching$$$$$$$$$$$$$$", peerAddr, response)
		}
	}
}

func (peer *NatPeer) punchAHole(targetId string) {

	peer.isApplier = true

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
		fmt.Println(err)
		return
	}

	if _, err := peer.keepAliveConn.Write(requestData); err != nil {
		fmt.Println(err)
		return
	}
}

func (peer *NatPeer) connectToPeers(response *nat_pb.NatConRes) {

	holeMsg := &nat_pb.Response{
		MsgType: nat_pb.NatMsgType_Ping,
		Pong: &nat_pb.NatPing{
			Ping: peer.peerID,
		},
	}

	data, err := proto.Marshal(holeMsg)
	if err != nil {
		fmt.Println("***************err hole message:->", err)
		return
	}

	//port, err := strconv.Atoi(response.PublicPort)
	//peerAddr := &net.UDPAddr{
	//	IP:   net.ParseIP(response.PublicIp),
	//	Port: port,
	//}

	p, err := shareport.DialUDP("udp4", "0.0.0.0:7001", response.PublicIp+":"+response.PublicPort)
	if err != nil {
		fmt.Println("******dial hole peer failed***", err)
		return
	}
	peer.holePunchConn = p

	go peer.readHoleMessage()

	for {
		no, err := peer.holePunchConn.Write(data)
		if err != nil {
			fmt.Println("*************punch conn*******failed to write hole message:-> ", err, no)
			break
		}

		fmt.Println("punch write done:->", no)

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
