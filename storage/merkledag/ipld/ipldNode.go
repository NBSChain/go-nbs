package ipld

import (
	"github.com/NBSChain/go-nbs/storage/core"
	"github.com/NBSChain/go-nbs/storage/merkledag/cid"
)

type NodeOption struct {
	Node DagNode
	Err  error
}

type Resolver interface {
	Resolve(path []string) (interface{}, []string, error)

	Tree(path string, depth int) []string
}

type DagNode interface {
	core.Block

	Resolver

	ResolveLink(path []string) (*DagLink, []string, error)

	Copy() DagNode

	Links() []*DagLink

	Size() (uint64, error)
}

type DagLink struct {
	Name string // utf8

	Size uint64

	Cid *cid.Cid
}
