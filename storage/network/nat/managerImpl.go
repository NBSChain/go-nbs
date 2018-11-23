package nat

import (
	"github.com/NBSChain/go-nbs/storage/network/denat"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/shareport"
	"github.com/NBSChain/go-nbs/utils"
	"net"
	"strconv"
	"time"
)

func NewNatManager(networkId string) *Manager {

	denat.GetDeNatSerIns().Setup(networkId)

	localPeers := ExternalIP()
	if len(localPeers) == 0 {
		logger.Panic("no available network")
	}

	logger.Debug("all network interfaces:", localPeers)

	natObj := &Manager{
		networkId: networkId,
		canServe:  make(chan bool),
		cache:     make(map[string]*HostBehindNat),
	}

	natObj.startNatService()

	go natObj.natServiceListening()

	go natObj.cacheManager()

	return natObj
}

func (nat *Manager) SetUpNatChannel(netNatAddr *nbsnet.NbsUdpAddr) error {

	port := strconv.Itoa(utils.GetConfig().NatChanSerPort)
	listener, err := shareport.ListenUDP("udp4", "0.0.0.0:"+port)
	if err != nil {
		logger.Warning("create share listening udp failed.")
		return err
	}

	serverIP := denat.GetDeNatSerIns().GetValidServer()
	client, err := shareport.DialUDP("udp4", "0.0.0.0:"+port, serverIP)
	if err != nil {
		logger.Warning("create share port dial udp connection failed.")
		return err
	}

	tunnel := &KATunnel{
		networkId:  nat.networkId,
		natAddr:    netNatAddr,
		closed:     make(chan bool),
		serverHub:  listener,
		kaConn:     client,
		sharedAddr: client.LocalAddr().String(),
		updateTime: time.Now(),
		proxyCache: make(map[string]*proxyConnItem),
	}

	go tunnel.runLoop()

	go tunnel.listening()

	go tunnel.readKeepAlive()

	go tunnel.connManage()

	nat.NatKATun = tunnel

	return nil
}

func (nat *Manager) WaitNatConfirm() chan bool {
	return nat.canServe
}

func (nat *Manager) PunchANatHole(lAddr, rAddr *nbsnet.NbsUdpAddr, connId string) (*net.UDPConn, error) {

	if err := nat.NatKATun.natHoleStep1InvitePeer(lAddr, rAddr, connId); err != nil {
		return nil, err
	}

	conn, err := nat.NatKATun.natHoleStep2Call(connId, rAddr)
	if err == nil {
		return conn, nil
	}

	return conn, err
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

			//TODO:: Support ip v6 later.
			if ip = ip.To4(); ip == nil {
				continue
			}

			ips = append(ips, ip.String())
		}
	}

	return ips
}
