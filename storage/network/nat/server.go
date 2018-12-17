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
	NotFundErr = fmt.Errorf("no such node behind nat device")
	logger     = utils.GetLogInstance()
)

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
}

func NewNatServer(networkId string) *Server {
	natObj := &Server{
		networkId: networkId,
		CanServe:  make(chan bool),
		cache:     make(map[string]*HostBehindNat),
	}

	natObj.initService()
	go natObj.receiveMsg()
	go natObj.timer()
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
	nbsnet.GetInstance().Register(nbsnet.NatBootReg, nat.checkWhoIsHe)
	nbsnet.GetInstance().Register(nbsnet.NatKeepAlive, nat.updateKATime)
	nbsnet.GetInstance().Register(nbsnet.NatReversInvite, nat.forwardInvite)
	nbsnet.GetInstance().Register(nbsnet.NatDigApply, nat.forwardDigApply)
	nbsnet.GetInstance().Register(nbsnet.NatDigConfirm, nat.forwardDigConfirm)
	nbsnet.GetInstance().Register(nbsnet.NatPingPong, nat.pong)
	nbsnet.GetInstance().Register(nbsnet.DrainOutOldKa, nat.checkKaTunnel)
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

		task := &nbsnet.Task{
			TypId: request.Typ,
			Param: &msgTask{
				message: request,
				addr:    peerAddr,
			},
		}

		nbsnet.GetInstance().Enqueue(task)
	}
}

func (nat *Server) timer() {

	for {
		select {
		case <-time.After(KeepAliveTime):
			task := &nbsnet.Task{
				TypId: nbsnet.DrainOutOldKa,
			}
			nbsnet.GetInstance().Enqueue(task)
		}
	}
}
