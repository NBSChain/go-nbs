package network

import (
	"context"
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/denat"
	"github.com/NBSChain/go-nbs/storage/network/nat"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/shareport"
	"github.com/NBSChain/go-nbs/thirdParty/account"
	"github.com/NBSChain/go-nbs/utils"
	"net"
	"sync"
	"time"
)

type nbsNetwork struct {
	networkId string
	ctx       context.Context
	CanServe  bool
	natServer *nat.Server
	natClient *nat.Client
	digTask   map[string]*connTask
	cmdRouter map[int]nat.CmdProcess
}

const (
	HolePunchTimeOut    = 4 * time.Second
	DigTryTimesOnNat    = 3
	ConnectionSeparator = "-"
)

var (
	CmdTaskErr = fmt.Errorf("convert cmmond task parameter err:->")
	once       sync.Once
	instance   *nbsNetwork
	logger     = utils.GetLogInstance()
)

/************************************************************************
*
*			public functions
*
*************************************************************************/
func GetInstance() Network {

	once.Do(func() {
		instance = newNetwork()
	})

	return instance
}

func newNetwork() *nbsNetwork {

	peerId := account.GetAccountInstance().GetPeerID()

	network := &nbsNetwork{
		ctx:       context.Background(),
		digTask:   make(map[string]*connTask),
		cmdRouter: make(map[int]nat.CmdProcess),
	}

	if peerId != "" {
		if err := network.StartUp(peerId); err != nil {
			panic(err)
		}
	} else {
		logger.Warning("no account right now, so the network is down")
	}

	return network
}

func (network *nbsNetwork) StartUp(peerId string) error {

	localPeers := nbsnet.ExternalIP()
	if len(localPeers) == 0 {
		logger.Panic("no available network")
	}

	logger.Debug("all network interfaces:->", localPeers)

	denat.GetDeNatSerIns().Setup(peerId)

	network.networkId = peerId

	network.natServer = nat.NewNatServer(network.networkId)

	if err := network.setupNatClient(peerId); err != nil {
		return err
	}

	return nil
}

func (network *nbsNetwork) GetNatAddr() (string, *nbsnet.NbsUdpAddr) {
	return network.networkId, network.natClient.NatAddr
}

func (network *nbsNetwork) DialUDP(nt string, localAddr, remoteAddr *net.UDPAddr) (*nbsnet.NbsUdpConn, error) {

	c, err := net.DialUDP(nt, localAddr, remoteAddr)
	if err != nil {
		return nil, err
	}
	Sid := c.LocalAddr().String() + ConnectionSeparator + remoteAddr.String()
	conn := nbsnet.NewNbsConn(c, Sid, nbsnet.CTypeNormal)
	return conn, nil
}

func (network *nbsNetwork) ListenUDP(nt string, lAddr *net.UDPAddr) (*nbsnet.NbsUdpConn, error) {

	var realConn *net.UDPConn
	var cType nbsnet.ConnType
	if network.CanServe {
		c, err := net.ListenUDP(nt, lAddr)
		if err != nil {
			return nil, err
		}
		realConn = c
		cType = nbsnet.CTypeNormal
	} else {
		c, err := shareport.ListenUDP(nt, lAddr.String())
		if err != nil {
			return nil, err
		}
		realConn = c
		cType = nbsnet.CTypeNatListen
	}

	conn := nbsnet.NewNbsConn(realConn, lAddr.String(), cType)
	return conn, nil
}

//TODO::bind local port and ip can't support right now.
func (network *nbsNetwork) Connect(lAddr, rAddr *nbsnet.NbsUdpAddr, toPort int) (*nbsnet.NbsUdpConn, error) {

	if lAddr == nil {
		lAddr = network.natClient.NatAddr
	}

	if rAddr.CanServe {
		return network.makeDirectConn(lAddr, rAddr, toPort)
	}

	var sessionID = lAddr.NetworkId + ConnectionSeparator + rAddr.NetworkId

	var realConn *net.UDPConn
	var connType nbsnet.ConnType
	if lAddr.CanServe {
		connType = nbsnet.CTypeNatSimplex
		c, err := network.invitePeerBehindNat(lAddr, rAddr, sessionID, toPort)
		if err != nil {
			return nil, err
		}
		realConn = c
	} else {
		c, t, err := network.punchANatHole(lAddr, rAddr, sessionID, toPort)
		if err != nil {
			return nil, err
		}
		connType = t
		realConn = c
	}
	conn := nbsnet.NewNbsConn(realConn, sessionID, connType)
	return conn, nil
}
func (network *nbsNetwork) runLoop() {
	network.cmdRouter[nat.CMDAnswerInvite] = network.answerInvite
	network.cmdRouter[nat.CMDDigOut] = network.digOut
	network.cmdRouter[nat.CMDDigSetup] = network.makeAHole

	for {
		select {
		case task := <-network.natClient.CmdTask:
			process, ok := network.cmdRouter[task.CmdType]
			if !ok {
				logger.Warning("unknown task type:->", task.CmdType)
				continue
			}
			if err := process(task.Params); err != nil {
				logger.Warning(err)
				continue
			}

		case <-network.natClient.Ctx.Done():
			logger.Warning("nat client process quit")
			network.natServer.CanServe = make(chan bool)
			go network.setupNatClient(network.networkId)
			return
		}
	}
}

func (network *nbsNetwork) setupNatClient(peerId string) error {
	defer close(network.natServer.CanServe)
	c, err := nat.NewNatClient(peerId, network.natServer.CanServe)
	if err != nil {
		return err
	}

	network.natClient = c
	network.CanServe = c.CanServer
	if !network.CanServe {
		go network.runLoop()
	}
	return nil
}
