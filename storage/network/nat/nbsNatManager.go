package nat

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/gogo/protobuf/proto"
	"net"
	"strconv"
	"sync"
)

var logger = utils.GetLogInstance()

const NetIoBufferSize = 1 << 11
const BootStrapNatServerTimeOutInSec = 6

type nbsNatManager struct {
	sync.Mutex
	natServer     *net.UDPConn
	natType       nat_pb.NatType
	publicAddress *net.UDPAddr
	privateIP     string
	networkId     string
}

//TODO::support multiple local ip address.
func NewNatManager(networkId string) Manager {

	localPeers := ExternalIP()
	if len(localPeers) == 0 {
		logger.Panic("no available network")
	}

	logger.Debug("all network interfaces:", localPeers)

	natObj := &nbsNatManager{
		networkId: networkId,
	}

	if !utils.GetConfig().NatServiceOff {

		natObj.startNatService()

		go natObj.natService()
	}

	return natObj
}

//TODO:: support ipv6 later.
func (nat *nbsNatManager) startNatService() {

	natServer, err := net.ListenUDP("udp4", &net.UDPAddr{
		Port: utils.GetConfig().NatServerPort,
	})

	if err != nil {
		logger.Panic("can't start nat natServer.", err)
	}

	nat.natServer = natServer
}

func (nat *nbsNatManager) natService() {

	logger.Info(">>>>>>Nat natServer start to listen......")

	for {
		peerAddr, request, err := nat.readNatRequest()
		if err != nil {
			logger.Error(err)
		}

		switch request.MsgType {
		case nat_pb.NatMsgType_BootStrapReg:
			if err := nat.bootNatResponse(request.BootRegReq, peerAddr); err != nil {
				logger.Error(err)
			}
		case nat_pb.NatMsgType_Ping:
			if err := nat.pong(request.Ping, peerAddr); err != nil {
				logger.Error(err)
			}
		}
	}
}

func (nat *nbsNatManager) readNatRequest() (*net.UDPAddr, *nat_pb.NatRequest, error) {

	data := make([]byte, NetIoBufferSize)

	n, peerAddr, err := nat.natServer.ReadFromUDP(data)
	if err != nil {
		logger.Warning("nat natServer read udp data failed:", err)
		return nil, nil, err
	}

	request := &nat_pb.NatRequest{}
	if err := proto.Unmarshal(data[:n], request); err != nil {
		logger.Warning("can't parse the nat request", err, peerAddr)
		return nil, nil, err
	}

	logger.Debug("request:", request)

	return peerAddr, request, nil
}

func (nat *nbsNatManager) bootNatResponse(request *nat_pb.BootNatRegReq, peerAddr *net.UDPAddr) error {

	response := &nat_pb.BootNatRegRes{}
	response.PublicIp = peerAddr.IP.String()
	response.PublicPort = fmt.Sprintf("%d", peerAddr.Port)
	response.Zone = peerAddr.Zone

	if peerAddr.IP.Equal(net.ParseIP(request.PrivateIp)) {

		response.NatType = nat_pb.NatType_NoNatDevice
	} else if strconv.Itoa(peerAddr.Port) == request.PrivatePort {

		response.NatType = nat_pb.NatType_ToBeChecked
		go nat.ping(peerAddr)

	} else {

		response.NatType = nat_pb.NatType_BehindNat
	}

	responseData, err := proto.Marshal(response)
	if err != nil {
		logger.Warning("failed to marshal nat response data", err)
		return err
	}

	if _, err := nat.natServer.WriteToUDP(responseData, peerAddr); err != nil {
		logger.Warning("failed to send nat response", err)
		return err
	}

	return nil
}

func (nat *nbsNatManager) confirmNatType() {
	nat.Lock()
	defer nat.Unlock()

	if nat.natType == nat_pb.NatType_BehindNat ||
		nat.natType == nat_pb.NatType_ToBeChecked {

		nat.natType = nat_pb.NatType_CanBeNatServer
	}
}

func ExternalIP() []string {

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil
	}

	var ips []string
	for _, face := range interfaces {

		if face.Flags&net.FlagUp == 0 ||
			face.Flags&net.FlagLoopback != 0 {
			continue
		}

		address, err := face.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range address {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			//TODO:: Support ip v6 lter.
			if ip = ip.To4(); ip == nil {
				continue
			}

			ips = append(ips, ip.String())
		}
	}

	return ips
}
