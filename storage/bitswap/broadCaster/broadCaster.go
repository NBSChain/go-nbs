package broadCaster

import (
	"github.com/NBSChain/go-nbs/storage/application/dataStore"
	"github.com/NBSChain/go-nbs/storage/merkledag/ipld"
	"github.com/NBSChain/go-nbs/storage/routing"
	"github.com/NBSChain/go-nbs/utils"
	"sync"
	"time"
)

//const MaxBroadCastCache	= 1 << 20
var logger 			= utils.GetLogInstance()
const KeysToBroadNoPerRound 	= 1 << 6
const ExchangeParamPrefix	= "keys_to_be_broadcast"
const MaxTimeToPutBlocks 	= 3
const RestTimeForWorking 	= 100
const RestTimeForIdle		= 5

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

type broadcastResult struct {
	sync.Mutex
	resultsQueue	map[string]ipld.DagNode
}

func NewBroadCaster() *BroadCaster {

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
*
*
********************************************************************/
func (broadcast *BroadCaster) BroadcastRunLoop()  {

	logger.Info("exchange layer start to work......")

	if err := broadcast.reloadBroadcastKeysToCache(); err != nil{
		logger.Error(err)
		return
	}

	for {
		time.Sleep(time.Millisecond * RestTimeForWorking)

		nodesWorkLoad := broadcast.popCache(KeysToBroadNoPerRound)
		if len(nodesWorkLoad) == 0{
			logger.Debug("no data to broadcast")
			time.Sleep(time.Second * RestTimeForIdle)
			continue
		}

		logger.Info("start to broad cast blocks to target peers ......")

		broadcast.startBroadCast(nodesWorkLoad)

		go broadcast.SyncCurrentCache()
	}
}

func (broadcast *BroadCaster) SyncCurrentCache(){

	broadcast.Lock()

	remainders := make([]string, 0, len(broadcast.broadcastCache))

	for key := range broadcast.broadcastCache{
		remainders = append(remainders, key)
	}

	broadcast.Unlock()

	broadcast.saveBroadcastKeysToStore(remainders)
}

func (broadcast *BroadCaster) startBroadCast(nodesWorkLoad map[string]ipld.DagNode)  {

	callbackQueue := &broadcastResult{
		resultsQueue:	make(map[string]ipld.DagNode),
	}

	var waitSignal sync.WaitGroup

	for key, node := range nodesWorkLoad{
		waitSignal.Add(1)
		go broadcast.sendOnNoe(key, node, &waitSignal, callbackQueue)
	}

	waitSignal.Wait()

	if len(callbackQueue.resultsQueue) == 0{
		return
	}

	broadcast.Lock()
	defer broadcast.Unlock()

	for key, node := range callbackQueue.resultsQueue{
		broadcast.broadcastCache[key] = node
	}
}

func (broadcast *BroadCaster) sendOnNoe(key string, node ipld.DagNode,
	waiter *sync.WaitGroup, callbackQueue *broadcastResult){

	defer  waiter.Done()

	errorChan := routing.GetInstance().PutValue(key, node.RawData())

	select {
		case err := <-errorChan:
			if err != nil{
				callbackQueue.pushItem(key, node)
				logger.Warning("saved data to network with error ", key, err)
			}

		case <-time.After(time.Second * MaxTimeToPutBlocks):
			callbackQueue.pushItem(key, node)
			logger.Warning("timeout, failed to put block onto network:key", key)
	}
}

/********************************************************************
*
*
*
********************************************************************/
func (queue *broadcastResult) pushItem(key string, node ipld.DagNode)  {

	queue.Lock()

	queue.resultsQueue[key] = node

	queue.Unlock()
}
