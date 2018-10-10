package bitswap

import (
	"github.com/NBSChain/go-nbs/storage/merkledag/cid"
	"github.com/NBSChain/go-nbs/storage/merkledag/ipld"
	"github.com/libp2p/go-libp2p-peer"
)

type Exchange interface {

	GetDagNode(*cid.Cid) (ipld.DagNode, error)

	GetDagNodes([]*cid.Cid) (<-chan ipld.DagNode, error)

	SaveToNetPeer([]ipld.DagNode) error
}

type SwapLedger interface {

	Score() float32

	Threshold() float32
}


type LedgerEngine interface {

	ReceiveData(fromNode peer.ID, data []byte) SwapLedger

	SupportData(toNode peer.ID, data []byte) SwapLedger

	GetLedger(nodeId peer.ID) SwapLedger
}
