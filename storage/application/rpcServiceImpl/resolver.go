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
	rootNode    ipld.DagNode
	currentNode chan ipld.DagNode
	readingErr  chan error
	parentUris  []string //TODO:: try to suport multi directory resolve.
}

//TODO:: rawData should be ok later.
func ReadStreamData(cidKey *cid.Cid, uris []string) (UrlResolver, error) {

	dagService := merkledag.GetDagInstance()
	node, err := dagService.Get(cidKey)
	if err != nil {
		return nil, err
	}

	resolver := &nbsUrlResolver{
		rootNode:    node,
		currentNode: make(chan ipld.DagNode),
		readingErr:  make(chan error),
		parentUris:  uris,
	}

	go resolver.createReader(node)

	return resolver, nil
}

func (resolver *nbsUrlResolver) createReader(rootNode ipld.DagNode) {

	resolver.traverseNode(rootNode)

	close(resolver.currentNode)

	resolver.readingErr <- io.EOF

	logger.Info("*************>>")
}

func (resolver *nbsUrlResolver) traverseNode(rootNode ipld.DagNode) {

	links := rootNode.Links()

	dagService := merkledag.GetDagInstance()

	if links == nil || len(links) == 0 {
		logger.Info("+++++++++++>", rootNode.String())
		resolver.currentNode <- rootNode
		return

	} else {
		for _, link := range links {
			node, err := dagService.Get(link.Cid)
			if err != nil {
				logger.Error(err)
				resolver.readingErr <- err
				return
			} else {
				resolver.traverseNode(node)
			}
		}
		logger.Debug(rootNode.String())
	}
}

func parseToBridgeNode(node ipld.DagNode) (*DagDataBridge, error) {

	bridgeNode := new(DagDataBridge)
	var ok bool
	bridgeNode.dag, ok = node.(*ipld.ProtoDagNode)
	if !ok {
		return nil, errors.New("only support protoDagNode right now. ")
	}

	bridgeNode.format = &unixfs_pb.Data{}

	err := proto.Unmarshal(node.Data(), bridgeNode.format)
	if err != nil {
		return nil, err
	}

	//TODO:: data_directory
	if bridgeNode.Type() != unixfs_pb.Data_File {
		return nil, ErrIsNotFileData
	}

	return bridgeNode, nil
}

func (resolver *nbsUrlResolver) Next() ([]byte, error) {

	select {

	case node, ok := <-resolver.currentNode:

		if !ok {
			return nil, io.EOF
		}

		bridgeNode, err := parseToBridgeNode(node)
		if err != nil {
			return nil, err
		}
		logger.Info("--------->", node.String())
		return bridgeNode.format.Data, err

	case err := <-resolver.readingErr:
		logger.Info(err.Error())
		logger.Info("*************>>")
		return nil, err
	}
}

func (resolver *nbsUrlResolver) Close() error {
	//TODO::
	return nil
}
