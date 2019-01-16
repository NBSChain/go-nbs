package dataStore

/*****************************************************************
*
*		DataStore interface and implements.
*
*****************************************************************/
type DataStore interface {
	Put(key string, value []byte) error

	Get(key string) ([]byte, error)

	Has(key string) (bool, error)

	Delete(key string) error

	Query(q Query) (Results, error)

	Batch() (Batch, error)
}

type Batch interface {
	Put(key string, val []byte) error

	Delete(key string) error

	Commit() error
}
