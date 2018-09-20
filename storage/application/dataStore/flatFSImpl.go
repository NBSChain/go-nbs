package dataStore

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	extension                  = ".data"
	diskUsageMessageTimeout    = 5 * time.Second
	diskUsageCheckpointPercent = 1.0
	diskUsageCheckpointTimeout = 2 * time.Second
)

var (
	DiskUsageFile 		= "diskUsage.cache"
	DiskUsageFilesAverage 	= 2000
	DiskUsageCalcTimeout 	= 5 * time.Minute
)
var (
	ErrDatastoreExists       = errors.New("datastore already exists")
	ErrDatastoreDoesNotExist = errors.New("datastore directory does not exist")
	ErrShardingFileMissing   = fmt.Errorf("%s file not found in datastore", SHARDING_FN)
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type diskUsageValue struct {
	DiskUsage int64        	`json:"diskUsage"`
	Accuracy  string	`json:"accuracy"`
}

type ShardFunc func(string) string

type opT int

type op struct {
	typ  opT
	key  Key
	tmp  string
	path string
	v    []byte
}

type opMap struct {
	ops sync.Map
}

type opResult struct {
	mu	sync.RWMutex
	success bool
	opMap 	*opMap
	name  	string
}

type FlatFileDataStore struct {
	diskUsage 	int64
	path 		string
	shardStr 	string
	getDir   	ShardFunc
	sync 		bool
	dirty       	bool
	storedValue 	diskUsageValue
	checkpointCh 	chan struct{}
	done         	chan struct{}
	opMap 		*opMap
}


func newFlatFileDataStore(path string, fun *ShardIdV1) (*Mount, error) {

	ds := &FlatFileDataStore{}
	return &Mount{
		prefix: 	NewKey("blocks"),
		dataStore:	ds,
	}, nil
}


func (fs *FlatFileDataStore) Put(key Key, value []byte) error{
	return nil
}

func (fs *FlatFileDataStore) Get(key Key) (value []byte, err error){
	return nil, nil
}

func (fs *FlatFileDataStore) Has(key Key) (exists bool, err error){
	return false, nil
}

func (fs *FlatFileDataStore) Delete(key Key) error{
	return nil
}

func (fs *FlatFileDataStore) Query(q Query) (Results, error){
	return nil, nil
}

func (fs *FlatFileDataStore) Batch() (Batch, error){
	return nil, nil
}