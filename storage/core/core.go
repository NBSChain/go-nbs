package core

import "github.com/NBSChain/go-nbs/storage/routing"

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
