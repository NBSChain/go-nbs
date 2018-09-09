package ipld

import (
	"github.com/NBSChain/go-nbs/storage/core/blocks"
	"github.com/NBSChain/go-nbs/storage/merkledag/cid"
)

/*****************************************************************
*
*		Interfaces
*
*****************************************************************/
type Resolver interface {
	Resolve(path []string) (interface{}, []string, error)

	Tree(path string, depth int) []string
}

type DagNode interface {
	blocks.Block

	Resolver

	ResolveLink(path []string) (*DagLink, []string, error)

	Copy() DagNode

	Links() []*DagLink

	Size() (int64, error)
}

/*****************************************************************
*
*		implements
*
*****************************************************************/

type DagLink struct {
	Name string // utf8
	Size uint64
	Cid  *cid.Cid
}

type ProtoDagNode struct {
	links   []*DagLink
	data    []byte
	encoded []byte
	cached  *cid.Cid
}

func (n *ProtoDagNode) SetData(d []byte) {
	n.encoded = nil
	n.cached = nil
	n.data = d
}

/*****************************************************************
*
*		blocks.Block Interface
*
*****************************************************************/
func (n *ProtoDagNode) RawData() []byte {
	return nil
}

func (n *ProtoDagNode) Cid() *cid.Cid {
	return nil
}

func (n *ProtoDagNode) Loggable() map[string]interface{} {
	return map[string]interface{}{
		"node": n.String(),
	}
}

func (n *ProtoDagNode) String() string {
	return n.Cid().String()
}

/*****************************************************************
*
*		Resolver Interface
*
*****************************************************************/
func (n *ProtoDagNode) Resolve(path []string) (interface{}, []string, error) {
	return nil, nil, nil
}

func (n *ProtoDagNode) Tree(path string, depth int) []string {
	return nil
}

/*****************************************************************
*
*		DagNode Interface
*
*****************************************************************/
func (n *ProtoDagNode) ResolveLink(path []string) (*DagLink, []string, error) {
	return &DagLink{}, nil, nil
}

func (n *ProtoDagNode) Copy() DagNode {
	return nil
}

func (n *ProtoDagNode) Links() []*DagLink {
	return nil
}

func (n *ProtoDagNode) Size() (int64, error) {
	return 0, nil
}
