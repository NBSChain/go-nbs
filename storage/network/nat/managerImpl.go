package nat

import (
	"fmt"
	das "github.com/NBSChain/go-nbs/storage/network/decentralizeNatSys"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/storage/network/shareport"
	"github.com/NBSChain/go-nbs/utils"
	"net"
	"strconv"
	"time"
)

//TODO::support multiple local ip address.
func NewNatManager(networkId string) *Manager {

	localPeers := ExternalIP()
	if len(localPeers) == 0 {
		logger.Panic("no available network")
	}

	logger.Debug("all network interfaces:", localPeers)

	natObj := &Manager{
		networkId:  networkId,
		canServe:   make(chan bool),
		cache:      make(map[string]*clientItem),
		dNatServer: das.NewDecentralizedNatServer(),
	}

	natObj.startNatService()

	go natObj.natServiceListening()

	go natObj.cacheManager()

	return natObj
}

//TODO:: we will replace this nat server by gossip protocol based nat server chose logic.
func (nat *Manager) FindWhoAmI() (address *net_pb.NbsAddress, err error) {

	config := utils.GetConfig()

	for _, serverIP := range config.NatServerIP {

		conn, err := nat.connectToNatServer(serverIP)
		if err != nil {
			logger.Error("can't know who am I", err)
			conn.Close()
			continue
		}
		conn.SetDeadline(time.Now().Add(time.Second * 3))

		localHost, err := nat.sendNatRequest(conn)
		if err != nil {
			logger.Error("failed to read nat response:", err)
			conn.Close()
			continue
		}

		response, err := nat.parseNatResponse(conn)
		if err != nil {
			logger.Debug("get NAT server info success.")
			conn.Close()
			continue
		}

		address = &net_pb.NbsAddress{
			PublicIp:     response.PublicIp,
			PrivateIp:    localHost,
			CanBeService: IsPublic(response.NatType),
		}

		if response.NatType == net_pb.NatType_ToBeChecked {

			select {
			case canServer := <-nat.canServe:
				address.CanBeService = canServer
			case <-time.After(time.Second * BootStrapNatServerTimeOutInSec / 2):
				address.CanBeService = false
			}
			close(nat.canServe)
		}

		return address, nil
	}

	return nil, fmt.Errorf("can't find available NAT server")
}

func (nat *Manager) NewKAChannel() (*KATunnel, error) {

	port := strconv.Itoa(utils.GetConfig().NatClientPort)
	natServer := nat.dNatServer.GossipNatServer()

	listener, err := shareport.ListenUDP("udp4", port)
	if err != nil {
		logger.Warning("create share listening udp failed.")
		return nil, err
	}

	client, err := shareport.DialUDP("udp4", "0.0.0.0:"+port, natServer)
	if err != nil {
		logger.Warning("create share port dial udp connection failed.")
		return nil, err
	}

	localAddr := client.LocalAddr().String()
	priIP, priPort, err := net.SplitHostPort(localAddr)

	channel := &KATunnel{
		closed:      make(chan bool),
		networkId:   nat.networkId,
		receiveHub:  listener,
		kaConn:      client,
		privateIP:   priIP,
		privatePort: priPort,
		updateTime:  time.Now(),
	}

	return channel, nil
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

func IsPublic(natType net_pb.NatType) bool {

	var canService bool
	switch natType {
	case net_pb.NatType_UnknownRES:
		canService = false

	case net_pb.NatType_NoNatDevice:
		canService = true

	case net_pb.NatType_BehindNat:
		canService = false

	case net_pb.NatType_CanBeNatServer:
		canService = true

	case net_pb.NatType_ToBeChecked:
		canService = false
	}

	return canService
}
