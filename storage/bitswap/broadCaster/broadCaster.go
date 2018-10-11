package broadCaster

import (
	"github.com/NBSChain/go-nbs/storage/application/dataStore"
	"github.com/NBSChain/go-nbs/storage/merkledag/ipld"
	"github.com/NBSChain/go-nbs/storage/routing"
	"github.com/NBSChain/go-nbs/utils"
	"sync"
	"time"
)

var logger 			= utils.GetLogInstance()
//const MaxBroadCastCache		= 1 << 20
const KeysToBroadNoPerRound 	= 1 << 6
const ExchangeParamPrefix	= "keys_to_be_broadcast"
const MaxTimeToPutBlocks 	= 3
const IdleTimeForRest		= 100

type BroadCaster struct {

	/*used to cache the bock data in memory and broad cast them later.*/
	sync.Mutex
	broadcastCache 		map[string]ipld.DagNode

	/*used to save task keys when nbs node is down.*/
	keystoreLock		sync.Mutex
	keyStoreService 	dataStore.DataStore

	/*used to get block data from local store.*/
	blockDataStore    	dataStore.DataStore
}

func NewBroadCaster()  *BroadCaster {

	keyStore  := dataStore.GetServiceDispatcher().ServiceByType(dataStore.ServiceTypeLocalParam)
	blockStore:= dataStore.GetServiceDispatcher().ServiceByType(dataStore.ServiceTypeBlock)

	return &BroadCaster{
		broadcastCache:  	make(map[string]ipld.DagNode),
		keyStoreService: 	keyStore,
		blockDataStore:  	blockStore,
	}
}

/********************************************************************
*
*		TODO:: Max size of broadcast cache
*
********************************************************************/
func (broadcast *BroadCaster) Cache(nodes []ipld.DagNode)  {

	keys := broadcast.pushCache(nodes)

	go broadcast.saveBroadcastKeysToStore(keys)
}

func (broadcast *BroadCaster) BroadcastRunLoop()  {

	logger.Info("exchange layer start to ")
	if err := broadcast.reloadBroadcastKeysToCache(); err != nil{
		logger.Error(err)
		return
	}

	for {
		time.Sleep(time.Millisecond * IdleTimeForRest)

		nodesWorkLoad := broadcast.popCache(KeysToBroadNoPerRound)
		if len(nodesWorkLoad) == 0{
			time.Sleep(time.Second)
			continue
		}

		logger.Info("start to broad cast blocks to target peers ......")

		broadcast.startBroadCast(nodesWorkLoad)

		go broadcast.syncCurrentCache()
	}
}

func (broadcast *BroadCaster) syncCurrentCache(){

	remainders := make([]string, len(broadcast.broadcastCache))

	for key := range broadcast.broadcastCache{
		remainders = append(remainders, key)
	}

	broadcast.saveBroadcastKeysToStore(remainders)
}

func (broadcast *BroadCaster) startBroadCast(nodesWorkLoad map[string]ipld.DagNode)  {
	callbackQueue := make(map[string]ipld.DagNode)

	var waitSignal sync.WaitGroup

	for _, node := range nodesWorkLoad{
		waitSignal.Add(1)
		go broadcast.sendOnNoe(node, &waitSignal, callbackQueue)
	}

	waitSignal.Wait()

	if len(callbackQueue) == 0{
		return
	}

	broadcast.Lock()
	defer broadcast.Unlock()

	for key, node := range callbackQueue{
		broadcast.broadcastCache[key] = node
	}
}



func (broadcast *BroadCaster) sendOnNoe(node ipld.DagNode,
	waiter *sync.WaitGroup, callbackQueue map[string]ipld.DagNode){

	defer  waiter.Done()

	key := node.Cid().KeyString()
	errorChan := routing.GetInstance().PutValue(key, node.RawData())

	select {
		case err := <-errorChan:
			if err != nil{
				callbackQueue[key] = node
				logger.Info("saved data to net work finished", key, err)
			}

		case <-time.After(time.Second * MaxTimeToPutBlocks):
			logger.Error("failed to put block onto network:key", key)
			callbackQueue[key] = node
	}
}
