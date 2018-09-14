package merkledag

import (
	"context"
	"errors"
	"github.com/NBSChain/go-nbs/storage/merkledag/ipld"
	"time"
)

const MaxNodes = 2 << 7
const MaxTimeToWaitInSeconds = 60

func NewBatch() *Batch {

	ctx, cancel := context.WithCancel(context.TODO())
	batch := &Batch{
		ctx:    ctx,
		cancel: cancel,
		nodes:  make(chan ipld.DagNode, MaxNodes),
	}

	go batch.run()

	return batch
}

type Batch struct {
	ctx           context.Context
	cancel        func()
	activeCommits int
	commitResult  error
	nodes         chan ipld.DagNode
}

func (batch *Batch) run() {

	dagService := GetDagInstance()

	for {
		select {

		case <-batch.ctx.Done():
			batch.commitResult = batch.ctx.Err()
			break

		case node := <-batch.nodes:
			err := dagService.Add(node)
			if err != nil {
				batch.commitResult = err
				batch.cancel()
				break
			}
		}
	}
}

func (batch *Batch) Add(node ipld.DagNode) error {

	if batch.commitResult != nil {
		return batch.commitResult
	}

	batch.nodes <- node

	return nil
}

func (batch *Batch) Commit() error {

	defer batch.cancel()

	if batch.commitResult != nil {
		return batch.commitResult
	}

	reminder := len(batch.nodes)
	if reminder == 0 {
		return nil
	}

	select {

	case <-time.After(time.Second * MaxTimeToWaitInSeconds):
		if batch.commitResult != nil || len(batch.nodes) > 0 {
			return errors.New("can't add dag node successfully. ")
		}
	}

	return nil
}
