package bitswap

import (
	"context"
	"github.com/NBSChain/go-nbs/storage/application/dataStore"
	"github.com/NBSChain/go-nbs/storage/merkledag/cid"
	"github.com/NBSChain/go-nbs/storage/merkledag/ipld"
	"github.com/NBSChain/go-nbs/utils"
	"sync"
)


const MaxBroadCastCache	= 1 << 20
var instance 		*bitSwap
var once 		sync.Once
var parentContext 	context.Context
var logger 		= utils.GetLogInstance()

func GetSwapInstance() Exchange {

	once.Do(func() {
		parentContext = context.Background()

		bs, err := newBitSwap()
		if err != nil {
			panic(err)
		}

		logger.Info("bitSwap start to run......\n")
		instance = bs
	})

	return instance
}

/*****************************************************************
*
*		DAGService interface implements.
*
*****************************************************************/
const ExchangeParamPrefix = "keys_to_be_broadcast"

func newBitSwap() (*bitSwap,error){

	ds := dataStore.GetServiceDispatcher().ServiceByType(dataStore.ServiceTypeLocalParam)

	return &bitSwap{
		broadcastCache:         make([]ipld.DagNode, 0, MaxBroadCastCache),
		broadcastDataStore: 	ds,
		workerSignal:		make(chan struct{}),
	}, nil
}

type bitSwap struct {

	sync.Mutex//broadcast cache locker.
	broadcastCache 		[]ipld.DagNode
	broadcastDataStore	dataStore.DataStore
	workerSignal		chan struct{}
}

func (bs *bitSwap) GetDagNode(*cid.Cid) (ipld.DagNode, error){
	return nil, nil
}

func (bs *bitSwap) GetDagNodes([]*cid.Cid) (<-chan ipld.DagNode, error){
	return nil, nil
}

func (bs *bitSwap) SaveToNetPeer(nodes []ipld.DagNode) error{

	bs.Lock()
	defer bs.Unlock()

	//TODO:: Max size of broadcast cache
	bs.broadcastCache = append(bs.broadcastCache, nodes...)

	keys := make([]string, len(nodes))

	for _, node := range nodes{
		keys = append(keys, node.Cid().String())
	}

	go bs.syncKeys(keys)

	return nil
}
