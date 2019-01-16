package memership

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/thirdParty/gossip/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"net"
	"time"
)

type ViewNode struct {
	nodeId        string
	probability   float64
	inAddr        *net.UDPAddr
	heartBeatTime time.Time
	expiredTime   time.Time
	outConn       *nbsnet.NbsUdpConn
	outAddr       *nbsnet.NbsUdpAddr
}

func (node *MemManager) newOutViewNode(host *pb.BasicHost, expire time.Time) (*ViewNode, error) {

	addr := nbsnet.ConvertFromGossipAddr(host)
	if oldItem, ok := node.PartialView[host.NetworkId]; ok { //TODO::check again.
		oldItem.expiredTime = expire
		oldItem.heartBeatTime = time.Now()
		oldItem.outAddr = addr
		logger.Debug("I've this item already :->", oldItem.nodeId)
		return oldItem, nil
	}

	port := utils.GetConfig().GossipCtrlPort
	conn, err := network.GetInstance().Connect(nil, addr, port)
	if err != nil {
		logger.Error("the contact failed to notify the subscriber:", err, addr, port)
		return nil, err
	}

	item := &ViewNode{
		nodeId:        host.NetworkId,
		probability:   node.meanProb(node.PartialView),
		outConn:       conn,
		outAddr:       addr,
		heartBeatTime: time.Now(),
		expiredTime:   expire,
	}

	node.PartialView[item.nodeId] = item
	go node.waitingWork(item)

	return item, nil
}

func (node *MemManager) newInViewNode(nodeId string, addr *net.UDPAddr) *ViewNode {

	view := &ViewNode{
		nodeId:        nodeId,
		inAddr:        addr,
		probability:   node.meanProb(node.InputView),
		heartBeatTime: time.Now(),
		outConn:       nil,
		outAddr:       nil,
	}
	node.InputView[nodeId] = view
	return view
}

func (node *MemManager) sendData(item *ViewNode, data []byte) error {

	if _, err := item.outConn.WriteWithSyn(data); err != nil {
		logger.Warning("item write data to peer err:->", err)
		node.removeFromView(item.nodeId, node.PartialView)
		return err
	}
	return nil
}

func (node *MemManager) send(item *ViewNode, pb proto.Message) error {
	data, err := proto.Marshal(pb)
	if err != nil {
		return err
	}
	return node.sendData(item, data)
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
					params: item.nodeId,
				},
			}
			return
		}
		msg := &pb.Gossip{}
		if err := proto.Unmarshal(buffer[:n], msg); err != nil {
			logger.Warning("unmarshal err:->", err, addr, n)
			continue
		}
		logger.Debug("receive ctrl msg:->", msg)
		task := &gossipTask{
			taskType: int(msg.MsgType),
		}
		task.msg = msg
		task.addr = addr
		node.taskQueue <- task
	}
}

func (item *ViewNode) String() string {

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
		"heartBeatTime",
		item.heartBeatTime.Format(format),
		"expiredTime",
		item.expiredTime.Format(format),
		"outAddr",
		outAddr,
	)
}
