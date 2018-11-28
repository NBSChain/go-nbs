// Code generated by protoc-gen-go. DO NOT EDIT.
// source: basicHost.proto

package pb

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

type BasicHost struct {
	CanServer            bool     `protobuf:"varint,1,opt,name=canServer,proto3" json:"canServer,omitempty"`
	NatServer            string   `protobuf:"bytes,2,opt,name=natServer,proto3" json:"natServer,omitempty"`
	NatIP                string   `protobuf:"bytes,3,opt,name=natIP,proto3" json:"natIP,omitempty"`
	NatPort              int32    `protobuf:"varint,4,opt,name=natPort,proto3" json:"natPort,omitempty"`
	PriIP                string   `protobuf:"bytes,5,opt,name=priIP,proto3" json:"priIP,omitempty"`
	PubIp                string   `protobuf:"bytes,6,opt,name=pubIp,proto3" json:"pubIp,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BasicHost) Reset()         { *m = BasicHost{} }
func (m *BasicHost) String() string { return proto.CompactTextString(m) }
func (*BasicHost) ProtoMessage()    {}
func (*BasicHost) Descriptor() ([]byte, []int) {
	return fileDescriptor_c96f4ee169cfb0c9, []int{0}
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

func init() {
	proto.RegisterType((*BasicHost)(nil), "pb.BasicHost")
}

func init() { proto.RegisterFile("basicHost.proto", fileDescriptor_c96f4ee169cfb0c9) }

var fileDescriptor_c96f4ee169cfb0c9 = []byte{
	// 176 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x44, 0x8f, 0xb1, 0x0e, 0x82, 0x30,
	0x10, 0x86, 0x53, 0x14, 0x94, 0x2e, 0xc6, 0xc6, 0xa1, 0x83, 0x43, 0xe3, 0xc4, 0xd4, 0xc5, 0x37,
	0x60, 0x92, 0xad, 0xc1, 0x27, 0x68, 0x09, 0x31, 0x0c, 0xf6, 0x2e, 0xed, 0xe9, 0x0b, 0xf9, 0xa2,
	0x86, 0x02, 0xe9, 0xf8, 0x7d, 0xdf, 0x3f, 0xdc, 0xf1, 0x93, 0xb3, 0x71, 0x1a, 0x1e, 0x10, 0x49,
	0x63, 0x00, 0x02, 0x51, 0xa0, 0xbb, 0xfd, 0x18, 0xaf, 0xdb, 0xcd, 0x8b, 0x2b, 0xaf, 0x07, 0xeb,
	0x9f, 0x63, 0xf8, 0x8e, 0x41, 0x32, 0xc5, 0x9a, 0x63, 0x9f, 0xc5, 0x5c, 0xbd, 0xa5, 0xb5, 0x16,
	0x8a, 0x35, 0x75, 0x9f, 0x85, 0xb8, 0xf0, 0xd2, 0x5b, 0xea, 0x8c, 0xdc, 0xa5, 0xb2, 0x80, 0x90,
	0xfc, 0xe0, 0x2d, 0x19, 0x08, 0x24, 0xf7, 0x8a, 0x35, 0x65, 0xbf, 0xe1, 0xbc, 0xc7, 0x30, 0x75,
	0x46, 0x96, 0xcb, 0x3e, 0x41, 0xb2, 0x1f, 0xd7, 0xa1, 0xac, 0x56, 0x3b, 0x43, 0xab, 0xf8, 0x79,
	0x80, 0xb7, 0xf6, 0x2e, 0xea, 0x17, 0xc4, 0x38, 0xa1, 0x46, 0xd7, 0xe6, 0xbb, 0x0d, 0x73, 0x55,
	0x7a, 0xe9, 0xfe, 0x0f, 0x00, 0x00, 0xff, 0xff, 0x8f, 0x68, 0xc9, 0x52, 0xe5, 0x00, 0x00, 0x00,
}
