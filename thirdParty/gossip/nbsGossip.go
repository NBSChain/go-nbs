package gossip

import (
	"github.com/NBSChain/go-nbs/utils"
	"sync"
)

type nbsGossip struct {
}

var instance *nbsGossip
var once sync.Once
var logger = utils.GetLogInstance()

func GetGossipInstance() BasicProtocol {

	once.Do(func() {
		instance = newNbsGossip()
	})

	return instance
}

func newNbsGossip() *nbsGossip {

	gossipObj := &nbsGossip{}

	return gossipObj
}
