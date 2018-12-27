package main

import (
	"github.com/NBSChain/go-nbs/storage/network/nat"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/storage/network/shareport"
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
	peerID string
	sync.Mutex
	ctrlChannel *net.TCPConn
	holeConn    *net.UDPConn
	lisConn     *net.UDPConn
	allMyHosts  []string
	netType     net_pb.NetWorkType
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
		allMyHosts:  make([]string, 0),
	}

	conn, err := shareport.ListenUDP("udp4", locServer)
	if err != nil {
		panic(err)
	}
	logger.Debug("listen at:->", conn.LocalAddr().String())

	locStr := conn.LocalAddr().String()
	request := &net_pb.NatMsg{
		Typ: nbsnet.NatKeepAlive,
		KeepAlive: &net_pb.KeepAlive{
			NodeId:  client.peerID,
			PriAddr: locStr,
		},
	}
	requestData, _ := proto.Marshal(request)

	hosts, ports := make(map[string]struct{}, 0), make(map[string]struct{}, 0)
	successHlepers := 0
	for _, addr := range ipHelpers {
		conn.SetDeadline(time.Now().Add(time.Second * 4))

		if _, err := conn.WriteToUDP(requestData, addr); err != nil {
			logger.Error("get all my ips err:->", err)
			continue
		}
		buffer := make([]byte, utils.NormalReadBuffer)
		n, peerAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			logger.Error("get ip err:->", err, addr.String())
			continue
		}

		msg := &net_pb.NatMsg{}
		proto.Unmarshal(buffer[:n], msg)

		if msg.KeepAlive == nil {
			logger.Error("this is not my want:->", msg, peerAddr)
			continue
		}

		pubHost := msg.KeepAlive.PubAddr
		ip, port, _ := net.SplitHostPort(pubHost)
		logger.Debug("get ip success:->", ip, port)
		hosts[ip] = struct{}{}
		ports[port] = struct{}{}
		successHlepers++
	}

	if successHlepers < 2 {
		panic("too small server to contact")
	}

	var networkType = nbsnet.SigIpSigPort
	if len(ports) == 1 {
		if len(hosts) > 1 {
			networkType = nbsnet.MulIpSigPort
			for ip := range hosts {
				client.allMyHosts = append(client.allMyHosts, ip)
			}
		}
	} else {
		panic("we don't support multi port model right now.")
	}

	client.netType = networkType
	logger.Info("network type:->", networkType)
	conn.SetDeadline(nat.NoTimeOut)
	client.lisConn = conn

	go client.ListenService()

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

			go peer.digOut(app)

			conn, err := shareport.DialUDP("udp4", locServer, natHelpServer.String())
			if err != nil {
				panic(err)
			}
			logger.Debug("receive dig application:->", app)
			ack := &net_pb.NatMsg{
				Typ: nbsnet.NatDigConfirm,
				DigConfirm: &net_pb.DigConfirm{
					TargetId: app.FromId,
					NtType:   peer.netType,
					PubIps:   peer.allMyHosts,
				},
			}

			data, _ := proto.Marshal(ack)

			if _, err := conn.Write(data); err != nil {
				panic(err)
			}
			conn.Close()
		case nbsnet.NatDigConfirm:

			ack := msg.DigConfirm
			logger.Debug("dig confirmed:->", ack)

			peer.Lock()
			digAddr := peer.holeConn.LocalAddr().(*net.UDPAddr)
			peer.holeConn.Close()
			peer.Unlock()

			ip, port, _ := nbsnet.SplitHostPort(ack.Public)
			if ack.NtType == nbsnet.SigIpSigPort {
				conn, err := net.DialUDP("udp4", digAddr, &net.UDPAddr{
					IP:   net.ParseIP(ip),
					Port: int(port),
				})
				if err != nil {
					panic(err)
				}
				logger.Debug("dial up success and send ka")
				go peer.blankKeepAlvie(conn)
				continue
			}

			ch := make(chan *net.UDPConn)
			for _, ips := range ack.PubIps {
				tarAddr := nbsnet.JoinHostPort(ips, port)
				go peer.findTheRightConn(digAddr.String(), tarAddr, ch)
			}

			select {
			case c := <-ch:
				logger.Debug("pick out the right conn and send ka:->", nbsnet.ConnString(c))
				go peer.blankKeepAlvie(c)

			case <-time.After(time.Second * 6):
				panic("failed dial up:->")
			}
		}
	}
}

func (peer *NatPeer) findTheRightConn(fromAddr, toAddr string, ch chan *net.UDPConn) {
	conn, err := shareport.DialUDP("udp4", fromAddr, toAddr)
	if err != nil {
		logger.Warning("dial up a bad conn:->", toAddr)
		return
	}
	msg := net_pb.NatMsg{
		Typ: nbsnet.NatBlankKA,
		Seq: time.Now().Unix(),
	}
	data, _ := proto.Marshal(&msg)
	if _, err := conn.Write(data); err != nil {
		logger.Warning("find a bad conn:->", toAddr)
		return
	}
	conn.SetDeadline(time.Now().Add(time.Second * 6))
	buffer := make([]byte, utils.NormalReadBuffer)
	if _, err := conn.Read(buffer); err != nil {
		logger.Warning("this conn failed:->", nbsnet.ConnString(conn), err)
		return
	}

	conn.SetDeadline(nat.NoTimeOut)
	ch <- conn
}

func (peer *NatPeer) blankKeepAlvie(conn *net.UDPConn) {

	for {
		msg := net_pb.NatMsg{
			Typ: nbsnet.NatBlankKA,
			Seq: time.Now().Unix(),
		}
		data, _ := proto.Marshal(&msg)

		if _, err := conn.Write(data); err != nil {
			panic(err)
		}
		logger.Debug("blank keep alive for hole tunnel:->", nbsnet.ConnString(conn))
		time.Sleep(time.Second * 5)
	}
}
func (peer *NatPeer) digOut(apply *net_pb.DigApply) {

	digMsg := &net_pb.NatMsg{
		Typ: nbsnet.NatDigOut,
		Seq: time.Now().Unix(),
	}
	data, _ := proto.Marshal(digMsg)

	if apply.NtType == nbsnet.SigIpSigPort {

		conn, err := shareport.DialUDP("udp4", locServer, apply.Public)
		if err != nil {
			panic(err)
		}

		for i := 0; i < 3; i++ {
			logger.Debug("signal ip send direct from me:->", nbsnet.ConnString(conn))
			if _, err := conn.Write(data); err != nil {
				panic(err)
			}
		}

		conn.Close()
		return
	}

	_, port, _ := nbsnet.SplitHostPort(apply.Public)
	for _, ip := range apply.PubIps {

		target := nbsnet.JoinHostPort(ip, port)

		conn, err := shareport.DialUDP("udp4", locServer, target)
		if err != nil {
			logger.Warning("dig out for multi ip target err:->", err)
			continue
		}

		for i := 0; i < 3; i++ {
			logger.Debug("multi ips send direct from me:->", nbsnet.ConnString(conn))
			if _, err := conn.Write(data); err != nil {
				panic(err)
			}
		}
		conn.Close()
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

	msg := &net_pb.NatMsg{
		Typ: nbsnet.NatDigApply,
		DigApply: &net_pb.DigApply{
			TargetId: targetId,
			FromId:   peer.peerID,
			NtType:   peer.netType,
			PubIps:   peer.allMyHosts,
		},
	}

	requestData, _ := proto.Marshal(msg)
	sConn, err := net.DialUDP("udp4", nil, natHelpServer)
	if err != nil {
		panic(err)
	}

	if _, err := sConn.Write(requestData); err != nil {
		panic(err)
	}
	peer.Lock()
	peer.holeConn = sConn
	peer.Unlock()

	logger.Debug("tel peer I want to make a connection:->", nbsnet.ConnString(sConn))
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
		logger.Debug("----listening service--punching success:->", peerAddr, msg)
		if _, err := peer.lisConn.WriteToUDP(buffer[:n], peerAddr); err != nil {
			panic(err)
		}
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
