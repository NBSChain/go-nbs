package main

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network"
	"github.com/NBSChain/go-nbs/thirdParty/account"
	"github.com/NBSChain/go-nbs/thirdParty/gossip"
	"github.com/NBSChain/go-nbs/utils"
	"os"
	"time"
)

type holeTest struct {
	conn *network.NbsUdpConn
}

func punchAHole() {

	ht := &holeTest{}

	nodeId := account.GetAccountInstance().GetPeerID()

	network.GetInstance().StartUp(nodeId)

	gossip.GetGossipInstance().Online(nodeId)

	fmt.Println("I started:->:", nodeId)

	var peerId string
	if len(os.Args) == 2 {
		peerId = os.Args[1]
		port := utils.GetConfig().GossipCtrlPort
		conn, err := network.GetInstance().Connect(nodeId, peerId, "", port)
		if err != nil {
			panic(err)
		}
		ht.conn = conn

		go ht.reading()
		go ht.writing()
	}

	<-make(chan struct{})
}

func (ht *holeTest) reading() {
	for {
		buffer := make([]byte, utils.NormalReadBuffer)
		n, peerAddr, err := ht.conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("read data from peer:", n, peerAddr)
	}
}

func (ht *holeTest) writing() {
	for {
		msg := []byte("sssssss")
		if _, err := ht.conn.Write(msg); err != nil {
			fmt.Println("write message error:", err)
		}

		time.Sleep(5 * time.Second)
	}
}
