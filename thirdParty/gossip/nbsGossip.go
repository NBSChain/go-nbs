package gossip

import (
	"github.com/NBSChain/go-nbs/utils"
	"sync"
)

var (
	instance *nbsGossip
	once     sync.Once
	logger   = utils.GetLogInstance()
)

type nbsGossip struct {
}

func GetGossipInstance() BasicProtocol {

	once.Do(func() {
		instance = newNbsGossip()
	})

	return instance
}

func newNbsGossip() *nbsGossip {

	gossipObj := &nbsGossip{}

	go gossipObj.registerToNetwork()

	return gossipObj
}

func (manager *nbsGossip) registerToNetwork() {

}

/*****************************************************************
*
*		interface implementations
*
*****************************************************************/
func (manager *nbsGossip) Publish(channel string, message []byte) error {
	return nil
}

func (manager *nbsGossip) Subscribe(channel string) error {
	return nil
}

func (manager *nbsGossip) AllPeers(channel string, depth int) ([]string, []string) {
	return nil, nil
}

func (manager *nbsGossip) AllMyTopics() []string {
	return nil
}

func (manager *nbsGossip) Unsubscribe(channel string) error {
	return nil
}
