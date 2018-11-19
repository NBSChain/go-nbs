package network

import (
	"context"
	"github.com/NBSChain/go-nbs/storage/network/nat"
	"github.com/NBSChain/go-nbs/utils"
	"sync"
)

type nbsNetwork struct {
	Context    context.Context
	natManager *nat.Manager
	netWorkId  string
}

var once sync.Once
var instance *nbsNetwork
var logger = utils.GetLogInstance()

func GetInstance() Network {

	once.Do(func() {
		instance = newNetwork()
	})

	return instance
}

func newNetwork() *nbsNetwork {

	network := &nbsNetwork{
		Context: context.Background(),
	}

	return network
}
