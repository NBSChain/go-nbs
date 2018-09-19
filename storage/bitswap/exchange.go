package bitswap

import (
	"context"
	"github.com/NBSChain/go-nbs/storage/merkledag/cid"
	"github.com/NBSChain/go-nbs/storage/merkledag/ipld"
	"github.com/NBSChain/go-nbs/utils"
	"gx/ipfs/QmWAzSEoqZ6xU6pu8yL8e5WaMb7wtbfbhhN4p1DknUPtr3/go-block-format"
	"io"
	"sync"
)

type Fetcher interface {
	GetBlock(*cid.Cid) (ipld.DagNode, error)
	GetBlocks([]*cid.Cid) (<-chan ipld.DagNode, error)
}

type Exchange interface {
	Fetcher

	HasBlock(blocks.Block) error

	IsOnline() bool

	io.Closer
}

var instance *bitSwap
var once sync.Once
var parentContext context.Context
var logger = utils.GetLogInstance()

func GetExInstance() Exchange {
	once.Do(func() {
		parentContext = context.Background()
		bs, err := newBitSwap()
		if err != nil {
			panic(err)
		}
		logger.Info("router start to run......\n")
		instance = bs
	})

	return instance
}

func newBitSwap() (*bitSwap,error){
	return &bitSwap{}, nil
}

type bitSwap struct {
}

func (bs *bitSwap) GetBlock(*cid.Cid) (ipld.DagNode, error){
	return nil, nil
}
func (bs *bitSwap) GetBlocks([]*cid.Cid) (<-chan ipld.DagNode, error){
	return nil, nil
}

func (bs *bitSwap) HasBlock(blocks.Block) error{
	return nil
}

func (bs *bitSwap) IsOnline() bool{
	return false
}
func (bs *bitSwap) Close() error{
	return nil
}