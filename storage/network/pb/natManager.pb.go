// Code generated by protoc-gen-go. DO NOT EDIT.
// source: natManager.proto

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
	return fileDescriptor_2f5cc5a30f7d4657, []int{0}
}

type BootReg struct {
	NodeId               string   `protobuf:"bytes,1,opt,name=nodeId,proto3" json:"nodeId,omitempty"`
	PrivateIp            string   `protobuf:"bytes,2,opt,name=privateIp,proto3" json:"privateIp,omitempty"`
	PrivatePort          string   `protobuf:"bytes,3,opt,name=privatePort,proto3" json:"privatePort,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BootReg) Reset()         { *m = BootReg{} }
func (m *BootReg) String() string { return proto.CompactTextString(m) }
func (*BootReg) ProtoMessage()    {}
func (*BootReg) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f5cc5a30f7d4657, []int{0}
}

func (m *BootReg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BootReg.Unmarshal(m, b)
}
func (m *BootReg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BootReg.Marshal(b, m, deterministic)
}
func (m *BootReg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BootReg.Merge(m, src)
}
func (m *BootReg) XXX_Size() int {
	return xxx_messageInfo_BootReg.Size(m)
}
func (m *BootReg) XXX_DiscardUnknown() {
	xxx_messageInfo_BootReg.DiscardUnknown(m)
}

var xxx_messageInfo_BootReg proto.InternalMessageInfo

func (m *BootReg) GetNodeId() string {
	if m != nil {
		return m.NodeId
	}
	return ""
}

func (m *BootReg) GetPrivateIp() string {
	if m != nil {
		return m.PrivateIp
	}
	return ""
}

func (m *BootReg) GetPrivatePort() string {
	if m != nil {
		return m.PrivatePort
	}
	return ""
}

type BootAnswer struct {
	NatType              NatType  `protobuf:"varint,1,opt,name=NatType,proto3,enum=net.pb.NatType" json:"NatType,omitempty"`
	PublicIp             string   `protobuf:"bytes,2,opt,name=PublicIp,proto3" json:"PublicIp,omitempty"`
	PublicPort           string   `protobuf:"bytes,3,opt,name=PublicPort,proto3" json:"PublicPort,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BootAnswer) Reset()         { *m = BootAnswer{} }
func (m *BootAnswer) String() string { return proto.CompactTextString(m) }
func (*BootAnswer) ProtoMessage()    {}
func (*BootAnswer) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f5cc5a30f7d4657, []int{1}
}

func (m *BootAnswer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BootAnswer.Unmarshal(m, b)
}
func (m *BootAnswer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BootAnswer.Marshal(b, m, deterministic)
}
func (m *BootAnswer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BootAnswer.Merge(m, src)
}
func (m *BootAnswer) XXX_Size() int {
	return xxx_messageInfo_BootAnswer.Size(m)
}
func (m *BootAnswer) XXX_DiscardUnknown() {
	xxx_messageInfo_BootAnswer.DiscardUnknown(m)
}

var xxx_messageInfo_BootAnswer proto.InternalMessageInfo

func (m *BootAnswer) GetNatType() NatType {
	if m != nil {
		return m.NatType
	}
	return NatType_UnknownRES
}

func (m *BootAnswer) GetPublicIp() string {
	if m != nil {
		return m.PublicIp
	}
	return ""
}

func (m *BootAnswer) GetPublicPort() string {
	if m != nil {
		return m.PublicPort
	}
	return ""
}

type DigApply struct {
	NatServer            string   `protobuf:"bytes,1,opt,name=NatServer,proto3" json:"NatServer,omitempty"`
	Public               string   `protobuf:"bytes,2,opt,name=Public,proto3" json:"Public,omitempty"`
	TargetId             string   `protobuf:"bytes,3,opt,name=TargetId,proto3" json:"TargetId,omitempty"`
	TargetPort           int32    `protobuf:"varint,4,opt,name=TargetPort,proto3" json:"TargetPort,omitempty"`
	SessionId            string   `protobuf:"bytes,5,opt,name=SessionId,proto3" json:"SessionId,omitempty"`
	FromId               string   `protobuf:"bytes,6,opt,name=FromId,proto3" json:"FromId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DigApply) Reset()         { *m = DigApply{} }
func (m *DigApply) String() string { return proto.CompactTextString(m) }
func (*DigApply) ProtoMessage()    {}
func (*DigApply) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f5cc5a30f7d4657, []int{2}
}

func (m *DigApply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DigApply.Unmarshal(m, b)
}
func (m *DigApply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DigApply.Marshal(b, m, deterministic)
}
func (m *DigApply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DigApply.Merge(m, src)
}
func (m *DigApply) XXX_Size() int {
	return xxx_messageInfo_DigApply.Size(m)
}
func (m *DigApply) XXX_DiscardUnknown() {
	xxx_messageInfo_DigApply.DiscardUnknown(m)
}

var xxx_messageInfo_DigApply proto.InternalMessageInfo

func (m *DigApply) GetNatServer() string {
	if m != nil {
		return m.NatServer
	}
	return ""
}

func (m *DigApply) GetPublic() string {
	if m != nil {
		return m.Public
	}
	return ""
}

func (m *DigApply) GetTargetId() string {
	if m != nil {
		return m.TargetId
	}
	return ""
}

func (m *DigApply) GetTargetPort() int32 {
	if m != nil {
		return m.TargetPort
	}
	return 0
}

func (m *DigApply) GetSessionId() string {
	if m != nil {
		return m.SessionId
	}
	return ""
}

func (m *DigApply) GetFromId() string {
	if m != nil {
		return m.FromId
	}
	return ""
}

type DigConfirm struct {
	Public               string   `protobuf:"bytes,1,opt,name=Public,proto3" json:"Public,omitempty"`
	SessionId            string   `protobuf:"bytes,2,opt,name=sessionId,proto3" json:"sessionId,omitempty"`
	TargetId             string   `protobuf:"bytes,3,opt,name=targetId,proto3" json:"targetId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DigConfirm) Reset()         { *m = DigConfirm{} }
func (m *DigConfirm) String() string { return proto.CompactTextString(m) }
func (*DigConfirm) ProtoMessage()    {}
func (*DigConfirm) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f5cc5a30f7d4657, []int{3}
}

func (m *DigConfirm) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DigConfirm.Unmarshal(m, b)
}
func (m *DigConfirm) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DigConfirm.Marshal(b, m, deterministic)
}
func (m *DigConfirm) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DigConfirm.Merge(m, src)
}
func (m *DigConfirm) XXX_Size() int {
	return xxx_messageInfo_DigConfirm.Size(m)
}
func (m *DigConfirm) XXX_DiscardUnknown() {
	xxx_messageInfo_DigConfirm.DiscardUnknown(m)
}

var xxx_messageInfo_DigConfirm proto.InternalMessageInfo

func (m *DigConfirm) GetPublic() string {
	if m != nil {
		return m.Public
	}
	return ""
}

func (m *DigConfirm) GetSessionId() string {
	if m != nil {
		return m.SessionId
	}
	return ""
}

func (m *DigConfirm) GetTargetId() string {
	if m != nil {
		return m.TargetId
	}
	return ""
}

type KeepAlive struct {
	NodeId               string   `protobuf:"bytes,1,opt,name=nodeId,proto3" json:"nodeId,omitempty"`
	LAddr                string   `protobuf:"bytes,2,opt,name=lAddr,proto3" json:"lAddr,omitempty"`
	PubIP                string   `protobuf:"bytes,3,opt,name=PubIP,proto3" json:"PubIP,omitempty"`
	PubPort              int32    `protobuf:"varint,4,opt,name=PubPort,proto3" json:"PubPort,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *KeepAlive) Reset()         { *m = KeepAlive{} }
func (m *KeepAlive) String() string { return proto.CompactTextString(m) }
func (*KeepAlive) ProtoMessage()    {}
func (*KeepAlive) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f5cc5a30f7d4657, []int{4}
}

func (m *KeepAlive) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KeepAlive.Unmarshal(m, b)
}
func (m *KeepAlive) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KeepAlive.Marshal(b, m, deterministic)
}
func (m *KeepAlive) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KeepAlive.Merge(m, src)
}
func (m *KeepAlive) XXX_Size() int {
	return xxx_messageInfo_KeepAlive.Size(m)
}
func (m *KeepAlive) XXX_DiscardUnknown() {
	xxx_messageInfo_KeepAlive.DiscardUnknown(m)
}

var xxx_messageInfo_KeepAlive proto.InternalMessageInfo

func (m *KeepAlive) GetNodeId() string {
	if m != nil {
		return m.NodeId
	}
	return ""
}

func (m *KeepAlive) GetLAddr() string {
	if m != nil {
		return m.LAddr
	}
	return ""
}

func (m *KeepAlive) GetPubIP() string {
	if m != nil {
		return m.PubIP
	}
	return ""
}

func (m *KeepAlive) GetPubPort() int32 {
	if m != nil {
		return m.PubPort
	}
	return 0
}

type PingPong struct {
	Ping                 string   `protobuf:"bytes,1,opt,name=ping,proto3" json:"ping,omitempty"`
	Pong                 string   `protobuf:"bytes,2,opt,name=pong,proto3" json:"pong,omitempty"`
	Nonce                string   `protobuf:"bytes,3,opt,name=nonce,proto3" json:"nonce,omitempty"`
	TTL                  int32    `protobuf:"varint,4,opt,name=TTL,proto3" json:"TTL,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PingPong) Reset()         { *m = PingPong{} }
func (m *PingPong) String() string { return proto.CompactTextString(m) }
func (*PingPong) ProtoMessage()    {}
func (*PingPong) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f5cc5a30f7d4657, []int{5}
}

func (m *PingPong) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PingPong.Unmarshal(m, b)
}
func (m *PingPong) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PingPong.Marshal(b, m, deterministic)
}
func (m *PingPong) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PingPong.Merge(m, src)
}
func (m *PingPong) XXX_Size() int {
	return xxx_messageInfo_PingPong.Size(m)
}
func (m *PingPong) XXX_DiscardUnknown() {
	xxx_messageInfo_PingPong.DiscardUnknown(m)
}

var xxx_messageInfo_PingPong proto.InternalMessageInfo

func (m *PingPong) GetPing() string {
	if m != nil {
		return m.Ping
	}
	return ""
}

func (m *PingPong) GetPong() string {
	if m != nil {
		return m.Pong
	}
	return ""
}

func (m *PingPong) GetNonce() string {
	if m != nil {
		return m.Nonce
	}
	return ""
}

func (m *PingPong) GetTTL() int32 {
	if m != nil {
		return m.TTL
	}
	return 0
}

type ReverseInvite struct {
	SessionId            string   `protobuf:"bytes,1,opt,name=sessionId,proto3" json:"sessionId,omitempty"`
	PubIp                string   `protobuf:"bytes,2,opt,name=pubIp,proto3" json:"pubIp,omitempty"`
	ToPort               int32    `protobuf:"varint,3,opt,name=toPort,proto3" json:"toPort,omitempty"`
	PeerId               string   `protobuf:"bytes,4,opt,name=peerId,proto3" json:"peerId,omitempty"`
	FromPort             string   `protobuf:"bytes,5,opt,name=fromPort,proto3" json:"fromPort,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReverseInvite) Reset()         { *m = ReverseInvite{} }
func (m *ReverseInvite) String() string { return proto.CompactTextString(m) }
func (*ReverseInvite) ProtoMessage()    {}
func (*ReverseInvite) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f5cc5a30f7d4657, []int{6}
}

func (m *ReverseInvite) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReverseInvite.Unmarshal(m, b)
}
func (m *ReverseInvite) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReverseInvite.Marshal(b, m, deterministic)
}
func (m *ReverseInvite) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReverseInvite.Merge(m, src)
}
func (m *ReverseInvite) XXX_Size() int {
	return xxx_messageInfo_ReverseInvite.Size(m)
}
func (m *ReverseInvite) XXX_DiscardUnknown() {
	xxx_messageInfo_ReverseInvite.DiscardUnknown(m)
}

var xxx_messageInfo_ReverseInvite proto.InternalMessageInfo

func (m *ReverseInvite) GetSessionId() string {
	if m != nil {
		return m.SessionId
	}
	return ""
}

func (m *ReverseInvite) GetPubIp() string {
	if m != nil {
		return m.PubIp
	}
	return ""
}

func (m *ReverseInvite) GetToPort() int32 {
	if m != nil {
		return m.ToPort
	}
	return 0
}

func (m *ReverseInvite) GetPeerId() string {
	if m != nil {
		return m.PeerId
	}
	return ""
}

func (m *ReverseInvite) GetFromPort() string {
	if m != nil {
		return m.FromPort
	}
	return ""
}

type ReverseInviteAck struct {
	SessionId            string   `protobuf:"bytes,1,opt,name=sessionId,proto3" json:"sessionId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReverseInviteAck) Reset()         { *m = ReverseInviteAck{} }
func (m *ReverseInviteAck) String() string { return proto.CompactTextString(m) }
func (*ReverseInviteAck) ProtoMessage()    {}
func (*ReverseInviteAck) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f5cc5a30f7d4657, []int{7}
}

func (m *ReverseInviteAck) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReverseInviteAck.Unmarshal(m, b)
}
func (m *ReverseInviteAck) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReverseInviteAck.Marshal(b, m, deterministic)
}
func (m *ReverseInviteAck) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReverseInviteAck.Merge(m, src)
}
func (m *ReverseInviteAck) XXX_Size() int {
	return xxx_messageInfo_ReverseInviteAck.Size(m)
}
func (m *ReverseInviteAck) XXX_DiscardUnknown() {
	xxx_messageInfo_ReverseInviteAck.DiscardUnknown(m)
}

var xxx_messageInfo_ReverseInviteAck proto.InternalMessageInfo

func (m *ReverseInviteAck) GetSessionId() string {
	if m != nil {
		return m.SessionId
	}
	return ""
}

type NatMsg struct {
	Typ                  MsgType           `protobuf:"varint,1,opt,name=Typ,proto3,enum=net.pb.MsgType" json:"Typ,omitempty"`
	Seq                  int64             `protobuf:"varint,2,opt,name=Seq,proto3" json:"Seq,omitempty"`
	BootReg              *BootReg          `protobuf:"bytes,3,opt,name=BootReg,proto3" json:"BootReg,omitempty"`
	KeepAlive            *KeepAlive        `protobuf:"bytes,4,opt,name=KeepAlive,proto3" json:"KeepAlive,omitempty"`
	DigApply             *DigApply         `protobuf:"bytes,5,opt,name=DigApply,proto3" json:"DigApply,omitempty"`
	DigConfirm           *DigConfirm       `protobuf:"bytes,6,opt,name=DigConfirm,proto3" json:"DigConfirm,omitempty"`
	PingPong             *PingPong         `protobuf:"bytes,7,opt,name=PingPong,proto3" json:"PingPong,omitempty"`
	ReverseInvite        *ReverseInvite    `protobuf:"bytes,8,opt,name=ReverseInvite,proto3" json:"ReverseInvite,omitempty"`
	ReverseInviteAck     *ReverseInviteAck `protobuf:"bytes,9,opt,name=ReverseInviteAck,proto3" json:"ReverseInviteAck,omitempty"`
	BootAnswer           *BootAnswer       `protobuf:"bytes,10,opt,name=BootAnswer,proto3" json:"BootAnswer,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *NatMsg) Reset()         { *m = NatMsg{} }
func (m *NatMsg) String() string { return proto.CompactTextString(m) }
func (*NatMsg) ProtoMessage()    {}
func (*NatMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f5cc5a30f7d4657, []int{8}
}

func (m *NatMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NatMsg.Unmarshal(m, b)
}
func (m *NatMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NatMsg.Marshal(b, m, deterministic)
}
func (m *NatMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NatMsg.Merge(m, src)
}
func (m *NatMsg) XXX_Size() int {
	return xxx_messageInfo_NatMsg.Size(m)
}
func (m *NatMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_NatMsg.DiscardUnknown(m)
}

var xxx_messageInfo_NatMsg proto.InternalMessageInfo

func (m *NatMsg) GetTyp() MsgType {
	if m != nil {
		return m.Typ
	}
	return MsgType_Error
}

func (m *NatMsg) GetSeq() int64 {
	if m != nil {
		return m.Seq
	}
	return 0
}

func (m *NatMsg) GetBootReg() *BootReg {
	if m != nil {
		return m.BootReg
	}
	return nil
}

func (m *NatMsg) GetKeepAlive() *KeepAlive {
	if m != nil {
		return m.KeepAlive
	}
	return nil
}

func (m *NatMsg) GetDigApply() *DigApply {
	if m != nil {
		return m.DigApply
	}
	return nil
}

func (m *NatMsg) GetDigConfirm() *DigConfirm {
	if m != nil {
		return m.DigConfirm
	}
	return nil
}

func (m *NatMsg) GetPingPong() *PingPong {
	if m != nil {
		return m.PingPong
	}
	return nil
}

func (m *NatMsg) GetReverseInvite() *ReverseInvite {
	if m != nil {
		return m.ReverseInvite
	}
	return nil
}

func (m *NatMsg) GetReverseInviteAck() *ReverseInviteAck {
	if m != nil {
		return m.ReverseInviteAck
	}
	return nil
}

func (m *NatMsg) GetBootAnswer() *BootAnswer {
	if m != nil {
		return m.BootAnswer
	}
	return nil
}

func init() {
	proto.RegisterEnum("net.pb.NatType", NatType_name, NatType_value)
	proto.RegisterType((*BootReg)(nil), "net.pb.BootReg")
	proto.RegisterType((*BootAnswer)(nil), "net.pb.BootAnswer")
	proto.RegisterType((*DigApply)(nil), "net.pb.DigApply")
	proto.RegisterType((*DigConfirm)(nil), "net.pb.DigConfirm")
	proto.RegisterType((*KeepAlive)(nil), "net.pb.KeepAlive")
	proto.RegisterType((*PingPong)(nil), "net.pb.PingPong")
	proto.RegisterType((*ReverseInvite)(nil), "net.pb.ReverseInvite")
	proto.RegisterType((*ReverseInviteAck)(nil), "net.pb.ReverseInviteAck")
	proto.RegisterType((*NatMsg)(nil), "net.pb.NatMsg")
}

func init() { proto.RegisterFile("natManager.proto", fileDescriptor_2f5cc5a30f7d4657) }

var fileDescriptor_2f5cc5a30f7d4657 = []byte{
	// 675 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x54, 0xdd, 0x6e, 0xd3, 0x4c,
	0x10, 0xfd, 0xdc, 0x34, 0x3f, 0x9e, 0x7c, 0x6d, 0xcd, 0xaa, 0x20, 0xab, 0xaa, 0x50, 0xf1, 0x55,
	0x41, 0x55, 0x40, 0xe1, 0x92, 0xab, 0xa4, 0x01, 0xc9, 0x82, 0x44, 0xd6, 0x26, 0xdc, 0x70, 0x51,
	0xc9, 0x8e, 0xa7, 0xae, 0x95, 0x74, 0xd7, 0xac, 0x37, 0xa9, 0xfa, 0x14, 0xbc, 0x0b, 0x4f, 0xc4,
	0xa3, 0xa0, 0x5d, 0xaf, 0x7f, 0x92, 0x02, 0x77, 0x73, 0xce, 0xec, 0xec, 0x99, 0xd9, 0x99, 0x59,
	0x70, 0x58, 0x28, 0xa7, 0x21, 0x0b, 0x13, 0x14, 0x83, 0x4c, 0x70, 0xc9, 0x49, 0x87, 0xa1, 0x1c,
	0x64, 0xd1, 0xd9, 0xff, 0x2c, 0xca, 0xa7, 0x79, 0x52, 0xb0, 0x5e, 0x08, 0xdd, 0x31, 0xe7, 0x92,
	0x62, 0x42, 0x5e, 0x40, 0x87, 0xf1, 0x18, 0xfd, 0xd8, 0xb5, 0x2e, 0xac, 0x4b, 0x9b, 0x1a, 0x44,
	0xce, 0xc1, 0xce, 0x44, 0xba, 0x0d, 0x25, 0xfa, 0x99, 0x7b, 0xa0, 0x5d, 0x35, 0x41, 0x2e, 0xa0,
	0x6f, 0x40, 0xc0, 0x85, 0x74, 0x5b, 0xda, 0xdf, 0xa4, 0xbc, 0x1c, 0x40, 0x49, 0x8c, 0x58, 0xfe,
	0x80, 0x82, 0xbc, 0x86, 0xee, 0x2c, 0x94, 0x8b, 0xc7, 0x0c, 0xb5, 0xcc, 0xf1, 0xf0, 0x64, 0x50,
	0x24, 0x36, 0x30, 0x34, 0x2d, 0xfd, 0xe4, 0x0c, 0x7a, 0xc1, 0x26, 0x5a, 0xa7, 0xcb, 0x4a, 0xb7,
	0xc2, 0xe4, 0x25, 0x40, 0x61, 0x37, 0x54, 0x1b, 0x8c, 0xf7, 0xd3, 0x82, 0xde, 0x24, 0x4d, 0x46,
	0x59, 0xb6, 0x7e, 0x54, 0x15, 0xcc, 0x42, 0x39, 0x47, 0xb1, 0x45, 0x61, 0x8a, 0xab, 0x09, 0x55,
	0x77, 0x11, 0x68, 0x44, 0x0c, 0x52, 0xf2, 0x8b, 0x50, 0x24, 0x28, 0xfd, 0xd8, 0x08, 0x54, 0x58,
	0xc9, 0x17, 0xb6, 0x96, 0x3f, 0xbc, 0xb0, 0x2e, 0xdb, 0xb4, 0xc1, 0x28, 0xc5, 0x39, 0xe6, 0x79,
	0xca, 0x99, 0x1f, 0xbb, 0xed, 0x42, 0xb1, 0x22, 0x94, 0xe2, 0x27, 0xc1, 0xef, 0xfd, 0xd8, 0xed,
	0x14, 0x8a, 0x05, 0xf2, 0x6e, 0x00, 0x26, 0x69, 0x72, 0xcd, 0xd9, 0x6d, 0x2a, 0xee, 0x1b, 0x79,
	0x59, 0x3b, 0x79, 0x9d, 0x83, 0x9d, 0x57, 0x77, 0x9b, 0x7e, 0x54, 0x84, 0xca, 0x5a, 0xee, 0x65,
	0x5d, 0x62, 0x2f, 0x05, 0xfb, 0x33, 0x62, 0x36, 0x5a, 0xa7, 0x5b, 0xfc, 0x6b, 0xbb, 0x4f, 0xa1,
	0xbd, 0x1e, 0xc5, 0xb1, 0x30, 0x57, 0x17, 0x40, 0xb1, 0xc1, 0x26, 0xf2, 0x03, 0x73, 0x67, 0x01,
	0x88, 0x0b, 0xdd, 0x60, 0x13, 0x35, 0xde, 0xa0, 0x84, 0xde, 0x37, 0xe8, 0x05, 0x29, 0x4b, 0x02,
	0xce, 0x12, 0x42, 0xe0, 0x30, 0x4b, 0x59, 0x62, 0x74, 0xb4, 0xad, 0x39, 0xce, 0x12, 0x23, 0xa2,
	0x6d, 0xa5, 0xc1, 0x38, 0x5b, 0x62, 0xa9, 0xa1, 0x01, 0x71, 0xa0, 0xb5, 0x58, 0x7c, 0x31, 0xf7,
	0x2b, 0xd3, 0xfb, 0x61, 0xc1, 0x11, 0xc5, 0x2d, 0x8a, 0x1c, 0x7d, 0xb6, 0x4d, 0x25, 0xee, 0x3e,
	0x89, 0xb5, 0xff, 0x24, 0xa7, 0xd0, 0xce, 0x36, 0x51, 0x35, 0x44, 0x05, 0x50, 0xf5, 0x4b, 0x5e,
	0x4d, 0x4f, 0x9b, 0x1a, 0xa4, 0xf8, 0x0c, 0x51, 0xf8, 0xb1, 0x96, 0xb4, 0xa9, 0x41, 0xea, 0x61,
	0x6f, 0x05, 0xbf, 0xd7, 0x11, 0x45, 0x47, 0x2b, 0xec, 0xbd, 0x03, 0x67, 0x27, 0xa1, 0xd1, 0x72,
	0xf5, 0xef, 0x9c, 0xbc, 0x5f, 0x2d, 0xe8, 0xcc, 0x42, 0x39, 0xcd, 0x13, 0xf2, 0x0a, 0x5a, 0x8b,
	0xc7, 0x6c, 0x7f, 0x1b, 0xa6, 0x79, 0xa2, 0xb7, 0x41, 0xf9, 0xd4, 0x1b, 0xcc, 0xf1, 0xbb, 0xce,
	0xbf, 0x45, 0x95, 0xa9, 0xd6, 0xc8, 0xec, 0xad, 0x4e, 0xbf, 0x5f, 0x07, 0x1a, 0x9a, 0x56, 0x7b,
	0xfd, 0xb6, 0xd1, 0x75, 0x5d, 0x53, 0x7f, 0xf8, 0xac, 0x3c, 0x5c, 0x39, 0x68, 0x63, 0x32, 0xae,
	0xea, 0xd5, 0xd1, 0x95, 0xf6, 0x87, 0x4e, 0x79, 0xbe, 0xe4, 0x69, 0xbd, 0x5c, 0xc3, 0xe6, 0xd0,
	0xea, 0x81, 0xee, 0x0f, 0x49, 0xe3, 0xbc, 0xf1, 0xd0, 0xe6, 0x68, 0x5f, 0xd5, 0xd3, 0xe1, 0x76,
	0x77, 0x15, 0x4a, 0x9e, 0xd6, 0xf3, 0xf3, 0x61, 0xaf, 0xdd, 0x6e, 0x4f, 0x87, 0x3c, 0x2f, 0x43,
	0x76, 0x9c, 0x74, 0x6f, 0x34, 0x26, 0x4f, 0x5b, 0xe3, 0xda, 0x3a, 0xde, 0xfd, 0x63, 0xfc, 0x68,
	0xb9, 0xa2, 0x4f, 0x9b, 0x39, 0x6c, 0xfe, 0x61, 0x2e, 0xec, 0x16, 0x59, 0x7b, 0x68, 0xe3, 0xd4,
	0x9b, 0x9b, 0xea, 0xa7, 0x23, 0xc7, 0x00, 0x5f, 0xd9, 0x8a, 0xf1, 0x07, 0x46, 0x3f, 0xce, 0x9d,
	0xff, 0xc8, 0x09, 0xf4, 0x67, 0x7c, 0x16, 0xca, 0x09, 0x6e, 0xd3, 0x25, 0x3a, 0x16, 0x39, 0x02,
	0x7b, 0x8c, 0x77, 0x29, 0x8b, 0x67, 0xa1, 0x74, 0x0e, 0x08, 0x81, 0xe3, 0xeb, 0x90, 0x8d, 0xb1,
	0xfa, 0xa4, 0x9c, 0x96, 0x8a, 0x59, 0xf0, 0x31, 0x5e, 0xdf, 0xe1, 0x72, 0x85, 0xb1, 0x73, 0x18,
	0x75, 0xf4, 0x0f, 0xfe, 0xfe, 0x77, 0x00, 0x00, 0x00, 0xff, 0xff, 0xa6, 0x6d, 0x5d, 0xb1, 0xeb,
	0x05, 0x00, 0x00,
}
