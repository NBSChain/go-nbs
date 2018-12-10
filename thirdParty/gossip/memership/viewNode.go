package memership

import (
	"github.com/NBSChain/go-nbs/storage/network"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/thirdParty/gossip/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"net"
	"time"
)

//TODO:: same msg forward times

type viewNode struct {
	nodeId      string
	probability float64
	inAddr      *net.UDPAddr
	updateTime  time.Time
	expiredTime time.Time
	outConn     *nbsnet.NbsUdpConn
	outAddr     *nbsnet.NbsUdpAddr
	manager     *MemManager
}

func (node *MemManager) newOutViewNode(host *pb.BasicHost, duration int64) (*viewNode, error) {

	addr := nbsnet.ConvertFromGossipAddr(host)
	port := utils.GetConfig().GossipCtrlPort

	conn, err := network.GetInstance().Connect(nil, addr, port)
	if err != nil {
		logger.Error("the contact failed to notify the subscriber:", err)
		return nil, err
	}

	item := &viewNode{
		nodeId:      host.NetworkId,
		outConn:     conn,
		outAddr:     addr,
		manager:     node,
		updateTime:  time.Now(),
		expiredTime: time.Now().Add(time.Duration(duration)),
	}

	node.partialView[item.nodeId] = item
	item.probability = 1 / float64(len(node.partialView))
	go item.waitingWork()

	return item, nil
}

func (node *MemManager) newInViewNode(nodeId string, addr *net.UDPAddr) *viewNode {

	view := &viewNode{
		nodeId:     nodeId,
		inAddr:     addr,
		manager:    node,
		updateTime: time.Now(),
	}

	node.inputView[nodeId] = view
	view.probability = 1 / float64(len(node.inputView))
	return view
}

func (item *viewNode) needUpdate() bool {
	return time.Now().Sub(item.updateTime) >= MemShipHeartBeat
}

func (item *viewNode) sendData(data []byte) error {

	if _, err := item.outConn.Write(data); err != nil {
		return err
	}

	item.updateTime = time.Now()

	return nil
}

func (item *viewNode) send(pb proto.Message) error {

	data, err := proto.Marshal(pb)

	if err != nil {
		return err
	}

	if _, err := item.outConn.Write(data); err != nil {
		return err
	}

	item.updateTime = time.Now()

	return nil
}

func (item *viewNode) waitingWork() {

	for {
		buffer := make([]byte, utils.NormalReadBuffer)
		n, err := item.outConn.Read(buffer)
		if err != nil {
			logger.Warning("node in view read err:->", err)
			break
		}
		msg := &pb.Gossip{}
		if err := proto.Unmarshal(buffer[:n], msg); err != nil {
			logger.Warning("unmarshal err:->", err)
			continue
		}

		addr := item.outConn.RealConn.RemoteAddr()
		task := &msgTask{
			msg:  msg,
			addr: addr.(*net.UDPAddr),
		}
		item.manager.taskQueue <- task
	}
}
