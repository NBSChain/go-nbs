package application

import (
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

	router := routing.GetInstance()
	router.Run()

	return nil
}

func NewNode() StorageNode {

	node := &NbsStorageNode{
		nodeId: "",
	}

	return node
}
