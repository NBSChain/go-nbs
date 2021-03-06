// Code generated by protoc-gen-go. DO NOT EDIT.
// source: account.proto

package account_pb

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

type Account struct {
	PeerID               string   `protobuf:"bytes,1,opt,name=peerID,proto3" json:"peerID,omitempty"`
	Encrypted            bool     `protobuf:"varint,2,opt,name=Encrypted,proto3" json:"Encrypted,omitempty"`
	PrivateKey           string   `protobuf:"bytes,3,opt,name=PrivateKey,proto3" json:"PrivateKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Account) Reset()         { *m = Account{} }
func (m *Account) String() string { return proto.CompactTextString(m) }
func (*Account) ProtoMessage()    {}
func (*Account) Descriptor() ([]byte, []int) {
	return fileDescriptor_8e28828dcb8d24f0, []int{0}
}

func (m *Account) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Account.Unmarshal(m, b)
}
func (m *Account) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Account.Marshal(b, m, deterministic)
}
func (m *Account) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Account.Merge(m, src)
}
func (m *Account) XXX_Size() int {
	return xxx_messageInfo_Account.Size(m)
}
func (m *Account) XXX_DiscardUnknown() {
	xxx_messageInfo_Account.DiscardUnknown(m)
}

var xxx_messageInfo_Account proto.InternalMessageInfo

func (m *Account) GetPeerID() string {
	if m != nil {
		return m.PeerID
	}
	return ""
}

func (m *Account) GetEncrypted() bool {
	if m != nil {
		return m.Encrypted
	}
	return false
}

func (m *Account) GetPrivateKey() string {
	if m != nil {
		return m.PrivateKey
	}
	return ""
}

func init() {
	proto.RegisterType((*Account)(nil), "account.pb.account")
}

func init() { proto.RegisterFile("account.proto", fileDescriptor_8e28828dcb8d24f0) }

var fileDescriptor_8e28828dcb8d24f0 = []byte{
	// 115 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4d, 0x4c, 0x4e, 0xce,
	0x2f, 0xcd, 0x2b, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x82, 0x73, 0x93, 0x94, 0xe2,
	0xb9, 0xd8, 0xa1, 0x3c, 0x21, 0x31, 0x2e, 0xb6, 0x82, 0xd4, 0xd4, 0x22, 0x4f, 0x17, 0x09, 0x46,
	0x05, 0x46, 0x0d, 0xce, 0x20, 0x28, 0x4f, 0x48, 0x86, 0x8b, 0xd3, 0x35, 0x2f, 0xb9, 0xa8, 0xb2,
	0xa0, 0x24, 0x35, 0x45, 0x82, 0x49, 0x81, 0x51, 0x83, 0x23, 0x08, 0x21, 0x20, 0x24, 0xc7, 0xc5,
	0x15, 0x50, 0x94, 0x59, 0x96, 0x58, 0x92, 0xea, 0x9d, 0x5a, 0x29, 0xc1, 0x0c, 0xd6, 0x89, 0x24,
	0x92, 0xc4, 0x06, 0xb6, 0xd3, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0xd3, 0x1e, 0xe9, 0x22, 0x84,
	0x00, 0x00, 0x00,
}
