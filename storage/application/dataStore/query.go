package dataStore

import "github.com/jbenet/goprocess"

type Entry struct {
	Key   string
	Value []byte
}

type Result struct {
	Entry

	Error error
}
type Query struct {
	Prefix   string
	Filters  []Filter
	Orders   []Order
	Limit    int
	Offset   int
	KeysOnly bool
}

type Results interface {
	Query() Query
	Next() <-chan Result
	NextSync() (Result, bool)
	Rest() ([]Entry, error)
	Close() error
	Process() goprocess.Process
}
