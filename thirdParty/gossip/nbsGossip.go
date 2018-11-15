package gossip

import (
	"github.com/NBSChain/go-nbs/thirdParty/gossip/memership"
	"github.com/NBSChain/go-nbs/utils"
	"sync"
)

var (
	instance *nbsGossip
	once     sync.Once
	logger   = utils.GetLogInstance()
)

type nbsGossip struct {
	peerId        string
	memberManager *memership.MemManager
}

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

func (manager *nbsGossip) StartUp(peerId string) error {

	manager.peerId = peerId

	memberNode := memership.NewMemberNode(peerId)
	memberNode.InitNode()
	manager.memberManager = memberNode

	logger.Info("gossip service start up......")

	return nil
}
