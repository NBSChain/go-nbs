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

const MsgPoolSize = 1 << 12

var (
	MsgConvertErr = fmt.Errorf("convert to message failed")
	NotFundErr    = fmt.Errorf("no such node behind nat device")
	logger        = utils.GetLogInstance()
)

type HostBehindNat struct {
	UpdateTIme time.Time
	PubAddr    *net.UDPAddr
	PriAddr    string
}

type Task struct {
	TypId  net_pb.MsgType
	Param  interface{}
	Result chan interface{}
}

type MsgProcess func(param *Task) error

type Server struct {
	sysNatServer *net.UDPConn
	networkId    string
	CanServe     chan bool
	cacheLock    sync.Mutex
	cache        map[string]*HostBehindNat
	msgQueue     chan *Task
	router       map[net_pb.MsgType]MsgProcess
}

func NewNatServer(networkId string) *Server {
	natObj := &Server{
		networkId: networkId,
		CanServe:  make(chan bool),
		cache:     make(map[string]*HostBehindNat),
		msgQueue:  make(chan *Task, MsgPoolSize),
		router:    make(map[net_pb.MsgType]MsgProcess),
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
	nat.router[nbsnet.NatBootReg] = nat.checkWhoIsHe
	nat.router[nbsnet.NatKeepAlive] = nat.updateKATime
	nat.router[nbsnet.NatReversInvite] = nat.forwardInvite
	nat.router[nbsnet.NatDigApply] = nat.forwardDigApply
	nat.router[nbsnet.NatDigConfirm] = nat.forwardDigConfirm
	nat.router[nbsnet.NatPingPong] = nat.pong
	nat.router[nbsnet.DrainOutOldKa] = nat.checkKaTunnel
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

		task := &Task{
			TypId: request.Typ,
			Param: &msgTask{
				message: request,
				addr:    peerAddr,
			},
		}

		nat.msgQueue <- task
	}
}

func (nat *Server) runLoop() {

	for {
		select {

		case task := <-nat.msgQueue:

			p, ok := nat.router[task.TypId]
			if !ok {
				logger.Warning("unknown task type:->", task.TypId)
				continue
			}

			if err := p(task); err != nil {
				logger.Warning("process task err:->", err)
				continue
			}
		}
	}
}

func (nat *Server) timer() {

	for {
		select {
		case <-time.After(KeepAliveTime):
			task := &Task{
				TypId: nbsnet.DrainOutOldKa,
			}
			nat.msgQueue <- task
		}
	}
}
