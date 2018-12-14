package nat

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/golang/protobuf/proto"
	"net"
	"strconv"
	"time"
)

const (
	DrainOutOldKa = 1
)

type taskProcess func(*natTask) error

type msgTask struct {
	addr    *net.UDPAddr
	message *net_pb.NatMsg
}
type innerTask struct {
	params []interface{}
	result chan interface{}
}
type natTask struct {
	taskType int
	msgTask
	innerTask
}

func (nat *Manager) checkWhoIsHe(task *natTask) error {
	peerAddr := task.addr
	request := task.message.BootReg

	response := &net_pb.BootAnswer{}
	response.PublicIp = peerAddr.IP.String()
	response.PublicPort = fmt.Sprintf("%d", peerAddr.Port)

	if peerAddr.IP.Equal(net.ParseIP(request.PrivateIp)) {
		response.NatType = net_pb.NatType_NoNatDevice
	} else if strconv.Itoa(peerAddr.Port) == request.PrivatePort {
		response.NatType = net_pb.NatType_ToBeChecked
		go nat.ping(peerAddr, request.NodeId)
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

	return nil
}

func (nat *Manager) updateKATime(task *natTask) error {

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

func (nat *Manager) forwardInvite(task *natTask) error {
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

func (nat *Manager) forwardDigApply(task *natTask) error {
	req := task.message.DigApply
	peerAddr := task.addr

	nat.cacheLock.Lock()
	defer nat.cacheLock.Unlock()

	toItem, ok := nat.cache[req.TargetId]
	if !ok {
		return NotFundErr
	}

	req.Public = peerAddr.String()

	logger.Debug("hole punch step2-2 forward dig out message:->", task.message.DigApply.Public)
	data, _ := proto.Marshal(task.message)
	if _, err := nat.sysNatServer.WriteToUDP(data, toItem.pubAddr); err != nil {
		return err
	}
	return nil
}

func (nat *Manager) forwardDigConfirm(task *natTask) error {
	ack := task.message.DigConfirm
	addr := task.addr

	nat.cacheLock.Lock()
	defer nat.cacheLock.Unlock()

	item, ok := nat.cache[ack.TargetId]
	if !ok {
		return NotFundErr
	}

	ack.Public = addr.String()

	logger.Debug("hole punch step2-6 forward dig out notification:->", task.message.DigConfirm.Public)
	data, _ := proto.Marshal(task.message)
	if _, err := nat.sysNatServer.WriteToUDP(data, item.pubAddr); err != nil {
		return err
	}
	return nil
}

func (nat *Manager) pong(task *natTask) error {
	nat.canServe <- true
	logger.Debug("I can serve as in public network.")
	return nil
}

func (nat *Manager) checkKaTunnel(task *natTask) error {

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
