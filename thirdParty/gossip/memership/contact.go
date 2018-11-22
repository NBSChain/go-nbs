package memership

import (
	"github.com/NBSChain/go-nbs/storage/network"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/thirdParty/gossip/pb"
	"github.com/gogo/protobuf/proto"
)

func (node *MemManager) proxyTheInitSub(request *pb.InitSub, applierAddr *nbsnet.NbsUdpAddr) {

	nextNodeId := node.randomContact()

	if node.peerId == nextNodeId {

		msg := &pb.Gossip{
			MessageType: pb.MsgType_reqContractAck,
			ContactRes: &pb.ReqContactACK{
				ApplierId:  request.NodeId,
				SupplierId: nextNodeId,
			},
		}

		//port := utils.GetConfig().GossipCtrlPort
		conn, err := network.GetInstance().Connect(nil, applierAddr)
		if err != nil {
			logger.Error("the contact failed to notify the subscriber:", err)
			return
		}
		defer conn.Close()

		msgData, err := proto.Marshal(msg)
		if err != nil {
			logger.Error("failed to marshal the contact init msg:", err)
			return
		}

		conn.Send(msgData)

	} else {
		//TODO:: forward this request to next contact server.
	}

}

/*
*4 Indirection mechanism for finding a contact node
Upon subscription( s, Counters, newSubscription) of a new subscriber s on a node ni
if ni is the initial contact then
Counters = 2 * Card􏰭PartialViewni 􏰮
􏰁Initialise the length of the walk to reach a random node􏰂
else
if Counters 􏰶􏰰 0 then
􏰁Normalize weight Wi j of n j 􏰃PartialView􏰂 for all n j 􏰃 PartialView do
Wi j 􏰰 Wi j Wout􏰭ni􏰮
end for
Choose n j 􏰃 PartialView with probability Wi j Decrement Counters;
Send(nj, s, Counters, newSubscription);
else
ni acts as the contact node and applies the basic SCAMP algorithm described in algorithm
1
end if end if
*/
func (node *MemManager) randomContact() string {
	//TODO::implement the indirect mechanism .
	return node.peerId
}
