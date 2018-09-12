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
		MaxSize:       8 << 20,
		MaxNodes:      128,
	}
}

type Batch struct {
	ctx    context.Context
	cancel func()

	activeCommits int
	err           error
	commitResults chan error

	nodes []ipld.DagNode
	size  int

	MaxSize  int
	MaxNodes int
}

func (t *Batch) Add(nd ipld.DagNode) error {
	if t.err != nil {
		return t.err
	}

	t.processResults()

	if t.err != nil {
		return t.err
	}

	t.nodes = append(t.nodes, nd)
	t.size += len(nd.RawData())

	if t.size > t.MaxSize || len(t.nodes) > t.MaxNodes {
		t.asyncCommit()
	}
	return t.err
}

func (t *Batch) processResults() {
	for t.activeCommits > 0 {
		select {
		case err := <-t.commitResults:
			t.activeCommits--
			if err != nil {
				t.setError(err)
				return
			}
		default:
			return
		}
	}
}

func (t *Batch) asyncCommit() {
	numBlocks := len(t.nodes)
	if numBlocks == 0 {
		return
	}
	if t.activeCommits >= ParallelBatchCommits {
		select {
		case err := <-t.commitResults:
			t.activeCommits--

			if err != nil {
				t.setError(err)
				return
			}
		case <-t.ctx.Done():
			t.setError(t.ctx.Err())
			return
		}
	}
	go func(ctx context.Context, b []ipld.DagNode, result chan error) {

		dagService := GetInstance()
		select {
		case result <- dagService.AddMany(b):
		case <-ctx.Done():
		}
	}(t.ctx, t.nodes, t.commitResults)

	t.activeCommits++
	t.nodes = make([]ipld.DagNode, 0, numBlocks)
	t.size = 0

	return
}

func (t *Batch) setError(err error) {
	t.err = err

	t.cancel()

	// Drain as much as we can without blocking.
loop:
	for {
		select {
		case <-t.commitResults:
		default:
			break loop
		}
	}

	// Be nice and cleanup. These can take a *lot* of memory.
	t.commitResults = nil
	t.ctx = nil
	t.nodes = nil
	t.size = 0
	t.activeCommits = 0
}

func (t *Batch) Commit() error {
	if t.err != nil {
		return t.err
	}

	t.asyncCommit()

loop:
	for t.activeCommits > 0 {
		select {
		case err := <-t.commitResults:
			t.activeCommits--
			if err != nil {
				t.setError(err)
				break loop
			}
		case <-t.ctx.Done():
			t.setError(t.ctx.Err())
			break loop
		}
	}

	return t.err
}
