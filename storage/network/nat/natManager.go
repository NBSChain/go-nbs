package nat

import (
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/utils"
	"net"
	"sync"
)

var logger = utils.GetLogInstance()

const NetIoBufferSize = 1 << 11
const BootStrapNatServerTimeOutInSec = 6

type NbsNatManager struct {
	sync.Mutex
	natServer     *net.UDPConn
	natType       nat_pb.NatType
	publicAddress *net.UDPAddr
	privateIP     string
	networkId     string
}
