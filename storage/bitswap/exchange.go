package bitswap

import (
	"github.com/NBSChain/go-nbs/storage/merkledag/cid"
	"github.com/NBSChain/go-nbs/storage/merkledag/ipld"
	"gx/ipfs/QmWAzSEoqZ6xU6pu8yL8e5WaMb7wtbfbhhN4p1DknUPtr3/go-block-format"
	"io"
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
