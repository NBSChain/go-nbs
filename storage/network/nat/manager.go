package nat

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/gogo/protobuf/proto"
	"net"
	"strconv"
	"sync"
	"time"
)

var logger = utils.GetLogInstance()

const NetIoBufferSize = 1 << 11
const BootStrapNatServerTimeOutInSec = 4

type ClientItem struct {
	nodeId     string
	pubIp      string
	pubPort    string
	priIp      string
	priPort    string
	canServer  bool
	kaAddr     *net.UDPAddr
	updateTIme time.Time
}

type Manager struct {
	cacheLock     sync.Mutex
	selfNatServer *net.UDPConn
	networkId     string
	canServe      chan bool
	NatKATun      *KATunnel
	cache         map[string]*ClientItem
	dNatServer    *DecentralizedNatServer
}

//TODO:: support ipv6 later.
func (nat *Manager) startNatService() {

	natServer, err := net.ListenUDP("udp4", &net.UDPAddr{
		Port: utils.GetConfig().NatServerPort,
	})

	if err != nil {
		logger.Panic("can't start nat selfNatServer.", err)
	}

	nat.selfNatServer = natServer
}

func (nat *Manager) natServiceListening() {

	logger.Info(">>>>>>Nat selfNatServer start to listen......")

	for {
		peerAddr, request, err := nat.readNatRequest()
		if err != nil {
			logger.Error(err)
			continue
		}

		switch request.MsgType {
		case net_pb.NatMsgType_BootStrapReg:
			if err = nat.checkWhoIsHe(request.BootRegReq, peerAddr); err != nil {
				logger.Error(err)
			}
		case net_pb.NatMsgType_Ping:
			if err = nat.pong(request.Ping, peerAddr); err != nil {
				logger.Error(err)
			}
		case net_pb.NatMsgType_Connect:
			if err = nat.invitePeers(request.ConnReq, peerAddr); err != nil {
				logger.Error(err)
			}
		case net_pb.NatMsgType_KeepAlive:
			if err = nat.updateKATime(request.KeepAlive, peerAddr); err != nil {
				logger.Error(err)
			}
		}

		if err != nil {
			nat.responseAnError(err, peerAddr)
		}
	}
}

func (nat *Manager) responseAnError(err error, peerAddr *net.UDPAddr) {
	response := &net_pb.Response{
		MsgType: net_pb.NatMsgType_error,
		Error: &net_pb.ErrorNotify{
			ErrMsg: err.Error(),
		},
	}

	resData, _ := proto.Marshal(response)
	nat.selfNatServer.WriteToUDP(resData, peerAddr)
}

func (nat *Manager) readNatRequest() (*net.UDPAddr, *net_pb.NatRequest, error) {

	data := make([]byte, NetIoBufferSize)

	n, peerAddr, err := nat.selfNatServer.ReadFromUDP(data)
	if err != nil {
		logger.Warning("nat selfNatServer read udp data failed:", err)
		return nil, nil, err
	}

	request := &net_pb.NatRequest{}
	if err := proto.Unmarshal(data[:n], request); err != nil {
		logger.Warning("can't parse the nat request", err, peerAddr)
		return nil, nil, err
	}

	logger.Debug("request:", request)

	return peerAddr, request, nil
}

func (nat *Manager) checkWhoIsHe(request *net_pb.BootNatRegReq, peerAddr *net.UDPAddr) error {

	response := &net_pb.BootNatRegRes{}
	response.PublicIp = peerAddr.IP.String()
	response.PublicPort = fmt.Sprintf("%d", peerAddr.Port)
	if peerAddr.IP.Equal(net.ParseIP(request.PrivateIp)) {

		response.NatType = net_pb.NatType_NoNatDevice
	} else if strconv.Itoa(peerAddr.Port) == request.PrivatePort {

		response.NatType = net_pb.NatType_ToBeChecked
		go nat.ping(peerAddr)

	} else {

		response.NatType = net_pb.NatType_BehindNat
	}

	pbRes := &net_pb.Response{
		MsgType:    net_pb.NatMsgType_BootStrapReg,
		BootRegRes: response,
	}

	pbResData, err := proto.Marshal(pbRes)
	if err != nil {
		logger.Warning("failed to marshal nat response data", err)
		return err
	}

	if _, err := nat.selfNatServer.WriteToUDP(pbResData, peerAddr); err != nil {
		logger.Warning("failed to send nat response", err)
		return err
	}

	item := &ClientItem{
		nodeId:     request.NodeId,
		pubIp:      response.PublicIp,
		pubPort:    response.PublicPort,
		priIp:      request.PrivateIp,
		priPort:    request.PrivatePort,
		canServer:  IsPublic(response.NatType),
		updateTIme: time.Now(),
	}

	nat.cache[request.NodeId] = item

	return nil
}

func (nat *Manager) cacheManager() {

	for {
		nat.cacheLock.Lock()

		currentClock := time.Now()
		for nodeId, item := range nat.cache {
			if item.canServer {
				continue
			}

			if currentClock.Sub(item.updateTIme) > time.Second*KeepAliveTimeOut {
				delete(nat.cache, nodeId)
			}
		}

		nat.cacheLock.Unlock()

		time.Sleep(time.Second * KeepAlive)
	}
}

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
		CanServe:    item.canServer,
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

func (nat *Manager) updateKATime(req *net_pb.NatKeepAlive, peerAddr *net.UDPAddr) error {
	nodeId := req.NodeId
	nat.cacheLock.Lock()
	defer nat.cacheLock.Unlock()

	if item, ok := nat.cache[nodeId]; !ok {
		item.updateTIme = time.Now()
		item.kaAddr = peerAddr
	}

	return nil
}
