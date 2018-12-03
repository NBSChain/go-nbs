package nat

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"net"
	"strconv"
	"time"
)

type taskProcess func(*MsgTask) (*net_pb.NatMsg, *net.UDPAddr, error)

type MsgTask struct {
	fromAddr *net.UDPAddr
	request  *net_pb.NatMsg
}

func (nat *Manager) checkWhoIsHe(task *MsgTask) (*net_pb.NatMsg, *net.UDPAddr, error) {
	peerAddr := task.fromAddr
	request := task.request.BootReg

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

	return pbRes, peerAddr, nil
}

func (nat *Manager) updateKATime(task *MsgTask) (*net_pb.NatMsg, *net.UDPAddr, error) {

	req := task.request.KeepAlive
	peerAddr := task.fromAddr
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

	return res, peerAddr, nil
}

func (nat *Manager) forwardInvite(task *MsgTask) (*net_pb.NatMsg, *net.UDPAddr, error) {
	invite := task.request.ReverseInvite

	nat.cacheLock.Lock()
	defer nat.cacheLock.Unlock()

	item, ok := nat.cache[invite.PeerId]
	if !ok {
		return nil, nil, NotFundErr
	}

	logger.Debug("Step3: forward notification to applier:", item.pubAddr)

	return task.request, item.pubAddr, nil
}

func (nat *Manager) forwardDigApply(task *MsgTask) (*net_pb.NatMsg, *net.UDPAddr, error) {
	req := task.request.DigApply
	peerAddr := task.fromAddr

	nat.cacheLock.Lock()
	defer nat.cacheLock.Unlock()

	toItem, ok := nat.cache[req.TargetId]
	if !ok {
		return nil, nil, NotFundErr
	}

	req.Public = peerAddr.String()

	logger.Debug("hole punch step2-2 forward dig out request:->", task.request.DigApply.Public)

	return task.request, toItem.pubAddr, nil
}

func (nat *Manager) forwardDigConfirm(task *MsgTask) (*net_pb.NatMsg, *net.UDPAddr, error) {
	ack := task.request.DigConfirm
	addr := task.fromAddr

	nat.cacheLock.Lock()
	defer nat.cacheLock.Unlock()

	item, ok := nat.cache[ack.TargetId]
	if !ok {
		return nil, nil, NotFundErr
	}

	ack.Public = addr.String()

	logger.Debug("hole punch step2-6 forward dig out notification:->", task.request.DigConfirm.Public)

	return task.request, item.pubAddr, nil
}
