// Code generated by protoc-gen-go. DO NOT EDIT.
// source: message.proto

package pb

import (
	fmt "fmt"
	pb "github.com/NBSChain/go-nbs/storage/network/pb"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type BasicHost struct {
	CanServer            bool     `protobuf:"varint,1,opt,name=CanServer,proto3" json:"CanServer,omitempty"`
	NatServer            string   `protobuf:"bytes,2,opt,name=NatServer,proto3" json:"NatServer,omitempty"`
	NatIP                string   `protobuf:"bytes,3,opt,name=NatIP,proto3" json:"NatIP,omitempty"`
	NatPort              int32    `protobuf:"varint,4,opt,name=NatPort,proto3" json:"NatPort,omitempty"`
	PriIP                string   `protobuf:"bytes,5,opt,name=PriIP,proto3" json:"PriIP,omitempty"`
	PubIp                string   `protobuf:"bytes,6,opt,name=PubIp,proto3" json:"PubIp,omitempty"`
	NetworkId            string   `protobuf:"bytes,7,opt,name=NetworkId,proto3" json:"NetworkId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BasicHost) Reset()         { *m = BasicHost{} }
func (m *BasicHost) String() string { return proto.CompactTextString(m) }
func (*BasicHost) ProtoMessage()    {}
func (*BasicHost) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{0}
}

func (m *BasicHost) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BasicHost.Unmarshal(m, b)
}
func (m *BasicHost) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BasicHost.Marshal(b, m, deterministic)
}
func (m *BasicHost) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BasicHost.Merge(m, src)
}
func (m *BasicHost) XXX_Size() int {
	return xxx_messageInfo_BasicHost.Size(m)
}
func (m *BasicHost) XXX_DiscardUnknown() {
	xxx_messageInfo_BasicHost.DiscardUnknown(m)
}

var xxx_messageInfo_BasicHost proto.InternalMessageInfo

func (m *BasicHost) GetCanServer() bool {
	if m != nil {
		return m.CanServer
	}
	return false
}

func (m *BasicHost) GetNatServer() string {
	if m != nil {
		return m.NatServer
	}
	return ""
}

func (m *BasicHost) GetNatIP() string {
	if m != nil {
		return m.NatIP
	}
	return ""
}

func (m *BasicHost) GetNatPort() int32 {
	if m != nil {
		return m.NatPort
	}
	return 0
}

func (m *BasicHost) GetPriIP() string {
	if m != nil {
		return m.PriIP
	}
	return ""
}

func (m *BasicHost) GetPubIp() string {
	if m != nil {
		return m.PubIp
	}
	return ""
}

func (m *BasicHost) GetNetworkId() string {
	if m != nil {
		return m.NetworkId
	}
	return ""
}

type Subscribe struct {
	SeqNo                int64      `protobuf:"varint,1,opt,name=SeqNo,proto3" json:"SeqNo,omitempty"`
	Duration             int64      `protobuf:"varint,2,opt,name=Duration,proto3" json:"Duration,omitempty"`
	Addr                 *BasicHost `protobuf:"bytes,3,opt,name=Addr,proto3" json:"Addr,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *Subscribe) Reset()         { *m = Subscribe{} }
func (m *Subscribe) String() string { return proto.CompactTextString(m) }
func (*Subscribe) ProtoMessage()    {}
func (*Subscribe) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{1}
}

func (m *Subscribe) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Subscribe.Unmarshal(m, b)
}
func (m *Subscribe) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Subscribe.Marshal(b, m, deterministic)
}
func (m *Subscribe) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Subscribe.Merge(m, src)
}
func (m *Subscribe) XXX_Size() int {
	return xxx_messageInfo_Subscribe.Size(m)
}
func (m *Subscribe) XXX_DiscardUnknown() {
	xxx_messageInfo_Subscribe.DiscardUnknown(m)
}

var xxx_messageInfo_Subscribe proto.InternalMessageInfo

func (m *Subscribe) GetSeqNo() int64 {
	if m != nil {
		return m.SeqNo
	}
	return 0
}

func (m *Subscribe) GetDuration() int64 {
	if m != nil {
		return m.Duration
	}
	return 0
}

func (m *Subscribe) GetAddr() *BasicHost {
	if m != nil {
		return m.Addr
	}
	return nil
}

type SynAck struct {
	SeqNo                int64    `protobuf:"varint,1,opt,name=SeqNo,proto3" json:"SeqNo,omitempty"`
	FromId               string   `protobuf:"bytes,2,opt,name=FromId,proto3" json:"FromId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SynAck) Reset()         { *m = SynAck{} }
func (m *SynAck) String() string { return proto.CompactTextString(m) }
func (*SynAck) ProtoMessage()    {}
func (*SynAck) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{2}
}

func (m *SynAck) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SynAck.Unmarshal(m, b)
}
func (m *SynAck) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SynAck.Marshal(b, m, deterministic)
}
func (m *SynAck) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SynAck.Merge(m, src)
}
func (m *SynAck) XXX_Size() int {
	return xxx_messageInfo_SynAck.Size(m)
}
func (m *SynAck) XXX_DiscardUnknown() {
	xxx_messageInfo_SynAck.DiscardUnknown(m)
}

var xxx_messageInfo_SynAck proto.InternalMessageInfo

func (m *SynAck) GetSeqNo() int64 {
	if m != nil {
		return m.SeqNo
	}
	return 0
}

func (m *SynAck) GetFromId() string {
	if m != nil {
		return m.FromId
	}
	return ""
}

type VoteContact struct {
	TTL                  int32      `protobuf:"varint,1,opt,name=TTL,proto3" json:"TTL,omitempty"`
	Subscribe            *Subscribe `protobuf:"bytes,2,opt,name=Subscribe,proto3" json:"Subscribe,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *VoteContact) Reset()         { *m = VoteContact{} }
func (m *VoteContact) String() string { return proto.CompactTextString(m) }
func (*VoteContact) ProtoMessage()    {}
func (*VoteContact) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{3}
}

func (m *VoteContact) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VoteContact.Unmarshal(m, b)
}
func (m *VoteContact) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VoteContact.Marshal(b, m, deterministic)
}
func (m *VoteContact) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VoteContact.Merge(m, src)
}
func (m *VoteContact) XXX_Size() int {
	return xxx_messageInfo_VoteContact.Size(m)
}
func (m *VoteContact) XXX_DiscardUnknown() {
	xxx_messageInfo_VoteContact.DiscardUnknown(m)
}

var xxx_messageInfo_VoteContact proto.InternalMessageInfo

func (m *VoteContact) GetTTL() int32 {
	if m != nil {
		return m.TTL
	}
	return 0
}

func (m *VoteContact) GetSubscribe() *Subscribe {
	if m != nil {
		return m.Subscribe
	}
	return nil
}

type HeartBeat struct {
	FromID               string   `protobuf:"bytes,1,opt,name=FromID,proto3" json:"FromID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HeartBeat) Reset()         { *m = HeartBeat{} }
func (m *HeartBeat) String() string { return proto.CompactTextString(m) }
func (*HeartBeat) ProtoMessage()    {}
func (*HeartBeat) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{4}
}

func (m *HeartBeat) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HeartBeat.Unmarshal(m, b)
}
func (m *HeartBeat) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HeartBeat.Marshal(b, m, deterministic)
}
func (m *HeartBeat) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HeartBeat.Merge(m, src)
}
func (m *HeartBeat) XXX_Size() int {
	return xxx_messageInfo_HeartBeat.Size(m)
}
func (m *HeartBeat) XXX_DiscardUnknown() {
	xxx_messageInfo_HeartBeat.DiscardUnknown(m)
}

var xxx_messageInfo_HeartBeat proto.InternalMessageInfo

func (m *HeartBeat) GetFromID() string {
	if m != nil {
		return m.FromID
	}
	return ""
}

type ArcReplace struct {
	FromId               string     `protobuf:"bytes,1,opt,name=FromId,proto3" json:"FromId,omitempty"`
	ToId                 string     `protobuf:"bytes,2,opt,name=ToId,proto3" json:"ToId,omitempty"`
	Addr                 *BasicHost `protobuf:"bytes,3,opt,name=Addr,proto3" json:"Addr,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *ArcReplace) Reset()         { *m = ArcReplace{} }
func (m *ArcReplace) String() string { return proto.CompactTextString(m) }
func (*ArcReplace) ProtoMessage()    {}
func (*ArcReplace) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{5}
}

func (m *ArcReplace) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ArcReplace.Unmarshal(m, b)
}
func (m *ArcReplace) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ArcReplace.Marshal(b, m, deterministic)
}
func (m *ArcReplace) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ArcReplace.Merge(m, src)
}
func (m *ArcReplace) XXX_Size() int {
	return xxx_messageInfo_ArcReplace.Size(m)
}
func (m *ArcReplace) XXX_DiscardUnknown() {
	xxx_messageInfo_ArcReplace.DiscardUnknown(m)
}

var xxx_messageInfo_ArcReplace proto.InternalMessageInfo

func (m *ArcReplace) GetFromId() string {
	if m != nil {
		return m.FromId
	}
	return ""
}

func (m *ArcReplace) GetToId() string {
	if m != nil {
		return m.ToId
	}
	return ""
}

func (m *ArcReplace) GetAddr() *BasicHost {
	if m != nil {
		return m.Addr
	}
	return nil
}

type ArcDrop struct {
	NodeId               string   `protobuf:"bytes,1,opt,name=nodeId,proto3" json:"nodeId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ArcDrop) Reset()         { *m = ArcDrop{} }
func (m *ArcDrop) String() string { return proto.CompactTextString(m) }
func (*ArcDrop) ProtoMessage()    {}
func (*ArcDrop) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{6}
}

func (m *ArcDrop) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ArcDrop.Unmarshal(m, b)
}
func (m *ArcDrop) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ArcDrop.Marshal(b, m, deterministic)
}
func (m *ArcDrop) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ArcDrop.Merge(m, src)
}
func (m *ArcDrop) XXX_Size() int {
	return xxx_messageInfo_ArcDrop.Size(m)
}
func (m *ArcDrop) XXX_DiscardUnknown() {
	xxx_messageInfo_ArcDrop.DiscardUnknown(m)
}

var xxx_messageInfo_ArcDrop proto.InternalMessageInfo

func (m *ArcDrop) GetNodeId() string {
	if m != nil {
		return m.NodeId
	}
	return ""
}

type Gossip struct {
	MsgType              pb.MsgType   `protobuf:"varint,1,opt,name=MsgType,proto3,enum=net.pb.MsgType" json:"MsgType,omitempty"`
	MsgId                string       `protobuf:"bytes,2,opt,name=MsgId,proto3" json:"MsgId,omitempty"`
	Subscribe            *Subscribe   `protobuf:"bytes,3,opt,name=Subscribe,proto3" json:"Subscribe,omitempty"`
	SubAck               *SynAck      `protobuf:"bytes,4,opt,name=SubAck,proto3" json:"SubAck,omitempty"`
	VoteContact          *VoteContact `protobuf:"bytes,5,opt,name=VoteContact,proto3" json:"VoteContact,omitempty"`
	VoteResult           *Subscribe   `protobuf:"bytes,6,opt,name=VoteResult,proto3" json:"VoteResult,omitempty"`
	HeartBeat            *HeartBeat   `protobuf:"bytes,7,opt,name=HeartBeat,proto3" json:"HeartBeat,omitempty"`
	SubConfirm           *SynAck      `protobuf:"bytes,8,opt,name=SubConfirm,proto3" json:"SubConfirm,omitempty"`
	VoteAck              *SynAck      `protobuf:"bytes,9,opt,name=VoteAck,proto3" json:"VoteAck,omitempty"`
	ArcReplace           *ArcReplace  `protobuf:"bytes,10,opt,name=ArcReplace,proto3" json:"ArcReplace,omitempty"`
	ArcDrop              *ArcDrop     `protobuf:"bytes,11,opt,name=ArcDrop,proto3" json:"ArcDrop,omitempty"`
	ReplaceAck           *Subscribe   `protobuf:"bytes,12,opt,name=ReplaceAck,proto3" json:"ReplaceAck,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Gossip) Reset()         { *m = Gossip{} }
func (m *Gossip) String() string { return proto.CompactTextString(m) }
func (*Gossip) ProtoMessage()    {}
func (*Gossip) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{7}
}

func (m *Gossip) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Gossip.Unmarshal(m, b)
}
func (m *Gossip) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Gossip.Marshal(b, m, deterministic)
}
func (m *Gossip) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Gossip.Merge(m, src)
}
func (m *Gossip) XXX_Size() int {
	return xxx_messageInfo_Gossip.Size(m)
}
func (m *Gossip) XXX_DiscardUnknown() {
	xxx_messageInfo_Gossip.DiscardUnknown(m)
}

var xxx_messageInfo_Gossip proto.InternalMessageInfo

func (m *Gossip) GetMsgType() pb.MsgType {
	if m != nil {
		return m.MsgType
	}
	return pb.MsgType_Error
}

func (m *Gossip) GetMsgId() string {
	if m != nil {
		return m.MsgId
	}
	return ""
}

func (m *Gossip) GetSubscribe() *Subscribe {
	if m != nil {
		return m.Subscribe
	}
	return nil
}

func (m *Gossip) GetSubAck() *SynAck {
	if m != nil {
		return m.SubAck
	}
	return nil
}

func (m *Gossip) GetVoteContact() *VoteContact {
	if m != nil {
		return m.VoteContact
	}
	return nil
}

func (m *Gossip) GetVoteResult() *Subscribe {
	if m != nil {
		return m.VoteResult
	}
	return nil
}

func (m *Gossip) GetHeartBeat() *HeartBeat {
	if m != nil {
		return m.HeartBeat
	}
	return nil
}

func (m *Gossip) GetSubConfirm() *SynAck {
	if m != nil {
		return m.SubConfirm
	}
	return nil
}

func (m *Gossip) GetVoteAck() *SynAck {
	if m != nil {
		return m.VoteAck
	}
	return nil
}

func (m *Gossip) GetArcReplace() *ArcReplace {
	if m != nil {
		return m.ArcReplace
	}
	return nil
}

func (m *Gossip) GetArcDrop() *ArcDrop {
	if m != nil {
		return m.ArcDrop
	}
	return nil
}

func (m *Gossip) GetReplaceAck() *Subscribe {
	if m != nil {
		return m.ReplaceAck
	}
	return nil
}

func init() {
	proto.RegisterType((*BasicHost)(nil), "pb.BasicHost")
	proto.RegisterType((*Subscribe)(nil), "pb.Subscribe")
	proto.RegisterType((*SynAck)(nil), "pb.SynAck")
	proto.RegisterType((*VoteContact)(nil), "pb.VoteContact")
	proto.RegisterType((*HeartBeat)(nil), "pb.HeartBeat")
	proto.RegisterType((*ArcReplace)(nil), "pb.ArcReplace")
	proto.RegisterType((*ArcDrop)(nil), "pb.ArcDrop")
	proto.RegisterType((*Gossip)(nil), "pb.gossip")
}

func init() { proto.RegisterFile("message.proto", fileDescriptor_33c57e4bae7b9afd) }

var fileDescriptor_33c57e4bae7b9afd = []byte{
	// 572 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x54, 0x4f, 0x8b, 0xd3, 0x40,
	0x14, 0xa7, 0xdb, 0x36, 0xdd, 0xbc, 0xb8, 0xbb, 0x32, 0x88, 0x84, 0xc5, 0xc3, 0x6e, 0x54, 0xa8,
	0x4a, 0x13, 0xac, 0xe0, 0xc1, 0x5b, 0xff, 0x20, 0x5b, 0xd8, 0x96, 0x32, 0x29, 0x5e, 0xbc, 0x38,
	0x93, 0x8c, 0xd9, 0xd0, 0x6d, 0x26, 0xce, 0x4c, 0x94, 0x3d, 0xfa, 0xd1, 0xfc, 0x66, 0x32, 0x93,
	0x34, 0x49, 0xa9, 0x8a, 0xb7, 0xfe, 0xfe, 0xcc, 0x7b, 0xbf, 0xbc, 0xf7, 0x28, 0x9c, 0xed, 0x98,
	0x94, 0x24, 0x61, 0x7e, 0x2e, 0xb8, 0xe2, 0xe8, 0x24, 0xa7, 0x97, 0x1f, 0x92, 0x54, 0xdd, 0x15,
	0xd4, 0x8f, 0xf8, 0x2e, 0x58, 0x4d, 0xc3, 0xd9, 0x1d, 0x49, 0xb3, 0x20, 0xe1, 0xa3, 0x8c, 0xca,
	0x40, 0x2a, 0x2e, 0x48, 0xc2, 0x82, 0x8c, 0xa9, 0x1f, 0x5c, 0x6c, 0x83, 0x9c, 0x06, 0x19, 0x95,
	0x4b, 0x99, 0x94, 0xef, 0xbd, 0x5f, 0x1d, 0xb0, 0xa7, 0x44, 0xa6, 0xd1, 0x0d, 0x97, 0x0a, 0x3d,
	0x03, 0x7b, 0x46, 0xb2, 0x90, 0x89, 0xef, 0x4c, 0xb8, 0x9d, 0xab, 0xce, 0xf0, 0x14, 0x37, 0x84,
	0x56, 0x57, 0x44, 0x55, 0xea, 0xc9, 0x55, 0x67, 0x68, 0xe3, 0x86, 0x40, 0x4f, 0xa0, 0xbf, 0x22,
	0x6a, 0xb1, 0x76, 0xbb, 0x46, 0x29, 0x01, 0x72, 0x61, 0xb0, 0x22, 0x6a, 0xcd, 0x85, 0x72, 0x7b,
	0x57, 0x9d, 0x61, 0x1f, 0xef, 0xa1, 0xf6, 0xaf, 0x45, 0xba, 0x58, 0xbb, 0xfd, 0xd2, 0x6f, 0x80,
	0x61, 0x0b, 0xba, 0xc8, 0x5d, 0xab, 0x62, 0x35, 0x30, 0x9d, 0xcb, 0x0f, 0x58, 0xc4, 0xee, 0xa0,
	0xea, 0xbc, 0x27, 0xbc, 0x2f, 0x60, 0x87, 0x05, 0x95, 0x91, 0x48, 0x29, 0xd3, 0x05, 0x42, 0xf6,
	0x6d, 0xc5, 0x4d, 0xfc, 0x2e, 0x2e, 0x01, 0xba, 0x84, 0xd3, 0x79, 0x21, 0x88, 0x4a, 0x79, 0x66,
	0x92, 0x77, 0x71, 0x8d, 0xd1, 0x35, 0xf4, 0x26, 0x71, 0x2c, 0x4c, 0x6e, 0x67, 0x7c, 0xe6, 0xe7,
	0xd4, 0xaf, 0x27, 0x82, 0x8d, 0xe4, 0xbd, 0x07, 0x2b, 0x7c, 0xc8, 0x26, 0xd1, 0xf6, 0x2f, 0xe5,
	0x9f, 0x82, 0xf5, 0x51, 0xf0, 0xdd, 0x22, 0xae, 0xc6, 0x52, 0x21, 0xef, 0x16, 0x9c, 0x4f, 0x5c,
	0xb1, 0x19, 0xcf, 0x14, 0x89, 0x14, 0x7a, 0x0c, 0xdd, 0xcd, 0xe6, 0xd6, 0x3c, 0xed, 0x63, 0xfd,
	0x13, 0xbd, 0x69, 0x45, 0x37, 0x6f, 0xab, 0x00, 0x35, 0x89, 0x1b, 0xdd, 0x7b, 0x0e, 0xf6, 0x0d,
	0x23, 0x42, 0x4d, 0x19, 0x51, 0x75, 0xcb, 0xb9, 0x29, 0xb7, 0x6f, 0x39, 0xf7, 0x3e, 0x03, 0x4c,
	0x44, 0x84, 0x59, 0x7e, 0x4f, 0x22, 0xd6, 0x0a, 0xd6, 0x76, 0xc5, 0x08, 0x41, 0x6f, 0xc3, 0xeb,
	0xb8, 0xe6, 0xf7, 0xff, 0xcc, 0xe1, 0x1a, 0x06, 0x13, 0x11, 0xcd, 0x05, 0xcf, 0x75, 0xe5, 0x8c,
	0xc7, 0xac, 0xa9, 0x5c, 0x22, 0xef, 0x67, 0x0f, 0xac, 0x84, 0x4b, 0x99, 0xe6, 0xe8, 0x15, 0x0c,
	0x96, 0x32, 0xd9, 0x3c, 0xe4, 0xcc, 0x78, 0xce, 0xc7, 0x17, 0x7e, 0xc6, 0x94, 0xae, 0x5b, 0xd1,
	0x78, 0xaf, 0xeb, 0xb1, 0x2e, 0x65, 0x52, 0x07, 0x2a, 0xc1, 0xe1, 0x74, 0xba, 0xff, 0x9e, 0x0e,
	0xf2, 0xc0, 0x0a, 0x0b, 0x3a, 0x89, 0xb6, 0xe6, 0xd0, 0x9c, 0x31, 0x18, 0xa7, 0xd9, 0x1a, 0xae,
	0x14, 0xf4, 0xf6, 0x60, 0x1f, 0xe6, 0xf2, 0x9c, 0xf1, 0x85, 0x36, 0xb6, 0x68, 0x7c, 0xb0, 0xb3,
	0x11, 0x80, 0x86, 0x98, 0xc9, 0xe2, 0x5e, 0x99, 0xab, 0x3c, 0x0a, 0xd1, 0x32, 0xe8, 0xc8, 0xf5,
	0x8e, 0xcc, 0xa5, 0x56, 0xee, 0x9a, 0xc4, 0xad, 0x1d, 0xbe, 0x06, 0x08, 0x0b, 0x3a, 0xe3, 0xd9,
	0xd7, 0x54, 0xec, 0xdc, 0xd3, 0xa3, 0xd8, 0x2d, 0x15, 0xbd, 0x80, 0x81, 0x6e, 0xa3, 0xbf, 0xcf,
	0x3e, 0x32, 0xee, 0x25, 0xe4, 0xb7, 0xb7, 0xef, 0x82, 0x31, 0x9e, 0x6b, 0x63, 0xc3, 0xe2, 0xf6,
	0x7d, 0xbc, 0xac, 0x17, 0xea, 0x3a, 0xc6, 0xec, 0x54, 0x66, 0x4d, 0xe1, 0x7a, 0xd9, 0x23, 0x80,
	0xea, 0x85, 0xee, 0xff, 0xe8, 0x8f, 0x43, 0x68, 0x0c, 0xd4, 0x32, 0xff, 0x2d, 0xef, 0x7e, 0x07,
	0x00, 0x00, 0xff, 0xff, 0xc3, 0x6d, 0x99, 0xda, 0xac, 0x04, 0x00, 0x00,
}
