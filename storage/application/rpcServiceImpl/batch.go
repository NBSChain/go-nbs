package rpcServiceImpl

import (
	"github.com/NBSChain/go-nbs/storage/merkledag"
	"github.com/NBSChain/go-nbs/storage/merkledag/ipld"
	"sync"
)

const MaxNodes	= DefaultLinksPerBlock

func NewBatch() *Batch {
	batch := &Batch{
		nodes:  make([]ipld.DagNode, 0, MaxNodes),
	}

	return batch
}

type Batch struct {
	wg		sync.WaitGroup
	commitResult  	error
	nodes         	[]ipld.DagNode
}


func (batch *Batch) Add(node ipld.DagNode) error {

	if batch.commitResult != nil {
		return batch.commitResult
	}

	batch.nodes = append(batch.nodes, node)
	if len(batch.nodes) >= MaxNodes{
		batch.subCommit()
		batch.nodes =  make([]ipld.DagNode, 0, MaxNodes)
	}

	logger.Warning(node.String())
	return nil
}

func (batch *Batch) Commit() error {

	if batch.commitResult != nil {
		return batch.commitResult
	}
	logger.Warning("commit and cancel")

	go batch.subCommit()

	batch.wg.Wait()

	return batch.commitResult
}


func (batch *Batch) subCommit(){

	logger.Warning("...start to subCommit......")

	reminder := len(batch.nodes)
	if reminder == 0 {
		return
	}

	dagService := merkledag.GetDagInstance()
	err := dagService.AddMany(batch.nodes)
	if err != nil{
		logger.Error(err)
		batch.commitResult = err
	}
}
