package main

import (
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/gogo/protobuf/proto"
	"net"
	"sync"
	"time"
)

const CtrlMsgPort = 8001
const HoleHelpPort = 8002

type HostBehindNat struct {
	UpdateTIme time.Time
	PubAddr    net.Addr
	PriAddr    string
	KAConn     *net.TCPConn
}

type NatServer struct {
	server     *net.TCPListener
	HoleHelper *net.UDPConn
	cacheLock  sync.Mutex
	natCache   map[string]*HostBehindNat
}

func NewServer() *NatServer {

	s, err := net.ListenTCP("tcp4", &net.TCPAddr{
		Port: CtrlMsgPort,
	})
	if err != nil {
		panic(err)
	}

	logger.Debug("control message server:->", s.Addr().String())

	server := &NatServer{
		server:   s,
		natCache: make(map[string]*HostBehindNat),
	}

	h, err := net.ListenUDP("udp4", &net.UDPAddr{
		Port: HoleHelpPort,
	})
	if err != nil {
		panic(err)
	}
	server.HoleHelper = h
	go server.HoleDigger()

	return server
}
func (s *NatServer) HoleDigger() {
	for {
		buffer := make([]byte, utils.NormalReadBuffer)
		n, peerAddr, err := s.HoleHelper.ReadFromUDP(buffer)
		if err != nil {
			panic(err)
		}

		msg := &net_pb.NatMsg{}
		proto.Unmarshal(buffer[:n], msg)
		s.cacheLock.Lock()

		switch msg.Typ {
		case nbsnet.NatDigApply:
			app := msg.DigApply
			item, ok := s.natCache[app.TargetId]
			if !ok {
				panic(app.TargetId)
			}
			app.Public = peerAddr.String()
			data, _ := proto.Marshal(msg)
			if _, err := item.KAConn.Write(data); err != nil {
				panic(err)
			}
			logger.Debug("forward apply from:->", app.Public)

		case nbsnet.NatDigConfirm:
			ack := msg.DigConfirm
			item, ok := s.natCache[ack.TargetId]
			if !ok {
				panic(ack.TargetId)
			}

			ack.Public = peerAddr.String()
			data, _ := proto.Marshal(msg)
			if _, err := item.KAConn.Write(data); err != nil {
				panic(err)
			}
			logger.Debug("forward confirm :->", ack.Public)

		case nbsnet.NatBlankKA:
			logger.Debug("get udp ka:->", peerAddr)
			//if _, err := s.HoleHelper.WriteToUDP(buffer[:n], peerAddr); err != nil {
			//	panic(err)
			//}
		case nbsnet.NatKeepAlive:
			ack := msg.KeepAlive
			ack.PubAddr = peerAddr.String()
			data, _ := proto.Marshal(msg)
			if _, err := s.HoleHelper.WriteToUDP(data, peerAddr); err != nil {
				logger.Debug("keep alive ack err:->", err)
			}
		}

		s.cacheLock.Unlock()
	}
}

func (s *NatServer) CtlMsg() {

	logger.Debug("start to run......")

	for {
		conn, err := s.server.AcceptTCP()
		if err != nil {
			logger.Debug(err.Error())
			continue
		}
		if err := conn.SetKeepAlive(true); err != nil {
			logger.Debug("set ka err:->", err)
			continue
		}

		go s.connThread(conn)
	}
}

func (s *NatServer) processMsg(msg *net_pb.NatMsg, conn *net.TCPConn) error {

	peerAddr := conn.RemoteAddr()

	switch msg.Typ {
	case nbsnet.NatKeepAlive:
		ack := msg.KeepAlive
		nodeId := ack.NodeId
		item, ok := s.natCache[nodeId]
		if !ok {
			item := &HostBehindNat{
				UpdateTIme: time.Now(),
				PubAddr:    conn.RemoteAddr(),
				KAConn:     conn,
				PriAddr:    msg.KeepAlive.PriAddr,
			}
			s.natCache[nodeId] = item
		} else {
			item.UpdateTIme = time.Now()
		}

		ack.PubAddr = peerAddr.String()
		data, _ := proto.Marshal(msg)
		if _, err := conn.Write(data); err != nil {
			logger.Debug("keep alive ack err:->", err)
			return err
		}

	case nbsnet.NatDigConfirm:
		ack := msg.DigConfirm
		item, ok := s.natCache[ack.TargetId]
		if !ok {
			panic(ok)
		}
		ack.Public = peerAddr.String()
		data, _ := proto.Marshal(msg)
		if _, err := item.KAConn.Write(data); err != nil {
			panic(err)
		}
		logger.Debug("forward this confirm to target:->", ack.TargetId, item.PubAddr, ack)

		if _, err := conn.Write(data); err != nil {
			panic(err)
		}

	case nbsnet.NatBlankKA:
		data, _ := proto.Marshal(msg)
		if _, err := conn.Write(data); err != nil {
			panic(err)
		}

	default:
		logger.Debug("unknown msg type")
	}
	return nil
}

func (s *NatServer) connThread(conn *net.TCPConn) {
	buffer := make([]byte, utils.NormalReadBuffer)
	defer conn.Close()
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			logger.Debug("read data from socket err:->", err)
			return
		}

		request := &net_pb.NatMsg{}
		if err := proto.Unmarshal(buffer[:n], request); err != nil {
			logger.Debug("can't parse the nat request:->", err)
			continue
		}

		logger.Debug("\nNat request:->", request, nbsnet.ConnString(conn))
		s.cacheLock.Lock()
		if err := s.processMsg(request, conn); err != nil {
			logger.Debug("process err:->", err)
		}
		s.cacheLock.Unlock()
	}
}
