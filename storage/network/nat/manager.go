package nat

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/gogo/protobuf/proto"
	"net"
	"strconv"
	"sync"
	"time"
)

var (
	NotFundErr = fmt.Errorf("no such node behind nat device")
	logger     = utils.GetLogInstance()
)

type HostBehindNat struct {
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
	cache        map[string]*HostBehindNat
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

		switch request.Typ {
		case nbsnet.NatBootReg:
			if err = nat.checkWhoIsHe(request.PayLoad, peerAddr); err != nil {
				logger.Error(err)
			}
		case nbsnet.NatPingPong:
			if err = nat.pong(request.PayLoad, peerAddr); err != nil {
				logger.Error(err)
			}
		case nbsnet.NatConnect:
			if err = nat.forwardDigRequest(request.PayLoad, peerAddr); err != nil {
				logger.Error(err)
			}
		case nbsnet.NatKeepAlive:
			if err = nat.updateKATime(request.PayLoad, peerAddr); err != nil {
				logger.Error(err)
			}
		case nbsnet.NatReversDig:
			if err = nat.forwardInvite(request.PayLoad, peerAddr); err != nil {
				logger.Error(err)
			}
		case nbsnet.NatConnectACK:
			if err = nat.forwardConnAck(request.PayLoad, peerAddr); err != nil {
				logger.Error(err)
			}
		}

		if err != nil {
			logger.Warning("nat server listening err:->", err)
		}
	}
}

func (nat *Manager) readNatRequest() (*net.UDPAddr, *net_pb.NatMsg, error) {

	data := make([]byte, utils.NormalReadBuffer)

	n, peerAddr, err := nat.sysNatServer.ReadFromUDP(data)
	if err != nil {
		logger.Warning("nat sysNatServer read udp data failed:", err)
		return nil, nil, err
	}

	request := &net_pb.NatMsg{}
	if err := proto.Unmarshal(data[:n], request); err != nil {
		logger.Warning("can't parse the nat request", err, peerAddr)
		return nil, nil, err
	}

	logger.Debug("request:", request, peerAddr)

	return peerAddr, request, nil
}

func (nat *Manager) checkWhoIsHe(data []byte, peerAddr *net.UDPAddr) error {

	request := &net_pb.BootReg{}
	if err := proto.Unmarshal(data, request); err != nil {
		return err
	}

	response := &net_pb.BootRegAck{}
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

	resData, _ := proto.Marshal(response)
	pbRes := &net_pb.NatMsg{
		Typ:     nbsnet.NatBootReg,
		Len:     int32(len(resData)),
		PayLoad: resData,
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
		select {
		case <-time.After(KeepAliveTime):
			nat.cacheLock.Lock()

			currentClock := time.Now()
			for nodeId, item := range nat.cache {

				if currentClock.Sub(item.updateTIme) > KeepAliveTimeOut {
					delete(nat.cache, nodeId)
				}
			}

			nat.cacheLock.Unlock()
		case <-nat.NatKATun.natChanged:
			logger.Warning("nat server ip changed.")
		}
	}
}

func (nat *Manager) updateKATime(data []byte, peerAddr *net.UDPAddr) error {
	req := &net_pb.NatKeepAlive{}
	if err := proto.Unmarshal(data, req); err != nil {
		return err
	}

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

	KeepAlive := &net_pb.NatKeepAlive{
		NodeId:  req.NodeId,
		PubIP:   peerAddr.IP.String(),
		PubPort: int32(peerAddr.Port),
	}

	kaData, _ := proto.Marshal(KeepAlive)
	res := &net_pb.NatMsg{
		Typ:     nbsnet.NatKeepAlive,
		Len:     int32(len(kaData)),
		PayLoad: kaData,
	}

	rawData, _ := proto.Marshal(res)
	if _, err := nat.sysNatServer.WriteToUDP(rawData, peerAddr); err != nil {
		return err
	}

	return nil
}

func (nat *Manager) forwardInvite(data []byte, peerAddr *net.UDPAddr) error {
	invite := &net_pb.ReverseInvite{}
	if err := proto.Unmarshal(data, invite); err != nil {
		return err
	}

	nat.cacheLock.Lock()
	defer nat.cacheLock.Unlock()

	item, ok := nat.cache[invite.PeerId]
	if !ok {
		return NotFundErr
	}

	res := &net_pb.NatMsg{
		Typ:     nbsnet.NatReversDig,
		Len:     int32(len(data)),
		PayLoad: data,
	}

	rawData, _ := proto.Marshal(res)

	if _, err := nat.sysNatServer.WriteTo(rawData, item.pubAddr); err != nil {
		return err
	}

	logger.Debug("Step3: forward notification to applier:", item.pubAddr)

	return nil
}

func (nat *Manager) forwardConnAck(bytes []byte, addr *net.UDPAddr) error {
	ack := &net_pb.NatConnectAck{}
	if err := proto.Unmarshal(bytes, ack); err != nil {
		return err
	}

	nat.cacheLock.Lock()
	defer nat.cacheLock.Unlock()

	item, ok := nat.cache[ack.TargetId]
	if !ok {
		return NotFundErr
	}

	ack.Public = addr.String()
	data, _, err := nbsnet.PackNatData(ack, nbsnet.NatConnectACK)
	if err != nil {
		return err
	}

	if _, err := nat.sysNatServer.WriteToUDP(data, item.pubAddr); err != nil {
		return err
	}

	return nil
}
