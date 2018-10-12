package fetcher

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/merkledag/cid"
	"github.com/NBSChain/go-nbs/storage/merkledag/ipld"
	"github.com/NBSChain/go-nbs/storage/routing"
	"github.com/libp2p/go-libp2p-peerstore"
)

const MaxWaitingKeySize		= 1 << 20
const MaxDepthForRouting  	= 20

type Fetcher struct {
	wantList	[]string
}

func NewRouterFetcher() *Fetcher{
	return &Fetcher{
		wantList:	make([]string, 0, MaxWaitingKeySize),
	}
}

func (getter *Fetcher)  GetNodeSync(cidObj *cid.Cid) (ipld.DagNode, error)  {

	router := routing.GetInstance()

	key := cid.CovertCidToDataStoreKey(cidObj)

	peers, err := router.FindPeer(key)
	if err != nil{
		return nil, err
	}

	data, peers, err := getter.findValueFromPeers(key, peers, MaxDepthForRouting)
	if data != nil{
		return ipld.Decode(data, cidObj)
	}

	return nil, fmt.Errorf("can't find the target block data")
}

func (getter *Fetcher) findValueFromPeers(key string,
	peers []peerstore.PeerInfo, depth int) ([]byte, []peerstore.PeerInfo, error){

	if depth -= 1; depth == 0{
		return nil, nil , fmt.Errorf("fail to find the data of key till end")
	}

	router := routing.GetInstance()

	for _, peerInfo := range peers{

		data, newPeers, err := router.GetValue(peerInfo, key)
		if err != nil{
			return nil, nil, err
		}

		if data != nil{
			return data, nil, nil
		}

		return getter.findValueFromPeers(key, newPeers, depth)
	}

	return nil, nil, nil
}