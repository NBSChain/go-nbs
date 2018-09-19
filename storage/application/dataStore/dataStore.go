package dataStore

type DataStore interface {
	Put(key Key, value []byte) error

	Get(key Key) (value []byte, err error)

	Has(key Key) (exists bool, err error)

	Delete(key Key) error

	Query(q Query) (Results, error)
}

type Batch interface {
	Put(key Key, val []byte) error

	Delete(key Key) error

	Commit() error
}
