package application

import (
	"context"
	"fmt"
	"github.com/NBSChain/go-nbs/storage/application/rpcService"
	"github.com/NBSChain/go-nbs/storage/core"
	"github.com/NBSChain/go-nbs/utils"
	"sync"
)

type NbsApplication struct {
	context context.Context
	node    core.StorageNode
}

var logger = utils.GetLogInstance()

var instance *NbsApplication
var once sync.Once

func GetInstance() Application {

	once.Do(func() {

		app, err := newApplication()
		if err != nil {
			panic(err)
		}
		fmt.Printf("--->Create application to run......\n")

		instance = app
	})

	return instance
}

func newApplication() (*NbsApplication, error) {

	return &NbsApplication{
		context: context.Background(),
		node:    core.NewNode(),
	}, nil
}

func (app *NbsApplication) Start() error {

	logger.Info("Application starting......")

	rpcService.StartCmdService()

	instance.node.Online()

	return nil
}
