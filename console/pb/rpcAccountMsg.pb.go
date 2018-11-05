// Code generated by protoc-gen-go. DO NOT EDIT.
// source: rpcAccountMsg.proto

package pb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
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

type AccountUnlockRequest struct {
	Password             string   `protobuf:"bytes,1,opt,name=password,proto3" json:"password,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AccountUnlockRequest) Reset()         { *m = AccountUnlockRequest{} }
func (m *AccountUnlockRequest) String() string { return proto.CompactTextString(m) }
func (*AccountUnlockRequest) ProtoMessage()    {}
func (*AccountUnlockRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_9e3dd6410332a063, []int{0}
}

func (m *AccountUnlockRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AccountUnlockRequest.Unmarshal(m, b)
}
func (m *AccountUnlockRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AccountUnlockRequest.Marshal(b, m, deterministic)
}
func (m *AccountUnlockRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccountUnlockRequest.Merge(m, src)
}
func (m *AccountUnlockRequest) XXX_Size() int {
	return xxx_messageInfo_AccountUnlockRequest.Size(m)
}
func (m *AccountUnlockRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AccountUnlockRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AccountUnlockRequest proto.InternalMessageInfo

func (m *AccountUnlockRequest) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type AccountUnlockResponse struct {
	Message              string   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AccountUnlockResponse) Reset()         { *m = AccountUnlockResponse{} }
func (m *AccountUnlockResponse) String() string { return proto.CompactTextString(m) }
func (*AccountUnlockResponse) ProtoMessage()    {}
func (*AccountUnlockResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_9e3dd6410332a063, []int{1}
}

func (m *AccountUnlockResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AccountUnlockResponse.Unmarshal(m, b)
}
func (m *AccountUnlockResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AccountUnlockResponse.Marshal(b, m, deterministic)
}
func (m *AccountUnlockResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccountUnlockResponse.Merge(m, src)
}
func (m *AccountUnlockResponse) XXX_Size() int {
	return xxx_messageInfo_AccountUnlockResponse.Size(m)
}
func (m *AccountUnlockResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AccountUnlockResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AccountUnlockResponse proto.InternalMessageInfo

func (m *AccountUnlockResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type CreateAccountRequest struct {
	Password             string   `protobuf:"bytes,1,opt,name=password,proto3" json:"password,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateAccountRequest) Reset()         { *m = CreateAccountRequest{} }
func (m *CreateAccountRequest) String() string { return proto.CompactTextString(m) }
func (*CreateAccountRequest) ProtoMessage()    {}
func (*CreateAccountRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_9e3dd6410332a063, []int{2}
}

func (m *CreateAccountRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateAccountRequest.Unmarshal(m, b)
}
func (m *CreateAccountRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateAccountRequest.Marshal(b, m, deterministic)
}
func (m *CreateAccountRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateAccountRequest.Merge(m, src)
}
func (m *CreateAccountRequest) XXX_Size() int {
	return xxx_messageInfo_CreateAccountRequest.Size(m)
}
func (m *CreateAccountRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateAccountRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateAccountRequest proto.InternalMessageInfo

func (m *CreateAccountRequest) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type CreateAccountResponse struct {
	Message              string   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateAccountResponse) Reset()         { *m = CreateAccountResponse{} }
func (m *CreateAccountResponse) String() string { return proto.CompactTextString(m) }
func (*CreateAccountResponse) ProtoMessage()    {}
func (*CreateAccountResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_9e3dd6410332a063, []int{3}
}

func (m *CreateAccountResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateAccountResponse.Unmarshal(m, b)
}
func (m *CreateAccountResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateAccountResponse.Marshal(b, m, deterministic)
}
func (m *CreateAccountResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateAccountResponse.Merge(m, src)
}
func (m *CreateAccountResponse) XXX_Size() int {
	return xxx_messageInfo_CreateAccountResponse.Size(m)
}
func (m *CreateAccountResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateAccountResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CreateAccountResponse proto.InternalMessageInfo

func (m *CreateAccountResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*AccountUnlockRequest)(nil), "pb.AccountUnlockRequest")
	proto.RegisterType((*AccountUnlockResponse)(nil), "pb.AccountUnlockResponse")
	proto.RegisterType((*CreateAccountRequest)(nil), "pb.CreateAccountRequest")
	proto.RegisterType((*CreateAccountResponse)(nil), "pb.CreateAccountResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// AccountTaskClient is the client API for AccountTask service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AccountTaskClient interface {
	AccountUnlock(ctx context.Context, in *AccountUnlockRequest, opts ...grpc.CallOption) (*AccountUnlockResponse, error)
	CreateAccount(ctx context.Context, in *CreateAccountRequest, opts ...grpc.CallOption) (*CreateAccountResponse, error)
}

type accountTaskClient struct {
	cc *grpc.ClientConn
}

func NewAccountTaskClient(cc *grpc.ClientConn) AccountTaskClient {
	return &accountTaskClient{cc}
}

func (c *accountTaskClient) AccountUnlock(ctx context.Context, in *AccountUnlockRequest, opts ...grpc.CallOption) (*AccountUnlockResponse, error) {
	out := new(AccountUnlockResponse)
	err := c.cc.Invoke(ctx, "/pb.AccountTask/AccountUnlock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountTaskClient) CreateAccount(ctx context.Context, in *CreateAccountRequest, opts ...grpc.CallOption) (*CreateAccountResponse, error) {
	out := new(CreateAccountResponse)
	err := c.cc.Invoke(ctx, "/pb.AccountTask/CreateAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AccountTaskServer is the server API for AccountTask service.
type AccountTaskServer interface {
	AccountUnlock(context.Context, *AccountUnlockRequest) (*AccountUnlockResponse, error)
	CreateAccount(context.Context, *CreateAccountRequest) (*CreateAccountResponse, error)
}

func RegisterAccountTaskServer(s *grpc.Server, srv AccountTaskServer) {
	s.RegisterService(&_AccountTask_serviceDesc, srv)
}

func _AccountTask_AccountUnlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AccountUnlockRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountTaskServer).AccountUnlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.AccountTask/AccountUnlock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountTaskServer).AccountUnlock(ctx, req.(*AccountUnlockRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountTask_CreateAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountTaskServer).CreateAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.AccountTask/CreateAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountTaskServer).CreateAccount(ctx, req.(*CreateAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _AccountTask_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.AccountTask",
	HandlerType: (*AccountTaskServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AccountUnlock",
			Handler:    _AccountTask_AccountUnlock_Handler,
		},
		{
			MethodName: "CreateAccount",
			Handler:    _AccountTask_CreateAccount_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rpcAccountMsg.proto",
}

func init() { proto.RegisterFile("rpcAccountMsg.proto", fileDescriptor_9e3dd6410332a063) }

var fileDescriptor_9e3dd6410332a063 = []byte{
	// 220 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2e, 0x2a, 0x48, 0x76,
	0x4c, 0x4e, 0xce, 0x2f, 0xcd, 0x2b, 0xf1, 0x2d, 0x4e, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17,
	0x62, 0x2a, 0x48, 0x52, 0x32, 0xe2, 0x12, 0x81, 0x8a, 0x87, 0xe6, 0xe5, 0xe4, 0x27, 0x67, 0x07,
	0xa5, 0x16, 0x96, 0xa6, 0x16, 0x97, 0x08, 0x49, 0x71, 0x71, 0x14, 0x24, 0x16, 0x17, 0x97, 0xe7,
	0x17, 0xa5, 0x48, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0xc1, 0xf9, 0x4a, 0x86, 0x5c, 0xa2, 0x68,
	0x7a, 0x8a, 0x0b, 0xf2, 0xf3, 0x8a, 0x53, 0x85, 0x24, 0xb8, 0xd8, 0x73, 0x53, 0x8b, 0x8b, 0x13,
	0xd3, 0x53, 0xa1, 0x7a, 0x60, 0x5c, 0x90, 0x35, 0xce, 0x45, 0xa9, 0x89, 0x25, 0xa9, 0x50, 0x8d,
	0x44, 0x5a, 0x83, 0xa6, 0x87, 0x90, 0x35, 0x46, 0x73, 0x19, 0xb9, 0xb8, 0xa1, 0xaa, 0x43, 0x12,
	0x8b, 0xb3, 0x85, 0xdc, 0xb8, 0x78, 0x51, 0x5c, 0x2a, 0x24, 0xa1, 0x57, 0x90, 0xa4, 0x87, 0xcd,
	0xc3, 0x52, 0x92, 0x58, 0x64, 0x20, 0xf6, 0x29, 0x31, 0x80, 0xcc, 0x41, 0x71, 0x0a, 0xc4, 0x1c,
	0x6c, 0x3e, 0x82, 0x98, 0x83, 0xd5, 0xdd, 0x4a, 0x0c, 0x4e, 0x9a, 0x5c, 0x42, 0xc9, 0xf9, 0xb9,
	0x7a, 0x79, 0x49, 0xc5, 0x7a, 0xc9, 0xf9, 0x79, 0xc5, 0xf9, 0x39, 0xa9, 0x7a, 0x05, 0x49, 0x4e,
	0x82, 0x41, 0x88, 0xc8, 0x81, 0x78, 0x24, 0x80, 0x31, 0x89, 0x0d, 0x1c, 0x47, 0xc6, 0x80, 0x00,
	0x00, 0x00, 0xff, 0xff, 0xaa, 0xee, 0x42, 0xe2, 0xba, 0x01, 0x00, 0x00,
}