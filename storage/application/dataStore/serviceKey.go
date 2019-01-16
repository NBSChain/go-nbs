package dataStore

import "strings"

type ServiceType int32

const RootServiceURL = "/"
const BLOCKServiceURL = "blocks"
const LocalParamKeyURL = "localParam"

const (
	ServiceTypeROOT       ServiceType = 0x01
	ServiceTypeBlock      ServiceType = 0x0101
	ServiceTypeLocalParam ServiceType = 0x0102
)

var ServiceDictionary = map[string]*RouterInfo{
	RootServiceURL: {
		code:      ServiceTypeROOT,
		subRouter: nil,
	},
	BLOCKServiceURL: {
		code:      ServiceTypeBlock,
		subRouter: nil,
	},
	LocalParamKeyURL: {
		code:      ServiceTypeLocalParam,
		subRouter: nil,
	},
}

type RouterInfo struct {
	code      ServiceType
	subRouter map[string]*RouterInfo
}

/*TODO::
check first level right now.
*Solve top level right now, I didn't find the necessary
*of complicate service routing.
*/
func NewServiceKey(fullPath string) ServiceType {

	if len(fullPath) <= 1 {
		return ServiceTypeROOT
	}

	if fullPath[0] == '/' {
		fullPath = fullPath[1:]
	}

	routers := strings.Split(fullPath, "/")

	topLevelServiceKey := routers[0]
	service := ServiceDictionary[topLevelServiceKey]
	if service != nil {
		return service.code
	}

	return ServiceTypeROOT
}
