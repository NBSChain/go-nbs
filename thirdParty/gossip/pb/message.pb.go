// Code generated by protoc-gen-go. DO NOT EDIT.
// source: message.proto

package pb

import (
	fmt "fmt"
	pb "github.com/NBSChain/go-nbs/storage/network/pb"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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
	SeqNo                int64                `protobuf:"varint,1,opt,name=SeqNo,proto3" json:"SeqNo,omitempty"`
	Expire               *timestamp.Timestamp `protobuf:"bytes,2,opt,name=Expire,proto3" json:"Expire,omitempty"`
	NodeId               string               `protobuf:"bytes,4,opt,name=nodeId,proto3" json:"nodeId,omitempty"`
	Addr                 *BasicHost           `protobuf:"bytes,3,opt,name=Addr,proto3" json:"Addr,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
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

func (m *Subscribe) GetExpire() *timestamp.Timestamp {
	if m != nil {
		return m.Expire
	}
	return nil
}

func (m *Subscribe) GetNodeId() string {
	if m != nil {
		return m.NodeId
	}
	return ""
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

type WeightUpdate struct {
	NodeId               string   `protobuf:"bytes,1,opt,name=nodeId,proto3" json:"nodeId,omitempty"`
	Weight               float64  `protobuf:"fixed64,2,opt,name=weight,proto3" json:"weight,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WeightUpdate) Reset()         { *m = WeightUpdate{} }
func (m *WeightUpdate) String() string { return proto.CompactTextString(m) }
func (*WeightUpdate) ProtoMessage()    {}
func (*WeightUpdate) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{7}
}

func (m *WeightUpdate) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WeightUpdate.Unmarshal(m, b)
}
func (m *WeightUpdate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WeightUpdate.Marshal(b, m, deterministic)
}
func (m *WeightUpdate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WeightUpdate.Merge(m, src)
}
func (m *WeightUpdate) XXX_Size() int {
	return xxx_messageInfo_WeightUpdate.Size(m)
}
func (m *WeightUpdate) XXX_DiscardUnknown() {
	xxx_messageInfo_WeightUpdate.DiscardUnknown(m)
}

var xxx_messageInfo_WeightUpdate proto.InternalMessageInfo

func (m *WeightUpdate) GetNodeId() string {
	if m != nil {
		return m.NodeId
	}
	return ""
}

func (m *WeightUpdate) GetWeight() float64 {
	if m != nil {
		return m.Weight
	}
	return 0
}

type Gossip struct {
	MsgType              pb.MsgType    `protobuf:"varint,1,opt,name=MsgType,proto3,enum=net.pb.MsgType" json:"MsgType,omitempty"`
	MsgId                string        `protobuf:"bytes,2,opt,name=MsgId,proto3" json:"MsgId,omitempty"`
	Subscribe            *Subscribe    `protobuf:"bytes,3,opt,name=Subscribe,proto3" json:"Subscribe,omitempty"`
	SubAck               *SynAck       `protobuf:"bytes,4,opt,name=SubAck,proto3" json:"SubAck,omitempty"`
	VoteContact          *VoteContact  `protobuf:"bytes,5,opt,name=VoteContact,proto3" json:"VoteContact,omitempty"`
	VoteResult           *Subscribe    `protobuf:"bytes,6,opt,name=VoteResult,proto3" json:"VoteResult,omitempty"`
	HeartBeat            *HeartBeat    `protobuf:"bytes,7,opt,name=HeartBeat,proto3" json:"HeartBeat,omitempty"`
	SubConfirm           *SynAck       `protobuf:"bytes,8,opt,name=SubConfirm,proto3" json:"SubConfirm,omitempty"`
	VoteAck              *SynAck       `protobuf:"bytes,9,opt,name=VoteAck,proto3" json:"VoteAck,omitempty"`
	ArcReplace           *ArcReplace   `protobuf:"bytes,10,opt,name=ArcReplace,proto3" json:"ArcReplace,omitempty"`
	ArcDrop              *ArcDrop      `protobuf:"bytes,11,opt,name=ArcDrop,proto3" json:"ArcDrop,omitempty"`
	ReplaceAck           *Subscribe    `protobuf:"bytes,12,opt,name=ReplaceAck,proto3" json:"ReplaceAck,omitempty"`
	OVWeight             *WeightUpdate `protobuf:"bytes,13,opt,name=OVWeight,proto3" json:"OVWeight,omitempty"`
	IVWeight             *WeightUpdate `protobuf:"bytes,14,opt,name=IVWeight,proto3" json:"IVWeight,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *Gossip) Reset()         { *m = Gossip{} }
func (m *Gossip) String() string { return proto.CompactTextString(m) }
func (*Gossip) ProtoMessage()    {}
func (*Gossip) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{8}
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

func (m *Gossip) GetOVWeight() *WeightUpdate {
	if m != nil {
		return m.OVWeight
	}
	return nil
}

func (m *Gossip) GetIVWeight() *WeightUpdate {
	if m != nil {
		return m.IVWeight
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
	proto.RegisterType((*WeightUpdate)(nil), "pb.WeightUpdate")
	proto.RegisterType((*Gossip)(nil), "pb.gossip")
}

func init() { proto.RegisterFile("message.proto", fileDescriptor_33c57e4bae7b9afd) }

var fileDescriptor_33c57e4bae7b9afd = []byte{
	// 674 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x94, 0xdd, 0x6e, 0x9b, 0x4a,
	0x10, 0xc7, 0x45, 0x6c, 0xe3, 0x78, 0xc8, 0x97, 0x56, 0x47, 0x11, 0x8a, 0x8e, 0x74, 0x12, 0xce,
	0x39, 0x52, 0xfa, 0x11, 0x50, 0x5d, 0xa9, 0x17, 0xbd, 0xa8, 0xe4, 0x38, 0xad, 0x82, 0x94, 0xb8,
	0xd6, 0xe2, 0xa6, 0x17, 0xbd, 0x5a, 0x60, 0x43, 0x50, 0x62, 0x96, 0xee, 0x2e, 0x4d, 0xf3, 0x16,
	0x7d, 0xa5, 0x3e, 0x48, 0xdf, 0xa5, 0xda, 0x65, 0xc1, 0x44, 0x69, 0xaa, 0xde, 0x31, 0xff, 0xf9,
	0x2d, 0xf3, 0x9f, 0xd9, 0x01, 0xd8, 0x5c, 0x52, 0x21, 0x48, 0x46, 0xfd, 0x92, 0x33, 0xc9, 0xd0,
	0x5a, 0x19, 0xef, 0xbd, 0xce, 0x72, 0x79, 0x55, 0xc5, 0x7e, 0xc2, 0x96, 0xc1, 0xec, 0x38, 0x9a,
	0x5e, 0x91, 0xbc, 0x08, 0x32, 0x76, 0x54, 0xc4, 0x22, 0x10, 0x92, 0x71, 0x92, 0xd1, 0xa0, 0xa0,
	0xf2, 0x96, 0xf1, 0xeb, 0xa0, 0x8c, 0x83, 0x22, 0x16, 0xe7, 0x22, 0xab, 0xcf, 0xef, 0xfd, 0x93,
	0x31, 0x96, 0xdd, 0xd0, 0x40, 0x47, 0x71, 0x75, 0x19, 0xc8, 0x7c, 0x49, 0x85, 0x24, 0xcb, 0xb2,
	0x06, 0xbc, 0xef, 0x16, 0x8c, 0x8e, 0x89, 0xc8, 0x93, 0x53, 0x26, 0x24, 0xfa, 0x1b, 0x46, 0x53,
	0x52, 0x44, 0x94, 0x7f, 0xa1, 0xdc, 0xb5, 0xf6, 0xad, 0xc3, 0x75, 0xbc, 0x12, 0x54, 0x76, 0x46,
	0xa4, 0xc9, 0xae, 0xed, 0x5b, 0x87, 0x23, 0xbc, 0x12, 0xd0, 0x5f, 0x30, 0x98, 0x11, 0x19, 0xce,
	0xdd, 0x9e, 0xce, 0xd4, 0x01, 0x72, 0x61, 0x38, 0x23, 0x72, 0xce, 0xb8, 0x74, 0xfb, 0xfb, 0xd6,
	0xe1, 0x00, 0x37, 0xa1, 0xe2, 0xe7, 0x3c, 0x0f, 0xe7, 0xee, 0xa0, 0xe6, 0x75, 0xa0, 0xd5, 0x2a,
	0x0e, 0x4b, 0xd7, 0x36, 0xaa, 0x0a, 0x74, 0xe5, 0xba, 0xc3, 0x30, 0x75, 0x87, 0xa6, 0x72, 0x23,
	0x78, 0xdf, 0x2c, 0x18, 0x45, 0x55, 0x2c, 0x12, 0x9e, 0xc7, 0x54, 0xbd, 0x21, 0xa2, 0x9f, 0x67,
	0x4c, 0xfb, 0xef, 0xe1, 0x3a, 0x40, 0x63, 0xb0, 0xdf, 0x7e, 0x2d, 0x73, 0x4e, 0xb5, 0x71, 0x67,
	0xbc, 0xe7, 0xd7, 0x93, 0xf1, 0x9b, 0xc9, 0xf8, 0x8b, 0x66, 0x32, 0xd8, 0x90, 0x68, 0x17, 0xec,
	0x82, 0xa5, 0x34, 0x4c, 0xb5, 0xf5, 0x11, 0x36, 0x11, 0x3a, 0x80, 0xfe, 0x24, 0x4d, 0xb9, 0x6e,
	0xd4, 0x19, 0x6f, 0xfa, 0x65, 0xec, 0xb7, 0x23, 0xc4, 0x3a, 0xe5, 0xbd, 0x02, 0x3b, 0xba, 0x2b,
	0x26, 0xc9, 0xf5, 0x23, 0x76, 0x76, 0xc1, 0x7e, 0xc7, 0xd9, 0x32, 0x4c, 0xcd, 0x1c, 0x4d, 0xe4,
	0x9d, 0x81, 0x73, 0xc1, 0x24, 0x9d, 0xb2, 0x42, 0x92, 0x44, 0xa2, 0x1d, 0xe8, 0x2d, 0x16, 0x67,
	0xfa, 0xe8, 0x00, 0xab, 0x47, 0xf4, 0xac, 0xd3, 0xaa, 0x69, 0x45, 0x1b, 0x68, 0x45, 0xbc, 0xca,
	0x7b, 0xff, 0xc2, 0xe8, 0x94, 0x12, 0x2e, 0x8f, 0x29, 0x91, 0x6d, 0xc9, 0x13, 0xfd, 0xba, 0xa6,
	0xe4, 0x89, 0xf7, 0x09, 0x60, 0xc2, 0x13, 0x4c, 0xcb, 0x1b, 0x92, 0xd0, 0x8e, 0xb1, 0x2e, 0x95,
	0x22, 0x04, 0xfd, 0x05, 0x6b, 0xed, 0xea, 0xe7, 0x3f, 0x99, 0xc3, 0x01, 0x0c, 0x27, 0x3c, 0x39,
	0xe1, 0xac, 0xec, 0x4c, 0xd3, 0xea, 0x4e, 0xd3, 0x7b, 0x03, 0x1b, 0x1f, 0x69, 0x9e, 0x5d, 0xc9,
	0x0f, 0x65, 0x4a, 0x24, 0x7d, 0x8c, 0x53, 0xfa, 0xad, 0xe6, 0xb4, 0x07, 0x0b, 0x9b, 0xc8, 0xfb,
	0xd1, 0x07, 0x3b, 0x63, 0x42, 0xe4, 0x25, 0x7a, 0x02, 0xc3, 0x73, 0x91, 0x2d, 0xee, 0x4a, 0xaa,
	0xcf, 0x6e, 0x8d, 0xb7, 0xfd, 0x82, 0x4a, 0xe5, 0xcb, 0xc8, 0xb8, 0xc9, 0xab, 0x6b, 0x39, 0x17,
	0x59, 0xdb, 0x50, 0x1d, 0xdc, 0x9f, 0x6e, 0xef, 0xf7, 0xd3, 0x45, 0x1e, 0xd8, 0x51, 0x15, 0x4f,
	0x92, 0x6b, 0xbd, 0x1e, 0xce, 0x18, 0x34, 0xa9, 0x6f, 0x1d, 0x9b, 0x0c, 0x7a, 0x71, 0xef, 0x3e,
	0xf5, 0xaa, 0x3b, 0xe3, 0x6d, 0x05, 0x76, 0x64, 0x7c, 0xef, 0xce, 0x8f, 0x00, 0x54, 0x88, 0xa9,
	0xa8, 0x6e, 0xa4, 0xfe, 0x0c, 0x1e, 0x98, 0xe8, 0x00, 0xca, 0x72, 0x7b, 0xc7, 0xfa, 0xd3, 0x30,
	0x74, 0x2b, 0xe2, 0xce, 0x0e, 0x3c, 0x05, 0x88, 0xaa, 0x78, 0xca, 0x8a, 0xcb, 0x9c, 0x2f, 0xdd,
	0xf5, 0x07, 0xb6, 0x3b, 0x59, 0xf4, 0x1f, 0x0c, 0x55, 0x19, 0xd5, 0xdf, 0xe8, 0x01, 0xd8, 0xa4,
	0x90, 0xdf, 0xdd, 0x1e, 0x17, 0x34, 0xb8, 0xa5, 0xc0, 0x95, 0x8a, 0xbb, 0xfb, 0xf5, 0x7f, 0xbb,
	0x10, 0xae, 0xa3, 0x61, 0xc7, 0xc0, 0x4a, 0xc2, 0xed, 0xb2, 0x1c, 0x01, 0x98, 0x13, 0xaa, 0xfe,
	0xc6, 0x2f, 0x87, 0xb0, 0x02, 0xd0, 0x73, 0x58, 0x7f, 0x7f, 0x51, 0x6f, 0x91, 0xbb, 0xa9, 0xe1,
	0x1d, 0x05, 0x77, 0xf7, 0x0a, 0xb7, 0x84, 0xa2, 0xc3, 0x86, 0xde, 0x7a, 0x8c, 0x6e, 0x88, 0xd8,
	0xd6, 0x7f, 0x88, 0x97, 0x3f, 0x03, 0x00, 0x00, 0xff, 0xff, 0xe9, 0x62, 0xfa, 0xea, 0x9a, 0x05,
	0x00, 0x00,
}
