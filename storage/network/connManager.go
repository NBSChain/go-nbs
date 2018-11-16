package network

import (
	"github.com/NBSChain/go-nbs/storage/network/nat"
	"sync"
)

type ConnManager struct {
	sync.Mutex
	natKATun *nat.KATunnel
	queue    map[string]*NbsUdpConn
}

func newConnManager() *ConnManager {
	cm := &ConnManager{
		queue: make(map[string]*NbsUdpConn),
	}

	go cm.runLoop()

	return cm
}

func (manager *ConnManager) put(conn *NbsUdpConn) {
	manager.Lock()
	defer manager.Unlock()

	if _, ok := manager.queue[conn.connId]; ok {
		return
	}

	manager.queue[conn.connId] = conn
}

func (manager *ConnManager) runLoop() {
}
