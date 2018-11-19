package nat

import (
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/gogo/protobuf/proto"
	"net"
	"sync"
	"time"
)

var logger = utils.GetLogInstance()

const NetIoBufferSize = 1 << 11
const BootStrapNatServerTimeOutInSec = 4

type hostBehindNat struct {
	addr       *MulAddr
	kaAddr     *net.UDPAddr
	updateTIme time.Time
}

type Manager struct {
	cacheLock     sync.Mutex
	selfNatServer *net.UDPConn
	networkId     string
	canServe      chan bool
	SelfAddr      *MulAddr
	NatKATun      *KATunnel
	cache         map[string]*hostBehindNat
	dNatServer    string
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
