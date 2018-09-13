package ipld

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/core/blocks"
	"github.com/NBSChain/go-nbs/storage/merkledag/cid"
	pb "github.com/NBSChain/go-nbs/storage/merkledag/pb"
	"github.com/multiformats/go-multihash"
	"sort"
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

var v0Cid = &cid.Cid{
	Version:  0,
	Code:     cid.DagProtobuf,
	HashType: multihash.SHA2_256,
	HashLen:  -1,
}

var v1Cid = &cid.Cid{
	Version:  1,
	Code:     cid.DagProtobuf,
	HashType: multihash.SHA2_256,
	HashLen:  -1,
}

type DagLink struct {
	Name string // utf8
	Size int64
	Cid  *cid.Cid
}

type LinkSlice []*DagLink

func (ls LinkSlice) Len() int           { return len(ls) }
func (ls LinkSlice) Swap(a, b int)      { ls[a], ls[b] = ls[b], ls[a] }
func (ls LinkSlice) Less(a, b int) bool { return ls[a].Name < ls[b].Name }

type ProtoDagNode struct {
	links   []*DagLink
	data    []byte
	encoded []byte
	cached  *cid.Cid
}

func MakeLink(n DagNode) (*DagLink, error) {

	s, err := n.Size()
	if err != nil {
		return nil, err
	}

	return &DagLink{
		Size: s,
		Cid:  n.Cid(),
	}, nil
}

/*****************************************************************
*
*		blocks.Block Interface
*
*****************************************************************/
func (node *ProtoDagNode) RawData() []byte {
	node.EncodeProtoBuf(false)
	return node.encoded
}

func (node *ProtoDagNode) Cid() *cid.Cid {

	if node.encoded != nil && node.cached.Hash != nil {
		return node.cached
	}

	err := node.EncodeProtoBuf(false)
	if err != nil {
		err = fmt.Errorf("invalid CID of length %d: %x: %v", len(node.RawData()), node.RawData(), err)
		panic(err)
	}

	return node.cached
}

func (node *ProtoDagNode) String() string {
	return node.Cid().String()
}

/*****************************************************************
*
*		Resolver Interface
*
*****************************************************************/
func (node *ProtoDagNode) Resolve(path []string) (interface{}, []string, error) {
	return nil, nil, nil
}

func (node *ProtoDagNode) Tree(path string, depth int) []string {
	return nil
}

/*****************************************************************
*
*		DagNode Interface
*
*****************************************************************/
func (node *ProtoDagNode) ResolveLink(path []string) (*DagLink, []string, error) {
	return &DagLink{}, nil, nil
}

func (node *ProtoDagNode) Copy() DagNode {
	return nil
}

func (node *ProtoDagNode) Links() []*DagLink {
	return nil
}

func (node *ProtoDagNode) Size() (int64, error) {

	err := node.EncodeProtoBuf(false)
	if err != nil {
		return 0, err
	}

	size := int64(len(node.encoded))

	for _, l := range node.links {
		size += l.Size
	}

	return size, nil
}

/*****************************************************************
*
*		Logic
*
*****************************************************************/
func (node *ProtoDagNode) SetData(d []byte) {
	node.encoded = nil
	node.cached.Hash = nil
	node.data = d
}

func (node *ProtoDagNode) AddNodeLink(name string, that DagNode) error {

	node.encoded = nil

	link, err := MakeLink(that)
	if err != nil {
		return err
	}
	link.Name = name

	node.AddRawLink(name, link)

	return nil
}

func (node *ProtoDagNode) EncodeProtoBuf(force bool) (err error) {

	if node.encoded != nil && !force {
		return nil
	}

	node.encoded, err = node.Marshal()
	if err != nil {
		return err
	}

	return node.cached.Sum(node.encoded)
}

func (node *ProtoDagNode) Marshal() ([]byte, error) {

	pbn := node.getPBNode()

	data, err := pbn.Marshal()
	if err != nil {
		return data, fmt.Errorf("marshal failed. %v", err)
	}
	return data, nil
}

func (node *ProtoDagNode) getPBNode() *pb.PBNode {

	pbn := &pb.PBNode{}
	if len(node.links) > 0 {
		pbn.Links = make([]*pb.PBLink, len(node.links))
	}

	sort.Stable(LinkSlice(node.links)) // keep links sorted
	for i, l := range node.links {
		pbn.Links[i] = &pb.PBLink{}
		pbn.Links[i].Name = l.Name
		pbn.Links[i].Tsize = l.Size

		if l.Cid != nil {
			pbn.Links[i].Hash = l.Cid.Bytes()
		}
	}

	if len(node.data) > 0 {
		pbn.Data = node.data
	}
	return pbn
}

func (node *ProtoDagNode) AddRawLink(name string, l *DagLink) error {

	node.encoded = nil
	node.links = append(node.links, &DagLink{
		Name: name,
		Size: l.Size,
		Cid:  l.Cid,
	})

	return nil
}

func NodeWithData(d []byte) *ProtoDagNode {

	return &ProtoDagNode{
		data: d,
	}
}

/*****************************************************************
*
*		ProtoDagNode implements.
*
*****************************************************************/
func (node *ProtoDagNode) AddChild(name string, child DagNode) error {

	//TODO:: really need to remove and add ?
	node.RemoveChild(name)

	return node.AddNodeLink(name, child)
}

func (node *ProtoDagNode) ForEachLink(func(*DagLink) error) error {
	return nil
}

func (node *ProtoDagNode) Find(string) (DagNode, error) {
	return nil, nil
}

func (node *ProtoDagNode) RemoveChild(string) error {
	return nil
}

func (node *ProtoDagNode) GetNode() (DagNode, error) {
	return node, nil
}
