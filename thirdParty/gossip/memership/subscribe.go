package memership

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/thirdParty/gossip/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/gogo/protobuf/proto"
	"net"
	"time"
)

const (
	InitSubscribeTimeOut = time.Second * 3
)

/*****************************************************************
*
*	member client functions about init subscribe request.
*
*****************************************************************/
func (node *MemManager) registerMySelf() error {

	servers := utils.GetConfig().GossipBootStrapIP

	var success bool

	for _, serverIp := range servers {

		conn, err := network.GetInstance().DialUDP("udp4",
			nil, &net.UDPAddr{
				IP:   net.ParseIP(serverIp),
				Port: utils.GetConfig().GossipCtrlPort,
			})

		if err != nil {
			logger.Error("dial to contract boot server failed:", err)
			goto CloseConn
		}

		conn.SetDeadline(time.Now().Add(InitSubscribeTimeOut))

		if err := node.acquireProxy(conn); err != nil {
			logger.Error("send init sub request failed:", err)
			goto CloseConn
		}

		if err := node.checkProxyValidation(conn); err == nil {
			logger.Debug("find gossip contact server success.", serverIp)
			success = true
			break
		}

	CloseConn:
		conn.Close()
	}

	if !success {
		return fmt.Errorf("failed to find a contract server")
	}

	return nil
}

func (node *MemManager) acquireProxy(conn *nbsnet.NbsUdpConn) error {

	msg := &pb.Gossip{
		MessageType: pb.MsgType_init,
		InitMsg: &pb.InitSub{
			Seq:    time.Now().Unix(),
			NodeId: node.nodeID,
			Addr:   nbsnet.ConvertToGossipAddr(conn.LocAddr),
		},
	}
	msgData, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	if _, err := conn.Write(msgData); err != nil {
		return err
	}

	return nil
}

func (node *MemManager) checkProxyValidation(conn *nbsnet.NbsUdpConn) error {

	buffer := make([]byte, utils.NormalReadBuffer)
	hasRead, err := conn.Read(buffer)
	if err != nil {
		return err
	}

	msg := &pb.Gossip{}
	if err := proto.Unmarshal(buffer[:hasRead], msg); err != nil {
		return err
	}

	if msg.MessageType != pb.MsgType_initACK {
		return fmt.Errorf("failed to send init sub request")
	}

	if msg.InitACK.SupplierID == node.nodeID {
		return fmt.Errorf("it's yourself")
	}

	return nil
}

/*****************************************************************
*
*	member server functions about init subscribe request.
*
*****************************************************************/
func (node *MemManager) confirmAndPrepare(request *pb.InitSub, peerAddr *net.UDPAddr) {

	message := &pb.Gossip{
		MessageType: pb.MsgType_initACK,
		InitACK: &pb.InitSubACK{
			Seq:        request.Seq,
			SupplierID: node.nodeID,
		},
	}

	msgData, _ := proto.Marshal(message)
	if _, err := node.serviceConn.WriteToUDP(msgData, peerAddr); err != nil {
		logger.Warning("failed to send init ack msg:", err)
		return
	}

	if node.nodeID == request.NodeId {
		return
	}

	task := innerTask{
		tType: ProxyInitSubRequest,
		param: make([]interface{}, 2),
	}

	rAddr := nbsnet.ConvertFromGossipAddr(request.Addr, peerAddr.IP.String())
	rAddr.NetworkId = request.NodeId

	task.param[0] = request
	task.param[1] = rAddr

	node.taskSignal <- task
}

func (node *MemManager) subToContract(ack *pb.ReqContactACK, addr *net.UDPAddr) {

	logger.Debug("gossip sub start:", ack, addr)

	_, ok := node.inputView[ack.SupplierID]
	if ok {
		logger.Info("duplicated sub confirm")
		return
	}

	item := &peerNodeItem{
		nodeId: ack.SupplierID,
		addr:   nbsnet.ConvertFromGossipAddr(ack.Supplier, addr.IP.String()),
	}

	node.inputView[ack.SupplierID] = item
	item.probability = 1 / float64(len(node.inputView))
}
