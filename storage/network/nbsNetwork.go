package network

import (
	"context"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-host"
	"sync"
)

type NbsNetwork struct {
	Host    host.Host
	Context context.Context
}

var once sync.Once
var instance *NbsNetwork
var logger = utils.GetLogInstance()

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
	logger.Info("Create host  %s\n", newHost.Addrs())

	network.Host = newHost

	return network
}

func (network *NbsNetwork) GetId() string {
	return string(network.Host.ID())
}
