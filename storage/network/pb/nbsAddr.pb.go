// Code generated by protoc-gen-go. DO NOT EDIT.
// source: nbsAddr.proto

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
	return fileDescriptor_b059c628fec3ac7c, []int{0}
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

func init() {
	proto.RegisterType((*NbsAddr)(nil), "net.pb.NbsAddr")
}

func init() { proto.RegisterFile("nbsAddr.proto", fileDescriptor_b059c628fec3ac7c) }

var fileDescriptor_b059c628fec3ac7c = []byte{
	// 159 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xcd, 0x4b, 0x2a, 0x76,
	0x4c, 0x49, 0x29, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0xcb, 0x4b, 0x2d, 0xd1, 0x2b,
	0x48, 0x52, 0x5a, 0xca, 0xc8, 0xc5, 0xee, 0x07, 0x91, 0x11, 0x92, 0xe1, 0xe2, 0xcc, 0x4b, 0x2d,
	0x29, 0xcf, 0x2f, 0xca, 0xf6, 0x4c, 0x91, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x42, 0x08, 0x80,
	0x64, 0x93, 0x13, 0xf3, 0x82, 0x53, 0x8b, 0xca, 0x52, 0x8b, 0x24, 0x98, 0x14, 0x18, 0x35, 0x38,
	0x82, 0x10, 0x02, 0x42, 0x22, 0x5c, 0xac, 0x05, 0xa5, 0x49, 0x9e, 0x05, 0x12, 0xcc, 0x60, 0x7d,
	0x10, 0x8e, 0x90, 0x04, 0x17, 0x7b, 0x41, 0x69, 0x52, 0x40, 0x7e, 0x51, 0x89, 0x04, 0x8b, 0x02,
	0xa3, 0x06, 0x6b, 0x10, 0x8c, 0x0b, 0x56, 0x5f, 0x94, 0xe9, 0x59, 0x20, 0xc1, 0x0a, 0x55, 0x0f,
	0xe2, 0x80, 0xd5, 0x17, 0x65, 0x82, 0xd5, 0xb3, 0x41, 0xd5, 0x43, 0xb8, 0x49, 0x6c, 0x60, 0x67,
	0x1b, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x34, 0x64, 0x16, 0xc3, 0xc7, 0x00, 0x00, 0x00,
}
