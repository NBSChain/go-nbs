package merkledag

import (
	"context"
	"github.com/NBSChain/go-nbs/storage/merkledag/cid"
	"github.com/NBSChain/go-nbs/storage/merkledag/ipld"
)

type NodeGetter interface {
	Get(context.Context, *cid.Cid) (ipld.DagNode, error)

	GetMany(context.Context, []*cid.Cid) <-chan *ipld.NodeOption
}

type LinkGetter interface {
	NodeGetter

	GetLinks(ctx context.Context, nd *cid.Cid) ([]*ipld.DagLink, error)
}

type DAGService interface {
	NodeGetter

	Add(context.Context, ipld.DagNode) error

	Remove(context.Context, *cid.Cid) error

	AddMany(context.Context, []ipld.DagNode) error

	RemoveMany(context.Context, []*cid.Cid) error
}
