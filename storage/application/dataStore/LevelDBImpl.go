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

func newLevelDB(opts *opt.Options) (DataStore, error) {

	var path = utils.GetConfig().LevelDBDir
	var err error
	var dataBase *leveldb.DB

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

	dateStore := &levelDBDataStore{
		dataBase:  dataBase,
		storePath: path,
	}

	return dateStore, nil
}

/*****************************************************************
*
*		The functions below are implements of interface
*
*****************************************************************/
//TODO:: What's the interface of following functions.
func (d *levelDBDataStore) QueryNew(query Query) (Results, error) {

	if len(query.Filters) > 0 ||
		len(query.Orders) > 0 ||
		query.Limit > 0 ||
		query.Offset > 0 {
		return d.QueryOrig(query)
	}
	var rnge *util.Range
	if query.Prefix != "" {
		rnge = util.BytesPrefix([]byte(query.Prefix))
	}

	iterator := d.dataBase.NewIterator(rnge, nil)

	return ResultsFromIterator(query, Iterator{
		Next: func() (Result, bool) {
			ok := iterator.Next()
			if !ok {
				return Result{}, false
			}

			key := string(iterator.Key())
			entry := Entry{Key: key}

			if !query.KeysOnly {
				buf := make([]byte, len(iterator.Value()))
				copy(buf, iterator.Value())
				entry.Value = buf
			}

			return Result{entry: entry}, true
		},
		Close: func() error {
			iterator.Release()
			return nil
		},
	}), nil
}

func (d *levelDBDataStore) QueryOrig(query Query) (Results, error) {

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

func (d *levelDBDataStore) runQuery(worker goprocess.Process, resultBuilder *ResultBuilder) {

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

		key := string(iterator.Key())
		entry := Entry{Key: key}

		if !resultBuilder.Query.KeysOnly {
			buf := make([]byte, len(iterator.Value()))
			copy(buf, iterator.Value())
			entry.Value = buf
		}

		select {
		case resultBuilder.Output <- Result{entry: entry}:
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

func (d *levelDBDataStore) DiskUsage() (uint64, error) {
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

func (d *levelDBDataStore) Close() (err error) {
	return d.dataBase.Close()
}

func (d *levelDBDataStore) IsThreadSafe() {}

/*****************************************************************
*
*		Batch interface and implements.
*
*****************************************************************/
type levelDBDataStore struct {
	dataBase  *leveldb.DB
	storePath string
}

func (d *levelDBDataStore) Put(key string, value []byte) (err error) {
	return d.dataBase.Put([]byte(key), value, nil)
}

func (d *levelDBDataStore) Get(key string) (value []byte, err error) {
	value, err = d.dataBase.Get([]byte(key), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return value, nil
}

func (d *levelDBDataStore) Has(key string) (exists bool, err error) {
	return d.dataBase.Has([]byte(key), nil)
}

func (d *levelDBDataStore) Delete(key string) (err error) {

	exists, err := d.dataBase.Has([]byte(key), nil)
	if !exists {
		return ErrNotFound
	} else if err != nil {
		return err
	}

	return d.dataBase.Delete([]byte(key), nil)
}

func (d *levelDBDataStore) Query(query Query) (Results, error) {
	return d.QueryNew(query)
}

func (d *levelDBDataStore) Batch() (Batch, error) {

	return &levelDBBatch{
		batch:    new(leveldb.Batch),
		database: d.dataBase,
	}, nil
}

/*****************************************************************
*
*		Batch interface and implements.
*
*****************************************************************/

type levelDBBatch struct {
	batch    *leveldb.Batch
	database *leveldb.DB
}

func (b *levelDBBatch) Put(key string, value []byte) error {
	b.batch.Put([]byte(key), value)
	return nil
}

func (b *levelDBBatch) Commit() error {
	return b.database.Write(b.batch, nil)
}

func (b *levelDBBatch) Delete(key string) error {
	b.batch.Delete([]byte(key))
	return nil
}
