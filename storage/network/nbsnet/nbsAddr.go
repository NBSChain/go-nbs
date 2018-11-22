package nbsnet

import (
	"github.com/NBSChain/go-nbs/storage/network/pb"
)

type NbsUdpAddr struct {
	NetworkId string
	CanServe  bool
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
