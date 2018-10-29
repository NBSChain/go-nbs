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
const MaxNatServerItem 	= 1 << 10
const SessionTimeOut 	= 24
const NetIoBufferSize	= 1 << 11

type natItem struct {
	nodeID		string
	privateInfo	*net.UDPAddr
	publicInfo	*net.UDPAddr
	updateTIme	time.Time
}

type nbsNat struct {
	sync.Mutex
	peers         map[string]natItem
	server        *net.UDPConn
	isPublic      bool
	publicAddress *net.UDPAddr
	privateIP     string
}

//TODO::support multiple local ip address.
func NewNatManager() NAT{
	localPeers := ExternalIP()
	if len(localPeers) == 0{
		logger.Panic("no available network")
	}

	natObj := &nbsNat{
		peers:     make(map[string]natItem),
		privateIP: localPeers[0],
	}

	natObj.startNatServer()

	natObj.registerToBootStrap()

	go natObj.natService()

	go natObj.runLoop()
	
	return natObj
}

func (nat *nbsNat)  startNatServer() {

	l, err := net.ListenUDP("udp",&net.UDPAddr{
		Port:utils.GetConfig().NatServerPort,
	})

	if err != nil{
		logger.Panic("can't start nat server.", err)
	}

	nat.server = l
}

func (nat *nbsNat) natService()  {

	for {
		data := make([]byte, NetIoBufferSize)

		n, peerAddr, err :=nat.server.ReadFromUDP(data)
		if err != nil{
			logger.Warning("nat server read udp data failed:", err)
			continue
		}

		request := &nat_pb.NatRequest{}
		if err := proto.Unmarshal(data[:n], request); err != nil{
			logger.Warning("can't parse the nat request", err, peerAddr)
			continue
		}

		response := &nat_pb.NatResponse{}
		if peerAddr.IP.Equal(net.ParseIP(request.PrivateIp)){
			response.IsAfterNat = false
		}else{
			nat.cacheItem(peerAddr, request)
			response.IsAfterNat 	= true
			response.PublicIp 	= peerAddr.IP.String()
			response.PublicPort 	= int32(peerAddr.Port)
			response.Zone 		= peerAddr.Zone
		}

		 responseData, err := proto.Marshal(response)
		if err != nil{
			logger.Warning("failed to marshal nat response data", err)
			continue
		}

		if _, err := nat.server.WriteToUDP(responseData, peerAddr); err != nil{
			logger.Warning("failed to send nat response", err)
			continue
		}
	}
}

//TODO:: set multiple servers to make it stronger.
func (nat *nbsNat) registerToBootStrap() error {

	config := utils.GetConfig()

	natServerAddr := &net.UDPAddr{
		IP:[]byte(config.NatServerIP),
		Port:config.NatServerPort,
	}

	connection, err := net.DialUDP("udp", &net.UDPAddr{
		Port:config.NatClientPort,
	}, natServerAddr)

	if err != nil{
		logger.Error("can't know who am I", err)
		return err
	}
	defer connection.Close()

	request := &nat_pb.NatRequest{
		NodeId:"",//TODO::
		PrivateIp:nat.privateIP,
		PrivatePort:int32(config.NatClientPort),
	}

	requestData, err := proto.Marshal(request)
	if err != nil{
		logger.Error("failed to marshal nat request", err)
		return err
	}

	if _, err := connection.WriteToUDP(requestData, natServerAddr); err != nil{
		logger.Error("failed to send nat request to server ", err)
		return err
	}

	responseData:= make([]byte, NetIoBufferSize)
	hasRead, _, err := connection.ReadFromUDP(responseData)
	if err != nil{
		logger.Error("failed to read nat response from server", err)
		return err
	}

	response := &nat_pb.NatResponse{}
	if err := proto.Unmarshal(responseData[:hasRead], response); err != nil{
		logger.Error("failed to unmarshal nat response data", err)
		return err
	}

	if response.IsAfterNat{
		nat.publicAddress = &net.UDPAddr{
			IP:[]byte(response.PublicIp),
			Port:int(response.PublicPort),
			Zone:response.Zone,
		}
	}

	return nil
}

func (nat *nbsNat) cacheItem(publicInfo *net.UDPAddr, privateInfo *nat_pb.NatRequest){

	item := natItem{
		nodeID: 	privateInfo.NodeId,
		privateInfo:	&net.UDPAddr{
			IP:	[]byte(privateInfo.PrivateIp),
			Port:	int(privateInfo.PrivatePort),
			Zone:	privateInfo.Zone,
		},
		publicInfo:	publicInfo,
		updateTIme:	time.Now(),
	}

	nat.Lock()
	nat.peers[item.nodeID] = item
	nat.Unlock()
}

func (nat *nbsNat) runLoop()  {
	for {
		if len(nat.peers) < MaxNatServerItem{
			time.Sleep(time.Second)
			logger.Debug("no much item to handle")
			continue
		}

		rightNow := time.Now()

		nat.Lock()
		for key, value := range nat.peers{
			if rightNow.Sub(value.updateTIme) > SessionTimeOut* time.Hour{
				delete(nat.peers, key)
			}
		}
		nat.Unlock()
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
			face.Flags&net.FlagLoopback != 0{
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

			ips = append(ips, ip.String())
		}
	}

	return ips
}