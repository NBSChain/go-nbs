package merkledag

import (
	"context"
	"github.com/NBSChain/go-nbs/storage/merkledag/cid"
	"github.com/NBSChain/go-nbs/storage/merkledag/ipld"
	"github.com/NBSChain/go-nbs/utils"
	"sync"
)

var instance *NbsDAGService
var once sync.Once
var parentContext context.Context
var logger = utils.GetLogInstance()

func GetDagInstance() DAGService {
	once.Do(func() {
		parentContext = context.Background()
		router, err := newNbsDagService()
		if err != nil {
			panic(err)
		}
		logger.Info("router start to run......\n")
		instance = router
	})

	return instance
}

type NbsDAGService struct {
	rehash     bool
	checkFirst bool
}

func newNbsDagService() (*NbsDAGService, error) {

	return &NbsDAGService{
		checkFirst: true,
		rehash:     false, //TODO:: I don't know the default value right now.
	}, nil
}

func (service *NbsDAGService) Get(*cid.Cid) (ipld.DagNode, error) {
	return nil, nil
}
func (service *NbsDAGService) GetMany([]*cid.Cid) <-chan *ipld.DagNode {
	return nil
}
func (service *NbsDAGService) Add(node ipld.DagNode) error {
	return nil
}
func (service *NbsDAGService) Remove(*cid.Cid) error {
	return nil
}
func (service *NbsDAGService) AddMany([]ipld.DagNode) error {
	return nil
}
func (service *NbsDAGService) RemoveMany([]*cid.Cid) error {
	return nil
}
