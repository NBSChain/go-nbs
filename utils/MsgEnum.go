package utils

const (
	unkown          int32 = iota
	NatError              = 1000
	NatBootReg            = 1001
	NatKeepAlive          = 1002
	NatConnect            = 1003
	NatPingPong           = 1004
	NatDigIn              = 1005
	NatDigOut             = 1006
	NatDigSuccess         = 1007
	NatReversDig          = 1008
	NatReversDigAck       = 1009
	NatPriDigSyn          = 1010
	NatPriDigAck          = 1011

	GspInitSub    = 2001
	GspInitSubACK = 2002
	GspRegContact = 2003
	GspContactAck = 2004
	GspForwardSub = 2005
	GspSubAck     = 2006
	GspHeartBeat  = 2007
)
