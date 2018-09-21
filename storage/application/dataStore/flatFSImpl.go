package dataStore

import (
	"errors"
	"fmt"
	"github.com/NBSChain/go-nbs/utils"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
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
	ErrDatastoreExists       = errors.New("dataStore already exists")
	ErrDatastoreDoesNotExist = errors.New("dataStore directory does not exist")
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
	typ  	opT
	key  	Key
	tmp  	string
	path 	string
	v    	[]byte
}

type opMap struct {
	ops sync.Map
}

type opResult struct {
	mu		sync.RWMutex
	success 	bool
	opMap 		*opMap
	name  		string
}

type FlatFileDataStore struct {
	diskUsage 	int64
	path 		string
	shardStr 	string
	getDir   	ShardFunc
	dirty       	bool
	storedValue 	diskUsageValue
	checkpointCh 	chan struct{}
	done         	chan struct{}
	opMap 		*opMap
}


func newFlatFileDataStore() (DataStore, error) {
	var path = utils.GetConfig().BlocksDir
	var shardIdV1 *ShardIdV1

	_, exist := utils.FileExists(path)
	if !exist{
		err := os.Mkdir(path, 0755)
		if err != nil{
			return nil, err
		}
		shardIdV1, err = ParseShardFunc(utils.GetConfig().ShardFun)

		err = WriteShardFunc(path, shardIdV1)
		if err != nil {
			return nil, err
		}
		err = WriteReadme(path, shardIdV1)

		if err != nil{
			return nil, err
		}
	}else{
		var err error
		shardIdV1, err = ReadShardFunc(path)
		if err != nil {
			return nil, err
		}
	}

	fileStore := &FlatFileDataStore{
		path:      	path,
		shardStr:  	shardIdV1.String(),
		getDir:    	shardIdV1.Func(),
		diskUsage: 	0,
		opMap:     	new(opMap),
	}

	err := fileStore.calculateDiskUsage()
	if err != nil {
		return nil, err
	}

	fileStore.checkpointCh = make(chan struct{}, 1)
	fileStore.done = make(chan struct{})
	go fileStore.checkpointLoop()

	return fileStore, nil
}

func (fs *FlatFileDataStore) calculateDiskUsage()  error{
	return nil
}

func (fs *FlatFileDataStore) checkpointLoop() {

	timerActive	:= true
	timer 		:= time.NewTimer(0)

	defer timer.Stop()

	for {
		select {
		case _, more := <-fs.checkpointCh:
			if !more {
			}
			if timerActive{

			}
		case <-timer.C:
			timerActive = false
		}
	}
}

func (fs *FlatFileDataStore) encode(key string) (dir, file string) {

	noSlash	:= key[len(BLOCKServiceURL) + 1:]

	dir 	= filepath.Join(fs.path, fs.getDir(noSlash))

	file 	= filepath.Join(dir, noSlash + extension)

	return dir, file
}

func (fs *FlatFileDataStore) decode(file string) (key string, ok bool) {

	if filepath.Ext(file) != extension {
		return "", false
	}

	name := file[:len(file)-len(extension)]

	return name, true
}

func (fs *FlatFileDataStore) makeDir(dir string) error {
	if _, exist := utils.FileExists(dir); exist{
		return nil
	}

	if err := os.Mkdir(dir, 0755); err != nil {
		return err
	}

	return nil
}


/*****************************************************************
*
*		DataStore interface and implements.
*
*****************************************************************/
func (fs *FlatFileDataStore) Put(key string, value []byte) error{

	dir, path := fs.encode(key)

	if err := fs.makeDir(dir); err != nil {
		return err
	}

	tmp, err := ioutil.TempFile(dir, "put-")
	if err != nil {
		return err
	}

	defer func() {
		tmp.Close()
		os.Remove(tmp.Name())
	}()

	if _, err := tmp.Write(value); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}

	err = os.Rename(tmp.Name(), path)
	if err != nil {
		return err
	}

	return nil
}

func (fs *FlatFileDataStore) Get(key string) (value []byte, err error){
	return nil, nil
}

func (fs *FlatFileDataStore) Has(key string) (exists bool, err error){
	return false, nil
}

func (fs *FlatFileDataStore) Delete(key string) error{
	return nil
}

func (fs *FlatFileDataStore) Query(q Query) (Results, error){
	return nil, nil
}

func (fs *FlatFileDataStore) Batch() (Batch, error){
	return nil, nil
}