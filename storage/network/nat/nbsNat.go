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

type nbsNat struct {
	sync.Mutex
	natServer     *net.UDPConn
	isPublic      bool
	publicAddress *net.UDPAddr
	privateIP     string
	networkId     string
}

//TODO::support multiple local ip address.
func NewNatManager(networkId string) NAT {

	localPeers := ExternalIP()
	if len(localPeers) == 0 {
		logger.Panic("no available network")
	}

	logger.Debug("all network interfaces:", localPeers)

	natObj := &nbsNat{
		privateIP: localPeers[0],
		networkId: networkId,
	}

	if !utils.GetConfig().NatServiceOff {

		natObj.startNatService()

		go natObj.natService()
	}

	return natObj
}

//TODO:: support ipv6 later.
func (nat *nbsNat) startNatService() {

	natServer, err := net.ListenUDP("udp4", &net.UDPAddr{
		Port: utils.GetConfig().NatServerPort,
	})

	if err != nil {
		logger.Panic("can't start nat natServer.", err)
	}

	nat.natServer = natServer
}

func (nat *nbsNat) natService() {

	logger.Info(">>>>>>Nat natServer start to listen......")

	for {
		data := make([]byte, NetIoBufferSize)

		n, peerAddr, err := nat.natServer.ReadFromUDP(data)
		if err != nil {
			logger.Warning("nat natServer read udp data failed:", err)
			continue
		}

		request := &nat_pb.NatRequest{}
		if err := proto.Unmarshal(data[:n], request); err != nil {
			logger.Warning("can't parse the nat request", err, peerAddr)
			continue
		}

		logger.Debug("get nat request from client:", request)

		response := &nat_pb.NatResponse{}
		if peerAddr.IP.Equal(net.ParseIP(request.PrivateIp)) {
			response.IsAfterNat = false
		} else {
			response.IsAfterNat = true
			response.PublicIp = peerAddr.IP.String()
			response.PublicPort = int32(peerAddr.Port)
			response.Zone = peerAddr.Zone
		}

		responseData, err := proto.Marshal(response)
		if err != nil {
			logger.Warning("failed to marshal nat response data", err)
			continue
		}

		if _, err := nat.natServer.WriteToUDP(responseData, peerAddr); err != nil {
			logger.Warning("failed to send nat response", err)
			continue
		}
	}
}

func (nat *nbsNat) connectToNatServer(serverIP string, localAddress *net.UDPAddr) (*net.UDPConn, error) {

	config := utils.GetConfig()
	natServerAddr := &net.UDPAddr{
		IP:   net.ParseIP(serverIP),
		Port: config.NatServerPort,
	}
	return net.DialUDP("udp", localAddress, natServerAddr)
}

func (nat *nbsNat) sendNatRequest(connection *net.UDPConn) error {

	request := &nat_pb.NatRequest{
		NodeId:      nat.networkId,
		PrivateIp:   nat.privateIP,
		PrivatePort: int32(utils.GetConfig().P2pListenPort),
	}

	requestData, err := proto.Marshal(request)
	if err != nil {
		logger.Error("failed to marshal nat request", err)
		return err
	}

	if no, err := connection.Write(requestData); err != nil || no == 0 {
		logger.Error("failed to send nat request to natServer ", err, no)
		return err
	}

	return nil
}

func (nat *nbsNat) parseNatResponse(connection *net.UDPConn) (*nat_pb.NatResponse, error) {

	responseData := make([]byte, NetIoBufferSize)
	hasRead, _, err := connection.ReadFromUDP(responseData)
	if err != nil {
		logger.Error("failed to read nat response from natServer", err)
		return nil, err
	}

	response := &nat_pb.NatResponse{}
	if err := proto.Unmarshal(responseData[:hasRead], response); err != nil {
		logger.Error("failed to unmarshal nat response data", err)
		return nil, err
	}

	logger.Debug("get response data from nat natServer:", response)

	if response.IsAfterNat {
		nat.publicAddress = &net.UDPAddr{
			IP:   net.ParseIP(response.PublicIp),
			Port: int(response.PublicPort),
			Zone: response.Zone,
		}
		nat.isPublic = false
	} else {
		nat.isPublic = true
		nat.publicAddress = nil
	}

	return response, nil
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

//TODO:: set multiple servers to make it stronger.
func (nat *nbsNat) FetchNatInfo() error {

	config := utils.GetConfig()

	for _, serverIP := range config.NatServerIP {

		//TIPS:: no need to bind local host and local port right now
		connection, err := nat.connectToNatServer(serverIP, nil)

		if err != nil {
			logger.Error("can't know who am I", err)
			continue
		}

		connection.SetDeadline(time.Now().Add(3 * time.Second))

		if err := nat.sendNatRequest(connection); err != nil {
			continue
		}

		_, err = nat.parseNatResponse(connection)
		if err == nil {
			break
		}
	}

	return nil
}
