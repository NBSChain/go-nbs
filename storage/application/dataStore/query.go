package dataStore

import "github.com/jbenet/goprocess"

type Entry struct {
	Key   	string
	Value 	[]byte
}

type Result struct {
	entry	Entry
	Error 	error
}
type Query struct {
	Prefix   	string
	Filters  	[]Filter
	Orders   	[]Order
	Limit    	int
	Offset   	int
	KeysOnly 	bool
}

type Results interface {
	Query() 	Query
	Next() 		<-chan Result
	NextSync() 	(Result, bool)
	Rest() 		([]Entry, error)
	Close() 	error
	Process() 	goprocess.Process
}

type results struct {
	query   	Query
	process 	goprocess.Process
	result  	<-chan Result
}

func (r *results) Next() <-chan Result {
	return r.result
}

func (r *results) NextSync() (Result, bool) {
	val, ok := <-r.result
	return val, ok
}

func (r *results) Rest() ([]Entry, error) {
	var entries []Entry

	for result := range r.result {
		if result.Error != nil {
			return entries, result.Error
		}
		entries = append(entries, result.entry)
	}
	<-r.process.Closed()
	return entries, nil
}

func (r *results) Process() goprocess.Process {
	return r.process
}

func (r *results) Close() error {
	return r.process.Close()
}

func (r *results) Query() Query {
	return r.query
}

type ResultBuilder struct {
	Query   	Query
	Process 	goprocess.Process
	Output  	chan Result
}

//TODO:: Refactoring result builder.
func (rb *ResultBuilder) Results() Results {
	return &results{
		query:   rb.Query,
		process: rb.Process,
		result:  rb.Output,
	}
}

const NormalBufSize 	= 1
const KeysOnlyBufSize 	= 128

func NewResultBuilder(query Query) *ResultBuilder {

	bufSize := NormalBufSize
	if query.KeysOnly {
		bufSize = KeysOnlyBufSize
	}

	resultBuilder := &ResultBuilder{
		Query:  query,
		Output: make(chan Result, bufSize),
	}

	resultBuilder.Process = goprocess.WithTeardown(func() error {
		close(resultBuilder.Output)
		return nil
	})

	return resultBuilder
}

func ResultsWithChan(query Query, resultChan <-chan Result) Results {

	resultBuilder := NewResultBuilder(query)

	resultBuilder.Process.Go(func(worker goprocess.Process) {
		for {
			select {
			case <-worker.Closing(): // client told us to close early
				return
			case result, more := <-resultChan:
				if !more {
					return
				}

				select {
				case resultBuilder.Output <- result:
				case <-worker.Closing():
					return
				}
			}
		}
	})

	go resultBuilder.Process.CloseAfterChildren()
	return resultBuilder.Results()
}

func ResultsWithEntries(query Query, entries []Entry) Results {
	resultBuilder := NewResultBuilder(query)

	resultBuilder.Process.Go(func(worker goprocess.Process) {
		for _, entry := range entries {
			select {
			case resultBuilder.Output <- Result{entry: entry}:
			case <-worker.Closing():
				return
			}
		}
		return
	})

	go resultBuilder.Process.CloseAfterChildren()
	return resultBuilder.Results()
}

func ResultsReplaceQuery(res Results, query Query) Results {

	switch r := res.(type) {
	case *results:
		return &results{query, r.process, r.result}

	case *resultsIterator:
		lr := r.legacyResults
		if lr != nil {
			lr = &results{query, lr.process, lr.result}
		}
		return &resultsIterator{query, r.next, r.close, lr}
	default:
		panic("unknown results type")
	}
}

func ResultsFromIterator(query Query, iterator Iterator) Results {
	if iterator.Close == nil {
		iterator.Close = noopClose
	}
	return &resultsIterator{
		query: 	query,
		next:  	iterator.Next,
		close: 	iterator.Close,
	}
}

func noopClose() error {
	return nil
}

type Iterator struct {
	Next  func() (Result, bool)
	Close func() error
}

type resultsIterator struct {
	query         Query
	next          func() (Result, bool)
	close         func() error
	legacyResults *results
}

func (r *resultsIterator) Next() <-chan Result {
	r.useLegacyResults()
	return r.legacyResults.Next()
}

func (r *resultsIterator) NextSync() (Result, bool) {
	if r.legacyResults != nil {
		return r.legacyResults.NextSync()
	} else {
		res, ok := r.next()
		if !ok {
			r.close()
		}
		return res, ok
	}
}

func (r *resultsIterator) Rest() ([]Entry, error) {
	var es []Entry
	for {
		e, ok := r.NextSync()
		if !ok {
			break
		}
		if e.Error != nil {
			return es, e.Error
		}
		es = append(es, e.entry)
	}
	return es, nil
}

func (r *resultsIterator) Process() goprocess.Process {
	r.useLegacyResults()
	return r.legacyResults.Process()
}

func (r *resultsIterator) Close() error {
	if r.legacyResults != nil {
		return r.legacyResults.Close()
	} else {
		return r.close()
	}
}

func (r *resultsIterator) Query() Query {
	return r.query
}

func (r *resultsIterator) useLegacyResults() {
	if r.legacyResults != nil {
		return
	}

	resultBuilder := NewResultBuilder(r.query)

	resultBuilder.Process.Go(func(worker goprocess.Process) {
		defer r.close()
		for {
			result, ok := r.next()
			if !ok {
				break
			}
			select {
			case resultBuilder.Output <- result:
			case <-worker.Closing(): // client told us to close early
				return
			}
		}
		return
	})

	go resultBuilder.Process.CloseAfterChildren()

	r.legacyResults = resultBuilder.Results().(*results)
}