package routing

import (
	"context"
	"github.com/libp2p/go-libp2p-peer"
	"github.com/libp2p/go-libp2p-peerstore"
)

type Routing interface {
	Ping(context.Context, peer.ID) error

	FindPeer(context.Context, peer.ID) (peerstore.PeerInfo, error)

	PutValue(context.Context, string, []byte) error

	GetValue(context.Context, string) ([]byte, error)

	Run()
}
