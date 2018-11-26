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

		if err := node.intSubStep1(conn); err != nil {
			logger.Error("send init sub request failed:", err)
			goto CloseConn
		}

		if err := node.intSubStep3(conn); err == nil {
			logger.Debug("find gossip contact server success.")
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

func (node *MemManager) intSubStep1(conn *nbsnet.NbsUdpConn) error {

	lAddr := conn.LocAddr

	msg := &pb.Gossip{
		MessageType: pb.MsgType_init,
		InitMsg: &pb.InitSub{
			NodeId:    node.peerId,
			CanServer: lAddr.CanServe,
			NatIP:     lAddr.NatIp,
			NatPort:   lAddr.NatPort,
			PriIP:     lAddr.PriIp,
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

func (node *MemManager) intSubStep3(conn *nbsnet.NbsUdpConn) error {

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

	return nil
}

/*****************************************************************
*
*	member server functions about init subscribe request.
*
*****************************************************************/
func (node *MemManager) intSubStep2(request *pb.InitSub, peerAddr *net.UDPAddr) {

	message := &pb.Gossip{
		MessageType: pb.MsgType_initACK,
		InitACK: &pb.InitSubACK{
			ContactId: node.peerId,
			ApplierId: request.NodeId,
		},
	}

	msgData, _ := proto.Marshal(message)
	if _, err := node.serviceConn.WriteToUDP(msgData, peerAddr); err != nil {
		logger.Warning("failed to send init ack msg:", err)
		return
	}

	task := innerTask{
		tType: ProxyInitSubRequest,
		param: make([]interface{}, 2),
	}

	rAddr := &nbsnet.NbsUdpAddr{
		NetworkId: request.NodeId,
		CanServe:  request.CanServer,
		NatServer: request.NatServer,
		NatIp:     request.NatIP,
		NatPort:   request.NatPort,
		PubIp:     peerAddr.IP.String(),
		PriIp:     request.PriIP,
	}
	task.param[0] = request
	task.param[1] = rAddr

	node.taskSignal <- task
}
