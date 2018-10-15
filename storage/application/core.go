package application

import (
	"github.com/NBSChain/go-nbs/storage/bitswap"
	"github.com/NBSChain/go-nbs/storage/merkledag"
	"github.com/NBSChain/go-nbs/storage/routing"
	"github.com/NBSChain/go-nbs/utils"
)

var logger = utils.GetLogInstance()

type StorageNode interface {
	Online() error
}

type NbsStorageNode struct {
	nodeId string
}

func (*NbsStorageNode) Online() error {

	//naming
	merkledag.GetDagInstance()

	bitswap.GetSwapInstance()

	routing.GetInstance()

	return nil
}

func NewNode() StorageNode {

	node := &NbsStorageNode{
		nodeId: "",
	}

	return node
}
