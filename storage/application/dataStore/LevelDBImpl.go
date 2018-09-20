package dataStore

import (
	"github.com/jbenet/goprocess"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/storage"
	"github.com/syndtr/goleveldb/leveldb/util"
	"os"
	"path/filepath"
)

func newLevelDB(path string, opts *Options) (*dataStore, error) {
	var nopts opt.Options
	if opts != nil {
		nopts = opt.Options(*opts)
	}

	var err error
	var db *leveldb.DB

	if path == "" {
		db, err = leveldb.Open(storage.NewMemStorage(), &nopts)
	} else {
		db, err = leveldb.OpenFile(path, &nopts)
		if errors.IsCorrupted(err) && !nopts.GetReadOnly() {
			db, err = leveldb.RecoverFile(path, &nopts)
		}
	}

	if err != nil {
		return nil, err
	}

	return &dataStore{
		DataBase: db,
		path:     path,
	}, nil
}

/*****************************************************************
*
*		Batch interface and implements.
*
*****************************************************************/
type dataStore struct {
	DataBase *leveldb.DB
	path     string
}

func (d *dataStore) Put(key Key, value []byte) (err error) {
	return d.DataBase.Put(key.Bytes(), value, nil)
}

func (d *dataStore) Get(key Key) (value []byte, err error) {
	val, err := d.DataBase.Get(key.Bytes(), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, err
		}
		return nil, err
	}
	return val, nil
}

func (d *dataStore) Has(key Key) (exists bool, err error) {
	return d.DataBase.Has(key.Bytes(), nil)
}

func (d *dataStore) Delete(key Key) (err error) {

	exists, err := d.DataBase.Has(key.Bytes(), nil)
	if !exists {
		return leveldb.ErrNotFound
	} else if err != nil {
		return err
	}
	return d.DataBase.Delete(key.Bytes(), nil)
}

func (d *dataStore) Query(q Query) (Results, error) {
	return d.QueryNew(q)
}

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
	i := d.DataBase.NewIterator(rnge, nil)
	return ResultsFromIterator(q, Iterator{
		Next: func() (Result, bool) {
			ok := i.Next()
			if !ok {
				return Result{}, false
			}
			k := string(i.Key())
			e := Entry{Key: k}

			if !q.KeysOnly {
				buf := make([]byte, len(i.Value()))
				copy(buf, i.Value())
				e.Value = buf
			}
			return Result{Entry: e}, true
		},
		Close: func() error {
			i.Release()
			return nil
		},
	}), nil
}

func (d *dataStore) QueryOrig(q Query) (Results, error) {
	// we can use multiple iterators concurrently. see:
	// https://godoc.org/github.com/syndtr/goleveldb/leveldb#DB.NewIterator
	// advance the iterator only if the reader reads
	//
	// run query in own sub-process tied to Results.Process(), so that
	// it waits for us to finish AND so that clients can signal to us
	// that resources should be reclaimed.
	qrb := NewResultBuilder(q)
	qrb.Process.Go(func(worker goprocess.Process) {
		d.runQuery(worker, qrb)
	})

	// go wait on the worker (without signaling close)
	go qrb.Process.CloseAfterChildren()

	// Now, apply remaining things (filters, order)
	qr := qrb.Results()
	for _, f := range q.Filters {
		qr = NaiveFilter(qr, f)
	}
	for _, o := range q.Orders {
		qr = NaiveOrder(qr, o)
	}
	return qr, nil
}

func (d *dataStore) runQuery(worker goprocess.Process, qrb *ResultBuilder) {

	var rnge *util.Range
	if qrb.Query.Prefix != "" {
		rnge = util.BytesPrefix([]byte(qrb.Query.Prefix))
	}
	i := d.DataBase.NewIterator(rnge, nil)
	defer i.Release()

	// advance iterator for offset
	if qrb.Query.Offset > 0 {
		for j := 0; j < qrb.Query.Offset; j++ {
			i.Next()
		}
	}

	// iterate, and handle limit, too
	for sent := 0; i.Next(); sent++ {
		// end early if we hit the limit
		if qrb.Query.Limit > 0 && sent >= qrb.Query.Limit {
			break
		}

		k := string(i.Key())
		e := Entry{Key: k}

		if !qrb.Query.KeysOnly {
			buf := make([]byte, len(i.Value()))
			copy(buf, i.Value())
			e.Value = buf
		}

		select {
		case qrb.Output <- Result{Entry: e}: // we sent it out
		case <-worker.Closing(): // client told us to end early.
			break
		}
	}

	if err := i.Error(); err != nil {
		select {
		case qrb.Output <- Result{Error: err}: // client read our error
		case <-worker.Closing(): // client told us to end.
			return
		}
	}
}

// DiskUsage returns the current disk size used by this levelDB.
// For in-mem datastores, it will return 0.
func (d *dataStore) DiskUsage() (uint64, error) {
	if d.path == "" { // in-mem
		return 0, nil
	}

	var du uint64

	err := filepath.Walk(d.path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		du += uint64(info.Size())
		return nil
	})

	if err != nil {
		return 0, err
	}

	return du, nil
}

// LevelDB needs to be closed.
func (d *dataStore) Close() (err error) {
	return d.DataBase.Close()
}

func (d *dataStore) IsThreadSafe() {}

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

type leveldbBatch struct {
	b  *leveldb.Batch
	db *leveldb.DB
}

func (d *dataStore) Batch() (Batch, error) {
	return &leveldbBatch{
		b:  new(leveldb.Batch),
		db: d.DataBase,
	}, nil
}

func (b *leveldbBatch) Put(key Key, value []byte) error {
	b.b.Put(key.Bytes(), value)
	return nil
}

func (b *leveldbBatch) Commit() error {
	return b.db.Write(b.b, nil)
}

func (b *leveldbBatch) Delete(key Key) error {
	b.b.Delete(key.Bytes())
	return nil
}