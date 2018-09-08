package core

import "github.com/NBSChain/go-nbs/storage/merkledag/cid"

type Block interface {
	RawData() []byte
	Cid() *cid.Cid
	String() string
	Loggable() map[string]interface{}
}

type BasicBlock struct {
	cid  *cid.Cid
	data []byte
}
