package bitswap

import (
	"context"
	"github.com/NBSChain/go-nbs/storage/merkledag/cid"
	"github.com/NBSChain/go-nbs/storage/merkledag/ipld"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/libp2p/go-libp2p-peer"
	"sync"
)

type Exchange interface {

	GetDagNode(*cid.Cid) (ipld.DagNode, error)

	GetDagNodes([]*cid.Cid) (<-chan ipld.DagNode, error)

	HasNode(ipld.DagNode) error
}

type SwapLedger interface {

	Score() float32

	Threshold() float32
}


type LedgerEngine interface {

	ReceiveData(fromNode peer.ID, data []byte) SwapLedger

	SupportData(toNode peer.ID, data []byte) SwapLedger

	GetLedger(nodeId peer.ID) SwapLedger
}


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

		logger.Info("bitSwap start to run......\n")
		instance = bs
	})

	return instance
}

func newBitSwap() (*bitSwap,error){
	return &bitSwap{}, nil
}

type bitSwap struct {
}

func (bs *bitSwap) GetDagNode(*cid.Cid) (ipld.DagNode, error){
	return nil, nil
}

func (bs *bitSwap) GetDagNodes([]*cid.Cid) (<-chan ipld.DagNode, error){
	return nil, nil
}

func (bs *bitSwap) HasNode(node ipld.DagNode) error{
	return nil
}

func (bs *bitSwap) Close() error{
	return nil
}