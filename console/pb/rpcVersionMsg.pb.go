// Code generated by protoc-gen-go. DO NOT EDIT.
// source: rpcVersionMsg.proto

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

type VersionRequest struct {
	CmdName              string   `protobuf:"bytes,1,opt,name=cmdName,proto3" json:"cmdName,omitempty"`
	Args                 []string `protobuf:"bytes,2,rep,name=args,proto3" json:"args,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VersionRequest) Reset()         { *m = VersionRequest{} }
func (m *VersionRequest) String() string { return proto.CompactTextString(m) }
func (*VersionRequest) ProtoMessage()    {}
func (*VersionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_691d697efa63b0aa, []int{0}
}

func (m *VersionRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VersionRequest.Unmarshal(m, b)
}
func (m *VersionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VersionRequest.Marshal(b, m, deterministic)
}
func (m *VersionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VersionRequest.Merge(m, src)
}
func (m *VersionRequest) XXX_Size() int {
	return xxx_messageInfo_VersionRequest.Size(m)
}
func (m *VersionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_VersionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_VersionRequest proto.InternalMessageInfo

func (m *VersionRequest) GetCmdName() string {
	if m != nil {
		return m.CmdName
	}
	return ""
}

func (m *VersionRequest) GetArgs() []string {
	if m != nil {
		return m.Args
	}
	return nil
}

type VersionResponse struct {
	Result               string   `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VersionResponse) Reset()         { *m = VersionResponse{} }
func (m *VersionResponse) String() string { return proto.CompactTextString(m) }
func (*VersionResponse) ProtoMessage()    {}
func (*VersionResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_691d697efa63b0aa, []int{1}
}

func (m *VersionResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VersionResponse.Unmarshal(m, b)
}
func (m *VersionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VersionResponse.Marshal(b, m, deterministic)
}
func (m *VersionResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VersionResponse.Merge(m, src)
}
func (m *VersionResponse) XXX_Size() int {
	return xxx_messageInfo_VersionResponse.Size(m)
}
func (m *VersionResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_VersionResponse.DiscardUnknown(m)
}

var xxx_messageInfo_VersionResponse proto.InternalMessageInfo

func (m *VersionResponse) GetResult() string {
	if m != nil {
		return m.Result
	}
	return ""
}

func init() {
	proto.RegisterType((*VersionRequest)(nil), "pb.VersionRequest")
	proto.RegisterType((*VersionResponse)(nil), "pb.VersionResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// VersionTaskClient is the client API for VersionTask service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type VersionTaskClient interface {
	SystemVersion(ctx context.Context, in *VersionRequest, opts ...grpc.CallOption) (*VersionResponse, error)
}

type versionTaskClient struct {
	cc *grpc.ClientConn
}

func NewVersionTaskClient(cc *grpc.ClientConn) VersionTaskClient {
	return &versionTaskClient{cc}
}

func (c *versionTaskClient) SystemVersion(ctx context.Context, in *VersionRequest, opts ...grpc.CallOption) (*VersionResponse, error) {
	out := new(VersionResponse)
	err := c.cc.Invoke(ctx, "/pb.VersionTask/SystemVersion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VersionTaskServer is the server API for VersionTask service.
type VersionTaskServer interface {
	SystemVersion(context.Context, *VersionRequest) (*VersionResponse, error)
}

func RegisterVersionTaskServer(s *grpc.Server, srv VersionTaskServer) {
	s.RegisterService(&_VersionTask_serviceDesc, srv)
}

func _VersionTask_SystemVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VersionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VersionTaskServer).SystemVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.VersionTask/SystemVersion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VersionTaskServer).SystemVersion(ctx, req.(*VersionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _VersionTask_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.VersionTask",
	HandlerType: (*VersionTaskServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SystemVersion",
			Handler:    _VersionTask_SystemVersion_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rpcVersionMsg.proto",
}

func init() { proto.RegisterFile("rpcVersionMsg.proto", fileDescriptor_691d697efa63b0aa) }

var fileDescriptor_691d697efa63b0aa = []byte{
	// 199 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x8f, 0xbf, 0x4b, 0xc0, 0x30,
	0x10, 0x85, 0x6d, 0x95, 0x4a, 0x4f, 0x54, 0xbc, 0x82, 0x14, 0xa7, 0xd2, 0xa9, 0x5d, 0x32, 0xe8,
	0xe6, 0xe0, 0xd0, 0xcd, 0x41, 0x91, 0x28, 0xee, 0x49, 0x3c, 0x8a, 0xd8, 0xfc, 0x30, 0x97, 0x0e,
	0xfe, 0xf7, 0x42, 0x8d, 0x4a, 0xb7, 0x7b, 0x1f, 0xf7, 0xb8, 0xfb, 0xa0, 0x89, 0xc1, 0xbc, 0x52,
	0xe4, 0x77, 0xef, 0x1e, 0x78, 0x16, 0x21, 0xfa, 0xe4, 0xb1, 0x0c, 0xba, 0xbf, 0x83, 0xb3, 0xcc,
	0x25, 0x7d, 0xae, 0xc4, 0x09, 0x5b, 0x38, 0x36, 0xf6, 0xed, 0x51, 0x59, 0x6a, 0x8b, 0xae, 0x18,
	0x6a, 0xf9, 0x1b, 0x11, 0xe1, 0x48, 0xc5, 0x99, 0xdb, 0xb2, 0x3b, 0x1c, 0x6a, 0xb9, 0xcd, 0xfd,
	0x08, 0xe7, 0x7f, 0x7d, 0x0e, 0xde, 0x31, 0xe1, 0x25, 0x54, 0x91, 0x78, 0x5d, 0x52, 0xee, 0xe7,
	0x74, 0x7d, 0x0f, 0x27, 0x79, 0xf5, 0x45, 0xf1, 0x07, 0xde, 0xc2, 0xe9, 0xf3, 0x17, 0x27, 0xb2,
	0x19, 0x22, 0x8a, 0xa0, 0xc5, 0xfe, 0x99, 0xab, 0x66, 0xc7, 0x7e, 0x0e, 0xf4, 0x07, 0xd3, 0x08,
	0x68, 0xbc, 0x15, 0x4e, 0xb3, 0x30, 0xde, 0xb1, 0x5f, 0x48, 0x04, 0x3d, 0x5d, 0xc8, 0x7f, 0x49,
	0x62, 0x56, 0x33, 0x3d, 0x15, 0xba, 0xda, 0x5c, 0x6f, 0xbe, 0x03, 0x00, 0x00, 0xff, 0xff, 0xf1,
	0x94, 0x6e, 0x2b, 0x02, 0x01, 0x00, 0x00,
}
