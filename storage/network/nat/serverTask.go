package nat

import (
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"net"
	"time"
)

type msgTask struct {
	addr    *net.UDPAddr
	message *net_pb.NatMsg
}

func (nat *Server) checkWhoIsHe(t *nbsnet.Task) error {
	task, ok := t.Param.(*msgTask)
	if !ok {
		return nbsnet.MsgConvertErr
	}

	peerAddr := task.addr
	request := task.message.BootReg

	response := &net_pb.BootAnswer{}
	response.PublicIp = peerAddr.IP.String()
	response.PublicPort = int32(peerAddr.Port)

	if peerAddr.IP.Equal(net.ParseIP(request.PrivateIp)) {
		response.NatType = net_pb.NatType_NoNatDevice
	} else if peerAddr.Port == int(request.PrivatePort) {
		if nat.networkId == request.NodeId {
			logger.Info("yeah, I'm the bootstrap node:->")
			response.NatType = net_pb.NatType_CanBeNatServer
		} else {
			response.NatType = net_pb.NatType_ToBeChecked
			go nat.ping(peerAddr, request.NodeId)
		}
	} else {
		response.NatType = net_pb.NatType_BehindNat
	}

	pbRes := &net_pb.NatMsg{
		Typ:        nbsnet.NatBootAnswer,
		BootAnswer: response,
	}
	data, _ := proto.Marshal(pbRes)
	if _, err := nat.sysNatServer.WriteToUDP(data, peerAddr); err != nil {
		return err
	}

	if response.NatType == net_pb.NatType_BehindNat ||
		response.NatType == net_pb.NatType_ToBeChecked {
		priAddr := nbsnet.JoinHostPort(request.PrivateIp, request.PrivatePort)
		item := &HostBehindNat{
			updateTIme: time.Now(),
			pubAddr:    peerAddr,
			priAddr:    priAddr,
		}
		nat.cache[request.NodeId] = item
	}

	return nil
}

func (nat *Server) updateKATime(t *nbsnet.Task) error {
	task, ok := t.Param.(*msgTask)
	if !ok {
		return nbsnet.MsgConvertErr
	}

	req := task.message.KeepAlive
	peerAddr := task.addr
	nodeId := req.NodeId

	nat.cacheLock.Lock()
	defer nat.cacheLock.Unlock()

	if item, ok := nat.cache[nodeId]; ok {
		item.updateTIme = time.Now()
		item.pubAddr = peerAddr
		item.priAddr = req.LAddr
	} else {
		item := &HostBehindNat{
			updateTIme: time.Now(),
			pubAddr:    peerAddr,
			priAddr:    req.LAddr,
		}
		nat.cache[nodeId] = item
	}

	res := &net_pb.NatMsg{
		Typ: nbsnet.NatKeepAlive,
		KeepAlive: &net_pb.KeepAlive{
			NodeId:  req.NodeId,
			PubIP:   peerAddr.IP.String(),
			PubPort: int32(peerAddr.Port),
		},
	}
	data, _ := proto.Marshal(res)
	if _, err := nat.sysNatServer.WriteToUDP(data, peerAddr); err != nil {
		return err
	}

	return nil
}

func (nat *Server) forwardInvite(t *nbsnet.Task) error {
	task, ok := t.Param.(*msgTask)
	if !ok {
		return nbsnet.MsgConvertErr
	}

	invite := task.message.ReverseInvite

	nat.cacheLock.Lock()
	defer nat.cacheLock.Unlock()

	item, ok := nat.cache[invite.PeerId]
	if !ok {
		return NotFundErr
	}

	logger.Debug("Step3: forward notification to applier:", item.pubAddr)
	data, _ := proto.Marshal(task.message)
	if _, err := nat.sysNatServer.WriteToUDP(data, item.pubAddr); err != nil {
		return err
	}
	return nil
}

func (nat *Server) forwardDigApply(t *nbsnet.Task) error {
	task, ok := t.Param.(*msgTask)
	if !ok {
		return nbsnet.MsgConvertErr
	}

	req := task.message.DigApply
	peerAddr := task.addr

	nat.cacheLock.Lock()
	defer nat.cacheLock.Unlock()

	toItem, ok := nat.cache[req.TargetId]
	if !ok {
		return NotFundErr
	}

	req.Public = peerAddr.String()

	logger.Info("hole punch step2-2 forward dig out message:->", task.message.DigApply.Public)
	data, _ := proto.Marshal(task.message)
	if _, err := nat.sysNatServer.WriteToUDP(data, toItem.pubAddr); err != nil {
		return err
	}
	return nil
}

func (nat *Server) forwardDigConfirm(t *nbsnet.Task) error {
	task, ok := t.Param.(*msgTask)
	if !ok {
		return nbsnet.MsgConvertErr
	}

	ack := task.message.DigConfirm
	addr := task.addr

	nat.cacheLock.Lock()
	defer nat.cacheLock.Unlock()

	item, ok := nat.cache[ack.TargetId]
	if !ok {
		return NotFundErr
	}

	ack.Public = addr.String()

	logger.Info("hole punch step2-6 forward dig out notification:->", task.message.DigConfirm.Public)
	data, _ := proto.Marshal(task.message)
	if _, err := nat.sysNatServer.WriteToUDP(data, item.pubAddr); err != nil {
		return err
	}
	return nil
}

func (nat *Server) pong(t *nbsnet.Task) error {
	nat.CanServe <- true
	logger.Debug("I can serve as in public network.")
	return nil
}

func (nat *Server) checkKaTunnel(t *nbsnet.Task) error {

	nat.cacheLock.Lock()
	defer nat.cacheLock.Unlock()

	currentClock := time.Now()
	for nodeId, item := range nat.cache {

		if currentClock.Sub(item.updateTIme) > KeepAliveTimeOut {
			logger.Warning("the nat item has been removed:->", nodeId)
			delete(nat.cache, nodeId)
		}
	}
	return nil
}

func (nat *Server) ping(peerAddr *net.UDPAddr, reqNodeId string) {

	conn, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   peerAddr.IP,
		Port: utils.GetConfig().NatServerPort,
	})

	if err != nil {
		logger.Warning("ping DialUDP err:->", err)
		return
	}
	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(BootStrapTimeOut / 2)); err != nil {
		logger.Warning("ping SetDeadline err:->", err)
		return
	}

	request := &net_pb.NatMsg{
		Typ: nbsnet.NatPingPong,
		PingPong: &net_pb.PingPong{
			Ping:  nat.networkId,
			Pong:  reqNodeId,
			Nonce: "", //TODO::security nonce
			TTL:   1,  //time to live
		},
	}
	reqData, _ := proto.Marshal(request)
	if _, err := conn.Write(reqData); err != nil {
		logger.Warning("ping Write err:->", err)
		return
	}
}
