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
	Conn       *net.TCPConn
	PriAddr    string
}

type Server struct {
	sysNatServer  *net.TCPListener
	holeDigHelper *net.UDPConn
	networkId     string
	CanServe      chan bool
	cacheLock     sync.Mutex
	cache         map[string]*HostBehindNat
}

func NewNatServer(networkId string) *Server {
	natObj := &Server{
		networkId: networkId,
		CanServe:  make(chan bool),
		cache:     make(map[string]*HostBehindNat),
	}

	natObj.initService()
	go natObj.holeHelper()
	go natObj.ctrlChServer()
	return natObj
}

//TODO:: support ipv6 later.
func (nat *Server) initService() {

	natServer, err := net.ListenTCP("tcp4", &net.TCPAddr{
		Port: utils.GetConfig().NatServerPort,
	})
	if err != nil {
		logger.Panic("can't start nat sysNatServer:->", err)
	}
	nat.sysNatServer = natServer

	h, err := net.ListenUDP("udp4", &net.UDPAddr{
		Port: utils.GetConfig().HolePuncherPort,
	})
	if err != nil {
		logger.Panic("can't start hole punch helper:->", err)
	}
	nat.holeDigHelper = h
}

func (nat *Server) ctrlChServer() {

	logger.Info(">>>>>>Nat sysNatServer start to listen......")

	for {
		conn, err := nat.sysNatServer.AcceptTCP()
		if err != nil {
			logger.Warning("nat sysNatServer accept err:->", err)
			continue
		}
		if err := conn.SetKeepAlive(true); err != nil {
			logger.Warning("set keep alive err:->", err)
			continue
		}

		go nat.processCtrlMsg(conn)
	}
}

func (nat *Server) forwardMsg(nodeId string, data []byte) error {

	nat.cacheLock.Lock()
	defer nat.cacheLock.Unlock()

	item, ok := nat.cache[nodeId]
	if !ok {
		return NotFundErr
	}
	logger.Debug("Step3: forward notification to applier:", item.Conn.RemoteAddr())
	if _, err := item.Conn.Write(data); err != nil {
		delete(nat.cache, nodeId)
	}
	return nil
}

func (nat *Server) holeHelper() {

	for {
		buffer := make([]byte, utils.NormalReadBuffer)
		n, peerAddr, err := nat.holeDigHelper.ReadFromUDP(buffer)
		if err != nil {
			logger.Error("this dig helper conn is bad:->", err)
			return
		}

		msg := net_pb.NatMsg{}
		proto.Unmarshal(buffer[:n], &msg)

		switch msg.Typ {
		case nbsnet.NatReversInvite:
			invite := msg.ReverseInvite
			err = nat.forwardMsg(invite.PeerId, buffer[:n])
			logger.Info("forward reverse invite request:->", invite)

		case nbsnet.NatDigApply:
			req := msg.DigApply
			req.Public = peerAddr.String()
			data, _ := proto.Marshal(&msg)
			err = nat.forwardMsg(req.TargetId, data)
			logger.Info("hole punch step2-2 forward dig out message:->", req.Public)

		case nbsnet.NatDigConfirm:
			ack := msg.DigConfirm
			ack.Public = peerAddr.String()
			data, _ := proto.Marshal(&msg)
			err = nat.forwardMsg(ack.TargetId, data)
			logger.Info("hole punch step2-6 forward dig out notification:->", ack.Public)

		}

		if err != nil {
			logger.Warning("failed to process dig msg:->", err)
		}
	}
}

func (nat *Server) processCtrlMsg(conn *net.TCPConn) {

	defer conn.Close()

	defer func() {
		nat.cacheLock.Lock()
		for nodeId, item := range nat.cache {
			if item.Conn == conn {
				delete(nat.cache, nodeId)
			}
		}
		nat.cacheLock.Unlock()
	}()

	for {
		data := make([]byte, utils.NormalReadBuffer)
		n, err := conn.Read(data)
		if err != nil {
			logger.Warning("read data from control channel err:->", err, conn.RemoteAddr().String())
			return
		}

		msg := &net_pb.NatMsg{}
		if err := proto.Unmarshal(data[:n], msg); err != nil {
			logger.Warning("can't parse the nat message", err)
			continue
		}

		logger.Debug("control channel message:->", msg)

		switch msg.Typ {

		case nbsnet.NatBootReg:
			err = nat.checkWhoIsHe(conn, msg.BootReg, msg.NetID)

		case nbsnet.NatPingPong:
			nat.CanServe <- true
			logger.Debug("I can serve as in public network.")

		case nbsnet.NatKeepAlive:
			err = nat.updateKATime(msg.KeepAlive, conn, msg.NetID)

		case nbsnet.NatCheckNetType:

			ip, port, _ := net.SplitHostPort(conn.RemoteAddr().String())
			msg.NatTypeCheck = &net_pb.IpV4{
				Ip:   ip,
				Port: port,
			}

			data, _ := proto.Marshal(msg)
			if _, err := conn.Write(data); err != nil {
				logger.Warning("write back nat info err:->", err)
			}

			return

		default:
			logger.Warning("this msg can't be processed success->", msg)
		}

		if err != nil {
			logger.Warning("control channel err :->", err)
		}
	}
}
