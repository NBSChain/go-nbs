package network

import "sync"

type ConnManager struct {
	sync.Mutex
	queue map[string]*NbsUdpConn
}

func newConnManager() *ConnManager {
	cm := &ConnManager{
		queue: make(map[string]*NbsUdpConn),
	}
	return cm
}

func (manager *ConnManager) KeepAliveConn(conn *NbsUdpConn) {
	manager.Lock()
	defer manager.Unlock()

	if _, ok := manager.queue[conn.connId]; ok {
		return
	}

	manager.queue[conn.connId] = conn
}

func (manager *ConnManager) runLoop() {

	for {

	}
}
