package nbsnet

import (
	"github.com/NBSChain/go-nbs/storage/network/pb"
)

const (
	NatBootReg         = net_pb.MsgType_NatBootReg
	NatKeepAlive       = net_pb.MsgType_NatKeepAlive
	NatDigApply        = net_pb.MsgType_NatDigApply
	NatPingPong        = net_pb.MsgType_NatPingPong
	NatDigOut          = net_pb.MsgType_NatDigOut
	NatReversInvite    = net_pb.MsgType_NatReversInvite
	NatReversInviteAck = net_pb.MsgType_NatReversInviteAck
	NatPriDigSyn       = net_pb.MsgType_NatPriDigSyn
	NatPriDigAck       = net_pb.MsgType_NatPriDigAck
	NatDigConfirm      = net_pb.MsgType_NatDigConfirm
	NatBootAnswer      = net_pb.MsgType_NatBootAnswer
	GspInitSub         = net_pb.MsgType_GspInitSub
	GspInitSubACK      = net_pb.MsgType_GspInitSubACK
	GspRegContact      = net_pb.MsgType_GspRegContact
	GspContactAck      = net_pb.MsgType_GspContactAck
	GspForwardSub      = net_pb.MsgType_GspForwardSub
	GspSubAck          = net_pb.MsgType_GspSubAck
	GspHeartBeat       = net_pb.MsgType_GspHeartBeat
)
