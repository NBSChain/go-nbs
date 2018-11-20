package nbsnet

import (
	"github.com/NBSChain/go-nbs/storage/network/nat"
	"sync"
)

type ConnManager struct {
	sync.Mutex
	natKATun *nat.KATunnel
	queue    map[string]*NbsUdpConn
}

func NewConnManager() *ConnManager {

	cm := &ConnManager{
		queue: make(map[string]*NbsUdpConn),
	}

	go cm.runLoop()

	return cm
}

func (manager *ConnManager) put(conn *NbsUdpConn) {
	manager.Lock()
	defer manager.Unlock()

	if _, ok := manager.queue[conn.ConnId]; ok {
		return
	}

	manager.queue[conn.ConnId] = conn
}

func (manager *ConnManager) Close(connId string) {
	delete(manager.queue, connId)
}

func (manager *ConnManager) runLoop() {

	//for {
	//	for id, conn := range manager.queue{
	//		if conn.isClosed{
	//
	//		}
	//	}
	//	time.Sleep(time.Second * 5)//TODO::
	//}
}
