package nat

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/denat"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/utils"
	"net"
	"time"
)

type MulAddr struct {
	canServe    bool
	addrID      string
	privateIP   string
	privatePort string
	publicIp    string
	publicPort  string
}

func NewNatManager(networkId string) *Manager {

	localPeers := ExternalIP()
	if len(localPeers) == 0 {
		logger.Panic("no available network")
	}

	logger.Debug("all network interfaces:", localPeers)

	decentralizedNatServer := denat.GetDNSInstance()

	natObj := &Manager{
		networkId:  networkId,
		canServe:   make(chan bool),
		cache:      make(map[string]*hostBehindNat),
		dNatServer: decentralizedNatServer.RandomNatSer(),
	}

	natObj.startNatService()

	go natObj.natServiceListening()

	go natObj.cacheManager()

	return natObj
}

func (nat *Manager) FindWhoAmI() (canServer bool, err error) {

	config := utils.GetConfig()

	//TODO:: we will replace this nat server by gossip protocol based nat server chose logic.
	for _, serverIP := range config.NatServerIP {

		conn, err := nat.connectToNatServer(serverIP)
		if err != nil {
			logger.Error("can't know who am I", err)
			conn.Close()
			continue
		}
		conn.SetDeadline(time.Now().Add(time.Second * 3))

		localHost, port, err := nat.sendNatRequest(conn)
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

		addr := &MulAddr{
			addrID:      nat.networkId,
			publicIp:    response.PublicIp,
			publicPort:  response.PublicPort,
			privatePort: port,
			privateIP:   localHost,
			canServe:    CanServe(response.NatType),
		}
		nat.SelfAddr = addr

		if response.NatType == net_pb.NatType_ToBeChecked {

			select {
			case c := <-nat.canServe:
				addr.canServe = c

			case <-time.After(time.Second * BootStrapNatServerTimeOutInSec / 2):
				addr.canServe = false
			}

			close(nat.canServe)
		}

		conn.Close()

		return addr.canServe, nil
	}

	return false, fmt.Errorf("can't find available NAT server")
}

func (nat *Manager) NewKAChannel() error {

	ka, err := nat.newKATunnel()
	if err != nil {
		logger.Warning("failed to create nat server ka channel.")
		return err
	}

	request := &net_pb.NatRequest{
		MsgType: net_pb.NatMsgType_BootStrapReg,
		BootRegReq: &net_pb.BootNatRegReq{
			NodeId:      nat.networkId,
			PrivateIp:   priIP,
			PrivatePort: priPort,
		},
	}

	if err := nat.registerPriHost(request); err != nil {
		return err
	}

	if err := ka.InitNatTunnel(); err != nil {
		logger.Warning("create NAT keep alive tunnel failed", err)
		return err
	}

	nat.NatKATun = ka

	return nil
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

func CanServe(natType net_pb.NatType) bool {

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
