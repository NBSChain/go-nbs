package gossip

import (
	"fmt"
	"github.com/NBSChain/go-nbs/thirdParty/account"
	"github.com/NBSChain/go-nbs/thirdParty/gossip/memership"
	"github.com/NBSChain/go-nbs/utils"
	"sync"
)

var (
	instance        *nbsGossip
	once            sync.Once
	logger          = utils.GetLogInstance()
	ServiceNotValid = fmt.Errorf("gossip service is not on right now")
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

	peerId := account.GetAccountInstance().GetPeerID()

	if peerId == "" {
		logger.Warning("no account right now, so the message gossip can't setup")
	} else {
		if err := gossipObj.Online(peerId); err != nil {
			logger.Error("gossip start failed:", err)
		}
	}

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

func (manager *nbsGossip) Offline() error {

	if manager.memberManager == nil {
		return fmt.Errorf("gossip is done already")
	}

	if err := manager.memberManager.DestroyNode(); err != nil {
		return err
	}
	manager.memberManager = nil
	return nil
}

func (manager *nbsGossip) Online(peerId string) error {
	if manager.memberManager != nil {
		return fmt.Errorf("gossip is already running")
	}

	manager.peerId = peerId

	memberNode := memership.NewMemberNode(peerId)
	manager.memberManager = memberNode
	if err := memberNode.InitNode(); err != nil {
		return err
	}

	logger.Info("gossip service start up......")

	return nil
}

func (manager *nbsGossip) IsOnline() bool {
	return manager.memberManager != nil
}

func (manager *nbsGossip) ShowInputViews() ([]*memership.ViewNode, error) {

	if manager.memberManager == nil {
		return nil, ServiceNotValid
	}

	var views []*memership.ViewNode
	for _, item := range manager.memberManager.InputView {
		views = append(views, item)
	}
	return views, nil
}

func (manager *nbsGossip) ShowOutputViews() ([]*memership.ViewNode, error) {
	if manager.memberManager == nil {
		return nil, ServiceNotValid
	}
	var views []*memership.ViewNode
	for _, item := range manager.memberManager.PartialView {
		views = append(views, item)
	}
	return views, nil
}
