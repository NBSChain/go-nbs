package memership

import (
	"github.com/NBSChain/go-nbs/storage/network"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/thirdParty/gossip/pb"
	"github.com/NBSChain/go-nbs/utils"
)

var (
	logger = utils.GetLogInstance()
)

type contractNode struct {
	isPublicNode bool
}

func newContractNode() *contractNode {

	natType := network.GetInstance().NatType()

	var canService bool
	switch natType {
	case nat_pb.NatType_UnknownRES:
		canService = false

	case nat_pb.NatType_NoNatDevice:
		canService = true

	case nat_pb.NatType_BehindNat:
		canService = false

	case nat_pb.NatType_CanBeNatServer:
		canService = true

	case nat_pb.NatType_ToBeChecked:
		canService = false
	}

	node := &contractNode{
		isPublicNode: canService,
	}

	return node
}

func (node *contractNode) initSub(sub *pb.InitSub) {

}
