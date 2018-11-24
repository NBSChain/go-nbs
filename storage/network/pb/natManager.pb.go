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

type NatMsgType int32

const (
	NatMsgType_UnknownReq    NatMsgType = 0
	NatMsgType_BootStrapReg  NatMsgType = 1
	NatMsgType_KeepAlive     NatMsgType = 2
	NatMsgType_Connect       NatMsgType = 3
	NatMsgType_Ping          NatMsgType = 4
	NatMsgType_DigIn         NatMsgType = 5
	NatMsgType_DigOut        NatMsgType = 6
	NatMsgType_ReverseDig    NatMsgType = 7
	NatMsgType_ReverseDigACK NatMsgType = 8
	NatMsgType_error         NatMsgType = 20
)

var NatMsgType_name = map[int32]string{
	0:  "UnknownReq",
	1:  "BootStrapReg",
	2:  "KeepAlive",
	3:  "Connect",
	4:  "Ping",
	5:  "DigIn",
	6:  "DigOut",
	7:  "ReverseDig",
	8:  "ReverseDigACK",
	20: "error",
}

var NatMsgType_value = map[string]int32{
	"UnknownReq":    0,
	"BootStrapReg":  1,
	"KeepAlive":     2,
	"Connect":       3,
	"Ping":          4,
	"DigIn":         5,
	"DigOut":        6,
	"ReverseDig":    7,
	"ReverseDigACK": 8,
	"error":         20,
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

func (m *BootNatRegReq) GetPrivateIp() string {
	if m != nil {
		return m.PrivateIp
	}
	return ""
}

type BootNatRegRes struct {
	NatType              NatType  `protobuf:"varint,1,opt,name=natType,proto3,enum=net.pb.NatType" json:"natType,omitempty"`
	NodeId               string   `protobuf:"bytes,2,opt,name=nodeId,proto3" json:"nodeId,omitempty"`
	PublicIp             string   `protobuf:"bytes,3,opt,name=publicIp,proto3" json:"publicIp,omitempty"`
	PublicPort           string   `protobuf:"bytes,4,opt,name=publicPort,proto3" json:"publicPort,omitempty"`
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

type NatConInvite struct {
	FromAddr             *NbsAddr `protobuf:"bytes,1,opt,name=fromAddr,proto3" json:"fromAddr,omitempty"`
	ToAddr               *NbsAddr `protobuf:"bytes,2,opt,name=toAddr,proto3" json:"toAddr,omitempty"`
	SessionId            string   `protobuf:"bytes,3,opt,name=sessionId,proto3" json:"sessionId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NatConInvite) Reset()         { *m = NatConInvite{} }
func (m *NatConInvite) String() string { return proto.CompactTextString(m) }
func (*NatConInvite) ProtoMessage()    {}
func (*NatConInvite) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f5cc5a30f7d4657, []int{2}
}

func (m *NatConInvite) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NatConInvite.Unmarshal(m, b)
}
func (m *NatConInvite) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NatConInvite.Marshal(b, m, deterministic)
}
func (m *NatConInvite) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NatConInvite.Merge(m, src)
}
func (m *NatConInvite) XXX_Size() int {
	return xxx_messageInfo_NatConInvite.Size(m)
}
func (m *NatConInvite) XXX_DiscardUnknown() {
	xxx_messageInfo_NatConInvite.DiscardUnknown(m)
}

var xxx_messageInfo_NatConInvite proto.InternalMessageInfo

func (m *NatConInvite) GetFromAddr() *NbsAddr {
	if m != nil {
		return m.FromAddr
	}
	return nil
}

func (m *NatConInvite) GetToAddr() *NbsAddr {
	if m != nil {
		return m.ToAddr
	}
	return nil
}

func (m *NatConInvite) GetSessionId() string {
	if m != nil {
		return m.SessionId
	}
	return ""
}

type NatKeepAlive struct {
	NodeId               string   `protobuf:"bytes,1,opt,name=nodeId,proto3" json:"nodeId,omitempty"`
	LAddr                string   `protobuf:"bytes,2,opt,name=lAddr,proto3" json:"lAddr,omitempty"`
	PubIP                string   `protobuf:"bytes,3,opt,name=PubIP,proto3" json:"PubIP,omitempty"`
	PubPort              int32    `protobuf:"varint,4,opt,name=PubPort,proto3" json:"PubPort,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NatKeepAlive) Reset()         { *m = NatKeepAlive{} }
func (m *NatKeepAlive) String() string { return proto.CompactTextString(m) }
func (*NatKeepAlive) ProtoMessage()    {}
func (*NatKeepAlive) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f5cc5a30f7d4657, []int{3}
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

func (m *NatKeepAlive) GetLAddr() string {
	if m != nil {
		return m.LAddr
	}
	return ""
}

func (m *NatKeepAlive) GetPubIP() string {
	if m != nil {
		return m.PubIP
	}
	return ""
}

func (m *NatKeepAlive) GetPubPort() int32 {
	if m != nil {
		return m.PubPort
	}
	return 0
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
	return fileDescriptor_2f5cc5a30f7d4657, []int{4}
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

type HoleDig struct {
	SessionId            string   `protobuf:"bytes,1,opt,name=sessionId,proto3" json:"sessionId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HoleDig) Reset()         { *m = HoleDig{} }
func (m *HoleDig) String() string { return proto.CompactTextString(m) }
func (*HoleDig) ProtoMessage()    {}
func (*HoleDig) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f5cc5a30f7d4657, []int{5}
}

func (m *HoleDig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HoleDig.Unmarshal(m, b)
}
func (m *HoleDig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HoleDig.Marshal(b, m, deterministic)
}
func (m *HoleDig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HoleDig.Merge(m, src)
}
func (m *HoleDig) XXX_Size() int {
	return xxx_messageInfo_HoleDig.Size(m)
}
func (m *HoleDig) XXX_DiscardUnknown() {
	xxx_messageInfo_HoleDig.DiscardUnknown(m)
}

var xxx_messageInfo_HoleDig proto.InternalMessageInfo

func (m *HoleDig) GetSessionId() string {
	if m != nil {
		return m.SessionId
	}
	return ""
}

type ErrorNotify struct {
	ErrMsg               string   `protobuf:"bytes,1,opt,name=errMsg,proto3" json:"errMsg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ErrorNotify) Reset()         { *m = ErrorNotify{} }
func (m *ErrorNotify) String() string { return proto.CompactTextString(m) }
func (*ErrorNotify) ProtoMessage()    {}
func (*ErrorNotify) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f5cc5a30f7d4657, []int{6}
}

func (m *ErrorNotify) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ErrorNotify.Unmarshal(m, b)
}
func (m *ErrorNotify) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ErrorNotify.Marshal(b, m, deterministic)
}
func (m *ErrorNotify) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ErrorNotify.Merge(m, src)
}
func (m *ErrorNotify) XXX_Size() int {
	return xxx_messageInfo_ErrorNotify.Size(m)
}
func (m *ErrorNotify) XXX_DiscardUnknown() {
	xxx_messageInfo_ErrorNotify.DiscardUnknown(m)
}

var xxx_messageInfo_ErrorNotify proto.InternalMessageInfo

func (m *ErrorNotify) GetErrMsg() string {
	if m != nil {
		return m.ErrMsg
	}
	return ""
}

type ReverseInvite struct {
	SessionId            string   `protobuf:"bytes,1,opt,name=sessionId,proto3" json:"sessionId,omitempty"`
	PubIp                string   `protobuf:"bytes,2,opt,name=pubIp,proto3" json:"pubIp,omitempty"`
	ToPort               int32    `protobuf:"varint,3,opt,name=toPort,proto3" json:"toPort,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReverseInvite) Reset()         { *m = ReverseInvite{} }
func (m *ReverseInvite) String() string { return proto.CompactTextString(m) }
func (*ReverseInvite) ProtoMessage()    {}
func (*ReverseInvite) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f5cc5a30f7d4657, []int{7}
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
	return fileDescriptor_2f5cc5a30f7d4657, []int{8}
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

type NatRequest struct {
	MsgType              NatMsgType        `protobuf:"varint,1,opt,name=msgType,proto3,enum=net.pb.NatMsgType" json:"msgType,omitempty"`
	BootRegReq           *BootNatRegReq    `protobuf:"bytes,2,opt,name=BootRegReq,proto3" json:"BootRegReq,omitempty"`
	ConnReq              *NatConInvite     `protobuf:"bytes,3,opt,name=connReq,proto3" json:"connReq,omitempty"`
	KeepAlive            *NatKeepAlive     `protobuf:"bytes,4,opt,name=keepAlive,proto3" json:"keepAlive,omitempty"`
	Ping                 *NatPing          `protobuf:"bytes,5,opt,name=ping,proto3" json:"ping,omitempty"`
	HoleMsg              *HoleDig          `protobuf:"bytes,6,opt,name=holeMsg,proto3" json:"holeMsg,omitempty"`
	Invite               *ReverseInvite    `protobuf:"bytes,7,opt,name=invite,proto3" json:"invite,omitempty"`
	InviteAck            *ReverseInviteAck `protobuf:"bytes,8,opt,name=inviteAck,proto3" json:"inviteAck,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *NatRequest) Reset()         { *m = NatRequest{} }
func (m *NatRequest) String() string { return proto.CompactTextString(m) }
func (*NatRequest) ProtoMessage()    {}
func (*NatRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f5cc5a30f7d4657, []int{9}
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

func (m *NatRequest) GetConnReq() *NatConInvite {
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

func (m *NatRequest) GetHoleMsg() *HoleDig {
	if m != nil {
		return m.HoleMsg
	}
	return nil
}

func (m *NatRequest) GetInvite() *ReverseInvite {
	if m != nil {
		return m.Invite
	}
	return nil
}

func (m *NatRequest) GetInviteAck() *ReverseInviteAck {
	if m != nil {
		return m.InviteAck
	}
	return nil
}

type NatResponse struct {
	MsgType              NatMsgType     `protobuf:"varint,1,opt,name=msgType,proto3,enum=net.pb.NatMsgType" json:"msgType,omitempty"`
	BootRegRes           *BootNatRegRes `protobuf:"bytes,2,opt,name=BootRegRes,proto3" json:"BootRegRes,omitempty"`
	ConnRes              *NatConInvite  `protobuf:"bytes,3,opt,name=connRes,proto3" json:"connRes,omitempty"`
	KeepAlive            *NatKeepAlive  `protobuf:"bytes,4,opt,name=keepAlive,proto3" json:"keepAlive,omitempty"`
	Pong                 *NatPing       `protobuf:"bytes,5,opt,name=pong,proto3" json:"pong,omitempty"`
	HoleMsg              *HoleDig       `protobuf:"bytes,6,opt,name=holeMsg,proto3" json:"holeMsg,omitempty"`
	Error                *ErrorNotify   `protobuf:"bytes,50,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *NatResponse) Reset()         { *m = NatResponse{} }
func (m *NatResponse) String() string { return proto.CompactTextString(m) }
func (*NatResponse) ProtoMessage()    {}
func (*NatResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f5cc5a30f7d4657, []int{10}
}

func (m *NatResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NatResponse.Unmarshal(m, b)
}
func (m *NatResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NatResponse.Marshal(b, m, deterministic)
}
func (m *NatResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NatResponse.Merge(m, src)
}
func (m *NatResponse) XXX_Size() int {
	return xxx_messageInfo_NatResponse.Size(m)
}
func (m *NatResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_NatResponse.DiscardUnknown(m)
}

var xxx_messageInfo_NatResponse proto.InternalMessageInfo

func (m *NatResponse) GetMsgType() NatMsgType {
	if m != nil {
		return m.MsgType
	}
	return NatMsgType_UnknownReq
}

func (m *NatResponse) GetBootRegRes() *BootNatRegRes {
	if m != nil {
		return m.BootRegRes
	}
	return nil
}

func (m *NatResponse) GetConnRes() *NatConInvite {
	if m != nil {
		return m.ConnRes
	}
	return nil
}

func (m *NatResponse) GetKeepAlive() *NatKeepAlive {
	if m != nil {
		return m.KeepAlive
	}
	return nil
}

func (m *NatResponse) GetPong() *NatPing {
	if m != nil {
		return m.Pong
	}
	return nil
}

func (m *NatResponse) GetHoleMsg() *HoleDig {
	if m != nil {
		return m.HoleMsg
	}
	return nil
}

func (m *NatResponse) GetError() *ErrorNotify {
	if m != nil {
		return m.Error
	}
	return nil
}

func init() {
	proto.RegisterEnum("net.pb.NatMsgType", NatMsgType_name, NatMsgType_value)
	proto.RegisterEnum("net.pb.NatType", NatType_name, NatType_value)
	proto.RegisterType((*BootNatRegReq)(nil), "net.pb.BootNatRegReq")
	proto.RegisterType((*BootNatRegRes)(nil), "net.pb.BootNatRegRes")
	proto.RegisterType((*NatConInvite)(nil), "net.pb.NatConInvite")
	proto.RegisterType((*NatKeepAlive)(nil), "net.pb.NatKeepAlive")
	proto.RegisterType((*NatPing)(nil), "net.pb.NatPing")
	proto.RegisterType((*HoleDig)(nil), "net.pb.HoleDig")
	proto.RegisterType((*ErrorNotify)(nil), "net.pb.ErrorNotify")
	proto.RegisterType((*ReverseInvite)(nil), "net.pb.ReverseInvite")
	proto.RegisterType((*ReverseInviteAck)(nil), "net.pb.ReverseInviteAck")
	proto.RegisterType((*NatRequest)(nil), "net.pb.natRequest")
	proto.RegisterType((*NatResponse)(nil), "net.pb.natResponse")
}

func init() { proto.RegisterFile("natManager.proto", fileDescriptor_2f5cc5a30f7d4657) }

var fileDescriptor_2f5cc5a30f7d4657 = []byte{
	// 774 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x55, 0x5d, 0x6f, 0xe3, 0x44,
	0x14, 0xdd, 0x7c, 0xd9, 0xc9, 0xcd, 0xa6, 0x3b, 0x0c, 0x05, 0x59, 0x2b, 0x84, 0x56, 0x46, 0x68,
	0xe9, 0x02, 0x11, 0x2a, 0x82, 0xf7, 0x7c, 0xac, 0x44, 0x54, 0x6a, 0xa2, 0x69, 0x78, 0x40, 0x48,
	0x48, 0x8e, 0x73, 0xeb, 0x8e, 0x92, 0xce, 0xb8, 0x9e, 0x49, 0x50, 0x5f, 0x79, 0x46, 0xbc, 0xf3,
	0xb3, 0xf8, 0x47, 0x68, 0xc6, 0x63, 0xc7, 0xae, 0x4a, 0xf9, 0x7a, 0xf3, 0xbd, 0xf7, 0xdc, 0xb9,
	0x67, 0xee, 0x39, 0x93, 0x00, 0x11, 0xb1, 0xbe, 0x8c, 0x45, 0x9c, 0x62, 0x3e, 0xce, 0x72, 0xa9,
	0x25, 0xf5, 0x04, 0xea, 0x71, 0xb6, 0x7e, 0x39, 0x12, 0x6b, 0x35, 0xd9, 0x6c, 0x5c, 0x3a, 0x4c,
	0x61, 0x34, 0x95, 0x52, 0x47, 0xb1, 0x66, 0x98, 0x32, 0xbc, 0xa3, 0xef, 0x83, 0x27, 0xe4, 0x06,
	0x17, 0x9b, 0xa0, 0xf5, 0xaa, 0xf5, 0xc9, 0x80, 0xb9, 0x88, 0xbe, 0x82, 0x61, 0x96, 0xf3, 0x43,
	0xac, 0x71, 0x29, 0x73, 0x1d, 0x74, 0x6c, 0xb1, 0x9e, 0xa2, 0x1f, 0xc0, 0xc0, 0x85, 0x8b, 0x2c,
	0x68, 0xdb, 0xfa, 0x31, 0x11, 0xfe, 0xd6, 0x6a, 0x4e, 0x52, 0xf4, 0x0c, 0x7c, 0x11, 0xeb, 0xd5,
	0x7d, 0x86, 0x76, 0xd4, 0xc9, 0xf9, 0x8b, 0x71, 0xc1, 0x71, 0x1c, 0x15, 0x69, 0x56, 0xd6, 0x6b,
	0xa4, 0xda, 0x0d, 0x52, 0x2f, 0xa1, 0x9f, 0xed, 0xd7, 0x3b, 0x9e, 0x2c, 0x32, 0xc7, 0xa8, 0x8a,
	0xe9, 0x87, 0x00, 0xc5, 0xb7, 0xe5, 0xdb, 0xb5, 0xd5, 0x5a, 0x26, 0xfc, 0xa5, 0x05, 0xcf, 0xa3,
	0x58, 0xcf, 0xa4, 0x58, 0x88, 0x03, 0xd7, 0x48, 0x3f, 0x85, 0xfe, 0x75, 0x2e, 0x6f, 0xcd, 0x72,
	0x2c, 0xa1, 0x61, 0x8d, 0x50, 0xb1, 0x33, 0x56, 0x01, 0xe8, 0x6b, 0xf0, 0xb4, 0xb4, 0xd0, 0xf6,
	0xe3, 0x50, 0x57, 0x36, 0x5b, 0x51, 0xa8, 0x14, 0x97, 0x62, 0xb1, 0x71, 0x1c, 0x8f, 0x89, 0x70,
	0x67, 0x39, 0x5c, 0x20, 0x66, 0x93, 0x1d, 0x3f, 0xe0, 0x5f, 0x6e, 0xff, 0x14, 0x7a, 0xbb, 0x6a,
	0xda, 0x80, 0x15, 0x81, 0xc9, 0x2e, 0xf7, 0xeb, 0xc5, 0xd2, 0x9d, 0x5b, 0x04, 0x34, 0x00, 0x7f,
	0xb9, 0x5f, 0x57, 0xb7, 0xee, 0xb1, 0x32, 0x0c, 0x7f, 0x00, 0x3f, 0x8a, 0xf5, 0x92, 0x8b, 0x94,
	0x52, 0xe8, 0x66, 0x5c, 0xa4, 0x6e, 0x8c, 0xfd, 0xb6, 0x39, 0x29, 0x52, 0x37, 0xc3, 0x7e, 0x9b,
	0x11, 0x42, 0x8a, 0x04, 0xcb, 0x11, 0x36, 0xa0, 0x04, 0x3a, 0xab, 0xd5, 0xb7, 0xee, 0x78, 0xf3,
	0x19, 0xbe, 0x06, 0xff, 0x1b, 0xb9, 0xc3, 0x39, 0x4f, 0x9b, 0x37, 0x6e, 0x3d, 0xbc, 0xf1, 0xc7,
	0x30, 0x7c, 0x9b, 0xe7, 0x32, 0x8f, 0xa4, 0xe6, 0xd7, 0xf7, 0xe6, 0xc2, 0x98, 0xe7, 0x97, 0xaa,
	0x64, 0xe2, 0xa2, 0xf0, 0x47, 0x18, 0x31, 0x3c, 0x60, 0xae, 0xd0, 0xa9, 0xf3, 0xe4, 0xa9, 0x86,
	0x66, 0xb6, 0x5f, 0x57, 0xbe, 0x2b, 0x02, 0x73, 0xb8, 0x96, 0x95, 0x5d, 0x7b, 0xcc, 0x45, 0xe1,
	0x17, 0x40, 0x1a, 0x87, 0x4f, 0x92, 0xed, 0xdf, 0xb0, 0xfe, 0xb5, 0x03, 0x20, 0x8c, 0x73, 0xef,
	0xf6, 0xa8, 0x34, 0xfd, 0x0c, 0xfc, 0x5b, 0x95, 0xd6, 0xac, 0x4b, 0x6b, 0xd6, 0xbd, 0x2c, 0x2a,
	0xac, 0x84, 0xd0, 0xaf, 0x00, 0x8c, 0xf3, 0x8b, 0x07, 0xe6, 0xfc, 0xf2, 0x5e, 0xd9, 0xd0, 0x78,
	0x7d, 0xac, 0x06, 0xa4, 0x63, 0xf0, 0x13, 0x29, 0x84, 0xe9, 0xe9, 0xd8, 0x9e, 0xd3, 0xda, 0x90,
	0xca, 0xb6, 0xac, 0x04, 0xd1, 0x73, 0x18, 0x6c, 0x4b, 0x23, 0x59, 0x69, 0x9a, 0x1d, 0x95, 0xc9,
	0xd8, 0x11, 0x46, 0x3f, 0x72, 0x36, 0xe8, 0x3d, 0x30, 0x71, 0xe1, 0x12, 0xe7, 0x8b, 0x33, 0xf0,
	0x6f, 0xe4, 0x0e, 0x8d, 0x48, 0x5e, 0x13, 0xe7, 0x24, 0x67, 0x65, 0x9d, 0x7e, 0x0e, 0x1e, 0xb7,
	0xb4, 0x02, 0xbf, 0x79, 0xcd, 0xc6, 0xbe, 0x99, 0x03, 0xd1, 0xaf, 0x61, 0xc0, 0x4b, 0x05, 0x82,
	0xbe, 0xed, 0x08, 0x1e, 0xed, 0x98, 0x24, 0x5b, 0x76, 0x84, 0x86, 0x7f, 0xb4, 0x61, 0x68, 0xe5,
	0x50, 0x99, 0x14, 0x0a, 0xff, 0x87, 0x1e, 0xea, 0x29, 0x3d, 0x54, 0x4d, 0x0f, 0x75, 0xd4, 0x43,
	0xfd, 0x13, 0x3d, 0xd4, 0x7f, 0xd6, 0x43, 0x3e, 0xa5, 0x87, 0xfc, 0x77, 0x7a, 0x9c, 0x41, 0x0f,
	0xcd, 0x6b, 0x0b, 0xce, 0x2d, 0xf0, 0xdd, 0x12, 0x58, 0x7b, 0x82, 0xac, 0x40, 0xbc, 0xf9, 0xbd,
	0x05, 0x70, 0xdc, 0x16, 0x3d, 0x01, 0xf8, 0x5e, 0x6c, 0x85, 0xfc, 0xd9, 0x78, 0x8b, 0x3c, 0xa3,
	0x04, 0x9e, 0x9b, 0x5d, 0x5c, 0xe9, 0x3c, 0xce, 0x18, 0xa6, 0xa4, 0x45, 0x47, 0x30, 0xa8, 0xee,
	0x40, 0xda, 0x74, 0x08, 0xfe, 0x4c, 0x0a, 0x81, 0x89, 0x26, 0x1d, 0xda, 0x87, 0xae, 0x21, 0x4c,
	0xba, 0x74, 0x00, 0xbd, 0x39, 0x4f, 0x17, 0x82, 0xf4, 0x28, 0x80, 0x37, 0xe7, 0xe9, 0x77, 0x7b,
	0x4d, 0x3c, 0x73, 0xbc, 0x13, 0x78, 0xce, 0x53, 0xe2, 0xd3, 0x77, 0xaa, 0xf7, 0x3e, 0xe7, 0xe9,
	0x64, 0x76, 0x41, 0xfa, 0xa6, 0xd3, 0x32, 0x23, 0xa7, 0x6f, 0x7e, 0xb2, 0x3f, 0x5c, 0x0f, 0x79,
	0xbd, 0xbd, 0x22, 0xcf, 0xe8, 0x0b, 0x18, 0x46, 0x32, 0x8a, 0xf5, 0x1c, 0x0f, 0x3c, 0xc1, 0x82,
	0xd6, 0x14, 0x6f, 0xb8, 0xd8, 0x44, 0xb1, 0x26, 0x6d, 0x4a, 0xe1, 0x64, 0x16, 0x8b, 0x29, 0x46,
	0xb1, 0xbe, 0xc2, 0xfc, 0x80, 0x39, 0xe9, 0x98, 0x9e, 0x95, 0x9c, 0xe2, 0xec, 0x06, 0x93, 0x2d,
	0x6e, 0x48, 0x77, 0xed, 0xd9, 0x3f, 0xc3, 0x2f, 0xff, 0x0c, 0x00, 0x00, 0xff, 0xff, 0xf2, 0xad,
	0x30, 0xdf, 0x37, 0x07, 0x00, 0x00,
}
