package merkledag

import (
	"context"
	"github.com/NBSChain/go-nbs/storage/merkledag/ipld"
	"runtime"
)

var ParallelBatchCommits = runtime.NumCPU() * 2

func NewBatch() *Batch {

	ctx, cancel := context.WithCancel(context.TODO())
	return &Batch{
		ctx:           ctx,
		cancel:        cancel,
		commitResults: make(chan error, ParallelBatchCommits),
		MaxNodes:      128,
	}
}

type Batch struct {
	ctx    context.Context
	cancel func()

	activeCommits int
	err           error
	commitResults chan error

	nodes    []ipld.DagNode
	MaxNodes int
}

func (t *Batch) Add(nd ipld.DagNode) error {

	if t.err != nil {
		return t.err
	}

	t.nodes = append(t.nodes, nd)

	if len(t.nodes) > t.MaxNodes {
		t.asyncCommit()
	}

	return t.err
}

func (t *Batch) addTask(ctx context.Context, b []ipld.DagNode, result chan error) {

	dagService := GetDagInstance()

	select {
	case result <- dagService.AddMany(b):
	case <-ctx.Done():
	}
}

func (t *Batch) asyncCommit() {

	numBlocks := len(t.nodes)
	if numBlocks == 0 {
		return
	}

	if t.activeCommits >= ParallelBatchCommits {
		select {
		case t.err = <-t.commitResults:
			t.activeCommits--

			if t.err != nil {
				t.cancel()
				return
			}
		case <-t.ctx.Done():
			t.err = t.ctx.Err()
			return
		}
	}

	go t.addTask(t.ctx, t.nodes, t.commitResults)

	t.activeCommits++
	t.nodes = make([]ipld.DagNode, 0, numBlocks)

	return
}

func (t *Batch) Commit() error {
	if t.err != nil {
		return t.err
	}

	t.asyncCommit()

	for t.activeCommits > 0 {
		select {
		case t.err = <-t.commitResults:
			t.activeCommits--

			if t.err != nil {
				t.cancel()
				return t.err
			}
		case <-t.ctx.Done():
			return t.ctx.Err()
		}
	}
}
