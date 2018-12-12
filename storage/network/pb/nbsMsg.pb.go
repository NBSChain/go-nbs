// Code generated by protoc-gen-go. DO NOT EDIT.
// source: nbsMsg.proto

package net_pb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type MsgType int32

const (
	MsgType_Error              MsgType = 0
	MsgType_NatBootReg         MsgType = 1001
	MsgType_NatKeepAlive       MsgType = 1002
	MsgType_NatDigApply        MsgType = 1003
	MsgType_NatPingPong        MsgType = 1004
	MsgType_NatDigOut          MsgType = 1006
	MsgType_NatReversInvite    MsgType = 1008
	MsgType_NatReversInviteAck MsgType = 1009
	MsgType_NatPriDigSyn       MsgType = 1010
	MsgType_NatPriDigAck       MsgType = 1011
	MsgType_NatDigConfirm      MsgType = 1012
	MsgType_NatBootAnswer      MsgType = 1013
	MsgType_GspSub             MsgType = 2001
	MsgType_GspSubACK          MsgType = 2002
	MsgType_GspVoteContact     MsgType = 2003
	MsgType_GspVoteResult      MsgType = 2004
	MsgType_GspVoteResAck      MsgType = 2005
	MsgType_GspIntroduce       MsgType = 2006
	MsgType_GspWelcome         MsgType = 2007
	MsgType_GspHeartBeat       MsgType = 2008
	MsgType_GspReplaceArc      MsgType = 2009
	MsgType_GspReplaceAck      MsgType = 2010
	MsgType_GspRemoveIVArc     MsgType = 2011
	MsgType_GspRemoveOVAcr     MsgType = 2012
	MsgType_GspUpdateOVWei     MsgType = 2014
	MsgType_GspUpdateIVWei     MsgType = 2015
)

var MsgType_name = map[int32]string{
	0:    "Error",
	1001: "NatBootReg",
	1002: "NatKeepAlive",
	1003: "NatDigApply",
	1004: "NatPingPong",
	1006: "NatDigOut",
	1008: "NatReversInvite",
	1009: "NatReversInviteAck",
	1010: "NatPriDigSyn",
	1011: "NatPriDigAck",
	1012: "NatDigConfirm",
	1013: "NatBootAnswer",
	2001: "GspSub",
	2002: "GspSubACK",
	2003: "GspVoteContact",
	2004: "GspVoteResult",
	2005: "GspVoteResAck",
	2006: "GspIntroduce",
	2007: "GspWelcome",
	2008: "GspHeartBeat",
	2009: "GspReplaceArc",
	2010: "GspReplaceAck",
	2011: "GspRemoveIVArc",
	2012: "GspRemoveOVAcr",
	2014: "GspUpdateOVWei",
	2015: "GspUpdateIVWei",
}

var MsgType_value = map[string]int32{
	"Error":              0,
	"NatBootReg":         1001,
	"NatKeepAlive":       1002,
	"NatDigApply":        1003,
	"NatPingPong":        1004,
	"NatDigOut":          1006,
	"NatReversInvite":    1008,
	"NatReversInviteAck": 1009,
	"NatPriDigSyn":       1010,
	"NatPriDigAck":       1011,
	"NatDigConfirm":      1012,
	"NatBootAnswer":      1013,
	"GspSub":             2001,
	"GspSubACK":          2002,
	"GspVoteContact":     2003,
	"GspVoteResult":      2004,
	"GspVoteResAck":      2005,
	"GspIntroduce":       2006,
	"GspWelcome":         2007,
	"GspHeartBeat":       2008,
	"GspReplaceArc":      2009,
	"GspReplaceAck":      2010,
	"GspRemoveIVArc":     2011,
	"GspRemoveOVAcr":     2012,
	"GspUpdateOVWei":     2014,
	"GspUpdateIVWei":     2015,
}

func (x MsgType) String() string {
	return proto.EnumName(MsgType_name, int32(x))
}

func (MsgType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_1ac27572d1ea9f39, []int{0}
}

func init() {
	proto.RegisterEnum("net.pb.MsgType", MsgType_name, MsgType_value)
}

func init() { proto.RegisterFile("nbsMsg.proto", fileDescriptor_1ac27572d1ea9f39) }

var fileDescriptor_1ac27572d1ea9f39 = []byte{
	// 368 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0xd2, 0xcb, 0x6e, 0xd4, 0x30,
	0x14, 0x06, 0x60, 0x24, 0xc4, 0x58, 0x75, 0x4b, 0x7d, 0x30, 0x48, 0xbc, 0x03, 0x8b, 0x6e, 0x78,
	0x82, 0x74, 0x8a, 0x42, 0x54, 0x35, 0x33, 0x4a, 0x21, 0x5d, 0x3b, 0xee, 0x21, 0xb2, 0x26, 0x63,
	0x5b, 0xf6, 0x49, 0xd0, 0x3c, 0x26, 0xf7, 0xab, 0x60, 0xcd, 0x45, 0x88, 0xfb, 0x1a, 0x79, 0x86,
	0x2c, 0x66, 0x96, 0xfe, 0xf4, 0xeb, 0x9c, 0xdf, 0x96, 0xf9, 0x81, 0x6d, 0xe2, 0x59, 0x6c, 0x8f,
	0x7c, 0x70, 0xe4, 0xe4, 0xc4, 0x22, 0x1d, 0xf9, 0xe6, 0xce, 0xdf, 0xab, 0x9c, 0x9d, 0xc5, 0xf6,
	0xc1, 0xca, 0xa3, 0xdc, 0xe3, 0xd7, 0xee, 0x85, 0xe0, 0x02, 0x5c, 0x91, 0x82, 0xf3, 0x52, 0xd1,
	0xb1, 0x73, 0x54, 0x61, 0x0b, 0x9f, 0x98, 0xbc, 0xc1, 0x0f, 0x4a, 0x45, 0xa7, 0x88, 0x3e, 0xeb,
	0xcc, 0x80, 0xf0, 0x99, 0x49, 0xe0, 0xfb, 0xa5, 0xa2, 0x13, 0xd3, 0x66, 0xde, 0x77, 0x2b, 0xf8,
	0x32, 0xca, 0xdc, 0xd8, 0x76, 0xee, 0x6c, 0x0b, 0x5f, 0x99, 0x3c, 0xe4, 0x7b, 0x9b, 0xcc, 0xac,
	0x27, 0xf8, 0xc6, 0xe4, 0x2d, 0x2e, 0x4a, 0x45, 0x15, 0x0e, 0x18, 0x62, 0x61, 0x07, 0x43, 0x08,
	0xdf, 0x99, 0xbc, 0xcd, 0xe5, 0x8e, 0x66, 0x7a, 0x01, 0x3f, 0xc6, 0xad, 0xf3, 0x60, 0x4e, 0x4c,
	0x7b, 0xbe, 0xb2, 0xf0, 0x73, 0x9b, 0x52, 0xea, 0x17, 0x93, 0x92, 0x5f, 0xdf, 0x2c, 0x99, 0x3a,
	0xfb, 0xc8, 0x84, 0x25, 0xfc, 0x1e, 0x2d, 0x5d, 0x20, 0xb3, 0xf1, 0x31, 0x06, 0xf8, 0xc3, 0xe4,
	0x3e, 0x9f, 0xe4, 0xd1, 0x9f, 0xf7, 0x0d, 0x3c, 0x11, 0xa9, 0xd9, 0xe6, 0x90, 0x4d, 0x4f, 0xe1,
	0xa9, 0x90, 0x37, 0xf9, 0x61, 0x1e, 0x7d, 0xed, 0x08, 0xa7, 0xce, 0x92, 0xd2, 0x04, 0xcf, 0x44,
	0x9a, 0xf2, 0x1f, 0x2b, 0x8c, 0x7d, 0x47, 0xf0, 0x7c, 0xc7, 0x52, 0x83, 0x17, 0x22, 0x95, 0xca,
	0xa3, 0x2f, 0x2c, 0x05, 0x77, 0xd9, 0x6b, 0x84, 0x97, 0x22, 0xbd, 0x60, 0x1e, 0xfd, 0x05, 0x76,
	0xda, 0x2d, 0x11, 0x5e, 0x8d, 0x99, 0xfb, 0xa8, 0x02, 0x1d, 0xa3, 0x22, 0x78, 0x3d, 0x8e, 0xaa,
	0xd0, 0x77, 0x4a, 0x63, 0x16, 0x34, 0xbc, 0xd9, 0x35, 0xbd, 0x80, 0xb7, 0x63, 0xb7, 0x0a, 0x97,
	0x6e, 0xc0, 0xa2, 0x4e, 0xc1, 0x77, 0xdb, 0x38, 0xab, 0x33, 0x1d, 0xe0, 0xfd, 0x88, 0x0f, 0xfd,
	0xa5, 0x22, 0x9c, 0xd5, 0x17, 0x68, 0xe0, 0xc3, 0x36, 0x16, 0x6b, 0xfc, 0x28, 0x9a, 0xc9, 0xfa,
	0x1f, 0xdc, 0xfd, 0x17, 0x00, 0x00, 0xff, 0xff, 0x30, 0xde, 0x58, 0xef, 0x17, 0x02, 0x00, 0x00,
}
