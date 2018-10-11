package broadCaster

import (
	"github.com/NBSChain/go-nbs/storage/bitswap/pb"
	"github.com/NBSChain/go-nbs/storage/merkledag/cid"
	"github.com/NBSChain/go-nbs/storage/merkledag/ipld"
	"github.com/golang/protobuf/proto"
)

func (broadcast *BroadCaster) loadSavedBroadcastKeys() (*bitswap_pb.BroadCastKey, error){

	savedKeys := &bitswap_pb.BroadCastKey{}

	bytesOfKeys, err := broadcast.broadcastKeyStore.Get(ExchangeParamPrefix)
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

func (broadcast *BroadCaster) saveBroadcastKeysToStore(savedKeys *bitswap_pb.BroadCastKey) error  {

	newBytesOfKeys, err := proto.Marshal(savedKeys)
	if err != nil{
		logger.Warning("failed(3) to get serialize broadcast keys.")
		return err
	}

	if err := broadcast.broadcastKeyStore.Put(ExchangeParamPrefix, newBytesOfKeys); err != nil{
		logger.Warning("failed(4) to save broadcast keys.")
		return err
	}

	return nil
}

func (broadcast *BroadCaster) SyncKeys(keys []string) {

	if len(keys) == 0{
		return
	}

	savedKeys, err := broadcast.loadSavedBroadcastKeys()
	if err != nil{
		logger.Warning("loadSavedBroadcastKeys failed")
		return
	}

	savedKeys.Key = append(savedKeys.Key, keys...)

	if err = broadcast.saveBroadcastKeysToStore(savedKeys); err != nil{

		logger.Warning("saveBroadcastKeysToStore failed")
		return
	}
}

func (broadcast *BroadCaster) reloadBroadcastKeysToCache() error{

	savedKeys, err := broadcast.loadSavedBroadcastKeys()

	if err != nil{
		return err
	}

	nodes := make([]ipld.DagNode, len(savedKeys.Key))

	for _, key := range savedKeys.Key{

		cid, err := cid.DsKeyToCid(key)
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

		nodes = append(nodes, node)
	}

	broadcast.cacheLock.Lock()
	broadcast.broadcastCache = append(broadcast.broadcastCache, nodes...)
	broadcast.cacheLock.Unlock()

	return nil
}
