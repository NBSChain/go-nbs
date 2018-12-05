package nbsnet

import (
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/thirdParty/gossip/pb"
	"net"
	"strconv"
)

//TODO:: refactoring this address setting.
type NbsUdpAddr struct {
	NetworkId string
	CanServe  bool
	NatServer string
	NatIp     string
	NatPort   int32
	PubIp     string
	PubPort   int32
	PriIp     string
	PriPort   int32
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

func SplitHostPort(addr string) (string, int32, error) {
	host, port, err := net.SplitHostPort(addr)
	intPort, _ := strconv.Atoi(port)

	return host, int32(intPort), err
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

func ConvertToGossipAddr(addr *NbsUdpAddr, nodeId string) *pb.BasicHost {
	return &pb.BasicHost{
		CanServer: addr.CanServe,
		NatServer: addr.NatServer,
		NatIP:     addr.NatIp,
		NatPort:   addr.NatPort,
		PriIP:     addr.PriIp,
		PubIp:     addr.PubIp,
		NetworkId: nodeId,
	}
}

func ConvertFromGossipAddr(addr *pb.BasicHost) *NbsUdpAddr {
	return &NbsUdpAddr{
		CanServe:  addr.CanServer,
		NatServer: addr.NatServer,
		NatIp:     addr.NatIP,
		NatPort:   addr.NatPort,
		PubIp:     addr.PubIp,
		PriIp:     addr.PriIP,
		NetworkId: addr.NetworkId,
	}
}
