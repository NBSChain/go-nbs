package bitswap

import (
	"github.com/NBSChain/go-nbs/storage/merkledag/ipld"
	"sync"
)

const KeysToBroadNoPerRound = 1 << 4

func (bs *bitSwap) broadcastRunLoop()  {

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


func (bs *bitSwap) doBroadCast()  {

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

func (bs *bitSwap) sendOnNoe(node ipld.DagNode, waiter *sync.WaitGroup){

}
