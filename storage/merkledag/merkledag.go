package merkledag

import (
	"context"
	"github.com/NBSChain/go-nbs/storage/merkledag/cid"
)

type NodeOption struct {
	Node Node
	Err  error
}

type NodeGetter interface {
	Get(context.Context, *cid.Cid) (Node, error)

	GetMany(context.Context, []*cid.Cid) <-chan *NodeOption
}

type LinkGetter interface {
	NodeGetter

	// TODO(ipfs/go-ipld-format#9): This should return []*cid.Cid

	GetLinks(ctx context.Context, nd *cid.Cid) ([]*Link, error)
}

type DAGService interface {
	NodeGetter

	Add(context.Context, Node) error

	Remove(context.Context, *cid.Cid) error

	AddMany(context.Context, []Node) error

	RemoveMany(context.Context, []*cid.Cid) error
}

type Resolver interface {
	Resolve(path []string) (interface{}, []string, error)

	Tree(path string, depth int) []string
}

type Node interface {
	Block

	Resolver

	ResolveLink(path []string) (*Link, []string, error)

	Copy() Node

	Links() []*Link

	Size() (uint64, error)
}

type Link struct {
	Name string // utf8

	Size uint64

	Cid *cid.Cid
}
