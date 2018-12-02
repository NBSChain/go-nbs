package nbsnet

import (
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/gogo/protobuf/proto"
	"time"
)

const (
	_               int32 = iota
	NatError              = 1000
	NatBootReg            = 1001
	NatKeepAlive          = 1002
	NatConnect            = 1003
	NatPingPong           = 1004
	NatDigOut             = 1006
	NatDigSuccess         = 1007
	NatReversDig          = 1008
	NatReversDigAck       = 1009
	NatPriDigSyn          = 1010
	NatPriDigAck          = 1011
	NatConnectACK         = 1012

	GspInitSub    = 2001
	GspInitSubACK = 2002
	GspRegContact = 2003
	GspContactAck = 2004
	GspForwardSub = 2005
	GspSubAck     = 2006
	GspHeartBeat  = 2007
)

func PackNatData(pb proto.Message, typ int32) ([]byte, *net_pb.NatMsg, error) {

	data, err := proto.Marshal(pb)
	if err != nil {
		return nil, nil, err
	}
	msg := &net_pb.NatMsg{
		Typ:     typ,
		Len:     int32(len(data)),
		Seq:     time.Now().Unix(),
		PayLoad: data,
	}
	rawData, err := proto.Marshal(msg)
	if err != nil {
		return nil, msg, err
	}
	return rawData, msg, nil
}

func PackEmptyNatData(typ int32) ([]byte, *net_pb.NatMsg, error) {
	msg := &net_pb.NatMsg{
		Typ: typ,
		Seq: time.Now().Unix(),
	}
	rawData, err := proto.Marshal(msg)
	if err != nil {
		return nil, nil, err
	}

	return rawData, msg, nil
}

func UnpackNatData(b []byte, pb proto.Message) (*net_pb.NatMsg, error) {
	natMsg := &net_pb.NatMsg{}
	if err := proto.Unmarshal(b, natMsg); err != nil {
		return nil, err
	}

	if natMsg.Len == 0 || pb == nil {
		return natMsg, nil
	}

	if err := proto.Unmarshal(natMsg.PayLoad, pb); err != nil {
		return natMsg, err
	}
	return natMsg, nil
}
