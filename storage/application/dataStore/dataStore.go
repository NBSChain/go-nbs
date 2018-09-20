package dataStore

import (
	"context"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"sort"
	"sync"
)

/*****************************************************************
*
*		DataStore interface and implements.
*
*****************************************************************/
type DataStore interface {
	Put(key Key, value []byte) error

	Get(key Key) (value []byte, err error)

	Has(key Key) (exists bool, err error)

	Delete(key Key) error

	Query(q Query) (Results, error)

	Batch() (Batch, error)
}

type Options opt.Options

var instance 		*MountDataStore
var once 		sync.Once
var parentContext 	context.Context
var logger 		= utils.GetLogInstance()

func GetDsInstance() DataStore {
	once.Do(func() {
		parentContext = context.Background()
		mounts, err := newMount()

		if err != nil {
			panic(err)
		}
		logger.Info("data store service start to run......\n")
		instance = mounts
	})

	return instance
}


type Mount struct {
	prefix    Key
	dataStore DataStore
}
type MountDataStore struct {
	mounts []Mount
}

//TODO:: Configurable this mount settings later.
func newMount() (*MountDataStore, error) {
	m := make([]Mount, 2)

	levelDb, err := newLevelDB("", &Options{
		Filter: filter.NewBloomFilter(10),
	})
	if err != nil{
		return nil ,err
	}
	m[0].dataStore = levelDb
	m[0].prefix = NewKey("/")


	flatFile, err := newFlatFileDataStore("", nil)
	if err != nil{
		return nil, err
	}
	m[1].dataStore = flatFile
	m[1].prefix = NewKey("blocks")

	sort.Slice(m, func(i, j int) bool { return m[i].prefix.String() > m[j].prefix.String() })

	return &MountDataStore{mounts: m}, nil
}



func (fs *MountDataStore) Put(key Key, value []byte) error{
	return nil
}

func (fs *MountDataStore) Get(key Key) (value []byte, err error){
	return nil, nil
}

func (fs *MountDataStore) Has(key Key) (exists bool, err error){
	return false, nil
}

func (fs *MountDataStore) Delete(key Key) error{
	return nil
}

func (fs *MountDataStore) Query(q Query) (Results, error){
	return nil, nil
}

func (fs *MountDataStore) Batch() (Batch, error){
	return nil, nil
}