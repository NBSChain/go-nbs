package application

import (
	"context"
	"fmt"
	"github.com/NBSChain/go-nbs/utils"
	"os"
	"sync"
)

type NbsApplication struct {
	Context context.Context
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
		Context: context.Background(),
	}, nil
}

func (*NbsApplication) AddFile(file *os.File) error {

	logger.Info("Application start to Add File", file)

	return nil
}
