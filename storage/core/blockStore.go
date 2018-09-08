package core

import (
	"context"
	"github.com/NBSChain/go-nbs/storage/merkledag/cid"
)

type BlockStore interface {
	DeleteBlock(*cid.Cid) error

	Has(*cid.Cid) (bool, error)

	Get(*cid.Cid) (Block, error)

	GetSize(*cid.Cid) (int, error)

	Put(Block) error

	PutMany([]Block) error

	AllKeysChan(ctx context.Context) (<-chan *cid.Cid, error)

	HashOnRead(enabled bool)
}
