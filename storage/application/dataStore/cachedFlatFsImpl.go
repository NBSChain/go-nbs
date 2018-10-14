package dataStore

import (
	"sync"
)
/*
//Hope some day I can find the right palce to use this bloom technology cause I love it so much.
//"github.com/AndreasBriese/bbloom"
//const HasBloomFilterSize	= 1 << 22
//const HasBloomFilterHashes	= 7
//bf := bbloom.New(float64(HasBloomFilterSize), float64(HasBloomFilterHashes))
*/


type cachedFlatFsDataStore struct {
	sync.Mutex
	cacheSet  map[string]struct{}
	dataStore DataStore
}

func NewBloomDataStore(ds DataStore) DataStore{

	return &cachedFlatFsDataStore{
		cacheSet:  make(map[string]struct{}),
		dataStore: ds,
	}
}

/*****************************************************************
*
*		DataStore interface and implements.
*
*****************************************************************/
func (bs *cachedFlatFsDataStore) Put(key string, value []byte) (err error){

	bs.Lock()
	defer bs.Unlock()
	if err = bs.dataStore.Put(key, value); err != nil{
		bs.cacheSet[key] = struct {}{}
	}

	return err
}

func (bs *cachedFlatFsDataStore) Get(key string) ([]byte, error){

	data, err := bs.dataStore.Get(key)

	if _, has := bs.cacheSet[key]; !has && err == nil {
		bs.Lock()
		bs.cacheSet[key] = struct{}{}
		bs.Unlock()
	}

	return data, err
}

func (bs *cachedFlatFsDataStore) Has(key string) (bool, error){

	if _, has := bs.cacheSet[key]; has{
		return true, nil
	}

	has, err := bs.dataStore.Has(key)
	if has {
		bs.Lock()
		bs.cacheSet[key] = struct{}{}
		bs.Unlock()
	}

	return has,err
}

func (bs *cachedFlatFsDataStore) Delete(key string) (err error){

	bs.Lock()
	defer bs.Unlock()
	if err = bs.dataStore.Delete(key); err == nil{
		delete(bs.cacheSet ,key)
	}

	return nil
}

func (bs *cachedFlatFsDataStore) Query(q Query) (Results, error){
	return bs.dataStore.Query(q)
}

func (bs *cachedFlatFsDataStore) Batch() (Batch, error){

	b, err := bs.dataStore.Batch()
	if err != nil{
		return nil, err
	}

	return &bloomBatch{
		parentCache: bs.cacheSet,
		tempCache:   make(map[string]struct{}),
		batch:       b,
	}, nil
}

/*****************************************************************
*
*		Batch interface and implements.
*
*****************************************************************/
type bloomBatch struct {
	batch       Batch
	parentCache map[string]struct{}
	tempCache   map[string]struct{}
}

func (bb *bloomBatch) Put(key string, val []byte) error{

	bb.tempCache[key] = struct{}{}
	return bb.batch.Put(key, val)
}

func (bb *bloomBatch) Delete(key string) error{
	delete(bb.tempCache, key)
	return bb.batch.Delete(key)
}

func (bb *bloomBatch) Commit() error{

	if err := bb.batch.Commit(); err != nil{
		return err
	}

	for key := range bb.tempCache{
		bb.parentCache[key] = struct{}{}
	}

	return nil
}
