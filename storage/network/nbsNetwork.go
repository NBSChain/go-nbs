package network

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
	//TODO:: replace this gx repository
	"gx/ipfs/QmNmJZL7FQySMtE2BQuLMuZg2EB2CLEunJJUSVSc9YnnbV/go-libp2p-host"
	"sync"
)

type NbsNetwork struct {
	Host    host.Host
	Context context.Context
}

var once sync.Once
var instance *NbsNetwork

func GetInstance() Network {

	once.Do(func() {
		instance = newNetwork()
	})

	return instance
}

func newNetwork() *NbsNetwork {

	network := &NbsNetwork{
		Context: context.Background(),
	}

	//--->convert gx version control to github one's
	newHost, err := libp2p.New(network.Context)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Create host  %s\n", newHost.Addrs())

	network.Host = newHost

	return network
}

func (network *NbsNetwork) GetId() string {
	return string(network.Host.ID())
}
