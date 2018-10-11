package bitswap

import (
	"context"
	"github.com/NBSChain/go-nbs/storage/bitswap/broadCaster"
	"github.com/NBSChain/go-nbs/storage/merkledag/cid"
	"github.com/NBSChain/go-nbs/storage/merkledag/ipld"
	"github.com/NBSChain/go-nbs/utils"
	"sync"
)

var instance 		*bitSwap
var once 		sync.Once
var parentContext 	context.Context
var logger 		= utils.GetLogInstance()

func GetSwapInstance() Exchange {

	once.Do(func() {
		parentContext = context.Background()

		bs, err := newBitSwap()
		if err != nil {
			panic(err)
		}

		go bs.broadCaster.BroadcastRunLoop()

		instance = bs
	})

	return instance
}

/*****************************************************************
*
*		DAGService interface implements.
*
*****************************************************************/
func newBitSwap() (*bitSwap,error){

	return &bitSwap{
		broadCaster:       broadCaster.NewBroadCaster(),
	}, nil
}

type bitSwap struct {
	broadCaster *broadCaster.BroadCaster
}

func (bs *bitSwap) GetDagNode(*cid.Cid) (ipld.DagNode, error){
	return nil, nil
}

func (bs *bitSwap) GetDagNodes([]*cid.Cid) (<-chan ipld.DagNode, error){
	return nil, nil
}

func (bs *bitSwap) SaveToNetPeer(nodes []ipld.DagNode) error{

	bs.broadCaster.Cache(nodes)

	return nil
}
