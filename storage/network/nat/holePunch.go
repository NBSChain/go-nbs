package nat

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/storage/network/shareport"
	"github.com/golang/protobuf/proto"
	"net"
)

func (tunnel *KATunnel) MakeANatConn(fromId, toId, connId string, port int) *ConnTask {

	sessionId := fromId + toId
	task := &ConnTask{
		sessionId: sessionId,
	}

	payload := &net_pb.NatConReq{
		FromPeerId: fromId,
		ToPeerId:   toId,
		ToPort:     int32(port),
		SessionId:  sessionId,
	}
	request := &net_pb.NatRequest{
		MsgType: net_pb.NatMsgType_Connect,
		ConnReq: payload,
	}

	reqData, err := proto.Marshal(request)
	if err != nil {
		logger.Error("failed to marshal the nat connect request", err)
		task.Err = err
		return task
	}

	if no, err := tunnel.kaConn.Write(reqData); err != nil || no == 0 {
		logger.Warning("nat channel keep alive message failed", err, no)
		task.Err = err
		return task
	}

	task.ConnCh = make(chan *net.UDPConn)

	tunnel.natTask[sessionId] = task

	return task
}

//TIPS::get peer's addr info and make a connection.
func (tunnel *KATunnel) punchAHole(response *net_pb.NatConRes) {
	sessionId := response.SessionId
	task, ok := tunnel.natTask[sessionId]
	if !ok {
		logger.Error("can't find the nat connection task")
		return
	}

	task.ProxyAddr = &net.UDPAddr{
		IP:   net.ParseIP(tunnel.privateIP),
		Port: int(response.TargetPort),
	}
	if response.IsCaller {
		task.CType = ConnTypeNat
	} else {
		task.CType = ConnTypeNatInverse
	}

	priConn, priErr := shareport.DialUDP("udp4", tunnel.privateIP+":"+tunnel.privatePort,
		response.PrivateIp+":"+response.PrivatePort)

	if priErr != nil {
		logger.Warning("failed to make a nat connection from private network while peer is behind nat")
	} else {
		task.ConnCh <- priConn
		return
	}

	pubConn, pubErr := shareport.DialUDP("udp4", tunnel.privateIP+":"+tunnel.privatePort,
		response.PublicIp+":"+response.PublicPort)
	if pubErr != nil {
		logger.Error("failed to make a nat connection from public network while peer is behind nat")
		err := fmt.Errorf("failed to make a nat connection while peer is behind nat.%s-%s", pubErr, priErr)
		task.Err = err
		task.ConnCh <- nil
		return
	}

	task.ConnCh <- pubConn

	delete(tunnel.natTask, task.sessionId)
}

/************************************************************************
*
*			server side
*
*************************************************************************/
//TODO::Find peers from nat gossip protocol
func (nat *Manager) invitePeers(req *net_pb.NatConReq, peerAddr *net.UDPAddr) error {
	nat.cacheLock.Lock()
	defer nat.cacheLock.Unlock()

	sessionId := req.SessionId
	fromItem, ok := nat.cache[req.FromPeerId]
	if !ok {
		return fmt.Errorf("the from peer id is not found")
	}

	toItem, ok := nat.cache[req.ToPeerId]
	if !ok {
		toItem = nat.dNatServer.SendConnInvite(fromItem, req.ToPeerId, sessionId, req.ToPort, false)
	} else {
		if err := nat.sendConnInvite(fromItem, toItem.kaAddr, sessionId, req.ToPort, false); err != nil {
			logger.Error("connect invite failed:", err)
		}
	}

	return nat.sendConnInvite(toItem, peerAddr, sessionId, req.ToPort, true)
}

func (nat *Manager) sendConnInvite(item *ClientItem, addr *net.UDPAddr, sessionId string, toPort int32, isCaller bool) error {

	connRes := &net_pb.NatConRes{
		PeerId:      item.nodeId,
		PublicIp:    item.pubIp,
		PublicPort:  item.pubPort,
		PrivateIp:   item.priIp,
		PrivatePort: item.priPort,
		SessionId:   sessionId,
		TargetPort:  toPort,
		IsCaller:    isCaller,
	}

	response := &net_pb.Response{
		MsgType: net_pb.NatMsgType_Connect,
		ConnRes: connRes,
	}

	toItemData, err := proto.Marshal(response)
	if err != nil {
		return err
	}

	if _, err := nat.selfNatServer.WriteToUDP(toItemData, item.kaAddr); err != nil {
		return err
	}
	return nil
}
