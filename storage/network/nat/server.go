package nat

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/gogo/protobuf/proto"
	"net"
	"sync"
	"time"
)

var (
	NotFundErr      = fmt.Errorf("no such node behind nat device")
	HandlerNotFound = fmt.Errorf("no taskhandler for this msg type")
	logger          = utils.GetLogInstance()
)

type HostBehindNat struct {
	updateTIme time.Time
	pubAddr    *net.UDPAddr
	priAddr    string
}

type Manager struct {
	sysNatServer *net.UDPConn
	networkId    string
	canServe     chan bool
	NatKATun     *KATunnel
	cacheLock    sync.Mutex
	cache        map[string]*HostBehindNat
	task         chan *natTask
	msgHandlers  map[int]taskProcess
}

//TODO:: support ipv6 later.
func (nat *Manager) initService() {

	natServer, err := net.ListenUDP("udp4", &net.UDPAddr{
		Port: utils.GetConfig().NatServerPort,
	})

	if err != nil {
		logger.Panic("can't start nat sysNatServer.", err)
	}

	nat.sysNatServer = natServer
	nat.msgHandlers[int(nbsnet.NatBootReg)] = nat.checkWhoIsHe
	nat.msgHandlers[int(nbsnet.NatKeepAlive)] = nat.updateKATime
	nat.msgHandlers[int(nbsnet.NatReversInvite)] = nat.forwardInvite
	nat.msgHandlers[int(nbsnet.NatDigApply)] = nat.forwardDigApply
	nat.msgHandlers[int(nbsnet.NatDigConfirm)] = nat.forwardDigConfirm
	nat.msgHandlers[int(nbsnet.NatPingPong)] = nat.pong
	nat.msgHandlers[DrainOutOldKa] = nat.checkKaTunnel
}

func (nat *Manager) MsgConsumer() {

}

func (nat *Manager) TaskReceiver() {

	logger.Info(">>>>>>Nat sysNatServer start to listen......")

	for {
		peerAddr, request, err := nat.readNatRequest()
		if err != nil {
			logger.Error(err)
			continue
		}

		task := &natTask{
			taskType: int(request.Typ),
		}

		task.message = request
		task.addr = peerAddr

		nat.task <- task
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
		logger.Warning("can't parse the nat message", err, peerAddr)
		return nil, nil, err
	}

	logger.Debug("message:", request, peerAddr)

	return peerAddr, request, nil
}

func (nat *Manager) RunLoop() {

	for {
		select {
		case task := <-nat.task:
			msgType := int(task.message.Typ)
			handler, ok := nat.msgHandlers[msgType]
			if !ok {
				logger.Warning(HandlerNotFound)
				continue
			}
			if err := handler(task); err != nil {
				logger.Warning("nat message proccess err :->", err)
			}
		}
	}
}

func (nat *Manager) timer() {

	for {
		select {
		case <-time.After(KeepAliveTime):
			nat.task <- &natTask{
				taskType: DrainOutOldKa,
			}
		}
	}
}
