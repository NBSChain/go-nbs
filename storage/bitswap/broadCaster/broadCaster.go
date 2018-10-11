package broadCaster

import (
	"fmt"
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

type BroadCaster struct {

	/*used to cache the bock data in memory and broad cast them later.*/
	sync.Mutex
	broadcastCache 		map[string]ipld.DagNode

	/*used to save task keys when nbs node is down.*/
	keystoreLock		sync.Mutex
	keyStoreService 	dataStore.DataStore

	/*used to get block data from local store.*/
	blockDataStore    	dataStore.DataStore

	/*used to run loop work to call back process.*/
	callbackQueue		map[string]ipld.DagNode
	callbackResult		map[string]error
	workerSignal      	chan struct{}
}

func NewBroadCaster()  *BroadCaster {

	keyStore  := dataStore.GetServiceDispatcher().ServiceByType(dataStore.ServiceTypeLocalParam)
	blockStore:= dataStore.GetServiceDispatcher().ServiceByType(dataStore.ServiceTypeBlock)

	return &BroadCaster{
		broadcastCache:  	make(map[string]ipld.DagNode),
		keyStoreService: 	keyStore,
		blockDataStore:  	blockStore,
		callbackQueue:   	make(map[string]ipld.DagNode),
		callbackResult:  	make(map[string]error),
		workerSignal:    	make(chan struct{}),
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

	broadcast.workerSignal<- struct{}{}
}

func (broadcast *BroadCaster) BroadcastRunLoop()  {

	logger.Info("exchange layer start to ")
	if err := broadcast.reloadBroadcastKeysToCache(); err != nil{
		logger.Error(err)
		return
	}

	for {
		select {
			case <-broadcast.workerSignal:
				broadcast.doBroadCast()
		}
	}
}


func (broadcast *BroadCaster) doBroadCast() {

	size := KeysToBroadNoPerRound
	if ok := len(broadcast.broadcastCache) < KeysToBroadNoPerRound; ok{
		size = len(broadcast.broadcastCache)
	}

	if size == 0{
		logger.Info("broad cast list is empty right now")
		return
	}


	broadcast.callbackResult = make(map[string]error)

	var waitSignal sync.WaitGroup

	broadcast.callbackQueue = broadcast.popCache(size)

	for _, node := range broadcast.callbackQueue{
		waitSignal.Add(1)
		go broadcast.sendOnNoe(node, &waitSignal)
	}

	waitSignal.Wait()

	remainders := make([]ipld.DagNode, size)
	for key, err := range broadcast.callbackResult{

		if err == nil{
			continue
		}

		node := broadcast.callbackQueue[key]
		remainders = append(remainders, node)
	}

	broadcast.Cache(remainders)
}

func (broadcast *BroadCaster) sendOnNoe(node ipld.DagNode, waiter *sync.WaitGroup){

	defer  waiter.Done()

	key := node.Cid().KeyString()

	errorChan := routing.GetInstance().PutValue(key, node.RawData())

	select {
		case err := <-errorChan:
			broadcast.callbackResult[key] = err
			logger.Info("saved data to net work finished", key, err)

		case <-time.After(time.Second * MaxTimeToPutBlocks):
			err := fmt.Errorf("failed to put block onto network:key=%s", key)
			broadcast.callbackResult[key] = err
	}
}
