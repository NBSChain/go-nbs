package nat

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/gogo/protobuf/proto"
	"net"
	"strconv"
)

var logger = utils.GetLogInstance()

const NetIoBufferSize = 1 << 11
const BootStrapNatServerTimeOutInSec = 4

type Manager struct {
	selfNatServer *net.UDPConn
	networkId     string
	canServe      chan bool
	dNatServer    *decentralizedNatServer
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

func (nat *Manager) runLoop() {

	logger.Info(">>>>>>Nat selfNatServer start to listen......")

	for {
		peerAddr, request, err := nat.readNatRequest()
		if err != nil {
			logger.Error(err)
		}

		switch request.MsgType {
		case net_pb.NatMsgType_BootStrapReg:
			if err := nat.bootNatResponse(request.BootRegReq, peerAddr); err != nil {
				logger.Error(err)
			}
		case net_pb.NatMsgType_Ping:
			if err := nat.pong(request.Ping, peerAddr); err != nil {
				logger.Error(err)
			}
		}
	}
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

func (nat *Manager) bootNatResponse(request *net_pb.BootNatRegReq, peerAddr *net.UDPAddr) error {

	response := &net_pb.BootNatRegRes{}
	response.PublicIp = peerAddr.IP.String()
	response.PublicPort = fmt.Sprintf("%d", peerAddr.Port)
	response.Zone = peerAddr.Zone

	if peerAddr.IP.Equal(net.ParseIP(request.PrivateIp)) {

		response.NatType = net_pb.NatType_NoNatDevice
	} else if strconv.Itoa(peerAddr.Port) == request.PrivatePort {

		response.NatType = net_pb.NatType_ToBeChecked
		go nat.ping(peerAddr)

	} else {

		response.NatType = net_pb.NatType_BehindNat
	}

	pbRes := &net_pb.Response{
		MsgType:net_pb.NatMsgType_BootStrapReg,
		BootRegRes:response,
	}

	pbResData, err := proto.Marshal(pbRes)
	if err != nil {
		logger.Warning("failed to marshal nat response data", err)
		return err
	}

	if _, err := nat.selfNatServer.WriteToUDP(pbResData, peerAddr); err != nil {
		logger.Warning("failed to send nat response", err)
		return err
	}

	return nil
}
