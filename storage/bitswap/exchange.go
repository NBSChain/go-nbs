package bitswap

import (
	"context"
	"github.com/NBSChain/go-nbs/storage/merkledag/cid"
	"github.com/NBSChain/go-nbs/storage/merkledag/ipld"
	"github.com/NBSChain/go-nbs/utils"
	"io"
	"sync"
)

type Fetcher interface {
	GetDagNode(*cid.Cid) (ipld.DagNode, error)
	GetDagNodes([]*cid.Cid) (<-chan ipld.DagNode, error)
}

type Exchange interface {

	Fetcher

	HasNode(ipld.DagNode) error

	IsOnline() bool

	io.Closer
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

func (bs *bitSwap) IsOnline() bool{
	return false
}
func (bs *bitSwap) Close() error{
	return nil
}