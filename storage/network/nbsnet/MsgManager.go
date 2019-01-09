package nbsnet

import (
	"github.com/NBSChain/go-nbs/storage/network/pb"
)

const (
	NatMsgBase         = net_pb.MsgType_NatBase
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
	NatBlankKA         = net_pb.MsgType_NatBlankKA
	NatCheckNetType    = net_pb.MsgType_NatCheckNetType
	NatFindPubIpSyn    = net_pb.MsgType_NatFindPubIpSyn
	NatFindPubIpACK    = net_pb.MsgType_NatFindPubIpACK
	NatBlankKAACK      = net_pb.MsgType_NatBlankKAACK
	NatEnd             = net_pb.MsgType_NatEnd

	GspBase        = net_pb.MsgType_GspBase
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
	GspInnerBase   = net_pb.MsgType_GspInnerBase
	GspEnd         = net_pb.MsgType_GspEnd

	SigIpSigPort = net_pb.NetWorkType_SigIpSigPort
	MulIpSigPort = net_pb.NetWorkType_MulIpSigPort
	MultiPort    = net_pb.NetWorkType_MultiPort
)
