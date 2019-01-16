package bitswap

import (
	"context"
	"github.com/NBSChain/go-nbs/storage/bitswap/broadCaster"
	"github.com/NBSChain/go-nbs/storage/bitswap/engine"
	"github.com/NBSChain/go-nbs/storage/bitswap/fetcher"
	"github.com/NBSChain/go-nbs/storage/merkledag/cid"
	"github.com/NBSChain/go-nbs/storage/merkledag/ipld"
	"github.com/NBSChain/go-nbs/utils"
	"sync"
)

var instance *bitSwap
var once sync.Once
var parentContext context.Context
var logger = utils.GetLogInstance()

func GetSwapInstance() Exchange {

	once.Do(func() {
		parentContext = context.Background()

		bs, err := newBitSwap()
		if err != nil {
			panic(err)
		}

		go bs.broadCaster.BroadcastRunLoop()
		go bs.wantManager.FetchRunLoop()
		instance = bs
	})

	return instance
}

/*****************************************************************
*
*		DAGService interface implements.
*
*****************************************************************/
func newBitSwap() (*bitSwap, error) {

	return &bitSwap{
		broadCaster:  broadCaster.NewBroadCaster(),
		wantManager:  fetcher.NewRouterFetcher(),
		ledgerEngine: engine.NewLedgerEngine(),
	}, nil
}

type bitSwap struct {
	broadCaster  *broadCaster.BroadCaster
	wantManager  *fetcher.Fetcher
	ledgerEngine LedgerEngine
}

func (bs *bitSwap) GetDagNode(cidObj *cid.Cid) (ipld.DagNode, error) {
	return bs.wantManager.GetNodeSync(cidObj)
}

func (bs *bitSwap) GetDagNodes(ctx context.Context, cidArr []*cid.Cid) <-chan fetcher.AsyncResult {
	return bs.wantManager.CacheRequest(ctx, cidArr)
}

func (bs *bitSwap) SaveToNetPeer(nodes map[string]ipld.DagNode) error {

	bs.broadCaster.PushCache(nodes)

	go bs.broadCaster.SyncCurrentCache()

	return nil
}

func (bs *bitSwap) GetLedgerEngine() LedgerEngine {
	return bs.ledgerEngine
}
