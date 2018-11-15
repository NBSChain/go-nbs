package memership

import "github.com/NBSChain/go-nbs/thirdParty/gossip/pb"

func (node *MemManager) proxyTheInitSub(request *pb.InitSub) {
	node.randomContact()
}
