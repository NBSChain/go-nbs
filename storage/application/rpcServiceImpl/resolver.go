package rpcServiceImpl

import (
	"github.com/NBSChain/go-nbs/storage/merkledag"
	"github.com/NBSChain/go-nbs/storage/merkledag/cid"
	"github.com/NBSChain/go-nbs/storage/merkledag/ipld"
	"io"
)

type UrlResolver interface {
	io.Closer
	Next() ([]byte, error)
}

type nbsUrlResolver struct {
	rootNode	ipld.DagNode
}


func ReadStreamData(cidKey *cid.Cid)  (UrlResolver, error){

	dagService := merkledag.GetDagInstance()
	node, err := dagService.Get(cidKey)

	if err != nil{
		return nil, err
	}

	return &nbsUrlResolver{
		rootNode:node,
	}, nil
}


func (resolver *nbsUrlResolver) Next() ([]byte, error)  {
	return nil, nil
}

func (resolver *nbsUrlResolver) Close() error{
	return nil
}