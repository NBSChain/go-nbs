package main

import (
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/storage/network/shareport"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"net"
	"os"
	"sync"
	"time"
)

var logger = utils.GetLogInstance()

type NatPeer struct {
	peerID        string
	keepAliveConn *net.TCPConn
	isApplier     bool
	startLock     sync.Mutex
	startConn     *net.UDPConn
	lisLock       sync.Mutex
	lisConn       *net.UDPConn
}

//var natServer = &net.TCPAddr{Port: CtrlMsgPort, IP: net.ParseIP("52.8.190.235")}
//var natHelpServer = &net.TCPAddr{Port: HoleHelpPort, IP: net.ParseIP("52.8.190.235")}

//var natHelpServer = &net.TCPAddr{Port: HoleHelpPort, IP: net.ParseIP("103.45.98.72")}
var natServer = &net.TCPAddr{Port: CtrlMsgPort, IP: net.ParseIP("192.168.103.101")}
var natHelpServer = &net.TCPAddr{Port: HoleHelpPort, IP: net.ParseIP("192.168.103.101")}
var locServer = "0.0.0.0:7001"

func NewPeer() *NatPeer {

	c1, err := net.DialTimeout("tcp4", natServer.String(), time.Second*2)
	if err != nil {
		panic(err)
	}

	logger.Debug("dialed----1--->", c1.LocalAddr().String(), c1.RemoteAddr())
	client := &NatPeer{
		keepAliveConn: c1.(*net.TCPConn),
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
			logger.Debug("get keep alive ack:->", msg)
		case nbsnet.NatDigApply:
			app := msg.DigApply
			logger.Debug("receive dig application:->", app)

			ack := &net_pb.NatMsg{
				Typ: nbsnet.NatDigConfirm,
				DigConfirm: &net_pb.DigConfirm{
					TargetId: app.FromId,
				},
			}

			data, _ := proto.Marshal(ack)
			conn, err := shareport.DialUDP("udp4", locServer, natHelpServer.String())
			if err != nil {
				panic(err)
			}

			if _, err := conn.Write(data); err != nil {
				panic(err)
			}

			go peer.readingDigOut(conn, "[222222]")

			ip, port, _ := nbsnet.SplitHostPort(app.Public)
			target := &net.UDPAddr{
				IP:   net.ParseIP(ip),
				Port: int(port),
			}
			go func() {
				for {
					logger.Debug("write to applier's peer:->", target)
					if _, err := peer.lisConn.WriteToUDP(data, target); err != nil {
						panic(err)
					}
					time.Sleep(time.Second * 2)
				}
			}()

		case nbsnet.NatDigConfirm:

			ack := msg.DigConfirm
			logger.Debug("dig confirmed:->", ack)
			//peer.startLock.Lock()
			//locAddr := peer.startConn.LocalAddr().String()
			//peer.startLock.Unlock()

			go func() {
				digMsg := &net_pb.NatMsg{
					Typ: nbsnet.NatDigOut,
					Seq: time.Now().Unix(),
				}
				data, _ := proto.Marshal(digMsg)
				ip, port, _ := nbsnet.SplitHostPort(ack.Public)
				target := &net.UDPAddr{
					IP:   net.ParseIP(ip),
					Port: int(port),
				}
				for {
					logger.Debug("send direct from me:->", target, peer.startConn.LocalAddr().String())
					if _, err := peer.startConn.WriteToUDP(data, target); err != nil {
						panic(err)
					}

					time.Sleep(time.Second * 2)
				}
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
			logger.Debug("---gun1 write---->", err, no)
			panic(err)
		}

		time.Sleep(time.Second * 60)
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

	conn, err := shareport.DialUDP("udp4", "", natHelpServer.String())
	if err != nil {
		panic(err)
	}
	if _, err := conn.Write(requestData); err != nil {
		logger.Debug(err)
		return
	}
	sConn, err := shareport.ListenUDP("udp4", conn.LocalAddr().String())
	if err != nil {
		panic(err)
	}
	conn.Close()

	peer.startLock.Lock()
	peer.startConn = sConn
	peer.startLock.Unlock()

	go peer.readingDigOut(sConn, "[111111]")

	logger.Debug("tel peer I want to make a connection:->", nbsnet.ConnString(conn))
}

func (peer *NatPeer) readingDigOut(conn *net.UDPConn, socketId string) {
	for {
		buffer := make([]byte, utils.NormalReadBuffer)
		n, peerAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			logger.Error("reading dig out err:->", err)
			return
		}
		msg := &net_pb.NatMsg{}
		proto.Unmarshal(buffer[:n], msg)
		logger.Infof("----%s-->get reading out message:%V->", socketId, peerAddr, msg)
	}
}

func (peer *NatPeer) ListenServiceBehindNat() {
	conn, err := shareport.ListenUDP("udp4", locServer)
	if err != nil {
		panic(err)
	}
	peer.lisConn = conn
	logger.Debug("listen at:->", peer.lisConn.LocalAddr().String())
	for {
		buffer := make([]byte, utils.NormalReadBuffer)
		n, peerAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			panic(err)
		}
		msg := &net_pb.NatMsg{}
		proto.Unmarshal(buffer[:n], msg)
		logger.Debug("----listening service---hole punching success:->", peerAddr, msg)
	}
}

func (peer *NatPeer) digDig(fromHost, targetHost string) {

	conn, err := shareport.DialUDP("udp4", fromHost, targetHost)
	if err != nil {
		panic(err)
	}
	go peer.readingDigOut(conn, "[333333]")

	digMsg := &net_pb.NatMsg{
		Typ: nbsnet.NatDigOut,
		Seq: time.Now().Unix(),
	}
	data, _ := proto.Marshal(digMsg)

	for {
		logger.Debug("dig a hole on peer's nat server:->", nbsnet.ConnString(conn))
		if _, err := conn.Write(data); err != nil {
			panic(err)
		}
		time.Sleep(time.Second * 2)
	}
}

func natTool() {

	if len(os.Args) < 2 {
		logger.Debug("input run mode -c or -s")
		os.Exit(1)
	}

	if os.Args[1] == "-c" {

		if len(os.Args) < 3 {
			logger.Debug("input this peer's ID")
			os.Exit(1)
		}

		client := NewPeer()

		go client.sendKA()
		go client.runLoop()

		go client.ListenServiceBehindNat()

		if len(os.Args) == 4 {
			client.punchAHole(os.Args[3])
		}

		<-make(chan struct{})

	} else if os.Args[1] == "-s" {
		server := NewServer()
		server.CtlMsg()

	} else {
		logger.Debug("unknown action")
	}
}
