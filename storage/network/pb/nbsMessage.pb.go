// Code generated by protoc-gen-go. DO NOT EDIT.
// source: nbsMessage.proto

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

type NbsMessage struct {
	Data                 []byte      `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	NbsAddr              *NbsAddress `protobuf:"bytes,2,opt,name=nbsAddr,proto3" json:"nbsAddr,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *NbsMessage) Reset()         { *m = NbsMessage{} }
func (m *NbsMessage) String() string { return proto.CompactTextString(m) }
func (*NbsMessage) ProtoMessage()    {}
func (*NbsMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_2df5997504409e8e, []int{0}
}

func (m *NbsMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NbsMessage.Unmarshal(m, b)
}
func (m *NbsMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NbsMessage.Marshal(b, m, deterministic)
}
func (m *NbsMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NbsMessage.Merge(m, src)
}
func (m *NbsMessage) XXX_Size() int {
	return xxx_messageInfo_NbsMessage.Size(m)
}
func (m *NbsMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_NbsMessage.DiscardUnknown(m)
}

var xxx_messageInfo_NbsMessage proto.InternalMessageInfo

func (m *NbsMessage) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *NbsMessage) GetNbsAddr() *NbsAddress {
	if m != nil {
		return m.NbsAddr
	}
	return nil
}

func init() {
	proto.RegisterType((*NbsMessage)(nil), "net.pb.NbsMessage")
}

func init() { proto.RegisterFile("nbsMessage.proto", fileDescriptor_2df5997504409e8e) }

var fileDescriptor_2df5997504409e8e = []byte{
	// 115 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0xc8, 0x4b, 0x2a, 0xf6,
	0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0xcb, 0x4b,
	0x2d, 0xd1, 0x2b, 0x48, 0x92, 0x02, 0xc9, 0x38, 0xa6, 0xa4, 0x14, 0xa5, 0x16, 0x17, 0x43, 0x64,
	0x94, 0xfc, 0xb8, 0xb8, 0xfc, 0xe0, 0xaa, 0x85, 0x84, 0xb8, 0x58, 0x52, 0x12, 0x4b, 0x12, 0x25,
	0x18, 0x15, 0x18, 0x35, 0x78, 0x82, 0xc0, 0x6c, 0x21, 0x1d, 0x2e, 0x76, 0xa8, 0x2e, 0x09, 0x26,
	0x05, 0x46, 0x0d, 0x6e, 0x23, 0x21, 0x3d, 0x88, 0x69, 0x7a, 0x7e, 0x70, 0xc3, 0x82, 0x60, 0x4a,
	0x92, 0xd8, 0xc0, 0xc6, 0x1a, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0xed, 0x97, 0x32, 0x66, 0x84,
	0x00, 0x00, 0x00,
}