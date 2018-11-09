package main

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/golang/protobuf/proto"
	"github.com/libp2p/go-reuseport"
	"net"
	"os"
	"strconv"
	"time"
)

type NatPeer struct {
	peerID      string
	conn        net.Conn
	privateIP   string
	privatePort string
	p2pConn     net.PacketConn
}

func NewPeer() *NatPeer {

	//47.52.172.234//52.8.190.235
	c, err := reuseport.Dial("udp4", "0.0.0.0:0", "52.8.190.235:8001")
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

	fmt.Println("dialed", c.LocalAddr(), c.RemoteAddr())

	l, err := reuseport.ListenPacket("udp4", host+":"+port)
	if err != nil {
		panic(err)
	}

	client.p2pConn = l

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

		peer.conn.SetDeadline(time.Now().Add(time.Second * 10))

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

		fmt.Println(response)

		switch response.MsgType {
		case nat_pb.NatMsgType_BootStrapReg:
			time.Sleep(10 * time.Second)
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

	if _, err := peer.conn.Write(requestData); err != nil {
		panic(err)
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
		fmt.Println("hole message:", err)
		return
	}

	dstPort, _ := strconv.Atoi(response.PublicPort)

	peerAddr := &net.UDPAddr{
		IP:   net.ParseIP(response.PublicIp),
		Port: dstPort,
	}

	go peer.p2pReader()

	for {
		var no int

		time.Sleep(3 * time.Second)

		if no, err = peer.p2pConn.WriteTo(data, peerAddr); err != nil || no == 0 {
			fmt.Println("failed to make p2p connection-> ", err, no)
			break
		}

		fmt.Println("\n\n---write data len:->", no, peer.p2pConn.LocalAddr().String(), peerAddr)

	}
}

func (peer *NatPeer) p2pReader() {

	for {

		time.Sleep(time.Second * 4)

		fmt.Println("-----------p2pReader----------------")

		readBuff := make([]byte, 2048)

		peer.p2pConn.SetReadDeadline(time.Now().Add(time.Second * 5))

		hasRead, peerAddr, err := peer.p2pConn.ReadFrom(readBuff)
		if err != nil {
			fmt.Println("****************hole message read:", err)
			continue
		}

		fmt.Println("****************hole message ->", hasRead, peerAddr)

		holeMsg := &nat_pb.Response{}

		proto.Unmarshal(readBuff[:hasRead], holeMsg)

		fmt.Println(holeMsg)
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
