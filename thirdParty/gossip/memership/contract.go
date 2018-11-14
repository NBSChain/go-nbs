package memership

import (
	"github.com/NBSChain/go-nbs/thirdParty/gossip/pb"
	"github.com/NBSChain/go-nbs/utils"
)

var (
	logger = utils.GetLogInstance()
)

type contractNode struct {
}

func newContractNode() *contractNode {

	node := &contractNode{}

	return node
}

func (node *contractNode) proxyInit(request *pb.InitSub) {

}
