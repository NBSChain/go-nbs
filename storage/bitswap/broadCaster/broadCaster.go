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
const MaxBroadCastCache		= 1 << 20
const KeysToBroadNoPerRound 	= 1 << 6
const ExchangeParamPrefix	= "keys_to_be_broadcast"
const MaxTimeToPutBlocks 	= 3

type BroadCaster struct {

	cacheLock		sync.Mutex
	broadcastCache    	[]ipld.DagNode
	keystoreLock		sync.Mutex
	broadcastKeyStore 	dataStore.DataStore
	blockDataStore    	dataStore.DataStore
	callbackQueue		map[string]ipld.DagNode
	callbackResult		map[string]error
	workerSignal      	chan struct{}
}

func NewBroadCaster()  *BroadCaster {

	keyStore  := dataStore.GetServiceDispatcher().ServiceByType(dataStore.ServiceTypeLocalParam)
	blockStore:= dataStore.GetServiceDispatcher().ServiceByType(dataStore.ServiceTypeBlock)

	return &BroadCaster{
		broadcastCache:    	make([]ipld.DagNode, 0, MaxBroadCastCache),
		broadcastKeyStore: 	keyStore,
		blockDataStore:		blockStore,
		callbackQueue:		make(map[string]ipld.DagNode),
		callbackResult:		make(map[string]error),
		workerSignal:      	make(chan struct{}),
	}
}

/********************************************************************
*
*		TODO:: Max size of broadcast cache
*
********************************************************************/
func (broadcast *BroadCaster) Cache(nodes []ipld.DagNode)  {
	broadcast.cacheLock.Lock()
	broadcast.broadcastCache = append(broadcast.broadcastCache, nodes...)
	broadcast.cacheLock.Unlock()
}

func (broadcast *BroadCaster) broadcastRunLoop()  {

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


func (broadcast *BroadCaster) doBroadCast()  {

	size := KeysToBroadNoPerRound

	if ok := len(broadcast.broadcastCache) < KeysToBroadNoPerRound; ok{
		size = len(broadcast.broadcastCache)
	}

	broadcast.cacheLock.Lock()
	nodes := broadcast.broadcastCache[:size]
	broadcast.broadcastCache = broadcast.broadcastCache[size:]
	broadcast.cacheLock.Unlock()

	var waitSignal sync.WaitGroup

	broadcast.callbackQueue = make(map[string]ipld.DagNode)
	broadcast.callbackResult = make(map[string]error)

	for _, node := range nodes{
		waitSignal.Add(1)
		go broadcast.sendOnNoe(node, &waitSignal)
	}

	waitSignal.Wait()

	remainder := make([]ipld.DagNode, size)
	for key, err := range broadcast.callbackResult{
		if err != nil{
			remindNode := broadcast.callbackQueue[key]
			remainder = append(remainder, remindNode)
		}
	}

	if len(remainder) == 0{
		return
	}

	broadcast.cacheLock.Lock()
	broadcast.broadcastCache = append(broadcast.broadcastCache, remainder...)
	broadcast.cacheLock.Unlock()

	//savedKeys := &bitswap_pb.BroadCastKey{
	//	Key:	broadcast.broadcastCache,
	//}
}

func (broadcast *BroadCaster) sendOnNoe(node ipld.DagNode, waiter *sync.WaitGroup){

	defer  waiter.Done()

	key := node.Cid().KeyString()
	broadcast.callbackQueue[key] = node

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
