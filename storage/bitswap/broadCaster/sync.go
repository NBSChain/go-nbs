package broadCaster

import (
	"github.com/NBSChain/go-nbs/storage/bitswap/pb"
	"github.com/NBSChain/go-nbs/storage/merkledag/cid"
	"github.com/NBSChain/go-nbs/storage/merkledag/ipld"
	"github.com/golang/protobuf/proto"
)

func (bs *BroadCaster) loadSavedBroadcastKeys() (*bitswap_pb.BroadCastKey, error){

	savedKeys := &bitswap_pb.BroadCastKey{}

	bytesOfKeys, err := bs.broadcastKeyStore.Get(ExchangeParamPrefix)
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

func (bs *BroadCaster) saveBroadcastKeysToStore(savedKeys *bitswap_pb.BroadCastKey) error  {

	newBytesOfKeys, err := proto.Marshal(savedKeys)
	if err != nil{
		logger.Warning("failed(3) to get serialize broadcast keys.")
		return err
	}

	if err := bs.broadcastKeyStore.Put(ExchangeParamPrefix, newBytesOfKeys); err != nil{
		logger.Warning("failed(4) to save broadcast keys.")
		return err
	}

	return nil
}

func (bs *BroadCaster) SyncKeys(keys []string) {

	if len(keys) == 0{
		return
	}

	savedKeys, err := bs.loadSavedBroadcastKeys()
	if err != nil{
		logger.Warning("loadSavedBroadcastKeys failed")
		return
	}

	savedKeys.Key = append(savedKeys.Key, keys...)

	if err = bs.saveBroadcastKeysToStore(savedKeys); err != nil{

		logger.Warning("saveBroadcastKeysToStore failed")
		return
	}
}

func (bs *BroadCaster) reloadBroadcastKeysToCache() error{

	savedKeys, err := bs.loadSavedBroadcastKeys()

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

		data, err := bs.blockDataStore.Get(key)
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

	bs.Lock()

	defer bs.Unlock()

	bs.broadcastCache = append(bs.broadcastCache, nodes...)

	return nil
}
