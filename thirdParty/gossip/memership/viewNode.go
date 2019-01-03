package memership

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/thirdParty/gossip/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"net"
	"sync"
	"time"
)

type ViewNode struct {
	sync.RWMutex
	nodeId      string
	probability float64
	inAddr      *net.UDPAddr
	updateTime  time.Time
	expiredTime time.Time
	outConn     *nbsnet.NbsUdpConn
	outAddr     *nbsnet.NbsUdpAddr
}

func (node *MemManager) newOutViewNode(host *pb.BasicHost, expire time.Time) (*ViewNode, error) {

	addr := nbsnet.ConvertFromGossipAddr(host)
	port := utils.GetConfig().GossipCtrlPort

	conn, err := network.GetInstance().Connect(nil, addr, port)
	if err != nil {
		logger.Error("the contact failed to notify the subscriber:", err, addr, port)
		return nil, err
	}

	item := &ViewNode{
		nodeId:      host.NetworkId,
		probability: node.meanProb(node.PartialView),
		outConn:     conn,
		outAddr:     addr,
		updateTime:  time.Now(),
		expiredTime: expire,
	}

	node.PartialView[item.nodeId] = item
	go node.waitingWork(item)

	return item, nil
}

func (node *MemManager) newInViewNode(nodeId string, addr *net.UDPAddr) *ViewNode {

	view := &ViewNode{
		nodeId:      nodeId,
		inAddr:      addr,
		probability: node.meanProb(node.InputView),
		updateTime:  time.Now(),
	}
	node.InputView[nodeId] = view
	return view
}

func (node *MemManager) sendData(item *ViewNode, data []byte) error {
	item.Lock()
	defer item.Unlock()

	if _, err := item.outConn.Write(data); err != nil {
		node.removeFromView(item, node.PartialView)
		return err
	}
	item.updateTime = time.Now()
	return nil
}

func (node *MemManager) send(item *ViewNode, pb proto.Message) error {
	item.Lock()
	defer item.Unlock()

	data, err := proto.Marshal(pb)

	if err != nil {
		return err
	}

	if _, err := item.outConn.Write(data); err != nil {
		node.removeFromView(item, node.PartialView)
		return err
	}
	item.updateTime = time.Now()

	return nil
}

func (node *MemManager) waitingWork(item *ViewNode) {

	for {
		buffer := make([]byte, utils.NormalReadBuffer)
		n, addr, err := item.outConn.ReadFromUDP(buffer)

		if err != nil {
			logger.Warning("node in view read err:->", err, item.nodeId)
			node.taskQueue <- &gossipTask{
				taskType: NodeFailed,
				innerTask: innerTask{
					params: item,
				},
			}
			return
		}
		msg := &pb.Gossip{}
		if err := proto.Unmarshal(buffer[:n], msg); err != nil {
			logger.Warning("unmarshal err:->", err, addr, n)
			continue
		}
		task := &gossipTask{
			taskType: int(msg.MsgType),
		}
		task.msg = msg
		task.addr = addr
		node.taskQueue <- task
	}
}

func (item *ViewNode) String() string {
	item.RLock()
	defer item.RUnlock()

	format := utils.GetConfig().SysTimeFormat

	var inAddr, outAddr string
	if item.inAddr != nil {
		inAddr = item.inAddr.String()
	}
	if item.outAddr != nil {
		outAddr = item.outAddr.String()
	}
	return fmt.Sprintf("------------%s------------\n"+
		"|%-15s:%20.2f|\n"+
		"|%-15s:%20s|\n"+
		"|%-15s:%20s|\n"+
		"|%-15s:%20s|\n"+
		"|%-15s:%20s|\n"+
		"-----------------------------------------------------------------------\n",
		item.nodeId,
		"probability",
		item.probability,
		"inAddr",
		inAddr,
		"updateTime",
		item.updateTime.Format(format),
		"expiredTime",
		item.expiredTime.Format(format),
		"outAddr",
		outAddr,
	)
}

func (node *MemManager) freshInputView(nodeId string) {
	if item, ok := node.InputView[nodeId]; ok {
		item.Lock()
		defer item.Unlock()
		logger.Debug("update input view item for msg receive :->", item.nodeId)
		item.updateTime = time.Now()
	}
}
