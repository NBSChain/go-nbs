package application

import (
	"context"
	"github.com/NBSChain/go-nbs/storage/bitswap"
	"github.com/NBSChain/go-nbs/storage/merkledag"
	"github.com/NBSChain/go-nbs/storage/network"
	"github.com/NBSChain/go-nbs/storage/routing"
	"github.com/NBSChain/go-nbs/thirdParty/account"
	"github.com/NBSChain/go-nbs/thirdParty/gossip"
	"github.com/NBSChain/go-nbs/utils"
	"sync"
)

var logger = utils.GetLogInstance()

type NbsApplication struct {
	context context.Context
	nodeId  string
}

var instance *NbsApplication
var once sync.Once

func GetInstance() Application {

	once.Do(func() {

		app, err := newApplication()
		if err != nil {
			panic(err)
		}

		logger.Info("--->Create application to run......\n")

		instance = app
	})

	return instance
}

func newApplication() (*NbsApplication, error) {

	acc := account.GetAccountInstance()

	return &NbsApplication{
		context: context.Background(),
		nodeId:  acc.GetPeerID(),
	}, nil
}

func (app *NbsApplication) GetNodeId() string {
	return app.nodeId
}

func (app *NbsApplication) Start() error {

	logger.Info("Application starting......")

	merkledag.GetDagInstance()

	bitswap.GetSwapInstance()

	routing.GetInstance()

	network.GetInstance()

	gossip.GetGossipInstance()

	return nil
}

func (app *NbsApplication) ReloadForNewAccount() error {

	app.nodeId = account.GetAccountInstance().GetPeerID()

	if err := network.GetInstance().StartUp(app.nodeId); err != nil {
		return err
	}

	return nil
}
