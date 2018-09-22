package dataStore

import (
	"context"
	"errors"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"sync"
)

var instance 		*ServiceRoutingMap
var once 		sync.Once
var parentContext 	context.Context
var logger 		= utils.GetLogInstance()
var ErrNotFound		= errors.New("dataStore: key not found")
var ErrInvalidType 	= errors.New("dataStore: invalid type error")
var ErrCommit		= errors.New("dataStore: batch commit failed")

func GetServiceDispatcher() *ServiceRoutingMap {

	once.Do(func() {
		parentContext = context.Background()
		mounts, err := newDispatcher()

		if err != nil {
			panic(err)
		}

		logger.Info("data store service start to run......\n")
		instance = mounts
	})

	return instance
}

type ServiceRoutingMap struct {
	serviceRouter	map[ServiceType]DataStore
}

//TODO:: Configurable this mount settings later.
func newDispatcher() (*ServiceRoutingMap, error) {

	serviceMap := new(ServiceRoutingMap)

	levelDbMount, err := newLevelDB( &opt.Options{
		Filter: filter.NewBloomFilter(10),
	})

	if err != nil{
		return nil ,err
	}
	serviceMap.addService(NewServiceKey(RootServiceURL), levelDbMount)


	flatFileDs, err := newFlatFileDataStore()
	if err != nil{
		return nil, err
	}
	serviceMap.addService(NewServiceKey(BLOCKServiceURL), flatFileDs)

	return serviceMap, nil
}

func (service *ServiceRoutingMap) addService(key ServiceType, item DataStore)  {
	service.serviceRouter[key] = item
}

func (service *ServiceRoutingMap) GetService(key string)  DataStore{

	serviceKey := NewServiceKey(key)

	return service.serviceRouter[serviceKey]
}
