package routing

import (
	"context"
	"github.com/NBSChain/go-nbs/storage/application/dataStore"
	"github.com/NBSChain/go-nbs/storage/network"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/libp2p/go-libp2p-peer"
	"github.com/libp2p/go-libp2p-peerstore"
	"sync"
)

type NbsDHT struct {
	peerId peer.ID
	localDataStore	dataStore.DataStore
}

var instance 		*NbsDHT
var once 		sync.Once
var parentContext 	context.Context
var logger 		= utils.GetLogInstance()

func GetInstance() Routing {
	once.Do(func() {
		parentContext = context.Background()
		router, err := newNbsDht()
		if err != nil {
			panic(err)
		}
		logger.Info("router start to run......\n")
		instance = router
	})

	return instance
}

func newNbsDht() (*NbsDHT, error) {

	network.GetInstance()

	distributeTable := &NbsDHT{
	}

	return distributeTable, nil
}

//----------->routing interface implementation<-----------//
func (*NbsDHT) Ping(peer peerstore.PeerInfo) Pong{
	return nil
}

func (*NbsDHT) FindPeer(key string) ([]peerstore.PeerInfo, error){
	return nil, nil
}

func (*NbsDHT) PutValue(key string, value []byte) chan error {
	return nil
}

func (*NbsDHT) GetValue(peer []peerstore.PeerInfo,  key string) ([]byte, []peerstore.PeerInfo, error) {
	return nil, nil, nil
}

func (dht *NbsDHT) saveData(key string, value []byte){
	dht.localDataStore.Put(key, value)
}