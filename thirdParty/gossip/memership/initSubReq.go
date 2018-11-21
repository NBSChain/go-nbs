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

		conn.SetDeadline(time.Now().Add(time.Second * 3))

		if err := node.sendInitSubRequest(conn); err != nil {
			logger.Error("send init sub request failed:", err)
			goto CloseConn
		}

		if err := node.readInitSubResponse(conn); err == nil {
			logger.Debug("find contract server success.")
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

func (node *MemManager) sendInitSubRequest(conn *nbsnet.NbsUdpConn) error {

	addr := network.GetInstance().GetAddress()
	msg := &pb.Gossip{
		MessageType: pb.MsgType_init,
		InitMsg: &pb.InitSub{
			NodeId:      addr.NetworkId,
			PublicIp:    addr.PubIp,
			PrivateIp:   addr.PriIp,
			CanBeServer: addr.CanServe,
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

func (node *MemManager) readInitSubResponse(conn *nbsnet.NbsUdpConn) error {

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
func (node *MemManager) initSubReqHandle(request *pb.InitSub, applierAddr *net.UDPAddr) {

	payLoad := &pb.InitSubACK{
		ContactId: node.peerId,
		ApplierId: request.NodeId,
	}

	message := &pb.Gossip{
		MessageType: pb.MsgType_initACK,
		InitACK:     payLoad,
	}

	msgData, err := proto.Marshal(message)
	if err != nil {
		logger.Warning("marshal init sub ack msg failed:", err)
		return
	}

	if _, err := node.serviceConn.WriteToUDP(msgData, applierAddr); err != nil {
		logger.Warning("failed to send init ack msg:", err)
		return
	}

	node.taskSignal <- innerTask{
		taskType: ProxyInitSubRequest,
		taskData: request,
	}
}
