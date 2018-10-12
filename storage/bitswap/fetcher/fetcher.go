package fetcher

import (
	"context"
	"fmt"
	"github.com/NBSChain/go-nbs/storage/merkledag/cid"
	"github.com/NBSChain/go-nbs/storage/merkledag/ipld"
	"github.com/NBSChain/go-nbs/storage/routing"
	"github.com/libp2p/go-libp2p-peerstore"
)

const MaxPeerEachSearch  	= 3
const MaxWaitingKeySize		= 1 << 20
const MaxDepthForRouting  	= 20
var ErrNotFound 		= fmt.Errorf("can't find the target block data")

type Fetcher struct {
	workCancel	context.CancelFunc
	findCtx		context.Context
	wantList	[]string
}

func NewRouterFetcher() *Fetcher{

	ctx, cal := context.WithCancel(context.Background())

	return &Fetcher{
		workCancel:	cal,
		findCtx:	ctx,
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
	peerSize := len(peers)
	if peerSize == 0{
		return nil, ErrNotFound
	}

	if peerSize > MaxPeerEachSearch{
		peerSize = MaxPeerEachSearch
	}

	targetPeers := peers[:peerSize]

	data, peers, err := getter.findValueFromPeers(key, targetPeers, MaxDepthForRouting)
	if data != nil{
		return ipld.Decode(data, cidObj)
	}

	return nil, ErrNotFound
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