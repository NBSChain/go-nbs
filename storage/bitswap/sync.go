package bitswap

import (
	"github.com/NBSChain/go-nbs/storage/bitswap/pb"
	"github.com/golang/protobuf/proto"
)

func (bs *bitSwap) loadSavedBroadcastKeys() (*bitswap_pb.BroadCastKey, error){

	savedKeys := &bitswap_pb.BroadCastKey{}

	bytesOfKeys, err := bs.broadcastDataStore.Get(ExchangeParamPrefix)
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

func (bs *bitSwap) saveBroadcastKeysToStore(savedKeys *bitswap_pb.BroadCastKey) error  {

	newBytesOfKeys, err := proto.Marshal(savedKeys)
	if err != nil{
		logger.Warning("failed(3) to get serialize broadcast keys.")
		return err
	}

	if err := bs.broadcastDataStore.Put(ExchangeParamPrefix, newBytesOfKeys); err != nil{
		logger.Warning("failed(4) to save broadcast keys.")
		return err
	}

	return nil
}

func (bs *bitSwap) syncKeys(keys []string) {

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

func (bs *bitSwap) reloadBroadcastKeysToCache() error{

	//savedKeys, err := bs.loadSavedBroadcastKeys()
	//
	//if err != nil{
	//	return err
	//}
	//
	//dagService := merkledag.GetDagInstance()
	//
	//nodes := make([]ipld.DagNode, len(savedKeys.Key))
	//
	//for _, key := range savedKeys.Key{
	//
	//	cid, err := cid.DsKeyToCid(key)
	//	if err != nil{
	//		logger.Warning(err)
	//		continue
	//	}
	//
	//	node, err := dagService.Get(cid)
	//	if err != nil{
	//		logger.Warning(err)
	//		continue
	//	}
	//
	//	nodes = append(nodes, node)
	//}
	//
	//bs.Lock()
	//
	//defer bs.Unlock()
	//
	//bs.broadcastCache = append(bs.broadcastCache, nodes...)

	return nil
}
