package nat

import (
	"github.com/NBSChain/go-nbs/utils"
	"net"
)

var logger = utils.GetLogInstance()

const NetIoBufferSize = 1 << 11
const BootStrapNatServerTimeOutInSec = 4

type NbsNatManager struct {
	natServer *net.UDPConn
	networkId string
	canServe  chan bool
}

type NbsAddress struct {
	PublicIP  string
	PrivateIp string
	IsInPub   bool
}
