package broadCaster

import (
	"github.com/NBSChain/go-nbs/storage/application/dataStore"
	"github.com/NBSChain/go-nbs/storage/merkledag/ipld"
	"github.com/NBSChain/go-nbs/storage/routing"
	"github.com/NBSChain/go-nbs/utils"
	"sync"
)

var logger 			= utils.GetLogInstance()
const MaxBroadCastCache		= 1 << 20
const KeysToBroadNoPerRound 	= 1 << 4
const ExchangeParamPrefix	= "keys_to_be_broadcast"

type BroadCaster struct {

	sync.Mutex//broadcast cache locker.
	broadcastCache    	[]ipld.DagNode
	broadcastKeyStore 	dataStore.DataStore
	blockDataStore    	dataStore.DataStore
	workerSignal      	chan struct{}
}

func NewBroadCaster()  *BroadCaster {

	keyStore  := dataStore.GetServiceDispatcher().ServiceByType(dataStore.ServiceTypeLocalParam)
	blockStore:= dataStore.GetServiceDispatcher().ServiceByType(dataStore.ServiceTypeBlock)

	return &BroadCaster{
		broadcastCache:    	make([]ipld.DagNode, 0, MaxBroadCastCache),
		broadcastKeyStore: 	keyStore,
		blockDataStore:		blockStore,
		workerSignal:      	make(chan struct{}),
	}
}

func (bs *BroadCaster) Cache(nodes []ipld.DagNode)  {
	bs.Lock()
	defer bs.Unlock()

	//TODO:: Max size of broadcast cache
	bs.broadcastCache = append(bs.broadcastCache, nodes...)
}


func (bs *BroadCaster) broadcastRunLoop()  {

	if err := bs.reloadBroadcastKeysToCache(); err != nil{
		logger.Error(err)
		return
	}

	for{
		select {
			case <-bs.workerSignal:
				bs.doBroadCast()
		}
	}
}

func (bs *BroadCaster) doBroadCast()  {

	size := KeysToBroadNoPerRound

	if ok := len(bs.broadcastCache) < KeysToBroadNoPerRound; ok{
		size = len(bs.broadcastCache)
	}

	nodes := bs.broadcastCache[:size]

	var waitSignal sync.WaitGroup

	for _, node := range nodes{
		waitSignal.Add(1)
		bs.sendOnNoe(node, &waitSignal)
	}

	waitSignal.Wait()

	//bs.broadcastCache = bs.broadcastCache[size:]
	//
	//bs.remove
}

func (bs *BroadCaster) sendOnNoe(node ipld.DagNode, waiter *sync.WaitGroup){

	router := routing.GetInstance()

	router.PutValue(node.Cid().KeyString(), node.RawData())
}
