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
		position:    -1,//-1 means start form self ,not sub nodes from links.
		parentUris:  uris,
	}
	resolver.links = make([]*ipld.DagLink, len(node.Links()))
	copy(resolver.links, node.Links())

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

	if resolver.position >= len(resolver.links){
		return nil, io.EOF
	}

	//TODO:: need to check what we can do if result is nil.
	result := resolver.currentNode.format.Data

	dataType := *resolver.currentNode.format.Type
	logger.Info("data type is :", dataType)

	resolver.position++
	if resolver.position >= len(resolver.links){
		resolver.currentNode = nil
		return result, nil
	}

	link := resolver.links[resolver.position]

	dagService := merkledag.GetDagInstance()
	node, err := dagService.Get(link.Cid)
	if err != nil{
		return nil, err
	}

	curNode, err := parseToBridgeNode(node)
	if err != nil{
		return nil ,err
	}

	resolver.currentNode = curNode

	return result, nil
}

func (resolver *nbsUrlResolver) Close() error{
	return nil
}