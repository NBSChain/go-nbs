package rpcService

import (
	"github.com/NBSChain/go-nbs/storage/merkledag/ipld"
	"sync"
	"time"
)

type Directory struct {
	name      string
	lock      sync.Mutex
	modTime   time.Time
	childDirs map[string]*Directory
	files     map[string]*ipld.DagNode
}
