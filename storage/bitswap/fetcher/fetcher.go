package fetcher

import (
	"context"
	"fmt"
	"github.com/NBSChain/go-nbs/storage/merkledag/cid"
	"github.com/NBSChain/go-nbs/storage/merkledag/ipld"
	"github.com/NBSChain/go-nbs/storage/routing"
	"github.com/libp2p/go-libp2p-peerstore"
	"sync"
	"time"
)

const MaxPeerEachSearch  	= 3
const MaxItemPerRound		= 10
const MaxIdleTime		= 100
const MaxDepthForRouting  	= 20
var ErrNotFound 		= fmt.Errorf("can't find the target block data")

type AsyncResult struct {
	node ipld.DagNode
	err  error
}

type wantItem struct {
	sessionID	uint
	resultNodeChan	chan AsyncResult
	waitingItem	[]*cid.Cid
}

type Fetcher struct {
	sync.Mutex
	sessionID	uint
	wantQueue	map[uint]*wantItem
}

func NewRouterFetcher() *Fetcher{

	return &Fetcher{
		wantQueue:	make(map[uint]*wantItem),
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
func (getter *Fetcher)  CacheRequest(ctx context.Context, cidArr []*cid.Cid) <-chan AsyncResult {

	resultNode := make(chan AsyncResult)

	getter.Lock()
	defer getter.Unlock()

	getter.sessionID += 1
	getter.wantQueue[getter.sessionID] = &wantItem{
		resultNodeChan:	resultNode,
		waitingItem:	cidArr,
		sessionID:	getter.sessionID,
	}

	return resultNode
}

func (getter *Fetcher) FetchRunLoop() {

	for{
		if len(getter.wantQueue) == 0{
			time.Sleep(time.Millisecond * MaxIdleTime)
			continue
		}

		noOneRound := MaxItemPerRound

		getter.Lock()
		for sessionID, item := range getter.wantQueue{
			if noOneRound -= 1; noOneRound < 0{
				break
			}

			go getter.getNodeArrayAsync(item)
			delete(getter.wantQueue, sessionID)

		}
		getter.Unlock()
	}
}

func (getter *Fetcher) getNodeArrayAsync(item *wantItem){

	for _, cidObj := range item.waitingItem{

		node, err := getter.GetNodeSync(cidObj)
		item.resultNodeChan<- AsyncResult{
			node: node,
			err:  err,
		}
	}
}