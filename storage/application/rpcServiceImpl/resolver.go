package rpcServiceImpl

import (
	"github.com/NBSChain/go-nbs/storage/merkledag/cid"
	"io"
)

type UrlResolver interface {
	io.Closer
	Next() ([]byte, error)
}

func ReadStreamData(cidKey *cid.Cid)  (UrlResolver, error){

	return nil, nil
}