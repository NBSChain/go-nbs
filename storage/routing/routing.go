package routing

import (
	"github.com/libp2p/go-libp2p-peerstore"
)

type Routing interface {

	Ping(peer peerstore.PeerInfo) Pong

	FindPeer(key string) (chan []peerstore.PeerInfo, error)//return k peers most closet to key

	PutValue(key string, value []byte) chan error

	GetValue(key string) (chan []byte, chan []peerstore.PeerInfo, error)//return value or k peers most closet to key
}

type Pong interface {
	Status() bool
	//TODO::more details
}

