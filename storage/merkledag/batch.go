package merkledag

import (
	"context"
	"github.com/NBSChain/go-nbs/storage/merkledag/ipld"
)

const MaxNodes	= 2 << 7

func NewBatch() *Batch {

	ctx, cancel := context.WithCancel(context.TODO())
	batch := &Batch{
		ctx:    ctx,
		cancel: cancel,
		nodes:  make([]ipld.DagNode, 0, MaxNodes),
	}

	return batch
}

type Batch struct {
	ctx           context.Context
	cancel        func()
	commitResult  error
	nodes         []ipld.DagNode
}


func (batch *Batch) Add(node ipld.DagNode) error {

	if batch.commitResult != nil {
		return batch.commitResult
	}

	batch.nodes = append(batch.nodes, node)
	if len(batch.nodes) >= MaxNodes{
		go batch.asyCommit()
		batch.nodes =  make([]ipld.DagNode, 0, MaxNodes)
	}

	return nil
}

func (batch *Batch) Commit() error {

	if batch.commitResult != nil {
		return batch.commitResult
	}

	defer batch.cancel()

	batch.asyCommit()

	return batch.commitResult
}


func (batch *Batch) asyCommit(){

	reminder := len(batch.nodes)
	if reminder == 0 {
		return
	}

	dagService := GetDagInstance()
	err := dagService.AddMany(batch.nodes)
	if err != nil{
		batch.commitResult = err
	}
}
