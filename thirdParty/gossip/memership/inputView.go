package memership

import (
	"github.com/NBSChain/go-nbs/thirdParty/gossip/pb"
	"net"
)

func (node *MemManager) subToContract(ack *pb.ReqContactACK, addr *net.UDPAddr) {
	logger.Info("gossip sub start:", ack, addr)
}
