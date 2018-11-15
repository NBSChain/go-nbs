package network

import (
	"context"
	"github.com/NBSChain/go-nbs/storage/network/nat"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/utils"
	"sync"
)

type nbsNetwork struct {
	Context    context.Context
	natManager *nat.NbsNatManager
	netWorkId  string
	addresses  *net_pb.NbsAddress
	peerStore  map[string]*net_pb.NbsAddress
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
		Context:   context.Background(),
		peerStore: make(map[string]*net_pb.NbsAddress),
	}

	return network
}
