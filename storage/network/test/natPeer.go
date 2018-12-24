package main

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/storage/network/shareport"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/NBSChain/go-nbs/utils/crypto"
	"github.com/golang/protobuf/proto"
	"net"
	"os"
	"sync"
	"time"
)

type NatPeer struct {
	peerID        string
	keepAliveConn *net.UDPConn
	isApplier     bool
	startLock     sync.Mutex
	startConn     *net.UDPConn
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
			conn, err := shareport.DialUDP("udp4", locServer, natServer.String())
			if err != nil {
				panic(err)
			}

			n, err := conn.Write(data)
			if err != nil {
				panic(err)
			}

			buf := make([]byte, utils.NormalReadBuffer)
			nn, err := conn.Read(buf)
			if err != nil {
				panic(err)
			}

			confirmAck := &net_pb.NatMsg{}
			proto.Unmarshal(buf[:nn], confirmAck)
			fmt.Println("my confirm message confirmed.", confirmAck)

			go peer.digDig(app.Public)

			go func() {
				keepAlive := &net_pb.NatMsg{
					Typ: nbsnet.NatBlankKA,
					ConnKA: &net_pb.ConnKA{
						KA: crypto.MD5SS(time.Now().String()),
					},
				}
				kaData, _ := proto.Marshal(keepAlive)
				if _, err := conn.Write(kaData); err != nil {
					panic(err)
				}
				time.Sleep(time.Second * 5)
				fmt.Println("keep the punch confirm channel alive:->", nbsnet.ConnString(conn))
			}()
			go func() {
				buffer := make([]byte, utils.NormalReadBuffer)
				n, err := conn.Read(buffer[:n])
				if err != nil {
					panic(err)
				}

				readMsg := &net_pb.NatMsg{}
				proto.Unmarshal(buffer[:n], readMsg)
				fmt.Println("222222211122212->", readMsg)

				conn.Close()
			}()

		case nbsnet.NatDigConfirm:

			ack := msg.DigConfirm
			fmt.Println("dig confirmed:->", ack)
			peer.startLock.Lock()
			locAddr := peer.startConn.LocalAddr().String()
			peer.startConn.Close()
			peer.startLock.Unlock()

			conn, err := shareport.DialUDP("udp4", locAddr, ack.Public)
			if err != nil {
				panic(err)
			}
			ka := &net_pb.NatMsg{
				Typ: nbsnet.NatBlankKA,
				ConnKA: &net_pb.ConnKA{
					KA: "=========this is a keep alive for nat connection:=>----------",
				},
			}

			data, _ := proto.Marshal(ka)
			fmt.Println("dial hole in back :->", nbsnet.ConnString(conn))
			if _, err := conn.Write(data); err != nil {
				panic(err)
			}

			go func() {
				buffer := make([]byte, utils.NormalReadBuffer)
				_, peerAddr, err := conn.ReadFromUDP(buffer)

				if err != nil {
					panic(err)
				}
				fmt.Println("-0---0000000-", peerAddr, buffer)
			}()

		}
	}
}

func (peer *NatPeer) sendKA() {
	locStr := peer.keepAliveConn.LocalAddr().String()
	request := &net_pb.NatMsg{
		Typ: nbsnet.NatKeepAlive,
		KeepAlive: &net_pb.KeepAlive{
			NodeId:  peer.peerID,
			PriAddr: locStr,
		},
	}
	requestData, _ := proto.Marshal(request)

	for {
		no, err := peer.keepAliveConn.Write(requestData)
		if err != nil || no == 0 {
			fmt.Println("---gun1 write---->", err, no)
			panic(err)
		}

		time.Sleep(time.Second * 5)
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

	conn, err := shareport.DialUDP("udp4", "", natServer.String())
	if err != nil {
		panic(err)
	}
	if _, err := conn.Write(requestData); err != nil {
		fmt.Println(err)
		return
	}
	peer.startLock.Lock()
	peer.startConn = conn
	peer.startLock.Unlock()

	fmt.Println("tel peer I want to make a connection:->", nbsnet.ConnString(conn))

	buffer := make([]byte, utils.NormalReadBuffer)
	n, err := conn.Read(buffer)
	if err != nil {
		panic(err)
	}
	readMsg := &net_pb.NatMsg{}
	proto.Unmarshal(buffer[:n], readMsg)
	fmt.Println("The server answer me:->", readMsg)

	for {
		keepAlive := &net_pb.NatMsg{
			Typ: nbsnet.NatBlankKA,
			ConnKA: &net_pb.ConnKA{
				KA: crypto.MD5SS(time.Now().String()),
			},
		}
		kaData, _ := proto.Marshal(keepAlive)
		if _, err := conn.Write(kaData); err != nil {
			panic(err)
		}
		time.Sleep(time.Second * 5)
		fmt.Println("keep the punch apply channel alive:->", nbsnet.ConnString(conn))
	}
}

func (peer *NatPeer) ListenHoleMsg() {
	conn, err := shareport.ListenUDP("udp4", locServer)
	if err != nil {
		panic(err)
	}
	for {
		buffer := make([]byte, utils.NormalReadBuffer)
		n, peerAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			panic(err)
		}
		msg := &net_pb.NatMsg{}
		proto.Unmarshal(buffer[:n], msg)
		fmt.Println("111111111---hole punching success:->", peerAddr, msg)
	}
}

func (peer *NatPeer) digDig(targetHost string) {

	conn, err := shareport.DialUDP("udp4", locServer, targetHost)
	if err != nil {
		panic(err)
	}

	digMsg := &net_pb.NatMsg{
		Typ: nbsnet.NatDigOut,
		Seq: time.Now().Unix(),
	}
	data, _ := proto.Marshal(digMsg)

	for i := 0; i < 3; i++ {
		fmt.Println("dig a hole on peer's nat server:->", nbsnet.ConnString(conn))
		if _, err := conn.Write(data); err != nil {
			panic(err)
		}
	}
	buffer := make([]byte, utils.NormalReadBuffer)
	n, err := conn.Read(buffer)
	if err != nil {
		panic(err)
	}

	readMsg := &net_pb.NatMsg{}
	proto.Unmarshal(buffer[:n], readMsg)
	fmt.Println("3333333333->", readMsg)
	conn.Close()
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

		go client.sendKA()
		go client.runLoop()

		go client.ListenHoleMsg()

		if len(os.Args) == 4 {
			client.punchAHole(os.Args[3])
		}

		<-make(chan struct{})

	} else if os.Args[1] == "-s" {
		server := NewServer()
		server.Processing()

	} else {
		fmt.Println("unknown action")
	}
}
