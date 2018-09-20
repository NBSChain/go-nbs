package dataStore

import (
	"context"
	"errors"
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

var instance 		*MountDataStore
var once 		sync.Once
var parentContext 	context.Context
var logger 		= utils.GetLogInstance()
var ErrNotFound		= errors.New("datastore: key not found")
var ErrInvalidType 	= errors.New("datastore: invalid type error")

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
	mounts 	[]Mount
}

//TODO:: Configurable this mount settings later.
func newMount() (*MountDataStore, error) {

	mounts := make([]Mount, 2)

	levelDbMount, err := newLevelDB( &opt.Options{
		Filter: filter.NewBloomFilter(10),
	})

	if err != nil{
		return nil ,err
	}

	mounts[0] = *levelDbMount


	flatFileMount, err := newFlatFileDataStore("", nil)
	if err != nil{
		return nil, err
	}

	mounts[1] = *flatFileMount


	sort.Slice(mounts, func(i, j int) bool { return mounts[i].prefix.String() > mounts[j].prefix.String() })

	return &MountDataStore{mounts: mounts}, nil
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