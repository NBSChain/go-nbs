package main

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/storage/network/shareport"
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
	//p2pConn     *net.UDPConn
	isApplier bool
}

var hostAddress = &net.UDPAddr{
	IP:   net.ParseIP("52.8.190.235"),
	Port: 8001,
}

func NewPeer() *NatPeer {

	c, err := shareport.ListenUDP("udp4", "0.0.0.0:7001")
	if err != nil {
		panic(err)
	}

	dialHost := c.LocalAddr().String()
	host, port, _ := net.SplitHostPort(dialHost)

	client := &NatPeer{
		conn:        c,
		peerID:      os.Args[2],
		privateIP:   host,
		privatePort: port,
	}

	fmt.Println("dialed", dialHost, c.RemoteAddr())

	//l, err := shareport.ListenUDP("udp4", "0.0.0.0:7001")
	//if err != nil {
	//	fmt.Println("********************dial share port failed:-> ", err)
	//	panic(err)
	//}
	//client.p2pConn = l
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

		peer.conn.SetDeadline(time.Now().Add(time.Second * 5))

		if no, err := peer.conn.WriteTo(requestData, hostAddress); err != nil || no == 0 {
			fmt.Println("failed to send nat request to natServer ", err, no)
			continue
		}

		responseData := make([]byte, 2048)
		hasRead, peerAddr, err := peer.conn.ReadFrom(responseData)
		if err != nil {
			fmt.Println("failed to read nat response from natServer", err)
			continue
		}

		response := &nat_pb.Response{}
		if err := proto.Unmarshal(responseData[:hasRead], response); err != nil {
			fmt.Println("failed to unmarshal nat response data", err)
			continue
		}

		fmt.Println("----->", peerAddr, response)

		switch response.MsgType {
		case nat_pb.NatMsgType_BootStrapReg:
			time.Sleep(10 * time.Second)
		case nat_pb.NatMsgType_Connect:
			go peer.connectToPeers(response.ConnRes)
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
		panic(err)
	}

	if _, err := peer.conn.WriteTo(requestData, hostAddress); err != nil {
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
		fmt.Println("********************err hole message:->", err)
		return
	}

	dstPort, _ := strconv.Atoi(response.PublicPort)

	peerAddr := &net.UDPAddr{
		IP:   net.ParseIP(response.PublicIp),
		Port: dstPort,
	}

	for {
		var no int

		if no, err = peer.conn.WriteTo(data, peerAddr); err != nil || no == 0 {
			fmt.Println("********************failed to make p2p connection:-> ", err, no)
			break
		}

		fmt.Println("\n\n********************write data len:->", no, peerAddr)

		if peer.isApplier {
			go peer.p2pReader()
			break
		} else {
			time.Sleep(3 * time.Second)
		}

	}

}

func (peer *NatPeer) p2pReader() {

	for {

		//time.Sleep(time.Second * 4)

		fmt.Println("********************start reading********************")

		readBuff := make([]byte, 2048)

		peer.conn.SetReadDeadline(time.Now().Add(time.Second * 5))

		hasRead, peerAddr, err := peer.conn.ReadFrom(readBuff)
		if err != nil {
			fmt.Println("****************reading from:->", err, peerAddr)
			continue
		}

		fmt.Println("****************has read :->", hasRead)

		holeMsg := &nat_pb.Response{}

		proto.Unmarshal(readBuff[:hasRead], holeMsg)

		fmt.Println("********************unmarshal:->", holeMsg)
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
