package fetcher

import (
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
	wantList      []string
}

func NewRouterFetcher() *Fetcher{

	return &Fetcher{
		wantList:	make([]string, 0, MaxWaitingKeySize),
	}
}


/*****************************************************************
*
*		interface GetDagNode implements.
*
*****************************************************************/

func (getter *Fetcher)  GetNodeSync(cidObj *cid.Cid) (ipld.DagNode, error)  {

	key := cid.CovertCidToDataStoreKey(cidObj)

	peers, err := routing.GetInstance().FindPeer(key)

	if err != nil{
		return nil, err
	}

	if len(peers) == 0{
		return nil, ErrNotFound
	}

	data, err := getter.findValueFromPeers(key, peers, MaxDepthForRouting)
	if err != nil{
		return nil, err
	}

	return ipld.Decode(data, cidObj)
}

func (getter *Fetcher) findValueFromPeers(key string,
	peers []peerstore.PeerInfo, depth int) ([]byte, error){

	if peerSize := len(peers); peerSize > MaxPeerEachSearch{
		peerSize = MaxPeerEachSearch
		peers = peers[:peerSize]
	}

	data, newPeers, err := routing.GetInstance().GetValue(peers, key)
	if err != nil {
		return nil, err
	}

	if data != nil{
		return data, nil
	}

	if depth -= 1; depth < 0 || len(newPeers) == 0{
		return nil, ErrNotFound
	}

	return getter.findValueFromPeers(key, newPeers, depth)
}


/*****************************************************************
*
*		interface GetDagNodes implements.
*
*****************************************************************/