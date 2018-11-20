package nat

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/storage/network/shareport"
	"github.com/golang/protobuf/proto"
	"net"
)

//TIPS::get peer's addr info and make a connection.
func (tunnel *KATunnel) punchAHole(response *net_pb.NatConRes) {
	sessionId := response.SessionId
	task, ok := tunnel.natTask[sessionId]
	if !ok {
		logger.Error("can't find the nat connection task")
		return
	}

	task.ProxyConn = &net.UDPAddr{
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
func (nat *Manager) notifyConnInvite(req *net_pb.NatConReq, peerAddr *net.UDPAddr) error {
	nat.cacheLock.Lock()
	defer nat.cacheLock.Unlock()

	req.

	sessionId := req.SessionId
	fromItem, ok := nat.cache[req.FromPeerId]
	if !ok {
		return fmt.Errorf("the from peer id is not found")
	}

	toItem, ok := nat.cache[req.ToPeerId]
	if !ok {
		toItem = nat.sysNatServer.SendConnInvite(fromItem, req.ToPeerId, sessionId, req.ToPort, false)
	} else {
		if err := nat.sendConnInvite(fromItem, toItem.pubAddr, sessionId, req.ToPort, false); err != nil {
			logger.Error("connect invite failed:", err)
		}
	}

	return nat.sendConnInvite(toItem, peerAddr, sessionId, req.ToPort, true)
}
