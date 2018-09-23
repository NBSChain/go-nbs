package merkledag

import (
	"github.com/NBSChain/go-nbs/storage/merkledag/cid"
	"github.com/NBSChain/go-nbs/storage/merkledag/ipld"
)

type NodeGetter interface {
	Get(*cid.Cid) (ipld.DagNode, error)

	GetMany([]*cid.Cid) <-chan *ipld.DagNode
}

type LinkGetter interface {
	NodeGetter

	GetLinks(nd *cid.Cid) ([]*ipld.DagLink, error)
}

type DAGService interface {
	NodeGetter

	Add(ipld.DagNode) error

	Remove(*cid.Cid) error

	AddMany([]ipld.DagNode) error

	RemoveMany([]*cid.Cid) error
}
