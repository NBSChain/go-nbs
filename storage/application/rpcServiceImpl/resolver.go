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
var ErrIsNotFileData = errors.New("this dag node is not a file")

type nbsUrlResolver struct {
	currentNode *DagDataBridge
	position    int
	links       []*ipld.DagLink
	parentUris  []string	//TODO:: try to suport multi directory resolve.
}

//TODO:: rawData should be ok later.
func ReadStreamData(cidKey *cid.Cid, uris []string)  (UrlResolver, error){

	dagService := merkledag.GetDagInstance()
	node, err := dagService.Get(cidKey)
	if err != nil{
		return nil, err
	}

	bridgeNode, err := parseToBridgeNode(node)
	if err != nil{
		return nil, err
	}

	resolver := &nbsUrlResolver{
		currentNode: bridgeNode,
		position:    0,
		parentUris:  uris,
	}

	if len(node.Links()) > 0 {
		resolver.links = make([]*ipld.DagLink, len(node.Links()))
		copy(resolver.links, node.Links())
	}

	return resolver, nil
}

func parseToBridgeNode(node ipld.DagNode) (*DagDataBridge, error)  {

	bridgeNode := new(DagDataBridge)
	var ok bool
	bridgeNode.dag, ok = node.(*ipld.ProtoDagNode)
	if !ok{
		return nil, errors.New("only support protoDagNode right now. ")
	}

	bridgeNode.format = &unixfs_pb.Data{}

	err := proto.Unmarshal(node.Data(), bridgeNode.format)
	if err != nil {
		return nil, err
	}

	//TODO:: data_directory
	if bridgeNode.Type() != unixfs_pb.Data_File{
		return nil, ErrIsNotFileData
	}

	return bridgeNode, nil
}


func (resolver *nbsUrlResolver) Next() ([]byte, error)  {

	//It's a leaf node
	if len(resolver.links) == 0{

		data := resolver.currentNode.format.Data

		if data == nil{
			return nil, io.EOF
		}else{
			resolver.currentNode.format.Data = nil
			return data, nil
		}
	}

	if resolver.position >= len(resolver.links){
		return nil, io.EOF
	}


	dagService := merkledag.GetDagInstance()

	link := resolver.links[resolver.position]
	node, err := dagService.Get(link.Cid)
	if err != nil{
		return nil, err
	}

	curNode, err := parseToBridgeNode(node)
	if err != nil{
		return nil ,err
	}

	resolver.currentNode = curNode
	resolver.position++

	return curNode.format.Data, nil
}

func (resolver *nbsUrlResolver) Close() error{
	//TODO::
	return nil
}