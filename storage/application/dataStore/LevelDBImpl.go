package dataStore

import (
	"github.com/NBSChain/go-nbs/utils"
	"github.com/jbenet/goprocess"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/storage"
	"github.com/syndtr/goleveldb/leveldb/util"
	"os"
	"path/filepath"
)

func newLevelDB(opts *opt.Options) (*Mount, error) {

	var path 	= utils.GetConfig().LevelDBDir
	var err 	error
	var dataBase 	*leveldb.DB

	if path == "" {
		dataBase, err = leveldb.Open(storage.NewMemStorage(), opts)
	} else {

		dataBase, err = leveldb.OpenFile(path, opts)
		if errors.IsCorrupted(err) && !opts.GetReadOnly() {
			dataBase, err = leveldb.RecoverFile(path, opts)
		}
	}

	if err != nil {
		return nil, err
	}

	dateStore := &dataStore{
		dataBase:  	dataBase,
		storePath: 	path,
	}

	return &Mount{
		prefix:		NewKey("/"),
		dataStore:	dateStore,
	}, nil
}
/*****************************************************************
*
*		The functions below are implements of interface
*
*****************************************************************/
//TODO:: What's the interface of following functions.
func (d *dataStore) QueryNew(q Query) (Results, error) {

	if len(q.Filters) > 0 ||
		len(q.Orders) > 0 ||
		q.Limit > 0 ||
		q.Offset > 0 {
		return d.QueryOrig(q)
	}
	var rnge *util.Range
	if q.Prefix != "" {
		rnge = util.BytesPrefix([]byte(q.Prefix))
	}

	iterator := d.dataBase.NewIterator(rnge, nil)

	return ResultsFromIterator(q, Iterator{
		Next: func() (Result, bool) {
			ok := iterator.Next()
			if !ok {
				return Result{}, false
			}

			key 	:= string(iterator.Key())
			entry 	:= Entry{Key: key}

			if !q.KeysOnly {
				buf := make([]byte, len(iterator.Value()))
				copy(buf, iterator.Value())
				entry.Value = buf
			}

			return Result{Entry: entry}, true
		},
		Close: func() error {
			iterator.Release()
			return nil
		},
	}), nil
}

func (d *dataStore) QueryOrig(query Query) (Results, error) {

	resultBuilder := NewResultBuilder(query)
	resultBuilder.Process.Go(func(worker goprocess.Process) {
		d.runQuery(worker, resultBuilder)
	})

	go resultBuilder.Process.CloseAfterChildren()

	queryResult := resultBuilder.Results()

	for _, filter := range query.Filters {
		queryResult = NaiveFilter(queryResult, filter)
	}

	for _, order := range query.Orders {
		queryResult = NaiveOrder(queryResult, order)
	}

	return queryResult, nil
}

func (d *dataStore) runQuery(worker goprocess.Process, resultBuilder *ResultBuilder) {

	var rnge *util.Range
	if resultBuilder.Query.Prefix != "" {
		rnge = util.BytesPrefix([]byte(resultBuilder.Query.Prefix))
	}

	iterator := d.dataBase.NewIterator(rnge, nil)
	defer iterator.Release()

	if resultBuilder.Query.Offset > 0 {
		for j := 0; j < resultBuilder.Query.Offset; j++ {
			iterator.Next()
		}
	}

	for sent := 0; iterator.Next(); sent++ {

		if resultBuilder.Query.Limit > 0 && sent >= resultBuilder.Query.Limit {
			break
		}

		key 	:= string(iterator.Key())
		entry 	:= Entry{Key: key}

		if !resultBuilder.Query.KeysOnly {
			buf := make([]byte, len(iterator.Value()))
			copy(buf, iterator.Value())
			entry.Value = buf
		}

		select {
		case resultBuilder.Output <- Result{Entry: entry}:
		case <-worker.Closing():
			break
		}
	}

	if err := iterator.Error(); err != nil {
		select {
		case resultBuilder.Output <- Result{Error: err}:
		case <-worker.Closing():
			return
		}
	}
}

func (d *dataStore) DiskUsage() (uint64, error) {
	if d.storePath == "" { // in-mem
		return 0, nil
	}

	var diskUsage uint64

	err := filepath.Walk(d.storePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		diskUsage += uint64(info.Size())
		return nil
	})

	if err != nil {
		return 0, err
	}

	return diskUsage, nil
}

func (d *dataStore) Close() (err error) {
	return d.dataBase.Close()
}

func (d *dataStore) IsThreadSafe() {}

/*****************************************************************
*
*		Batch interface and implements.
*
*****************************************************************/
type dataStore struct {
	dataBase  	*leveldb.DB
	storePath 	string
}

func (d *dataStore) Put(key Key, value []byte) (err error) {
	return d.dataBase.Put(key.Bytes(), value, nil)
}

func (d *dataStore) Get(key Key) (value []byte, err error) {
	value, err = d.dataBase.Get(key.Bytes(), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return value, nil
}

func (d *dataStore) Has(key Key) (exists bool, err error) {
	return d.dataBase.Has(key.Bytes(), nil)
}

func (d *dataStore) Delete(key Key) (err error) {

	exists, err := d.dataBase.Has(key.Bytes(), nil)
	if !exists {
		return ErrNotFound
	} else if err != nil {
		return err
	}
	return d.dataBase.Delete(key.Bytes(), nil)
}

func (d *dataStore) Query(q Query) (Results, error) {
	return d.QueryNew(q)
}

func (d *dataStore) Batch() (Batch, error) {

	return &levelDBBatch{
		batch:    	new(leveldb.Batch),
		database: 	d.dataBase,
	}, nil
}
/*****************************************************************
*
*		Batch interface and implements.
*
*****************************************************************/

type Batch interface {

	Put(key Key, val []byte) error

	Delete(key Key) error

	Commit() error
}

type levelDBBatch struct {
	batch    	*leveldb.Batch
	database 	*leveldb.DB
}

func (b *levelDBBatch) Put(key Key, value []byte) error {
	b.batch.Put(key.Bytes(), value)
	return nil
}

func (b *levelDBBatch) Commit() error {
	return b.database.Write(b.batch, nil)
}

func (b *levelDBBatch) Delete(key Key) error {
	b.batch.Delete(key.Bytes())
	return nil
}