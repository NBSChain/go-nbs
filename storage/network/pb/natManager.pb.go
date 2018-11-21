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
	NatMsgType_UnknownReq   NatMsgType = 0
	NatMsgType_BootStrapReg NatMsgType = 1
	NatMsgType_KeepAlive    NatMsgType = 2
	NatMsgType_Connect      NatMsgType = 3
	NatMsgType_Ping         NatMsgType = 4
	NatMsgType_holeAck      NatMsgType = 5
	NatMsgType_error        NatMsgType = 20
)

var NatMsgType_name = map[int32]string{
	0:  "UnknownReq",
	1:  "BootStrapReg",
	2:  "KeepAlive",
	3:  "Connect",
	4:  "Ping",
	5:  "holeAck",
	20: "error",
}

var NatMsgType_value = map[string]int32{
	"UnknownReq":   0,
	"BootStrapReg": 1,
	"KeepAlive":    2,
	"Connect":      3,
	"Ping":         4,
	"holeAck":      5,
	"error":        20,
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
	CType                int32    `protobuf:"varint,4,opt,name=cType,proto3" json:"cType,omitempty"`
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

func (m *NatConInvite) GetCType() int32 {
	if m != nil {
		return m.CType
	}
	return 0
}

type NbsAddr struct {
	NetworkId            string   `protobuf:"bytes,1,opt,name=networkId,proto3" json:"networkId,omitempty"`
	CanServer            bool     `protobuf:"varint,2,opt,name=canServer,proto3" json:"canServer,omitempty"`
	PubIp                string   `protobuf:"bytes,3,opt,name=pubIp,proto3" json:"pubIp,omitempty"`
	PubPort              int32    `protobuf:"varint,4,opt,name=pubPort,proto3" json:"pubPort,omitempty"`
	PriIp                string   `protobuf:"bytes,5,opt,name=priIp,proto3" json:"priIp,omitempty"`
	PriPort              int32    `protobuf:"varint,6,opt,name=priPort,proto3" json:"priPort,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NbsAddr) Reset()         { *m = NbsAddr{} }
func (m *NbsAddr) String() string { return proto.CompactTextString(m) }
func (*NbsAddr) ProtoMessage()    {}
func (*NbsAddr) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f5cc5a30f7d4657, []int{3}
}

func (m *NbsAddr) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NbsAddr.Unmarshal(m, b)
}
func (m *NbsAddr) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NbsAddr.Marshal(b, m, deterministic)
}
func (m *NbsAddr) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NbsAddr.Merge(m, src)
}
func (m *NbsAddr) XXX_Size() int {
	return xxx_messageInfo_NbsAddr.Size(m)
}
func (m *NbsAddr) XXX_DiscardUnknown() {
	xxx_messageInfo_NbsAddr.DiscardUnknown(m)
}

var xxx_messageInfo_NbsAddr proto.InternalMessageInfo

func (m *NbsAddr) GetNetworkId() string {
	if m != nil {
		return m.NetworkId
	}
	return ""
}

func (m *NbsAddr) GetCanServer() bool {
	if m != nil {
		return m.CanServer
	}
	return false
}

func (m *NbsAddr) GetPubIp() string {
	if m != nil {
		return m.PubIp
	}
	return ""
}

func (m *NbsAddr) GetPubPort() int32 {
	if m != nil {
		return m.PubPort
	}
	return 0
}

func (m *NbsAddr) GetPriIp() string {
	if m != nil {
		return m.PriIp
	}
	return ""
}

func (m *NbsAddr) GetPriPort() int32 {
	if m != nil {
		return m.PriPort
	}
	return 0
}

type NatKeepAlive struct {
	NodeId               string   `protobuf:"bytes,1,opt,name=nodeId,proto3" json:"nodeId,omitempty"`
	LAddr                string   `protobuf:"bytes,2,opt,name=lAddr,proto3" json:"lAddr,omitempty"`
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

func (m *NatKeepAlive) GetLAddr() string {
	if m != nil {
		return m.LAddr
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

type HoleAck struct {
	SessionId            string   `protobuf:"bytes,1,opt,name=sessionId,proto3" json:"sessionId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HoleAck) Reset()         { *m = HoleAck{} }
func (m *HoleAck) String() string { return proto.CompactTextString(m) }
func (*HoleAck) ProtoMessage()    {}
func (*HoleAck) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f5cc5a30f7d4657, []int{6}
}

func (m *HoleAck) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HoleAck.Unmarshal(m, b)
}
func (m *HoleAck) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HoleAck.Marshal(b, m, deterministic)
}
func (m *HoleAck) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HoleAck.Merge(m, src)
}
func (m *HoleAck) XXX_Size() int {
	return xxx_messageInfo_HoleAck.Size(m)
}
func (m *HoleAck) XXX_DiscardUnknown() {
	xxx_messageInfo_HoleAck.DiscardUnknown(m)
}

var xxx_messageInfo_HoleAck proto.InternalMessageInfo

func (m *HoleAck) GetSessionId() string {
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
	return fileDescriptor_2f5cc5a30f7d4657, []int{7}
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

type NatRequest struct {
	MsgType              NatMsgType     `protobuf:"varint,1,opt,name=msgType,proto3,enum=net.pb.NatMsgType" json:"msgType,omitempty"`
	BootRegReq           *BootNatRegReq `protobuf:"bytes,2,opt,name=BootRegReq,proto3" json:"BootRegReq,omitempty"`
	ConnReq              *NatConInvite  `protobuf:"bytes,3,opt,name=connReq,proto3" json:"connReq,omitempty"`
	KeepAlive            *NatKeepAlive  `protobuf:"bytes,4,opt,name=keepAlive,proto3" json:"keepAlive,omitempty"`
	Ping                 *NatPing       `protobuf:"bytes,5,opt,name=ping,proto3" json:"ping,omitempty"`
	HoleAck              *HoleAck       `protobuf:"bytes,6,opt,name=holeAck,proto3" json:"holeAck,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *NatRequest) Reset()         { *m = NatRequest{} }
func (m *NatRequest) String() string { return proto.CompactTextString(m) }
func (*NatRequest) ProtoMessage()    {}
func (*NatRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f5cc5a30f7d4657, []int{8}
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

func (m *NatRequest) GetHoleAck() *HoleAck {
	if m != nil {
		return m.HoleAck
	}
	return nil
}

type Response struct {
	MsgType              NatMsgType     `protobuf:"varint,1,opt,name=msgType,proto3,enum=net.pb.NatMsgType" json:"msgType,omitempty"`
	BootRegRes           *BootNatRegRes `protobuf:"bytes,2,opt,name=BootRegRes,proto3" json:"BootRegRes,omitempty"`
	ConnRes              *NatConInvite  `protobuf:"bytes,3,opt,name=connRes,proto3" json:"connRes,omitempty"`
	KeepAlive            *NatKeepAlive  `protobuf:"bytes,4,opt,name=keepAlive,proto3" json:"keepAlive,omitempty"`
	Pong                 *NatPing       `protobuf:"bytes,5,opt,name=pong,proto3" json:"pong,omitempty"`
	Error                *ErrorNotify   `protobuf:"bytes,6,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f5cc5a30f7d4657, []int{9}
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

func (m *Response) GetConnRes() *NatConInvite {
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

func (m *Response) GetError() *ErrorNotify {
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
	proto.RegisterType((*NbsAddr)(nil), "net.pb.NbsAddr")
	proto.RegisterType((*NatKeepAlive)(nil), "net.pb.NatKeepAlive")
	proto.RegisterType((*NatPing)(nil), "net.pb.NatPing")
	proto.RegisterType((*HoleAck)(nil), "net.pb.HoleAck")
	proto.RegisterType((*ErrorNotify)(nil), "net.pb.ErrorNotify")
	proto.RegisterType((*NatRequest)(nil), "net.pb.NatRequest")
	proto.RegisterType((*Response)(nil), "net.pb.response")
}

func init() { proto.RegisterFile("natManager.proto", fileDescriptor_2f5cc5a30f7d4657) }

var fileDescriptor_2f5cc5a30f7d4657 = []byte{
	// 711 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x55, 0x4d, 0x6f, 0xd3, 0x4c,
	0x10, 0x6e, 0x3e, 0x1c, 0x27, 0x93, 0x7e, 0x58, 0xfb, 0xe6, 0x7d, 0x15, 0xbd, 0x42, 0xa8, 0x32,
	0x42, 0xa5, 0x05, 0xe5, 0x50, 0xc4, 0x8d, 0x4b, 0x13, 0x2a, 0x11, 0x41, 0xad, 0x6a, 0x1b, 0x0e,
	0x5c, 0x90, 0x1c, 0x67, 0xea, 0x5a, 0x49, 0x77, 0xdd, 0xdd, 0x4d, 0xaa, 0xfe, 0x09, 0xce, 0xdc,
	0x51, 0xff, 0x27, 0xda, 0x0f, 0x3b, 0x4e, 0x05, 0x15, 0x82, 0xdb, 0xce, 0xcc, 0x33, 0x9e, 0x67,
	0x9f, 0x79, 0x36, 0x81, 0x80, 0xc5, 0xea, 0x2c, 0x66, 0x71, 0x8a, 0x62, 0x90, 0x0b, 0xae, 0x38,
	0x69, 0x31, 0x54, 0x83, 0x7c, 0x1a, 0xa6, 0xb0, 0x33, 0xe4, 0x5c, 0x45, 0xb1, 0xa2, 0x98, 0x52,
	0xbc, 0x21, 0xff, 0x41, 0x8b, 0xf1, 0x19, 0x8e, 0x67, 0xfd, 0xda, 0x7e, 0xed, 0x45, 0x87, 0xba,
	0x88, 0xec, 0x43, 0x37, 0x17, 0xd9, 0x2a, 0x56, 0x78, 0xce, 0x85, 0xea, 0x37, 0x4c, 0xb1, 0x9a,
	0x22, 0x4f, 0xa0, 0xe3, 0xc2, 0x71, 0xde, 0xaf, 0x9b, 0xfa, 0x3a, 0x11, 0x7e, 0xad, 0x6d, 0x4e,
	0x92, 0xe4, 0x10, 0x7c, 0x16, 0xab, 0xc9, 0x5d, 0x8e, 0x66, 0xd4, 0xee, 0xf1, 0xde, 0xc0, 0x92,
	0x1a, 0x44, 0x36, 0x4d, 0x8b, 0x7a, 0x85, 0x54, 0x7d, 0x83, 0xd4, 0xff, 0xd0, 0xce, 0x97, 0xd3,
	0x45, 0x96, 0x8c, 0x73, 0xc7, 0xa8, 0x8c, 0xc9, 0x53, 0x00, 0x7b, 0x36, 0x7c, 0x9b, 0xa6, 0x5a,
	0xc9, 0x84, 0xdf, 0x6a, 0xb0, 0x1d, 0xc5, 0x6a, 0xc4, 0xd9, 0x98, 0xad, 0x32, 0x85, 0xe4, 0x25,
	0xb4, 0x2f, 0x05, 0xbf, 0x3e, 0x99, 0xcd, 0x84, 0x21, 0xd4, 0xad, 0x10, 0x9a, 0x4a, 0x9d, 0xa6,
	0x25, 0x80, 0x1c, 0x40, 0x4b, 0x71, 0x03, 0xad, 0xff, 0x1c, 0xea, 0xca, 0x5a, 0x15, 0x89, 0x52,
	0x66, 0x9c, 0x8d, 0x67, 0x8e, 0xe3, 0x3a, 0x41, 0x7a, 0xe0, 0x25, 0x46, 0x01, 0xcd, 0xcf, 0xa3,
	0x36, 0x08, 0xef, 0x6b, 0xe0, 0xbb, 0xef, 0xe8, 0x7e, 0x86, 0xea, 0x96, 0x8b, 0x79, 0xb9, 0x92,
	0x75, 0x42, 0x57, 0x93, 0x98, 0x5d, 0xa0, 0x58, 0xa1, 0x65, 0xd2, 0xa6, 0xeb, 0x84, 0xfe, 0x7a,
	0xbe, 0x9c, 0x96, 0xda, 0xd8, 0x80, 0xf4, 0xc1, 0xcf, 0x97, 0xd3, 0x52, 0x15, 0x8f, 0x16, 0xa1,
	0xc1, 0x8b, 0x6c, 0x9c, 0xf7, 0x3d, 0x87, 0xd7, 0x81, 0xc1, 0x8b, 0xcc, 0xe0, 0x5b, 0x0e, 0x6f,
	0xc3, 0xf0, 0xad, 0x51, 0xf0, 0x03, 0x62, 0x7e, 0xb2, 0xc8, 0x56, 0xf8, 0x4b, 0xef, 0xf4, 0xc0,
	0x5b, 0x94, 0x5a, 0x75, 0xa8, 0x0d, 0xc2, 0xcf, 0xe0, 0x47, 0xb1, 0x3a, 0xcf, 0x58, 0x4a, 0x08,
	0x34, 0xf3, 0x8c, 0xa5, 0xae, 0xcd, 0x9c, 0x4d, 0x8e, 0xb3, 0xd4, 0xf5, 0x98, 0xb3, 0xfe, 0x10,
	0xe3, 0x2c, 0xc1, 0xe2, 0x42, 0x26, 0x20, 0x01, 0x34, 0x26, 0x93, 0x8f, 0xee, 0x32, 0xfa, 0x18,
	0x1e, 0x80, 0xff, 0x9e, 0x2f, 0xf0, 0x24, 0x99, 0x6f, 0xea, 0x5f, 0x7b, 0xa0, 0x7f, 0xf8, 0x1c,
	0xba, 0xa7, 0x42, 0x70, 0x11, 0x71, 0x95, 0x5d, 0xde, 0xe9, 0x0b, 0xa0, 0x10, 0x67, 0xb2, 0x60,
	0xe2, 0xa2, 0xf0, 0xbe, 0x0e, 0x60, 0x8c, 0x7b, 0xb3, 0x44, 0xa9, 0xc8, 0x2b, 0xf0, 0xaf, 0x65,
	0x5a, 0x71, 0x2e, 0xa9, 0x38, 0xf7, 0xcc, 0x56, 0x68, 0x01, 0x21, 0x6f, 0x00, 0xb4, 0xf1, 0xed,
	0xfb, 0x72, 0x76, 0xf9, 0xb7, 0x68, 0xd8, 0x78, 0x7c, 0xb4, 0x02, 0x24, 0x03, 0xf0, 0x13, 0xce,
	0x98, 0xee, 0x69, 0x98, 0x9e, 0x5e, 0x65, 0x48, 0xe9, 0x5a, 0x5a, 0x80, 0xc8, 0x31, 0x74, 0xe6,
	0xc5, 0x26, 0x8c, 0x16, 0x9b, 0x1d, 0xe5, 0x96, 0xe8, 0x1a, 0x46, 0x9e, 0x39, 0xdd, 0xbd, 0x07,
	0x1e, 0xb6, 0x6b, 0x71, 0x8b, 0x38, 0x04, 0xff, 0xca, 0x8a, 0x69, 0xf6, 0x5f, 0xc1, 0x39, 0x8d,
	0x69, 0x51, 0x0f, 0xbf, 0xd7, 0xa1, 0x2d, 0x50, 0xe6, 0x9c, 0x49, 0xfc, 0x0b, 0x95, 0xe4, 0x63,
	0x2a, 0xc9, 0x8a, 0x4a, 0x72, 0xad, 0x92, 0xfc, 0x1d, 0x95, 0xe4, 0x1f, 0xab, 0xc4, 0x1f, 0x53,
	0x89, 0x1b, 0x95, 0x3c, 0xd4, 0x4e, 0x72, 0x1a, 0xfd, 0x53, 0xa0, 0x2a, 0xf6, 0xa2, 0x16, 0x71,
	0xb4, 0x30, 0x66, 0x72, 0x0a, 0x90, 0x5d, 0x80, 0x4f, 0x6c, 0xce, 0xf8, 0xad, 0xde, 0x62, 0xb0,
	0x45, 0x02, 0xd8, 0xd6, 0xf7, 0xbb, 0x50, 0x22, 0xce, 0x29, 0xa6, 0x41, 0x8d, 0xec, 0x40, 0xa7,
	0xe4, 0x15, 0xd4, 0x49, 0x17, 0xfc, 0x11, 0x67, 0x0c, 0x13, 0x15, 0x34, 0x48, 0x1b, 0x9a, 0x9a,
	0x44, 0xd0, 0xd4, 0x69, 0xb7, 0x86, 0xc0, 0x23, 0x1d, 0xc7, 0x26, 0xe8, 0x1d, 0x7d, 0x31, 0xcf,
	0xec, 0xe1, 0xa8, 0xd3, 0x8b, 0x60, 0x8b, 0xec, 0x41, 0x37, 0xe2, 0x51, 0xac, 0xde, 0xe1, 0x2a,
	0x4b, 0xd0, 0x4e, 0x1a, 0xe2, 0x55, 0xc6, 0x66, 0x51, 0xac, 0x82, 0x3a, 0x21, 0xb0, 0x3b, 0x8a,
	0xd9, 0x10, 0xa3, 0x58, 0xd9, 0x5f, 0x94, 0xa0, 0xa1, 0x7b, 0x26, 0x7c, 0x88, 0xa3, 0x2b, 0x4c,
	0xe6, 0x38, 0x0b, 0x9a, 0xd3, 0x96, 0xf9, 0x43, 0x79, 0xfd, 0x23, 0x00, 0x00, 0xff, 0xff, 0x35,
	0xac, 0x1f, 0xf6, 0x64, 0x06, 0x00, 0x00,
}
