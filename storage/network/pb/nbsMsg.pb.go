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
	MsgType_GspInitSub         MsgType = 2001
	MsgType_GspInitSubACK      MsgType = 2002
	MsgType_GspRegContact      MsgType = 2003
	MsgType_GspContactAck      MsgType = 2004
	MsgType_GspForwardSub      MsgType = 2005
	MsgType_GspSubAck          MsgType = 2006
	MsgType_GspHeartBeat       MsgType = 2007
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
	2001: "GspInitSub",
	2002: "GspInitSubACK",
	2003: "GspRegContact",
	2004: "GspContactAck",
	2005: "GspForwardSub",
	2006: "GspSubAck",
	2007: "GspHeartBeat",
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
	"GspInitSub":         2001,
	"GspInitSubACK":      2002,
	"GspRegContact":      2003,
	"GspContactAck":      2004,
	"GspForwardSub":      2005,
	"GspSubAck":          2006,
	"GspHeartBeat":       2007,
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
	// 294 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0xd0, 0x4b, 0x4e, 0x42, 0x31,
	0x14, 0x06, 0x60, 0x63, 0x22, 0x0d, 0x15, 0xbd, 0xc7, 0xc6, 0xc4, 0x3d, 0x38, 0x60, 0xe2, 0x0a,
	0x2e, 0xa0, 0x48, 0x08, 0x57, 0x02, 0x6e, 0xa0, 0x17, 0x8f, 0x4d, 0x03, 0xb6, 0x4d, 0x7b, 0x80,
	0xb0, 0x4c, 0xdf, 0x6b, 0xf0, 0x19, 0x9f, 0x73, 0x53, 0xb1, 0x31, 0x3a, 0xfd, 0x72, 0xf2, 0xff,
	0xfd, 0xcb, 0x6b, 0xa6, 0x0c, 0xbd, 0xa0, 0xea, 0xce, 0x5b, 0xb2, 0xa2, 0x62, 0x90, 0xea, 0xae,
	0xdc, 0x7d, 0x5a, 0xe5, 0xac, 0x17, 0xd4, 0xf1, 0xc2, 0xa1, 0xa8, 0xf2, 0xb5, 0x7d, 0xef, 0xad,
	0x87, 0x15, 0x91, 0x71, 0x5e, 0x48, 0x6a, 0x58, 0x4b, 0x03, 0x54, 0x70, 0xc7, 0xc4, 0x16, 0xaf,
	0x15, 0x92, 0xba, 0x88, 0x2e, 0x9f, 0xe8, 0x19, 0xc2, 0x3d, 0x13, 0xc0, 0xd7, 0x0b, 0x49, 0x2d,
	0xad, 0x72, 0xe7, 0x26, 0x0b, 0x78, 0x48, 0xd2, 0xd7, 0x46, 0xf5, 0xad, 0x51, 0xf0, 0xc8, 0xc4,
	0x26, 0xaf, 0x2e, 0x6f, 0x8e, 0xa6, 0x04, 0xcf, 0x4c, 0x6c, 0xf3, 0xac, 0x90, 0x34, 0xc0, 0x19,
	0xfa, 0xd0, 0x31, 0x33, 0x4d, 0x08, 0x2f, 0x4c, 0xec, 0x70, 0xf1, 0x4f, 0xf3, 0xd1, 0x18, 0x5e,
	0x53, 0x6b, 0xdf, 0xeb, 0x96, 0x56, 0xc3, 0x85, 0x81, 0xb7, 0xbf, 0x14, 0xaf, 0xde, 0x99, 0x10,
	0x7c, 0x63, 0x59, 0xd2, 0xb4, 0xe6, 0x54, 0xfb, 0x33, 0xf8, 0x48, 0x16, 0x07, 0xe4, 0x26, 0xcc,
	0xd1, 0xc3, 0x27, 0x8b, 0xa3, 0xda, 0xc1, 0x75, 0x8c, 0xa6, 0xe1, 0xb4, 0x84, 0xf3, 0x2c, 0x1e,
	0xfd, 0x42, 0xde, 0xec, 0xc2, 0x45, 0xb2, 0x01, 0xc6, 0x30, 0x92, 0x23, 0x82, 0xcb, 0x64, 0x3f,
	0x10, 0x4b, 0xaf, 0x92, 0x1d, 0x58, 0x3f, 0x97, 0xfe, 0x24, 0xe6, 0x5d, 0x67, 0x71, 0x6d, 0x3b,
	0xb8, 0x98, 0x35, 0x1a, 0xc3, 0x4d, 0x16, 0xdf, 0xda, 0x0e, 0xee, 0x10, 0xa5, 0xa7, 0x06, 0x4a,
	0x82, 0xdb, 0xac, 0xac, 0x7c, 0x7f, 0xff, 0xde, 0x57, 0x00, 0x00, 0x00, 0xff, 0xff, 0x62, 0x82,
	0xcd, 0x2f, 0x8e, 0x01, 0x00, 0x00,
}