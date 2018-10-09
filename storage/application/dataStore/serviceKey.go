package dataStore

import "strings"
type ServiceType int32

const (
	ServiceTypeROOT  ServiceType = 0x01
	ServiceTypeBlock ServiceType = 0x0101
)

var ServiceDictionary  = map[string]*RouterInfo{
	"/":{
		key:       "/",
		code:      ServiceTypeROOT,
		subRouter: nil,
	},
	"blocks":{
		key:       "blocks",
		code:      ServiceTypeBlock,
		subRouter: nil,
	},
}

type RouterInfo struct {
	key		string
	code		ServiceType
	subRouter	map[string]*RouterInfo
}

//TODO::check first level right now.
func NewServiceKey(fullPath string) ServiceType{

	if len(fullPath) <= 1{
		return ServiceTypeROOT
	}

	if fullPath[0] == '/' {
		fullPath = fullPath[1:]
	}

	routers := strings.Split(fullPath, "/")

	/*TODO::
	*Solve top level right now, I didn't find the necessary
	*of complicate service routing.
	*/

	topLevelServiceKey := routers[0]
	service := ServiceDictionary[topLevelServiceKey]
	if service != nil{
		return service.code
	}

	return ServiceTypeROOT
}
