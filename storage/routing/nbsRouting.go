package routing

import (
	"context"
	"fmt"
	"github.com/W-B-S/nbs-node/storage/network"
	"github.com/libp2p/go-libp2p-peer"
	"github.com/libp2p/go-libp2p-peerstore"
	"sync"
)

type NbsDHT struct {
	peerId peer.ID
}

var instance *NbsDHT
var once sync.Once
var parentContext context.Context

func GetInstance() Routing {
	once.Do(func() {
		parentContext = context.Background()
		router, err := newNbsDht()
		if err != nil {
			panic(err)
		}
		fmt.Printf("router start to run......\n")
		instance = router
	})

	return instance
}

func newNbsDht() (*NbsDHT, error) {

	network := network.GetInstance()

	distributeTable := &NbsDHT{
		peerId: peer.ID(network.GetId()),
	}

	return distributeTable, nil
}

//----------->routing interface implementation<-----------//
func (*NbsDHT) Ping(context.Context, peer.ID) error {
	return nil
}

func (*NbsDHT) FindPeer(context.Context, peer.ID) (peerstore.PeerInfo, error) {
	return peerstore.PeerInfo{}, nil
}

func (*NbsDHT) PutValue(context.Context, string, []byte) error {
	return nil
}

func (*NbsDHT) GetValue(context.Context, string) ([]byte, error) {
	return nil, nil
}

func (router *NbsDHT) Run() {

	select {
	case <-parentContext.Done():
		fmt.Printf("routing node done!\n")
	}
}
