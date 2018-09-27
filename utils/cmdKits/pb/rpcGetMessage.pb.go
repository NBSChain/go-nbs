// Code generated by protoc-gen-go. DO NOT EDIT.
// source: rpcGetMessage.proto

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

type GetRequest struct {
	DataUri              string   `protobuf:"bytes,1,opt,name=dataUri,proto3" json:"dataUri,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetRequest) Reset()         { *m = GetRequest{} }
func (m *GetRequest) String() string { return proto.CompactTextString(m) }
func (*GetRequest) ProtoMessage()    {}
func (*GetRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4017a178f96b5b2b, []int{0}
}

func (m *GetRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetRequest.Unmarshal(m, b)
}
func (m *GetRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetRequest.Marshal(b, m, deterministic)
}
func (m *GetRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetRequest.Merge(m, src)
}
func (m *GetRequest) XXX_Size() int {
	return xxx_messageInfo_GetRequest.Size(m)
}
func (m *GetRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetRequest proto.InternalMessageInfo

func (m *GetRequest) GetDataUri() string {
	if m != nil {
		return m.DataUri
	}
	return ""
}

type GetResponse struct {
	Content              []byte   `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetResponse) Reset()         { *m = GetResponse{} }
func (m *GetResponse) String() string { return proto.CompactTextString(m) }
func (*GetResponse) ProtoMessage()    {}
func (*GetResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4017a178f96b5b2b, []int{1}
}

func (m *GetResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetResponse.Unmarshal(m, b)
}
func (m *GetResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetResponse.Marshal(b, m, deterministic)
}
func (m *GetResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetResponse.Merge(m, src)
}
func (m *GetResponse) XXX_Size() int {
	return xxx_messageInfo_GetResponse.Size(m)
}
func (m *GetResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetResponse proto.InternalMessageInfo

func (m *GetResponse) GetContent() []byte {
	if m != nil {
		return m.Content
	}
	return nil
}

func init() {
	proto.RegisterType((*GetRequest)(nil), "pb.GetRequest")
	proto.RegisterType((*GetResponse)(nil), "pb.GetResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// GetTaskClient is the client API for GetTask service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type GetTaskClient interface {
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (GetTask_GetClient, error)
}

type getTaskClient struct {
	cc *grpc.ClientConn
}

func NewGetTaskClient(cc *grpc.ClientConn) GetTaskClient {
	return &getTaskClient{cc}
}

func (c *getTaskClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (GetTask_GetClient, error) {
	stream, err := c.cc.NewStream(ctx, &_GetTask_serviceDesc.Streams[0], "/pb.GetTask/Get", opts...)
	if err != nil {
		return nil, err
	}
	x := &getTaskGetClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type GetTask_GetClient interface {
	Recv() (*GetResponse, error)
	grpc.ClientStream
}

type getTaskGetClient struct {
	grpc.ClientStream
}

func (x *getTaskGetClient) Recv() (*GetResponse, error) {
	m := new(GetResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetTaskServer is the server API for GetTask service.
type GetTaskServer interface {
	Get(*GetRequest, GetTask_GetServer) error
}

func RegisterGetTaskServer(s *grpc.Server, srv GetTaskServer) {
	s.RegisterService(&_GetTask_serviceDesc, srv)
}

func _GetTask_Get_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GetTaskServer).Get(m, &getTaskGetServer{stream})
}

type GetTask_GetServer interface {
	Send(*GetResponse) error
	grpc.ServerStream
}

type getTaskGetServer struct {
	grpc.ServerStream
}

func (x *getTaskGetServer) Send(m *GetResponse) error {
	return x.ServerStream.SendMsg(m)
}

var _GetTask_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.GetTask",
	HandlerType: (*GetTaskServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Get",
			Handler:       _GetTask_Get_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "rpcGetMessage.proto",
}

func init() { proto.RegisterFile("rpcGetMessage.proto", fileDescriptor_4017a178f96b5b2b) }

var fileDescriptor_4017a178f96b5b2b = []byte{
	// 191 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x8f, 0xcd, 0x4a, 0xc3, 0x40,
	0x14, 0x46, 0x1d, 0x05, 0x83, 0xd7, 0x3f, 0x18, 0x37, 0x41, 0x10, 0x24, 0x0b, 0x15, 0x17, 0x97,
	0xa0, 0xf8, 0x02, 0xd9, 0xcc, 0x42, 0x04, 0x19, 0xda, 0x07, 0x98, 0x99, 0x5c, 0x4a, 0x68, 0x33,
	0x33, 0xcd, 0xbd, 0x79, 0xff, 0xd2, 0x94, 0xb4, 0x5d, 0x1e, 0x38, 0x7c, 0x1f, 0x07, 0x9e, 0x86,
	0x1c, 0x0c, 0xc9, 0x1f, 0x31, 0xbb, 0x15, 0x61, 0x1e, 0x92, 0x24, 0x7d, 0x99, 0x7d, 0xf5, 0x06,
	0x60, 0x48, 0x2c, 0x6d, 0x47, 0x62, 0xd1, 0x25, 0x14, 0xad, 0x13, 0xb7, 0x1c, 0xba, 0x52, 0xbd,
	0xaa, 0x8f, 0x1b, 0x3b, 0x63, 0xf5, 0x0e, 0xb7, 0x93, 0xc7, 0x39, 0x45, 0xa6, 0xbd, 0x18, 0x52,
	0x14, 0x8a, 0x32, 0x89, 0x77, 0x76, 0xc6, 0xaf, 0x1f, 0x28, 0x0c, 0xc9, 0xc2, 0xf1, 0x5a, 0x7f,
	0xc2, 0x95, 0x21, 0xd1, 0x0f, 0x98, 0x3d, 0x9e, 0x4e, 0x9e, 0x1f, 0x8f, 0x7c, 0x18, 0xab, 0x2e,
	0x6a, 0xd5, 0xd4, 0xf0, 0x12, 0x52, 0x8f, 0xd1, 0x33, 0xc6, 0xd4, 0x12, 0x8e, 0xd2, 0x6d, 0x18,
	0x43, 0xdf, 0xfe, 0x76, 0xc2, 0x98, 0x7d, 0x73, 0x6f, 0xcf, 0x0b, 0xfe, 0x95, 0xbf, 0x9e, 0x22,
	0xbe, 0x77, 0x01, 0x00, 0x00, 0xff, 0xff, 0xd9, 0xe1, 0x9e, 0x0e, 0xdb, 0x00, 0x00, 0x00,
}
