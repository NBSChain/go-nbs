package nat

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"net"
	"sync"
	"time"
)

var (
	NotFundErr      = fmt.Errorf("no such node behind nat device")
	HandlerNotFound = fmt.Errorf("no taskhandler for this msg type")
	logger          = utils.GetLogInstance()
)

const MsgPoolSize = 1 << 10

type HostBehindNat struct {
	updateTIme time.Time
	pubAddr    *net.UDPAddr
	priAddr    string
}

type Server struct {
	sysNatServer *net.UDPConn
	networkId    string
	CanServe     chan bool
	cacheLock    sync.Mutex
	cache        map[string]*HostBehindNat
	taskQueue    chan *natTask
	taskRouter   map[int]taskProcess
}

func NewNatServer(networkId string) *Server {
	natObj := &Server{
		networkId:  networkId,
		CanServe:   make(chan bool),
		cache:      make(map[string]*HostBehindNat),
		taskQueue:  make(chan *natTask, MsgPoolSize),
		taskRouter: make(map[int]taskProcess),
	}

	natObj.initService()
	go natObj.receiveMsg()
	go natObj.timer()
	go natObj.runLoop()

	return natObj
}

//TODO:: support ipv6 later.
func (nat *Server) initService() {

	natServer, err := net.ListenUDP("udp4", &net.UDPAddr{
		Port: utils.GetConfig().NatServerPort,
	})

	if err != nil {
		logger.Panic("can't start nat sysNatServer.", err)
	}

	nat.sysNatServer = natServer
	nat.taskRouter[int(nbsnet.NatBootReg)] = nat.checkWhoIsHe
	nat.taskRouter[int(nbsnet.NatKeepAlive)] = nat.updateKATime
	nat.taskRouter[int(nbsnet.NatReversInvite)] = nat.forwardInvite
	nat.taskRouter[int(nbsnet.NatDigApply)] = nat.forwardDigApply
	nat.taskRouter[int(nbsnet.NatDigConfirm)] = nat.forwardDigConfirm
	nat.taskRouter[int(nbsnet.NatPingPong)] = nat.pong
	nat.taskRouter[DrainOutOldKa] = nat.checkKaTunnel
}

func (nat *Server) receiveMsg() {

	logger.Info(">>>>>>Nat sysNatServer start to listen......")

	for {
		data := make([]byte, utils.NormalReadBuffer)

		n, peerAddr, err := nat.sysNatServer.ReadFromUDP(data)
		if err != nil {
			logger.Warning("nat sysNatServer read udp data failed:", err)
			continue
		}

		request := &net_pb.NatMsg{}
		if err := proto.Unmarshal(data[:n], request); err != nil {
			logger.Warning("can't parse the nat message", err, peerAddr)
			continue
		}

		logger.Debug("message:", request, peerAddr)

		task := &natTask{
			taskType: int(request.Typ),
		}

		task.message = request
		task.addr = peerAddr
		nat.taskQueue <- task
	}
}

func (nat *Server) runLoop() {

	for {
		select {
		case task := <-nat.taskQueue:
			handler, ok := nat.taskRouter[task.taskType]
			if !ok {
				logger.Warning(HandlerNotFound)
				continue
			}
			if err := handler(task); err != nil {
				logger.Warning("nat message process err :->", err)
			}
		}
	}
}

func (nat *Server) timer() {

	for {
		select {
		case <-time.After(KeepAliveTime):
			nat.taskQueue <- &natTask{
				taskType: DrainOutOldKa,
			}
		}
	}
}
