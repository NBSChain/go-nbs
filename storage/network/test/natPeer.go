package main

import (
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

var logger = utils.GetLogInstance()

const CtrlMsgPort = 8001
const HoleHelpPort = 8002

type NatPeer struct {
	peerID      string
	ctrlChannel *net.TCPConn
	startConn   *net.UDPConn
	lisConn     *net.UDPConn
	hostLock    sync.Mutex
	allMyHosts  map[string]struct{}
}

//var natServer = &net.TCPAddr{Port: CtrlMsgPort, IP: net.ParseIP("52.8.190.235")}
//var natHelpServer = &net.UDPAddr{Port: HoleHelpPort, IP: net.ParseIP("52.8.190.235")}

var natServer = &net.TCPAddr{Port: CtrlMsgPort, IP: net.ParseIP("103.45.98.72")}
var natHelpServer = &net.UDPAddr{Port: HoleHelpPort, IP: net.ParseIP("103.45.98.72")}

var ipHelpers = []*net.UDPAddr{&net.UDPAddr{Port: HoleHelpPort, IP: net.ParseIP("103.45.98.72")},
	&net.UDPAddr{Port: HoleHelpPort, IP: net.ParseIP("52.8.190.235")},
	&net.UDPAddr{Port: HoleHelpPort, IP: net.ParseIP("47.52.172.234")},
	&net.UDPAddr{Port: HoleHelpPort, IP: net.ParseIP("13.57.241.215")}}

//var natServer = &net.TCPAddr{Port: CtrlMsgPort, IP: net.ParseIP("192.168.103.101")}
//var natHelpServer = &net.UDPAddr{Port: HoleHelpPort, IP: net.ParseIP("192.168.103.101")}

var locServer = "0.0.0.0:7001"

func NewPeer() *NatPeer {

	c1, err := net.DialTimeout("tcp4", natServer.String(), time.Second*2)
	if err != nil {
		panic(err)
	}

	logger.Debug("dialed----1--->", c1.LocalAddr().String(), c1.RemoteAddr())
	client := &NatPeer{
		ctrlChannel: c1.(*net.TCPConn),
		peerID:      os.Args[2],
		allMyHosts:  make(map[string]struct{}),
	}

	conn, err := net.ListenUDP("udp4", &net.UDPAddr{
		Port: 7001,
	})
	if err != nil {
		panic(err)
	}
	client.lisConn = conn
	logger.Debug("listen at:->", client.lisConn.LocalAddr().String())
	go client.ListenService()
	locStr := client.lisConn.LocalAddr().String()
	request := &net_pb.NatMsg{
		Typ: nbsnet.NatKeepAlive,
		KeepAlive: &net_pb.KeepAlive{
			NodeId:  client.peerID,
			PriAddr: locStr,
		},
	}
	requestData, _ := proto.Marshal(request)

	for _, addr := range ipHelpers {
		if _, err := conn.WriteToUDP(requestData, addr); err != nil {
			logger.Error("get all my ips err:->", err)
			continue
		}
	}

	return client
}

func (peer *NatPeer) runLoop() {
	for {
		buffer := make([]byte, utils.NormalReadBuffer)
		n, err := peer.ctrlChannel.Read(buffer)
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

			ips := make([]string, 0)
			peer.hostLock.Lock()
			for ip := range peer.allMyHosts {
				ips = append(ips, ip)
			}
			peer.hostLock.Unlock()
			ack := &net_pb.NatMsg{
				Typ: nbsnet.NatDigConfirm,
				DigConfirm: &net_pb.DigConfirm{
					TargetId: app.FromId,
					AllIps:   ips,
				},
			}

			data, _ := proto.Marshal(ack)
			if _, err := peer.lisConn.WriteToUDP(data, natHelpServer); err != nil {
				panic(err)
			}
			go peer.udpKA(peer.lisConn)

			go peer.dididididid(peer.lisConn, app.Public, app.AllIps)

		case nbsnet.NatDigConfirm:

			ack := msg.DigConfirm
			logger.Debug("dig confirmed:->", ack)

			go peer.dididididid(peer.startConn, ack.Public, ack.AllIps)
		}
	}
}
func (peer *NatPeer) dididididid(conn *net.UDPConn, public string, allIps []string) {
	digMsg := &net_pb.NatMsg{
		Typ: nbsnet.NatDigOut,
		Seq: time.Now().Unix(),
	}
	data, _ := proto.Marshal(digMsg)

	for _, ip := range allIps {

		_, port, _ := nbsnet.SplitHostPort(public)
		target := &net.UDPAddr{
			IP:   net.ParseIP(ip),
			Port: int(port),
		}

		for {
			logger.Debug("send direct from me:->", target, conn.LocalAddr().String())
			if _, err := conn.WriteToUDP(data, target); err != nil {
				panic(err)
			}

			time.Sleep(time.Second * 2)
		}
	}
}

func (peer *NatPeer) sendKA() {
	locStr := peer.ctrlChannel.LocalAddr().String()
	request := &net_pb.NatMsg{
		Typ: nbsnet.NatKeepAlive,
		KeepAlive: &net_pb.KeepAlive{
			NodeId:  peer.peerID,
			PriAddr: locStr,
		},
	}
	requestData, _ := proto.Marshal(request)
	for {
		no, err := peer.ctrlChannel.Write(requestData)
		if err != nil || no == 0 {
			logger.Debug("---gun1 write---->", err, no)
			panic(err)
		}

		time.Sleep(time.Second * 60)
	}
}

func (peer *NatPeer) punchAHole(targetId string) {
	ips := make([]string, 0)
	peer.hostLock.Lock()
	for ip := range peer.allMyHosts {
		ips = append(ips, ip)
	}
	peer.hostLock.Unlock()
	msg := &net_pb.NatMsg{
		Typ: nbsnet.NatDigApply,
		DigApply: &net_pb.DigApply{
			TargetId: targetId,
			FromId:   peer.peerID,
			AllIps:   ips,
		},
	}
	requestData, _ := proto.Marshal(msg)

	sConn, err := net.ListenUDP("udp4", &net.UDPAddr{
		Port: 11223,
	})
	if err != nil {
		panic(err)
	}

	if _, err := sConn.WriteToUDP(requestData, natHelpServer); err != nil {
		panic(err)
	}

	peer.startConn = sConn

	go peer.udpKA(sConn)
	go peer.readingDigOut(sConn, "[111111]")

	logger.Debug("tel peer I want to make a connection:->", sConn.LocalAddr().String())
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
		logger.Infof("----%s-->get reading out message:%s->", socketId, peerAddr, msg)

		//if _, err := conn.WriteToUDP(buffer[:n], peerAddr); err != nil {
		//	panic(err)
		//}
	}
}

func (peer *NatPeer) ListenService() {

	for {
		buffer := make([]byte, utils.NormalReadBuffer)
		n, peerAddr, err := peer.lisConn.ReadFromUDP(buffer)
		if err != nil {
			panic(err)
		}
		msg := &net_pb.NatMsg{}
		proto.Unmarshal(buffer[:n], msg)
		logger.Debug("----listening service---hole punching success:->", peerAddr, msg)
		switch msg.Typ {
		case nbsnet.NatKeepAlive:
			ka := msg.KeepAlive
			ip, _, _ := nbsnet.SplitHostPort(ka.PubAddr)
			peer.hostLock.Lock()
			peer.allMyHosts[ip] = struct{}{}
			peer.hostLock.Unlock()
		default:
			logger.Warning("maybe success:--->")
		}
	}
}

func (peer *NatPeer) udpKA(conn *net.UDPConn) {

	digMsg := &net_pb.NatMsg{
		Typ: nbsnet.NatBlankKA,
		Seq: time.Now().Unix(),
	}
	data, _ := proto.Marshal(digMsg)

	for {
		logger.Debug("udp ka:->", conn.LocalAddr().String())
		if _, err := conn.WriteToUDP(data, natHelpServer); err != nil {
			panic(err)
		}
		time.Sleep(time.Second * 2)
	}
}

var probeServer = []*net.UDPAddr{
	&net.UDPAddr{
		IP:   net.ParseIP("103.45.98.72"),
		Port: 8002,
	},
	&net.UDPAddr{
		IP:   net.ParseIP("52.8.190.235"),
		Port: 8002,
	},
	&net.UDPAddr{
		IP:   net.ParseIP("47.52.172.234"),
		Port: 8002,
	},
	&net.UDPAddr{
		IP:   net.ParseIP("13.57.241.215"),
		Port: 8002,
	},
}

type Probe struct {
	probeConn *net.UDPConn
}

func NewProbe() *Probe {
	return &Probe{}
}

func (p *Probe) Probe() {

	idx_1, _ := strconv.Atoi(os.Args[2])
	s_1 := probeServer[idx_1]

	conn, err := net.ListenUDP("udp4", &net.UDPAddr{
		Port: 7002,
	})
	if err != nil {
		panic(err)
	}
	logger.Debug("start liten:->", conn.LocalAddr().String())
	p.probeConn = conn

	go p.Reading()

	locStr := conn.LocalAddr().String()
	request := &net_pb.NatMsg{
		Typ: nbsnet.NatKeepAlive,
		KeepAlive: &net_pb.KeepAlive{
			NodeId:  "xxx---probe---XXX",
			PriAddr: locStr,
		},
	}
	requestData, _ := proto.Marshal(request)

	if _, err := conn.WriteToUDP(requestData, s_1); err != nil {
		panic(err)
	}
	logger.Debug("write one server:->", s_1.String())
	if len(os.Args) > 3 {
		idx_2, _ := strconv.Atoi(os.Args[3])
		s_2 := probeServer[idx_2]
		if _, err := conn.WriteToUDP(requestData, s_2); err != nil {
			panic(err)
		}
		logger.Debug("write second server:->", s_2.String())

	}
}
func (p *Probe) Reading() {

	for {
		buffer := make([]byte, utils.NormalReadBuffer)
		n, peerAdd, err := p.probeConn.ReadFromUDP(buffer)
		if err != nil {
			panic(err)
		}
		msg := &net_pb.NatMsg{}
		proto.Unmarshal(buffer[:n], msg)
		logger.Debug("probe reading :->", peerAdd, msg)
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
		go client.ListenService()

		if len(os.Args) == 4 {
			client.punchAHole(os.Args[3])
		}

		<-make(chan struct{})

	} else if os.Args[1] == "-s" {
		server := NewServer()
		server.CtlMsg()

	} else if os.Args[1] == "-p" {
		probe := NewProbe()
		probe.Probe()
		<-make(chan struct{})

	} else {
		logger.Debug("unknown action")
	}
}
