package nbsnet

import (
	"github.com/NBSChain/go-nbs/storage/network/pb"
)

const (
	NatBootReg = net_pb.MsgType_NatBootReg

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

	GspSub         = net_pb.MsgType_GspSub
	GspSubACK      = net_pb.MsgType_GspSubACK
	GspVoteContact = net_pb.MsgType_GspVoteContact
	GspVoteResult  = net_pb.MsgType_GspVoteResult
	GspVoteResAck  = net_pb.MsgType_GspVoteResAck
	GspIntroduce   = net_pb.MsgType_GspIntroduce
	GspWelcome     = net_pb.MsgType_GspWelcome
	GspHeartBeat   = net_pb.MsgType_GspHeartBeat
	GspReplaceArc  = net_pb.MsgType_GspReplaceArc
	GspReplaceAck  = net_pb.MsgType_GspReplaceAck
	GspRemoveIVArc = net_pb.MsgType_GspRemoveIVArc
	GspRemoveOVAcr = net_pb.MsgType_GspRemoveOVAcr
	GspUpdateOVWei = net_pb.MsgType_GspUpdateOVWei
	GspUpdateIVWei = net_pb.MsgType_GspUpdateIVWei
)
