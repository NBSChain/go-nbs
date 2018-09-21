package dataStore
/*****************************************************************
*
*		DataStore interface and implements.
*
*****************************************************************/
const RootServiceURL 	= "/"
const BLOCKServiceURL 	= "/blocks"


type DataStore interface {
	Put(key string, value []byte) error

	Get(key string) (value []byte, err error)

	Has(key string) (exists bool, err error)

	Delete(key string) error

	Query(q Query) (Results, error)

	Batch() (Batch, error)
}


type Batch interface {

	Put(key string, val []byte) error

	Delete(key string) error

	Commit() error
}
