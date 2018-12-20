package memership

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network"
	"github.com/NBSChain/go-nbs/storage/network/nat"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/thirdParty/gossip/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/gogo/protobuf/proto"
	"net"
	"time"
)

/*****************************************************************
*
*	member client functions about init subscribe request.
*
*****************************************************************/
func (node *MemManager) RegisterMySelf() error {

	servers := utils.GetConfig().GossipBootStrapIP

	var success bool

	for _, serverIp := range servers {

		logger.Debug("start to find gossip server:->", serverIp)

		conn, err := network.GetInstance().DialUDP("udp4",
			nil, &net.UDPAddr{
				IP:   net.ParseIP(serverIp),
				Port: utils.GetConfig().GossipCtrlPort,
			})

		if err != nil {
			logger.Warning("dial to contract boot server failed:->", err)
			goto CloseConn
		}

		if err := conn.SetDeadline(time.Now().Add(SubscribeTimeOut)); err != nil {
			logger.Warning("set outConn time out err:->", err)
			goto CloseConn
		}

		if err := node.acquireProxy(conn); err != nil {
			logger.Warning("send init sub request failed:", err)
			goto CloseConn
		}

		if err := node.checkProxyValidation(conn); err == nil {
			logger.Info("find gossip contact server success.", serverIp)
			success = true
			break
		}

	CloseConn:
		conn.Close()
	}

	if !success {
		if node.isBootNode {
			logger.Info("I'm a boot strap node and alone now")
			return nil
		}
		return fmt.Errorf("failed to find a contract server")
	}

	return nil
}

func (node *MemManager) acquireProxy(conn *nbsnet.NbsUdpConn) error {
	msg := &pb.Gossip{
		MsgType: nbsnet.GspSub,
		Subscribe: &pb.Subscribe{
			SeqNo:  1,
			Expire: time.Now().Add(DefaultSubExpire).Unix(),
			NodeId: node.nodeID,
			Addr:   nbsnet.ConvertToGossipAddr(conn.LocAddr, node.nodeID),
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
		logger.Warning("read contact proxy result err:->", err)
		return err
	}

	msg := &pb.Gossip{}
	if err := proto.Unmarshal(buffer[:hasRead], msg); err != nil {
		logger.Debug("it's not sub ack:->", err)
		return err
	}

	if msg.MsgType != nbsnet.GspSubACK {
		return fmt.Errorf("failed to send init sub request")
	}

	if msg.SubAck.FromId == node.nodeID {
		node.isBootNode = true
		return fmt.Errorf("it's yourself")
	}

	logger.Info("He will proxy our sub:->", conn.String())

	return nil
}

/*****************************************************************
*
*	member server functions about init subscribe request.
*
*****************************************************************/
func (node *MemManager) firstInitSub(task *gossipTask) error {

	subReq := task.msg.Subscribe
	peerAddr := task.addr

	subReq.SeqNo++
	message := &pb.Gossip{
		MsgType: nbsnet.GspSubACK,
		SubAck: &pb.SynAck{
			SeqNo:  subReq.SeqNo,
			FromId: node.nodeID,
		},
	}

	msgData, _ := proto.Marshal(message)
	if _, err := node.serviceConn.WriteToUDP(msgData, peerAddr); err != nil {
		logger.Warning("failed to send init ack msg:", err)
		return err
	}

	if node.nodeID == subReq.NodeId {
		logger.Info("it's yourself.")
		return nil
	}

	counter := 2 * len(node.PartialView)
	return node.asContactProxy(subReq, counter)
}

func (node *MemManager) subToContract(task *gossipTask) error {

	result := task.msg.VoteResult
	nodeId := result.NodeId
	expire := time.Unix(result.Expire, 0)
	item, ok := node.InputView[nodeId]
	if ok {
		logger.Info("duplicated sub confirm")
		item.expiredTime = expire
		return nil
	}

	logger.Debug("get contact node:->", result, task.addr)

	item, err := node.newOutViewNode(result.Addr, expire)
	if err != nil {
		logger.Error("sub to contact node:->", err)
		return err
	}

	node.newInViewNode(nodeId, task.addr)
	msg := &pb.Gossip{
		MsgType: nbsnet.GspVoteResAck,
		VoteAck: &pb.SynAck{
			SeqNo:  result.SeqNo + 1,
			FromId: node.nodeID,
		},
	}

	return item.send(msg)
}

func (node *MemManager) subAccepted(task *gossipTask) error {
	ack := task.msg.SubConfirm
	_, ok := node.InputView[ack.FromId]
	if ok {
		return fmt.Errorf("duplicated sub accepted")
	}

	node.newInViewNode(ack.FromId, task.addr)
	return nil
}

func (node *MemManager) Resub() error {

	if node.isBootNode {
		logger.Info("I'm the boot node, so maybe it's normal situation")
		return nil
	}

	if len(node.PartialView) == 0 {
		logger.Debug("register myself because of no partial view in my cache")
		return node.RegisterMySelf()
	}

	item := node.choseRandomInPartialView()
	logger.Debug("I am alone and need to subscribe to random node:->", item.nodeId)

	if err := node.acquireProxy(item.outConn); err != nil {
		logger.Warning("send reSub request err:->", err)
		node.removeFromView(item, node.PartialView)
		return err
	}

	if err := item.outConn.SetDeadline(time.Now().Add(SubscribeTimeOut)); err != nil {
		logger.Warning("set outConn time out when reSub err:->", err)
		node.removeFromView(item, node.PartialView)
		return err
	}

	if err := node.checkProxyValidation(item.outConn); err != nil {
		logger.Warning("can't reSub to node err:->", item.nodeId, err)
		node.removeFromView(item, node.PartialView)
		return err
	}

	if err := item.outConn.SetDeadline(nat.NoTimeOut); err != nil {
		logger.Warning("set outConn time out when reSub err:->", err)
		node.removeFromView(item, node.PartialView)
		return err
	}

	return nil
}

func (node *MemManager) reSubAckConfirm(task *gossipTask) error {
	logger.Debug("he will solve our reSub request:->", task.addr)
	return nil
}
