package broadCaster

import (
	"github.com/NBSChain/go-nbs/storage/bitswap/pb"
	"github.com/NBSChain/go-nbs/storage/merkledag/cid"
	"github.com/NBSChain/go-nbs/storage/merkledag/ipld"
	"github.com/golang/protobuf/proto"
)

func (broadcast *BroadCaster) loadSavedBroadcastKeys() (*bitswap_pb.BroadCastKey, error){
	broadcast.keystoreLock.Lock()
	defer broadcast.keystoreLock.Unlock()

	savedKeys := &bitswap_pb.BroadCastKey{}

	bytesOfKeys, err := broadcast.keyStoreService.Get(ExchangeParamPrefix)
	if err != nil{
		logger.Warning("failed(1) to get saved broadcast keys.")
		return nil, err
	}

	if err := proto.UnmarshalMerge(bytesOfKeys, savedKeys); err != nil{

		logger.Warning("failed(2) to get parse broadcast keys.")
		return nil, err
	}

	return savedKeys, nil
}

func (broadcast *BroadCaster) saveBroadcastKeysToStore(keys []string) error  {

	if len(keys) == 0{
		return nil
	}

	broadcast.keystoreLock.Lock()
	defer broadcast.keystoreLock.Unlock()

	keyStores := &bitswap_pb.BroadCastKey{Keys:keys}

	newBytesOfKeys, err := proto.Marshal(keyStores)
	if err != nil{
		logger.Warning("failed(3) to get serialize broadcast keys.")
		return err
	}

	if err := broadcast.keyStoreService.Put(ExchangeParamPrefix, newBytesOfKeys); err != nil{
		logger.Warning("failed(4) to save broadcast keys.")
		return err
	}

	return nil
}

func (broadcast *BroadCaster) reloadBroadcastKeysToCache() error{

	savedKeys, err := broadcast.loadSavedBroadcastKeys()
	if err != nil{
		return err
	}

	broadcast.Lock()
	defer broadcast.Unlock()

	for _, key := range savedKeys.Keys{

		cid, err := cid.CovertDataStoreKeyToCid(key)
		if err != nil{
			logger.Warningf("convert string to cid object failed:%s", err.Error())
			continue
		}

		data, err := broadcast.blockDataStore.Get(key)
		if err != nil{
			logger.Warningf("get data from local store failed:%s", err.Error())
			continue
		}

		node, err := ipld.Decode(data, cid)
		if err != nil{
			logger.Warningf("decode local data failed:%s", err.Error())
			continue
		}

		broadcast.broadcastCache[key] = node
	}

	return nil
}

func (broadcast *BroadCaster) PushCache(nodes map[string]ipld.DagNode){

	//TODO:: Max size of broadcast cache
	if len(nodes) == 0{
		return
	}

	broadcast.Lock()
	defer broadcast.Unlock()

	for key, node := range nodes{
		broadcast.broadcastCache[key] = node
	}
}

func (broadcast *BroadCaster) popCache(size int) map[string]ipld.DagNode{

	broadcast.Lock()
	defer broadcast.Unlock()

	nodes := make(map[string]ipld.DagNode)

	for key, node := range broadcast.broadcastCache{
		nodes[key] = node

		delete(broadcast.broadcastCache, key)

		if size -= 1; size <= 0{
			break
		}
	}

	return nodes
}