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

type hostBehindNat struct {
	updateTIme time.Time
	pubAddr    *net.UDPAddr
	priAddr    string
}

type Manager struct {
	cacheLock    sync.Mutex
	sysNatServer *net.UDPConn
	networkId    string
	canServe     chan bool
	NatKATun     *KATunnel
	cache        map[string]*hostBehindNat
}

//TODO:: support ipv6 later.
func (nat *Manager) startNatService() {

	natServer, err := net.ListenUDP("udp4", &net.UDPAddr{
		Port: utils.GetConfig().NatServerPort,
	})

	if err != nil {
		logger.Panic("can't start nat sysNatServer.", err)
	}

	nat.sysNatServer = natServer
}

func (nat *Manager) natServiceListening() {

	logger.Info(">>>>>>Nat sysNatServer start to listen......")

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
			if err = nat.notifyConnInvite(request.ConnReq, peerAddr); err != nil {
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
	nat.sysNatServer.WriteToUDP(resData, peerAddr)
}

func (nat *Manager) readNatRequest() (*net.UDPAddr, *net_pb.NatRequest, error) {

	data := make([]byte, utils.NormalReadBuffer)

	n, peerAddr, err := nat.sysNatServer.ReadFromUDP(data)
	if err != nil {
		logger.Warning("nat sysNatServer read udp data failed:", err)
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

	if _, err := nat.sysNatServer.WriteToUDP(pbResData, peerAddr); err != nil {
		logger.Warning("failed to send nat response", err)
		return err
	}

	return nil
}

/************************************************************************
*
*			server side
*
*************************************************************************/
func (nat *Manager) cacheManager() {

	for {
		nat.cacheLock.Lock()

		currentClock := time.Now()
		for nodeId, item := range nat.cache {

			if currentClock.Sub(item.updateTIme) > time.Second*KeepAliveTimeOut {
				delete(nat.cache, nodeId)
			}
		}

		nat.cacheLock.Unlock()

		time.Sleep(KeepAliveTime)
	}
}

func (nat *Manager) updateKATime(req *net_pb.NatKeepAlive, peerAddr *net.UDPAddr) error {

	nodeId := req.NodeId

	nat.cacheLock.Lock()
	defer nat.cacheLock.Unlock()

	if item, ok := nat.cache[nodeId]; ok {
		item.updateTIme = time.Now()
		item.pubAddr = peerAddr
		item.priAddr = req.LAddr
	} else {
		item := &hostBehindNat{
			updateTIme: time.Now(),
			pubAddr:    peerAddr,
			priAddr:    req.LAddr,
		}
		nat.cache[nodeId] = item
	}

	return nil
}
