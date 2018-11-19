package nat

import (
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/storage/network/shareport"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"net"
	"strconv"
	"time"
)

const (
	KeepAlive        = 15
	KeepAliveTimeOut = 45
)

type ConnType int8

const (
	_ ConnType = iota
	ConnTypeNormal
	ConnTypeNat
	ConnTypeNatInverse
)

type ConnTask struct {
	Err       error
	sessionId string
	CType     ConnType
	ProxyAddr *net.UDPAddr
	ConnCh    chan *net.UDPConn
}

type KATunnel struct {
	closed     chan bool
	receiveHub *net.UDPConn
	kaConn     *net.UDPConn
	updateTime time.Time
	natTask    map[string]*ConnTask
}

func (nat *Manager) newKATunnel() (*KATunnel, error) {

	port := strconv.Itoa(utils.GetConfig().NatClientPort)
	listener, err := shareport.ListenUDP("udp4", port)
	if err != nil {
		logger.Warning("create share listening udp failed.")
		return nil, err
	}

	client, err := shareport.DialUDP("udp4", "0.0.0.0:"+port, nat.dNatServer)
	if err != nil {
		logger.Warning("create share port dial udp connection failed.")
		return nil, err
	}

	localAddr := client.LocalAddr().String()
	priIP, priPort, err := net.SplitHostPort(localAddr)
	nat.SelfAddr.privatePort = priPort
	nat.SelfAddr.privateIP = priIP

	tunnel := &KATunnel{
		closed:     make(chan bool),
		receiveHub: listener,
		kaConn:     client,
		updateTime: time.Now(),
		natTask:    make(map[string]*ConnTask),
	}

	request := &net_pb.NatRequest{
		MsgType: net_pb.NatMsgType_BootStrapReg,
		BootRegReq: &net_pb.BootNatRegReq{
			NodeId:      nat.networkId,
			PrivateIp:   priIP,
			PrivatePort: priPort,
		},
	}

	if err := nat.registerPriHost(request); err != nil {
		return nil, err
	}

	return tunnel, nil
}

func (nat *Manager) registerPriHost(request *net_pb.NatRequest) error {

	reqData, err := proto.Marshal(request)
	if err != nil {
		logger.Warning("failed to marshal nat keep alive message", err)
		return err
	}

	if no, err := nat.NatKATun.kaConn.Write(reqData); err != nil || no == 0 {
		logger.Warning("nat channel keep alive message failed", err, no)
		return err
	}

	if err := nat.readRegResponse(); err != nil {
		logger.Warning("failed to read channel initialize response.")
		return err
	}

	return nil
}

func (tunnel *KATunnel) Close() {

	tunnel.receiveHub.Close()
	tunnel.kaConn.Close()
	tunnel.closed <- true
}

func (tunnel *KATunnel) runLoop() {

	for {
		select {
		case <-time.After(time.Second * KeepAlive):
			tunnel.sendKeepAlive()
		case <-tunnel.closed:
			logger.Info("keep alive channel is closed.")
			return
		}
	}
}

func (tunnel *KATunnel) sendKeepAlive() error {

	request := &net_pb.NatRequest{
		MsgType: net_pb.NatMsgType_KeepAlive,
		KeepAlive: &net_pb.NatKeepAlive{
			NodeId: tunnel.networkId,
		},
	}

	requestData, err := proto.Marshal(request)
	if err != nil {
		logger.Warning("failed to marshal keep alive request:", err)
		return err
	}

	logger.Debug("keep alive channel start")

	if no, err := tunnel.kaConn.Write(requestData); err != nil || no == 0 {
		logger.Warning("failed to send keep alive channel message:", err)
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

		time.Sleep(time.Second * KeepAlive)
	}
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
