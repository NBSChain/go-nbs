package application

import (
	"context"
	"github.com/NBSChain/go-nbs/storage/bitswap"
	"github.com/NBSChain/go-nbs/storage/merkledag"
	"github.com/NBSChain/go-nbs/storage/network"
	"github.com/NBSChain/go-nbs/storage/routing"
	"github.com/NBSChain/go-nbs/thirdParty/account"
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

	if app.nodeId != "" {
		//TIPS:: if no account exist, we can't identify the network. so there are no connections right now.
		network.GetInstance().StartUp(app.nodeId)
	} else {
		logger.Warning("no account now so the network is down now")
	}

	return nil
}

func (app *NbsApplication) ReloadForNewAccount() error {

	app.nodeId = account.GetAccountInstance().GetPeerID()

	network.GetInstance().StartUp(app.nodeId)

	return nil
}
