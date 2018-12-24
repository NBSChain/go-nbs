package main

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/storage/network/shareport"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"net"
	"os"
	"time"
)

type NatPeer struct {
	peerID        string
	keepAliveConn *net.UDPConn
	isApplier     bool
	waitStr       string
	listenConn    *net.UDPConn
}

var natServer = &net.UDPAddr{Port: NatServerTestPort, IP: net.ParseIP("52.8.190.235")}

//var natServer = &net.UDPAddr{Port: NatServerTestPort, IP: net.ParseIP("192.168.103.101")}
var locServer = "0.0.0.0:7001"

func NewPeer() *NatPeer {

	c1, err := net.DialUDP("udp4", nil, natServer)
	if err != nil {
		panic(err)
	}

	fmt.Println("dialed----1--->", c1.LocalAddr().String(), c1.RemoteAddr())
	client := &NatPeer{
		keepAliveConn: c1,
		peerID:        os.Args[2],
	}

	return client
}

func (peer *NatPeer) runLoop() {

	go peer.sendKA(peer.keepAliveConn)
	go peer.Listening()

	if len(os.Args) == 4 {
		go peer.punchAHole(os.Args[3])
	}

	for {
		buffer := make([]byte, utils.NormalReadBuffer)
		n, err := peer.keepAliveConn.Read(buffer)
		if err != nil {
			panic(err)
		}
		msg := &net_pb.NatMsg{}
		proto.Unmarshal(buffer[:n], msg)
		switch msg.Typ {
		case nbsnet.NatKeepAlive:
		case nbsnet.NatDigApply:
			app := msg.DigApply
			fmt.Println("receive dig application:->", app)

			ack := &net_pb.NatMsg{
				Typ: nbsnet.NatDigConfirm,
				DigConfirm: &net_pb.DigConfirm{
					TargetId: app.FromId,
				},
			}

			data, _ := proto.Marshal(ack)

			if _, err := peer.listenConn.WriteToUDP(data, natServer); err != nil {
				panic(err)
			}

			go peer.digDig(app.Public)

		case nbsnet.NatDigConfirm:

			ack := msg.DigConfirm
			fmt.Println("dig confirmed:->", ack)
			locAddr := peer.waitStr
			conn, err := shareport.DialUDP("udp4", locAddr, ack.Public)
			if err != nil {
				panic(err)
			}

			go func() {
				for {
					fmt.Println("dial hole in back :->", nbsnet.ConnString(conn))
					if _, err := conn.Write(buffer[:n]); err != nil {
						panic(err)
					}
					time.Sleep(time.Second)
				}
			}()

			go func() {
				buffer := make([]byte, utils.NormalReadBuffer)
				if _, err := conn.Read(buffer); err != nil {
					panic(err)
				}
				fmt.Println("-0-----", buffer)
			}()

		}
	}
}

func (peer *NatPeer) sendKA(conn *net.UDPConn) {
	locStr := conn.LocalAddr().String()
	for {
		request := &net_pb.NatMsg{
			Typ: nbsnet.NatKeepAlive,
			KeepAlive: &net_pb.KeepAlive{
				NodeId:  peer.peerID,
				PriAddr: locStr,
			},
		}
		requestData, _ := proto.Marshal(request)

		if no, err := conn.Write(requestData); err != nil || no == 0 {
			fmt.Println("---gun1 write---->", err, no)
			panic(err)
		}

		time.Sleep(time.Second * 10)
	}
}

func (peer *NatPeer) punchAHole(targetId string) {
	msg := &net_pb.NatMsg{
		Typ: nbsnet.NatDigApply,
		DigApply: &net_pb.DigApply{
			TargetId: targetId,
			FromId:   peer.peerID,
		},
	}
	requestData, _ := proto.Marshal(msg)

	conn, err := net.DialUDP("udp4", nil, natServer)
	if err != nil {
		panic(err)
	}
	if _, err := conn.Write(requestData); err != nil {
		fmt.Println(err)
		return
	}
	peer.waitStr = conn.LocalAddr().String()

	fmt.Println("tel peer I want to make a connection:->", nbsnet.ConnString(conn))
}

func (peer *NatPeer) Listening2(conn *net.UDPConn) {
	for {
		buffer := make([]byte, utils.NormalReadBuffer)
		n, err := conn.Read(buffer)
		if err != nil {
			panic(err)
		}
		msg := &net_pb.NatMsg{}
		proto.Unmarshal(buffer[:n], msg)
		println("2222222hole punching success:->", msg)

		conn.Write(buffer)
	}
}

func (peer *NatPeer) Listening() {

	lisConn, err := shareport.ListenUDP("udp4", locServer)
	if err != nil {
		panic(err)
	}

	peer.listenConn = lisConn

	for {
		buffer := make([]byte, utils.NormalReadBuffer)
		n, peerAddr, err := lisConn.ReadFromUDP(buffer)
		if err != nil {
			panic(err)
		}
		msg := &net_pb.NatMsg{}
		proto.Unmarshal(buffer[:n], msg)
		println("111111111hole punching success:->", msg, peerAddr)
	}
}

func (peer *NatPeer) digDig(targetHost string) {

	digMsg := &net_pb.NatMsg{
		Typ: nbsnet.NatDigOut,
		Seq: time.Now().Unix(),
	}
	data, _ := proto.Marshal(digMsg)

	host, port, _ := nbsnet.SplitHostPort(targetHost)
	addr := &net.UDPAddr{
		IP:   net.ParseIP(host),
		Port: int(port),
	}
	for i := 0; i < 5; i++ {
		println("dig a hole on peer's nat server:->", peer.listenConn.LocalAddr().String())
		if _, err := peer.listenConn.WriteToUDP(data, addr); err != nil {
			panic(err)
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func natTool() {

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
		client.runLoop()
	} else if os.Args[1] == "-s" {
		server := NewServer()
		server.Processing()

	} else {
		fmt.Println("unknown action")
	}
}
