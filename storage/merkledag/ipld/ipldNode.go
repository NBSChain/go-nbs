package ipld

import (
	"errors"
	"fmt"
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

type Block interface {
	Data() []byte
	RawData() []byte
	Cid() *cid.Cid
	String() string
}

type DagNode interface {
	Block

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
	Size int64
	Cid  *cid.Cid
}

var v0Cid = cid.Cid{
	Version:  0,
	Code:     cid.DagProtobuf,
	HashType: multihash.SHA2_256,
	HashLen:  -1,
}
var v1Cid = cid.Cid{
	Version:  1,
	Code:     cid.DagProtobuf,
	HashType: multihash.SHA2_256,
	HashLen:  -1,
}

type LinkSlice []*DagLink

func (ls LinkSlice) Len() int           { return len(ls) }
func (ls LinkSlice) Swap(a, b int)      { ls[a], ls[b] = ls[b], ls[a] }
func (ls LinkSlice) Less(a, b int) bool { return ls[a].Name < ls[b].Name }

type ProtoDagNode struct {
	links     []*DagLink
	data      []byte
	encoded   []byte
	cached    *cid.Cid
	linkCache map[string]int
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
*		Block Interface
*
*****************************************************************/
func (n *ProtoDagNode) RawData() []byte {
	n.EncodeProtoBuf(false)
	return n.encoded
}

func (n *ProtoDagNode) Data() []byte {
	return n.data
}

func (n *ProtoDagNode) Cid() *cid.Cid {

	if n.encoded != nil && n.cached != nil {
		return n.cached
	}

	err := n.EncodeProtoBuf(false)
	if err != nil {
		err = fmt.Errorf("invalid CID of length %d: %x: %v", len(n.RawData()), n.RawData(), err)
		panic(err)
	}

	return n.cached
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
	return n.links
}

func (n *ProtoDagNode) Size() (int64, error) {

	err := n.EncodeProtoBuf(false)
	if err != nil {
		return 0, err
	}

	size := int64(len(n.encoded))

	for _, l := range n.links {
		size += l.Size
	}

	return size, nil
}

/*****************************************************************
*
*		Logic
*
*****************************************************************/
func (n *ProtoDagNode) SetData(d []byte) {
	n.encoded = nil
	n.cached = nil
	n.data = d
}

func (n *ProtoDagNode) AddNodeLink(name string, that DagNode) error {

	n.encoded = nil

	link, err := MakeLink(that)
	if err != nil {
		return err
	}
	link.Name = name

	n.AddRawLink(name, link)

	return nil
}

func (n *ProtoDagNode) EncodeProtoBuf(force bool) error {

	if n.encoded != nil && !force {
		return nil
	}

	var err error
	n.encoded, err = n.Marshal()
	if err != nil {
		return err
	}

	//TODO:: add V2 cid version.
	v1 := v0Cid
	n.cached = &v1

	return n.cached.Sum(n.encoded)
}

func (n *ProtoDagNode) Marshal() ([]byte, error) {

	pbn := n.getPBNode()

	data, err := pbn.Marshal()
	if err != nil {
		return data, fmt.Errorf("marshal failed. %v", err)
	}
	return data, nil
}

func (n *ProtoDagNode) getPBNode() *pb.PBNode {

	pbn := &pb.PBNode{}
	if len(n.links) > 0 {
		pbn.Links = make([]*pb.PBLink, len(n.links))
	}

	sort.Stable(LinkSlice(n.links)) // keep links sorted
	for i, l := range n.links {
		pbn.Links[i] = &pb.PBLink{}
		pbn.Links[i].Name = l.Name
		pbn.Links[i].Tsize = l.Size

		if l.Cid != nil {
			pbn.Links[i].Hash = l.Cid.Bytes()
		}
	}

	if len(n.data) > 0 {
		pbn.Data = n.data
	}
	return pbn
}

func (n *ProtoDagNode) AddRawLink(name string, l *DagLink) error {

	n.encoded = nil
	n.links = append(n.links, &DagLink{
		Name: name,
		Size: l.Size,
		Cid:  l.Cid,
	})

	n.linkCache[name] = len(n.links) - 1

	return nil
}

func NodeWithData(d []byte) *ProtoDagNode {

	return &ProtoDagNode{
		data:      d,
		linkCache: make(map[string]int),
	}
}
func NewNode() *ProtoDagNode {
	return &ProtoDagNode{
		linkCache: make(map[string]int),
	}
}

func (n *ProtoDagNode) unmarshal(encoded []byte) error {
	var pbn pb.PBNode
	if err := pbn.Unmarshal(encoded); err != nil {
		return fmt.Errorf("unmarshal failed. %v", err)
	}

	pbnl := pbn.GetLinks()
	n.links = make([]*DagLink, len(pbnl))
	for i, l := range pbnl {
		n.links[i] = &DagLink{Name: l.GetName(), Size: l.GetTsize()}
		c, err := cid.Cast(l.GetHash())
		if err != nil {
			return fmt.Errorf("link hash #%d is not valid multihash. %v", i, err)
		}
		n.links[i].Cid = c
	}
	sort.Stable(LinkSlice(n.links))

	n.data = pbn.GetData()
	n.encoded = encoded
	return nil
}

/*****************************************************************
*
*		ProtoDagNode implements.
*
*****************************************************************/
func (n *ProtoDagNode) AddChild(name string, child DagNode) error {

	n.RemoveChild(name)

	return n.AddNodeLink(name, child)
}

func (n *ProtoDagNode) ForEachLink(func(*DagLink) error) error {
	return nil
}

func (n *ProtoDagNode) Find(string) (DagNode, error) {
	return nil, nil
}

//TODO:: need to test.
func (n *ProtoDagNode) RemoveChild(name string) error {

	if index, ok := n.linkCache[name]; ok {

		n.links = n.links[:index]
		n.linkCache = make(map[string]int)

		for newIdx, l := range n.links {
			n.linkCache[l.Name] = newIdx
		}

		return nil
	}

	return errors.New("can't find the dag link:" + name)
}
