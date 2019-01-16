package gossip

import (
	"fmt"
	"github.com/NBSChain/go-nbs/thirdParty/account"
	"github.com/NBSChain/go-nbs/thirdParty/gossip/memership"
	"github.com/NBSChain/go-nbs/thirdParty/gossip/message"
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
	msgManager    *message.MsgManager
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
			logger.Error("gossip online failed:->", err)
		}
	}

	return gossipObj
}

/*****************************************************************
*
*		interface implementations
*
*****************************************************************/
func (manager *nbsGossip) Publish(channel string, msg string) error {
	if manager.memberManager == nil {
		return fmt.Errorf("gossip isn't online right now")
	}

	manager.memberManager.FanOut(channel, []byte(msg), message.MsgTypePlainTxt)
	return nil
}

func (manager *nbsGossip) Subscribe(channel string) (chan *message.MsgEntity, error) {
	return manager.msgManager.NewSub(channel)
}

func (manager *nbsGossip) AllPeers(channel string, depth int) ([]string, []string) {
	return nil, nil
}

func (manager *nbsGossip) AllMyTopics() []string {
	return nil
}

func (manager *nbsGossip) Unsubscribe(channel string) {
	manager.msgManager.CancelSub(channel)

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
	manager.msgManager = message.NewMsgManager()
	memberNode.AppMsgHub = manager.msgManager.MsgReceiver

	manager.memberManager = memberNode
	if err := memberNode.InitNode(); err != nil {
		return err
	}

	if err := memberNode.RegisterMySelf(); err != nil {
		logger.Warning("failed to register myself:->", err)
		return err
	}
	logger.Info("gossip service start up......")

	return nil
}

func (manager *nbsGossip) IsOnline() bool {
	return manager.memberManager != nil
}

func (manager *nbsGossip) ShowInputViews() ([]string, error) {

	if manager.memberManager == nil {
		return nil, ServiceNotValid
	}

	return manager.memberManager.GetViewsInfo(manager.memberManager.InputView), nil
}

func (manager *nbsGossip) ShowOutputViews() ([]string, error) {
	if manager.memberManager == nil {
		return nil, ServiceNotValid
	}

	return manager.memberManager.GetViewsInfo(manager.memberManager.PartialView), nil
}

func (manager *nbsGossip) ClearInputViews() int {
	if manager.memberManager == nil {
		return 0
	}
	return manager.memberManager.RemoveViewsInfo(manager.memberManager.InputView)
}

func (manager *nbsGossip) ClearOutputViews() int {
	if manager.memberManager == nil {
		return 0
	}
	return manager.memberManager.RemoveViewsInfo(manager.memberManager.PartialView)
}
