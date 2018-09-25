package rpcServiceImpl

import (
	"errors"
	"github.com/NBSChain/go-nbs/storage/application/pb"
	"github.com/NBSChain/go-nbs/storage/merkledag"
	"github.com/NBSChain/go-nbs/storage/merkledag/cid"
	"github.com/NBSChain/go-nbs/storage/merkledag/ipld"
	"github.com/gogo/protobuf/proto"
	"io"
)

type UrlResolver interface {
	io.Closer
	Next() ([]byte, error)
}
var ErrIsNotData = errors.New("this dag node is not a file")

type nbsUrlResolver struct {
	rootNode	*DagDataBridge//ipld.DagNode
	position	int
	links		[]*ipld.DagLink
}


func ReadStreamData(cidKey *cid.Cid)  (UrlResolver, error){

	dagService := merkledag.GetDagInstance()
	node, err := dagService.Get(cidKey)

	if err != nil{
		return nil, err
	}

	bridgeNode := new(DagDataBridge)
	bridgeNode.dag = node.(*ipld.ProtoDagNode)
	bridgeNode.format = &unixfs_pb.Data{}

	err = proto.Unmarshal(node.Data(), bridgeNode.format)
	if err != nil {
		return nil, err
	}

	return &nbsUrlResolver{
		rootNode:	bridgeNode,
		links:		node.Links(),
		position:	0,
	}, nil
}


func (resolver *nbsUrlResolver) Next() ([]byte, error)  {

	//TODO:: rawData should be ok later.
	if resolver.rootNode.Type() != TFile{
		return nil, ErrIsNotData
	}

	dagService := merkledag.GetDagInstance()

	if resolver.position >= len(resolver.links){
		return nil, io.EOF
	}

	link := resolver.links[resolver.position]
	resolver.position++

	node, err := dagService.Get(link.Cid)
	if err != nil{
		return nil, err
	}

	return node.Data(), nil
}

func (resolver *nbsUrlResolver) Close() error{
	return nil
}