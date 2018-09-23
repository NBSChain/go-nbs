package application

import (
	"context"
	"github.com/NBSChain/go-nbs/storage/application/rpcService"
	"sync"
)

type NbsApplication struct {
	context context.Context
	node    StorageNode
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

	return &NbsApplication{
		context: context.Background(),
		node:    NewNode(),
	}, nil
}

func (app *NbsApplication) Start() error {

	logger.Info("Application starting......")

	rpcService.StartCmdService()

	instance.node.Online()

	return nil
}
