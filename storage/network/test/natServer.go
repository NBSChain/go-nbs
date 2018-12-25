package main

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/gogo/protobuf/proto"
	"net"
	"time"
)

const NatServerTestPort = 8001

type HostBehindNat struct {
	UpdateTIme time.Time
	PubAddr    net.Addr
	PriAddr    string
	KAConn     *net.TCPConn
}

type NatServer struct {
	server   *net.TCPListener
	natCache map[string]*HostBehindNat
}

func NewServer() *NatServer {

	s, err := net.ListenTCP("tcp4", &net.TCPAddr{
		Port: NatServerTestPort,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(s.Addr().String())

	server := &NatServer{
		server:   s,
		natCache: make(map[string]*HostBehindNat),
	}

	return server
}

func (s *NatServer) Processing() {

	fmt.Println("start to run......")

	for {
		conn, err := s.server.AcceptTCP()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		if err := conn.SetKeepAlive(true); err != nil {
			fmt.Println("set ka err:->", err)
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
			fmt.Println("keep alive ack err:->", err)
			return err
		}

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
		fmt.Println("forward apply from:->", app.Public)

		if _, err := conn.Write(data); err != nil {
			panic(err)
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
		fmt.Println("forward this confirm to target:->", ack.TargetId, item.PubAddr, ack)

		if _, err := conn.Write(data); err != nil {
			panic(err)
		}

	case nbsnet.NatBlankKA:
		data, _ := proto.Marshal(msg)
		if _, err := conn.Write(data); err != nil {
			panic(err)
		}

	default:
		fmt.Println("unknown msg type")
	}
	return nil
}

func (s *NatServer) connThread(conn *net.TCPConn) {
	buffer := make([]byte, utils.NormalReadBuffer)
	defer conn.Close()
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("read data from socket err:->", err)
			continue
		}

		request := &net_pb.NatMsg{}
		if err := proto.Unmarshal(buffer[:n], request); err != nil {
			fmt.Println("can't parse the nat request:->", err)
			continue
		}

		fmt.Println("\nNat request:->", request, nbsnet.ConnString(conn))

		if err := s.processMsg(request, conn); err != nil {
			fmt.Println("process err:->", err)
		}
	}
}
