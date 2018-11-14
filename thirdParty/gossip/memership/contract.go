package memership

import (
	"github.com/NBSChain/go-nbs/thirdParty/gossip/pb"
	"github.com/NBSChain/go-nbs/utils"
)

var (
	logger = utils.GetLogInstance()
)

type ContractNode struct {
}

func newContractNode() *ContractNode {

	node := &ContractNode{}

	return node
}

func (node *ContractNode) proxyInit(request *pb.InitSub) {

}
