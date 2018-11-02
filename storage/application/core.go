package application

import (
	"github.com/NBSChain/go-nbs/storage/bitswap"
	"github.com/NBSChain/go-nbs/storage/merkledag"
	"github.com/NBSChain/go-nbs/storage/network"
	"github.com/NBSChain/go-nbs/storage/routing"
	"github.com/NBSChain/go-nbs/thirdParty/account"
	"github.com/NBSChain/go-nbs/utils"
)

var logger = utils.GetLogInstance()

type StorageNode interface {
	Online() error
}

type NbsStorageNode struct {
	nodeId string
}

func (node *NbsStorageNode) Online() error {

	//naming

	merkledag.GetDagInstance()

	bitswap.GetSwapInstance()

	routing.GetInstance()

	network.GetInstance().StartUp(node.nodeId)

	return nil
}

func newNode() StorageNode {

	acc := account.GetAccountInstance()

	node := &NbsStorageNode{
		nodeId: acc.GetPeerID(),
	}

	return node
}
