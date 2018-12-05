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
}

func newOutViewNode(sub *pb.Subscribe, views map[string]*viewNode) (*viewNode, error) {

	addr := nbsnet.ConvertFromGossipAddr(sub.Addr)
	port := utils.GetConfig().GossipCtrlPort

	conn, err := network.GetInstance().Connect(nil, addr, port)
	if err != nil {
		logger.Error("the contact failed to notify the subscriber:", err)
		return nil, err
	}

	item := &viewNode{
		nodeId:      sub.Addr.NetworkId,
		outConn:     conn,
		updateTime:  time.Now(),
		expiredTime: time.Now().Add(time.Duration(sub.Duration)),
	}

	views[item.nodeId] = item
	item.probability = 1 / float64(len(views))
	updateProbability(views)

	return item, nil
}

func newInViewNode(nodeId string, addr *net.UDPAddr, views map[string]*viewNode) *viewNode {

	view := &viewNode{
		nodeId:     nodeId,
		inAddr:     addr,
		updateTime: time.Now(),
	}

	views[nodeId] = view
	view.probability = 1 / float64(len(views))
	updateProbability(views)

	return view
}

//TODO:: broadcast probability???
func updateProbability(view map[string]*viewNode) {

	var summerOut float64
	for _, item := range view {
		summerOut += item.probability
	}

	for _, item := range view {
		item.probability = item.probability / summerOut
	}
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
		if _, err := item.outConn.Read(buffer); err != nil {
			logger.Warning("node in view read err:->", err)
		}
	}
}
