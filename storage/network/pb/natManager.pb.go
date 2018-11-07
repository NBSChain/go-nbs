// Code generated by protoc-gen-go. DO NOT EDIT.
// source: natManager.proto

package nat_pb

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

type NatMsgType int32

const (
	NatMsgType_UnknownReq   NatMsgType = 0
	NatMsgType_BootStrapReg NatMsgType = 1
	NatMsgType_KeepAlive    NatMsgType = 2
	NatMsgType_Connect      NatMsgType = 3
	NatMsgType_Ping         NatMsgType = 4
)

var NatMsgType_name = map[int32]string{
	0: "UnknownReq",
	1: "BootStrapReg",
	2: "KeepAlive",
	3: "Connect",
	4: "Ping",
}

var NatMsgType_value = map[string]int32{
	"UnknownReq":   0,
	"BootStrapReg": 1,
	"KeepAlive":    2,
	"Connect":      3,
	"Ping":         4,
}

func (x NatMsgType) String() string {
	return proto.EnumName(NatMsgType_name, int32(x))
}

func (NatMsgType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_2f5cc5a30f7d4657, []int{0}
}

type NatType int32

const (
	NatType_UnknownRES     NatType = 0
	NatType_NoNatDevice    NatType = 1
	NatType_BehindNat      NatType = 2
	NatType_CanBeNatServer NatType = 3
	NatType_ToBeChecked    NatType = 4
)

var NatType_name = map[int32]string{
	0: "UnknownRES",
	1: "NoNatDevice",
	2: "BehindNat",
	3: "CanBeNatServer",
	4: "ToBeChecked",
}

var NatType_value = map[string]int32{
	"UnknownRES":     0,
	"NoNatDevice":    1,
	"BehindNat":      2,
	"CanBeNatServer": 3,
	"ToBeChecked":    4,
}

func (x NatType) String() string {
	return proto.EnumName(NatType_name, int32(x))
}

func (NatType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_2f5cc5a30f7d4657, []int{1}
}

type BootNatRegReq struct {
	NodeId               string   `protobuf:"bytes,1,opt,name=nodeId,proto3" json:"nodeId,omitempty"`
	PrivatePort          string   `protobuf:"bytes,3,opt,name=privatePort,proto3" json:"privatePort,omitempty"`
	Zone                 string   `protobuf:"bytes,4,opt,name=zone,proto3" json:"zone,omitempty"`
	PrivateIp            string   `protobuf:"bytes,2,opt,name=privateIp,proto3" json:"privateIp,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BootNatRegReq) Reset()         { *m = BootNatRegReq{} }
func (m *BootNatRegReq) String() string { return proto.CompactTextString(m) }
func (*BootNatRegReq) ProtoMessage()    {}
func (*BootNatRegReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f5cc5a30f7d4657, []int{0}
}

func (m *BootNatRegReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BootNatRegReq.Unmarshal(m, b)
}
func (m *BootNatRegReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BootNatRegReq.Marshal(b, m, deterministic)
}
func (m *BootNatRegReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BootNatRegReq.Merge(m, src)
}
func (m *BootNatRegReq) XXX_Size() int {
	return xxx_messageInfo_BootNatRegReq.Size(m)
}
func (m *BootNatRegReq) XXX_DiscardUnknown() {
	xxx_messageInfo_BootNatRegReq.DiscardUnknown(m)
}

var xxx_messageInfo_BootNatRegReq proto.InternalMessageInfo

func (m *BootNatRegReq) GetNodeId() string {
	if m != nil {
		return m.NodeId
	}
	return ""
}

func (m *BootNatRegReq) GetPrivatePort() string {
	if m != nil {
		return m.PrivatePort
	}
	return ""
}

func (m *BootNatRegReq) GetZone() string {
	if m != nil {
		return m.Zone
	}
	return ""
}

func (m *BootNatRegReq) GetPrivateIp() string {
	if m != nil {
		return m.PrivateIp
	}
	return ""
}

type BootNatRegRes struct {
	NatType              NatType  `protobuf:"varint,1,opt,name=natType,proto3,enum=nat.pb.NatType" json:"natType,omitempty"`
	NodeId               string   `protobuf:"bytes,2,opt,name=nodeId,proto3" json:"nodeId,omitempty"`
	PublicIp             string   `protobuf:"bytes,3,opt,name=publicIp,proto3" json:"publicIp,omitempty"`
	PublicPort           string   `protobuf:"bytes,4,opt,name=publicPort,proto3" json:"publicPort,omitempty"`
	Zone                 string   `protobuf:"bytes,5,opt,name=zone,proto3" json:"zone,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BootNatRegRes) Reset()         { *m = BootNatRegRes{} }
func (m *BootNatRegRes) String() string { return proto.CompactTextString(m) }
func (*BootNatRegRes) ProtoMessage()    {}
func (*BootNatRegRes) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f5cc5a30f7d4657, []int{1}
}

func (m *BootNatRegRes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BootNatRegRes.Unmarshal(m, b)
}
func (m *BootNatRegRes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BootNatRegRes.Marshal(b, m, deterministic)
}
func (m *BootNatRegRes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BootNatRegRes.Merge(m, src)
}
func (m *BootNatRegRes) XXX_Size() int {
	return xxx_messageInfo_BootNatRegRes.Size(m)
}
func (m *BootNatRegRes) XXX_DiscardUnknown() {
	xxx_messageInfo_BootNatRegRes.DiscardUnknown(m)
}

var xxx_messageInfo_BootNatRegRes proto.InternalMessageInfo

func (m *BootNatRegRes) GetNatType() NatType {
	if m != nil {
		return m.NatType
	}
	return NatType_UnknownRES
}

func (m *BootNatRegRes) GetNodeId() string {
	if m != nil {
		return m.NodeId
	}
	return ""
}

func (m *BootNatRegRes) GetPublicIp() string {
	if m != nil {
		return m.PublicIp
	}
	return ""
}

func (m *BootNatRegRes) GetPublicPort() string {
	if m != nil {
		return m.PublicPort
	}
	return ""
}

func (m *BootNatRegRes) GetZone() string {
	if m != nil {
		return m.Zone
	}
	return ""
}

type NatConReq struct {
	FromPeerId           string   `protobuf:"bytes,1,opt,name=FromPeerId,proto3" json:"FromPeerId,omitempty"`
	ToPeerId             string   `protobuf:"bytes,2,opt,name=ToPeerId,proto3" json:"ToPeerId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NatConReq) Reset()         { *m = NatConReq{} }
func (m *NatConReq) String() string { return proto.CompactTextString(m) }
func (*NatConReq) ProtoMessage()    {}
func (*NatConReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f5cc5a30f7d4657, []int{2}
}

func (m *NatConReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NatConReq.Unmarshal(m, b)
}
func (m *NatConReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NatConReq.Marshal(b, m, deterministic)
}
func (m *NatConReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NatConReq.Merge(m, src)
}
func (m *NatConReq) XXX_Size() int {
	return xxx_messageInfo_NatConReq.Size(m)
}
func (m *NatConReq) XXX_DiscardUnknown() {
	xxx_messageInfo_NatConReq.DiscardUnknown(m)
}

var xxx_messageInfo_NatConReq proto.InternalMessageInfo

func (m *NatConReq) GetFromPeerId() string {
	if m != nil {
		return m.FromPeerId
	}
	return ""
}

func (m *NatConReq) GetToPeerId() string {
	if m != nil {
		return m.ToPeerId
	}
	return ""
}

type NatConRes struct {
	PeerId               string   `protobuf:"bytes,2,opt,name=peerId,proto3" json:"peerId,omitempty"`
	PublicIp             string   `protobuf:"bytes,3,opt,name=publicIp,proto3" json:"publicIp,omitempty"`
	PublicPort           string   `protobuf:"bytes,4,opt,name=publicPort,proto3" json:"publicPort,omitempty"`
	PrivateIp            string   `protobuf:"bytes,5,opt,name=privateIp,proto3" json:"privateIp,omitempty"`
	PrivatePort          string   `protobuf:"bytes,6,opt,name=privatePort,proto3" json:"privatePort,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NatConRes) Reset()         { *m = NatConRes{} }
func (m *NatConRes) String() string { return proto.CompactTextString(m) }
func (*NatConRes) ProtoMessage()    {}
func (*NatConRes) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f5cc5a30f7d4657, []int{3}
}

func (m *NatConRes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NatConRes.Unmarshal(m, b)
}
func (m *NatConRes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NatConRes.Marshal(b, m, deterministic)
}
func (m *NatConRes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NatConRes.Merge(m, src)
}
func (m *NatConRes) XXX_Size() int {
	return xxx_messageInfo_NatConRes.Size(m)
}
func (m *NatConRes) XXX_DiscardUnknown() {
	xxx_messageInfo_NatConRes.DiscardUnknown(m)
}

var xxx_messageInfo_NatConRes proto.InternalMessageInfo

func (m *NatConRes) GetPeerId() string {
	if m != nil {
		return m.PeerId
	}
	return ""
}

func (m *NatConRes) GetPublicIp() string {
	if m != nil {
		return m.PublicIp
	}
	return ""
}

func (m *NatConRes) GetPublicPort() string {
	if m != nil {
		return m.PublicPort
	}
	return ""
}

func (m *NatConRes) GetPrivateIp() string {
	if m != nil {
		return m.PrivateIp
	}
	return ""
}

func (m *NatConRes) GetPrivatePort() string {
	if m != nil {
		return m.PrivatePort
	}
	return ""
}

type NatKeepAlive struct {
	NodeId               string   `protobuf:"bytes,1,opt,name=nodeId,proto3" json:"nodeId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NatKeepAlive) Reset()         { *m = NatKeepAlive{} }
func (m *NatKeepAlive) String() string { return proto.CompactTextString(m) }
func (*NatKeepAlive) ProtoMessage()    {}
func (*NatKeepAlive) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f5cc5a30f7d4657, []int{4}
}

func (m *NatKeepAlive) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NatKeepAlive.Unmarshal(m, b)
}
func (m *NatKeepAlive) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NatKeepAlive.Marshal(b, m, deterministic)
}
func (m *NatKeepAlive) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NatKeepAlive.Merge(m, src)
}
func (m *NatKeepAlive) XXX_Size() int {
	return xxx_messageInfo_NatKeepAlive.Size(m)
}
func (m *NatKeepAlive) XXX_DiscardUnknown() {
	xxx_messageInfo_NatKeepAlive.DiscardUnknown(m)
}

var xxx_messageInfo_NatKeepAlive proto.InternalMessageInfo

func (m *NatKeepAlive) GetNodeId() string {
	if m != nil {
		return m.NodeId
	}
	return ""
}

type NatPing struct {
	Ping                 string   `protobuf:"bytes,1,opt,name=ping,proto3" json:"ping,omitempty"`
	Pong                 string   `protobuf:"bytes,2,opt,name=pong,proto3" json:"pong,omitempty"`
	Nonce                string   `protobuf:"bytes,3,opt,name=nonce,proto3" json:"nonce,omitempty"`
	TTL                  int32    `protobuf:"varint,4,opt,name=TTL,proto3" json:"TTL,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NatPing) Reset()         { *m = NatPing{} }
func (m *NatPing) String() string { return proto.CompactTextString(m) }
func (*NatPing) ProtoMessage()    {}
func (*NatPing) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f5cc5a30f7d4657, []int{5}
}

func (m *NatPing) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NatPing.Unmarshal(m, b)
}
func (m *NatPing) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NatPing.Marshal(b, m, deterministic)
}
func (m *NatPing) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NatPing.Merge(m, src)
}
func (m *NatPing) XXX_Size() int {
	return xxx_messageInfo_NatPing.Size(m)
}
func (m *NatPing) XXX_DiscardUnknown() {
	xxx_messageInfo_NatPing.DiscardUnknown(m)
}

var xxx_messageInfo_NatPing proto.InternalMessageInfo

func (m *NatPing) GetPing() string {
	if m != nil {
		return m.Ping
	}
	return ""
}

func (m *NatPing) GetPong() string {
	if m != nil {
		return m.Pong
	}
	return ""
}

func (m *NatPing) GetNonce() string {
	if m != nil {
		return m.Nonce
	}
	return ""
}

func (m *NatPing) GetTTL() int32 {
	if m != nil {
		return m.TTL
	}
	return 0
}

type NatRequest struct {
	MsgType              NatMsgType     `protobuf:"varint,1,opt,name=msgType,proto3,enum=nat.pb.NatMsgType" json:"msgType,omitempty"`
	BootRegReq           *BootNatRegReq `protobuf:"bytes,2,opt,name=BootRegReq,proto3" json:"BootRegReq,omitempty"`
	ConnReq              *NatConReq     `protobuf:"bytes,3,opt,name=connReq,proto3" json:"connReq,omitempty"`
	KeepAlive            *NatKeepAlive  `protobuf:"bytes,4,opt,name=keepAlive,proto3" json:"keepAlive,omitempty"`
	Ping                 *NatPing       `protobuf:"bytes,5,opt,name=ping,proto3" json:"ping,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *NatRequest) Reset()         { *m = NatRequest{} }
func (m *NatRequest) String() string { return proto.CompactTextString(m) }
func (*NatRequest) ProtoMessage()    {}
func (*NatRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f5cc5a30f7d4657, []int{6}
}

func (m *NatRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NatRequest.Unmarshal(m, b)
}
func (m *NatRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NatRequest.Marshal(b, m, deterministic)
}
func (m *NatRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NatRequest.Merge(m, src)
}
func (m *NatRequest) XXX_Size() int {
	return xxx_messageInfo_NatRequest.Size(m)
}
func (m *NatRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_NatRequest.DiscardUnknown(m)
}

var xxx_messageInfo_NatRequest proto.InternalMessageInfo

func (m *NatRequest) GetMsgType() NatMsgType {
	if m != nil {
		return m.MsgType
	}
	return NatMsgType_UnknownReq
}

func (m *NatRequest) GetBootRegReq() *BootNatRegReq {
	if m != nil {
		return m.BootRegReq
	}
	return nil
}

func (m *NatRequest) GetConnReq() *NatConReq {
	if m != nil {
		return m.ConnReq
	}
	return nil
}

func (m *NatRequest) GetKeepAlive() *NatKeepAlive {
	if m != nil {
		return m.KeepAlive
	}
	return nil
}

func (m *NatRequest) GetPing() *NatPing {
	if m != nil {
		return m.Ping
	}
	return nil
}

type Response struct {
	MsgType              NatMsgType     `protobuf:"varint,1,opt,name=msgType,proto3,enum=nat.pb.NatMsgType" json:"msgType,omitempty"`
	BootRegRes           *BootNatRegRes `protobuf:"bytes,2,opt,name=BootRegRes,proto3" json:"BootRegRes,omitempty"`
	ConnRes              *NatConRes     `protobuf:"bytes,3,opt,name=connRes,proto3" json:"connRes,omitempty"`
	KeepAlive            *NatKeepAlive  `protobuf:"bytes,4,opt,name=keepAlive,proto3" json:"keepAlive,omitempty"`
	Pong                 *NatPing       `protobuf:"bytes,5,opt,name=pong,proto3" json:"pong,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f5cc5a30f7d4657, []int{7}
}

func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (m *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(m, src)
}
func (m *Response) XXX_Size() int {
	return xxx_messageInfo_Response.Size(m)
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

func (m *Response) GetMsgType() NatMsgType {
	if m != nil {
		return m.MsgType
	}
	return NatMsgType_UnknownReq
}

func (m *Response) GetBootRegRes() *BootNatRegRes {
	if m != nil {
		return m.BootRegRes
	}
	return nil
}

func (m *Response) GetConnRes() *NatConRes {
	if m != nil {
		return m.ConnRes
	}
	return nil
}

func (m *Response) GetKeepAlive() *NatKeepAlive {
	if m != nil {
		return m.KeepAlive
	}
	return nil
}

func (m *Response) GetPong() *NatPing {
	if m != nil {
		return m.Pong
	}
	return nil
}

func init() {
	proto.RegisterEnum("nat.pb.NatMsgType", NatMsgType_name, NatMsgType_value)
	proto.RegisterEnum("nat.pb.NatType", NatType_name, NatType_value)
	proto.RegisterType((*BootNatRegReq)(nil), "nat.pb.BootNatRegReq")
	proto.RegisterType((*BootNatRegRes)(nil), "nat.pb.BootNatRegRes")
	proto.RegisterType((*NatConReq)(nil), "nat.pb.NatConReq")
	proto.RegisterType((*NatConRes)(nil), "nat.pb.NatConRes")
	proto.RegisterType((*NatKeepAlive)(nil), "nat.pb.NatKeepAlive")
	proto.RegisterType((*NatPing)(nil), "nat.pb.NatPing")
	proto.RegisterType((*NatRequest)(nil), "nat.pb.NatRequest")
	proto.RegisterType((*Response)(nil), "nat.pb.response")
}

func init() { proto.RegisterFile("natManager.proto", fileDescriptor_2f5cc5a30f7d4657) }

var fileDescriptor_2f5cc5a30f7d4657 = []byte{
	// 578 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x54, 0xcb, 0x6a, 0x14, 0x41,
	0x14, 0x4d, 0xcf, 0x33, 0x73, 0x27, 0x8f, 0xf2, 0x12, 0x65, 0x10, 0x09, 0xa1, 0x05, 0xd1, 0x28,
	0xb3, 0x18, 0xf1, 0x03, 0x9c, 0xf1, 0x41, 0xd0, 0x34, 0xa1, 0xd3, 0x2e, 0xdc, 0x08, 0x95, 0x9e,
	0x4b, 0xa7, 0x49, 0x52, 0x55, 0xe9, 0xaa, 0x8c, 0xa8, 0x1f, 0x23, 0xf8, 0x8f, 0xe2, 0x56, 0xaa,
	0xfa, 0x31, 0x35, 0x91, 0x64, 0x61, 0x76, 0xf7, 0xd5, 0x75, 0xcf, 0x39, 0xf7, 0xd0, 0xc0, 0x04,
	0x37, 0x87, 0x5c, 0xf0, 0x8c, 0x8a, 0xb1, 0x2a, 0xa4, 0x91, 0xd8, 0x13, 0xdc, 0x8c, 0xd5, 0x49,
	0xf8, 0x03, 0x36, 0xa7, 0x52, 0x9a, 0x88, 0x9b, 0x98, 0xb2, 0x98, 0x2e, 0xf1, 0x01, 0xf4, 0x84,
	0x9c, 0xd3, 0xc1, 0x7c, 0x14, 0xec, 0x05, 0x4f, 0x07, 0x71, 0x95, 0xe1, 0x1e, 0x0c, 0x55, 0x91,
	0x2f, 0xb8, 0xa1, 0x23, 0x59, 0x98, 0x51, 0xdb, 0x35, 0xfd, 0x12, 0x22, 0x74, 0xbe, 0x4b, 0x41,
	0xa3, 0x8e, 0x6b, 0xb9, 0x18, 0x1f, 0xc1, 0xa0, 0x1a, 0x39, 0x50, 0xa3, 0x96, 0x6b, 0x2c, 0x0b,
	0xe1, 0xaf, 0x60, 0x75, 0xbb, 0xc6, 0x67, 0xd0, 0x17, 0xdc, 0x24, 0xdf, 0x14, 0xb9, 0xf5, 0x5b,
	0x93, 0xed, 0x71, 0x09, 0x74, 0x1c, 0x95, 0xe5, 0xb8, 0xee, 0x7b, 0x40, 0x5b, 0x2b, 0x40, 0x1f,
	0xc2, 0xba, 0xba, 0x3a, 0x39, 0xcf, 0xd3, 0x03, 0x55, 0xa1, 0x6c, 0x72, 0xdc, 0x05, 0x28, 0x63,
	0xc7, 0xa1, 0x04, 0xea, 0x55, 0x1a, 0x0a, 0xdd, 0x25, 0x85, 0xf0, 0x3d, 0x0c, 0x22, 0x6e, 0x66,
	0x52, 0x58, 0x75, 0x76, 0x01, 0xde, 0x15, 0xf2, 0xe2, 0x88, 0xa8, 0x68, 0x14, 0xf2, 0x2a, 0x76,
	0x79, 0x22, 0xab, 0x6e, 0x09, 0xab, 0xc9, 0xc3, 0x9f, 0xc1, 0xf2, 0x25, 0x6d, 0xe1, 0x2b, 0x7f,
	0xae, 0xca, 0xee, 0x04, 0x7f, 0x45, 0xed, 0xee, 0x35, 0xb5, 0xaf, 0x5f, 0xb0, 0xf7, 0xcf, 0x05,
	0xc3, 0x27, 0xb0, 0x11, 0x71, 0xf3, 0x81, 0x48, 0xbd, 0x3e, 0xcf, 0x17, 0x74, 0x93, 0x17, 0xc2,
	0xcf, 0xd0, 0x8f, 0xb8, 0x39, 0xca, 0x45, 0x66, 0x15, 0x53, 0xb9, 0xc8, 0xaa, 0x01, 0x17, 0xbb,
	0x9a, 0x14, 0x59, 0x45, 0xcc, 0xc5, 0xb8, 0x03, 0x5d, 0x21, 0x45, 0x4a, 0x15, 0xa7, 0x32, 0x41,
	0x06, 0xed, 0x24, 0xf9, 0xe8, 0x98, 0x74, 0x63, 0x1b, 0x86, 0x7f, 0x02, 0x00, 0x67, 0x87, 0xcb,
	0x2b, 0xd2, 0x06, 0x5f, 0x40, 0xff, 0x42, 0x67, 0x9e, 0x1f, 0xd0, 0xf3, 0xc3, 0x61, 0xd9, 0x89,
	0xeb, 0x11, 0x7c, 0x05, 0x60, 0xed, 0x54, 0x3a, 0xd9, 0xad, 0x1f, 0x4e, 0xee, 0xd7, 0x1f, 0xac,
	0xd8, 0x3c, 0xf6, 0x06, 0xf1, 0x39, 0xf4, 0x53, 0x29, 0xec, 0x7d, 0x1d, 0xba, 0xe1, 0xe4, 0x9e,
	0xb7, 0xa4, 0x3c, 0x7c, 0x5c, 0x4f, 0xe0, 0x04, 0x06, 0x67, 0xb5, 0x40, 0x0e, 0xf8, 0x70, 0xb2,
	0xe3, 0x8d, 0x37, 0xe2, 0xc5, 0xcb, 0x31, 0x7c, 0x5c, 0x89, 0xd4, 0x75, 0xe3, 0xbe, 0xa5, 0xad,
	0x86, 0xa5, 0x6a, 0xe1, 0xef, 0x00, 0xd6, 0x0b, 0xd2, 0x4a, 0x0a, 0x4d, 0x77, 0xe0, 0xad, 0x6f,
	0xe3, 0xad, 0x3d, 0xde, 0x7a, 0xc9, 0x5b, 0xdf, 0xc4, 0x5b, 0xd7, 0xbc, 0xf5, 0x7f, 0xf3, 0x96,
	0xb7, 0xf1, 0x96, 0x22, 0xdb, 0x4f, 0xdc, 0xc1, 0x2b, 0x4e, 0xb8, 0x05, 0xf0, 0x49, 0x9c, 0x09,
	0xf9, 0xd5, 0x8a, 0xcd, 0xd6, 0x90, 0xc1, 0x86, 0x45, 0x7c, 0x6c, 0x0a, 0xae, 0x62, 0xca, 0x58,
	0x80, 0x9b, 0x30, 0x68, 0x96, 0xb1, 0x16, 0x0e, 0xa1, 0x3f, 0x93, 0x42, 0x50, 0x6a, 0x58, 0x1b,
	0xd7, 0xa1, 0x63, 0x5f, 0x66, 0x9d, 0xfd, 0x2f, 0xce, 0xa2, 0xd7, 0x9f, 0x7c, 0x7b, 0xcc, 0xd6,
	0x70, 0x1b, 0x86, 0x91, 0x8c, 0xb8, 0x79, 0x43, 0x8b, 0x3c, 0xa5, 0xf2, 0xc5, 0x29, 0x9d, 0xe6,
	0x62, 0x1e, 0x71, 0xc3, 0x5a, 0x88, 0xb0, 0x35, 0xe3, 0x62, 0x4a, 0x11, 0x37, 0xc7, 0x54, 0x2c,
	0xa8, 0x60, 0x6d, 0xfb, 0x4d, 0x22, 0xa7, 0x34, 0x3b, 0xa5, 0xf4, 0x8c, 0xe6, 0xac, 0x73, 0xd2,
	0x73, 0xbf, 0xd1, 0x97, 0x7f, 0x03, 0x00, 0x00, 0xff, 0xff, 0x14, 0x87, 0x45, 0x92, 0x5a, 0x05,
	0x00, 0x00,
}
