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

func GetInstance() DAGService {
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
}

func newNbsDagService() (*NbsDAGService, error) {

	return &NbsDAGService{}, nil
}

func (service *NbsDAGService) Get(context.Context, *cid.Cid) (ipld.DagNode, error) {
	return nil, nil
}
func (service *NbsDAGService) GetMany(context.Context, []*cid.Cid) <-chan *ipld.DagNode {
	return nil
}
func (service *NbsDAGService) Add(context.Context, ipld.DagNode) error {
	return nil
}
func (service *NbsDAGService) Remove(context.Context, *cid.Cid) error {
	return nil
}
func (service *NbsDAGService) AddMany(context.Context, []ipld.DagNode) error {
	return nil
}
func (service *NbsDAGService) RemoveMany(context.Context, []*cid.Cid) error {
	return nil
}
