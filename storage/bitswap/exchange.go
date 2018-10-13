package bitswap

import (
	"context"
	"github.com/NBSChain/go-nbs/storage/bitswap/fetcher"
	"github.com/NBSChain/go-nbs/storage/merkledag/cid"
	"github.com/NBSChain/go-nbs/storage/merkledag/ipld"
	"github.com/libp2p/go-libp2p-peer"
)



type Exchange interface {

	GetDagNode(*cid.Cid) (ipld.DagNode, error)

	GetDagNodes(context.Context, []*cid.Cid) (<-chan fetcher.AsyncResult)

	SaveToNetPeer(map[string]ipld.DagNode) error

	GetLedgerEngine() LedgerEngine
}

type SwapLedger interface {

	Score() float64

	Threshold() float64
}


type LedgerEngine interface {

	ReceiveData(fromNode peer.ID, data []byte) SwapLedger

	SupportData(toNode peer.ID, data []byte) SwapLedger

	GetLedger(nodeId peer.ID) SwapLedger
}
