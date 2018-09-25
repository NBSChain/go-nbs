// Code generated by protoc-gen-go. DO NOT EDIT.
// source: rpcAddMsg.proto

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

type FileType int32

const (
	FileType_FILE       FileType = 0
	FileType_LARGEFILE  FileType = 1
	FileType_DIRECTORY  FileType = 2
	FileType_SYSTEMLINK FileType = 3
)

var FileType_name = map[int32]string{
	0: "FILE",
	1: "LARGEFILE",
	2: "DIRECTORY",
	3: "SYSTEMLINK",
}

var FileType_value = map[string]int32{
	"FILE":       0,
	"LARGEFILE":  1,
	"DIRECTORY":  2,
	"SYSTEMLINK": 3,
}

func (x FileType) String() string {
	return proto.EnumName(FileType_name, int32(x))
}

func (FileType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_f5acc0f250d5865c, []int{0}
}

type FileChunk struct {
	Content              []byte   `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FileChunk) Reset()         { *m = FileChunk{} }
func (m *FileChunk) String() string { return proto.CompactTextString(m) }
func (*FileChunk) ProtoMessage()    {}
func (*FileChunk) Descriptor() ([]byte, []int) {
	return fileDescriptor_f5acc0f250d5865c, []int{0}
}

func (m *FileChunk) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FileChunk.Unmarshal(m, b)
}
func (m *FileChunk) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FileChunk.Marshal(b, m, deterministic)
}
func (m *FileChunk) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FileChunk.Merge(m, src)
}
func (m *FileChunk) XXX_Size() int {
	return xxx_messageInfo_FileChunk.Size(m)
}
func (m *FileChunk) XXX_DiscardUnknown() {
	xxx_messageInfo_FileChunk.DiscardUnknown(m)
}

var xxx_messageInfo_FileChunk proto.InternalMessageInfo

func (m *FileChunk) GetContent() []byte {
	if m != nil {
		return m.Content
	}
	return nil
}

type AddRequest struct {
	FileName             string   `protobuf:"bytes,1,opt,name=fileName,proto3" json:"fileName,omitempty"`
	FullPath             string   `protobuf:"bytes,2,opt,name=fullPath,proto3" json:"fullPath,omitempty"`
	FileSize             int64    `protobuf:"varint,3,opt,name=fileSize,proto3" json:"fileSize,omitempty"`
	FileType             FileType `protobuf:"varint,4,opt,name=fileType,proto3,enum=pb.FileType" json:"fileType,omitempty"`
	SplitterSize         int32    `protobuf:"varint,5,opt,name=splitterSize,proto3" json:"splitterSize,omitempty"`
	FileData             []byte   `protobuf:"bytes,6,opt,name=fileData,proto3" json:"fileData,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddRequest) Reset()         { *m = AddRequest{} }
func (m *AddRequest) String() string { return proto.CompactTextString(m) }
func (*AddRequest) ProtoMessage()    {}
func (*AddRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f5acc0f250d5865c, []int{1}
}

func (m *AddRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddRequest.Unmarshal(m, b)
}
func (m *AddRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddRequest.Marshal(b, m, deterministic)
}
func (m *AddRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddRequest.Merge(m, src)
}
func (m *AddRequest) XXX_Size() int {
	return xxx_messageInfo_AddRequest.Size(m)
}
func (m *AddRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AddRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AddRequest proto.InternalMessageInfo

func (m *AddRequest) GetFileName() string {
	if m != nil {
		return m.FileName
	}
	return ""
}

func (m *AddRequest) GetFullPath() string {
	if m != nil {
		return m.FullPath
	}
	return ""
}

func (m *AddRequest) GetFileSize() int64 {
	if m != nil {
		return m.FileSize
	}
	return 0
}

func (m *AddRequest) GetFileType() FileType {
	if m != nil {
		return m.FileType
	}
	return FileType_FILE
}

func (m *AddRequest) GetSplitterSize() int32 {
	if m != nil {
		return m.SplitterSize
	}
	return 0
}

func (m *AddRequest) GetFileData() []byte {
	if m != nil {
		return m.FileData
	}
	return nil
}

type AddResponse struct {
	Message              string   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	SessionId            string   `protobuf:"bytes,2,opt,name=sessionId,proto3" json:"sessionId,omitempty"`
	Name                 string   `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Hash                 string   `protobuf:"bytes,4,opt,name=hash,proto3" json:"hash,omitempty"`
	Bytes                int64    `protobuf:"varint,5,opt,name=bytes,proto3" json:"bytes,omitempty"`
	Size                 string   `protobuf:"bytes,6,opt,name=size,proto3" json:"size,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddResponse) Reset()         { *m = AddResponse{} }
func (m *AddResponse) String() string { return proto.CompactTextString(m) }
func (*AddResponse) ProtoMessage()    {}
func (*AddResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_f5acc0f250d5865c, []int{2}
}

func (m *AddResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddResponse.Unmarshal(m, b)
}
func (m *AddResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddResponse.Marshal(b, m, deterministic)
}
func (m *AddResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddResponse.Merge(m, src)
}
func (m *AddResponse) XXX_Size() int {
	return xxx_messageInfo_AddResponse.Size(m)
}
func (m *AddResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AddResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AddResponse proto.InternalMessageInfo

func (m *AddResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *AddResponse) GetSessionId() string {
	if m != nil {
		return m.SessionId
	}
	return ""
}

func (m *AddResponse) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *AddResponse) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

func (m *AddResponse) GetBytes() int64 {
	if m != nil {
		return m.Bytes
	}
	return 0
}

func (m *AddResponse) GetSize() string {
	if m != nil {
		return m.Size
	}
	return ""
}

func init() {
	proto.RegisterEnum("pb.FileType", FileType_name, FileType_value)
	proto.RegisterType((*FileChunk)(nil), "pb.FileChunk")
	proto.RegisterType((*AddRequest)(nil), "pb.AddRequest")
	proto.RegisterType((*AddResponse)(nil), "pb.AddResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// AddTaskClient is the client API for AddTask service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AddTaskClient interface {
	AddFile(ctx context.Context, in *AddRequest, opts ...grpc.CallOption) (*AddResponse, error)
	TransLargeFile(ctx context.Context, opts ...grpc.CallOption) (AddTask_TransLargeFileClient, error)
}

type addTaskClient struct {
	cc *grpc.ClientConn
}

func NewAddTaskClient(cc *grpc.ClientConn) AddTaskClient {
	return &addTaskClient{cc}
}

func (c *addTaskClient) AddFile(ctx context.Context, in *AddRequest, opts ...grpc.CallOption) (*AddResponse, error) {
	out := new(AddResponse)
	err := c.cc.Invoke(ctx, "/pb.AddTask/AddFile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *addTaskClient) TransLargeFile(ctx context.Context, opts ...grpc.CallOption) (AddTask_TransLargeFileClient, error) {
	stream, err := c.cc.NewStream(ctx, &_AddTask_serviceDesc.Streams[0], "/pb.AddTask/TransLargeFile", opts...)
	if err != nil {
		return nil, err
	}
	x := &addTaskTransLargeFileClient{stream}
	return x, nil
}

type AddTask_TransLargeFileClient interface {
	Send(*FileChunk) error
	CloseAndRecv() (*AddResponse, error)
	grpc.ClientStream
}

type addTaskTransLargeFileClient struct {
	grpc.ClientStream
}

func (x *addTaskTransLargeFileClient) Send(m *FileChunk) error {
	return x.ClientStream.SendMsg(m)
}

func (x *addTaskTransLargeFileClient) CloseAndRecv() (*AddResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(AddResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// AddTaskServer is the server API for AddTask service.
type AddTaskServer interface {
	AddFile(context.Context, *AddRequest) (*AddResponse, error)
	TransLargeFile(AddTask_TransLargeFileServer) error
}

func RegisterAddTaskServer(s *grpc.Server, srv AddTaskServer) {
	s.RegisterService(&_AddTask_serviceDesc, srv)
}

func _AddTask_AddFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AddTaskServer).AddFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.AddTask/AddFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AddTaskServer).AddFile(ctx, req.(*AddRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AddTask_TransLargeFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(AddTaskServer).TransLargeFile(&addTaskTransLargeFileServer{stream})
}

type AddTask_TransLargeFileServer interface {
	SendAndClose(*AddResponse) error
	Recv() (*FileChunk, error)
	grpc.ServerStream
}

type addTaskTransLargeFileServer struct {
	grpc.ServerStream
}

func (x *addTaskTransLargeFileServer) SendAndClose(m *AddResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *addTaskTransLargeFileServer) Recv() (*FileChunk, error) {
	m := new(FileChunk)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _AddTask_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.AddTask",
	HandlerType: (*AddTaskServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddFile",
			Handler:    _AddTask_AddFile_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "TransLargeFile",
			Handler:       _AddTask_TransLargeFile_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "rpcAddMsg.proto",
}

func init() { proto.RegisterFile("rpcAddMsg.proto", fileDescriptor_f5acc0f250d5865c) }

var fileDescriptor_f5acc0f250d5865c = []byte{
	// 421 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x92, 0xdf, 0x6a, 0xdb, 0x30,
	0x14, 0xc6, 0xab, 0x38, 0x49, 0xe3, 0xb3, 0x24, 0x0d, 0x62, 0x17, 0xa6, 0x6c, 0x10, 0x0c, 0x03,
	0x33, 0x86, 0x19, 0xdd, 0x5e, 0x20, 0x69, 0x93, 0x11, 0x9a, 0x76, 0x45, 0xf1, 0x4d, 0x2f, 0xe5,
	0x48, 0x4b, 0x4c, 0x6d, 0x59, 0xf3, 0x51, 0x2e, 0xba, 0x17, 0xd9, 0x3b, 0xed, 0xa9, 0x86, 0xe4,
	0x3f, 0x61, 0xb0, 0xbb, 0xf3, 0x3b, 0xdf, 0x39, 0x92, 0xbe, 0xa3, 0x03, 0x57, 0x95, 0xde, 0x2f,
	0x84, 0x78, 0xc0, 0x43, 0xac, 0xab, 0xd2, 0x94, 0xb4, 0xa7, 0xd3, 0xf0, 0x03, 0xf8, 0xeb, 0x2c,
	0x97, 0xb7, 0xc7, 0x93, 0x7a, 0xa1, 0x01, 0x5c, 0xee, 0x4b, 0x65, 0xa4, 0x32, 0x01, 0x99, 0x93,
	0x68, 0xcc, 0x5a, 0x0c, 0xff, 0x10, 0x80, 0x85, 0x10, 0x4c, 0xfe, 0x3c, 0x49, 0x34, 0xf4, 0x1a,
	0x46, 0x3f, 0xb2, 0x5c, 0x3e, 0xf2, 0x42, 0xba, 0x4a, 0x9f, 0x75, 0xec, 0xb4, 0x53, 0x9e, 0x3f,
	0x71, 0x73, 0x0c, 0x7a, 0x8d, 0xd6, 0x70, 0xdb, 0xb7, 0xcb, 0x7e, 0xc9, 0xc0, 0x9b, 0x93, 0xc8,
	0x63, 0x1d, 0xd3, 0xa8, 0xd6, 0x92, 0x57, 0x2d, 0x83, 0xfe, 0x9c, 0x44, 0xd3, 0x9b, 0x71, 0xac,
	0xd3, 0x78, 0xdd, 0xe4, 0x58, 0xa7, 0xd2, 0x10, 0xc6, 0xa8, 0xf3, 0xcc, 0x18, 0x59, 0xb9, 0x93,
	0x06, 0x73, 0x12, 0x0d, 0xd8, 0x3f, 0xb9, 0xf6, 0xa6, 0x3b, 0x6e, 0x78, 0x30, 0x74, 0x5e, 0x3a,
	0x0e, 0x7f, 0x13, 0x78, 0xe3, 0xcc, 0xa0, 0x2e, 0x15, 0x4a, 0x6b, 0xbb, 0x90, 0x88, 0xfc, 0xd0,
	0x9a, 0x69, 0x91, 0xbe, 0x03, 0x1f, 0x25, 0x62, 0x56, 0xaa, 0x8d, 0x68, 0xcc, 0x9c, 0x13, 0x94,
	0x42, 0x5f, 0xd9, 0x09, 0x78, 0x4e, 0x70, 0xb1, 0xcd, 0x1d, 0x39, 0x1e, 0x9d, 0x03, 0x9f, 0xb9,
	0x98, 0xbe, 0x85, 0x41, 0xfa, 0x6a, 0x24, 0xba, 0x87, 0x7a, 0xac, 0x06, 0x5b, 0x89, 0xf6, 0xf5,
	0xc3, 0xba, 0xd2, 0xc6, 0x1f, 0x97, 0x30, 0x6a, 0xfd, 0xd2, 0x11, 0xf4, 0xd7, 0x9b, 0xed, 0x6a,
	0x76, 0x41, 0x27, 0xe0, 0x6f, 0x17, 0xec, 0xdb, 0xca, 0x21, 0xb1, 0x78, 0xb7, 0x61, 0xab, 0xdb,
	0xe4, 0x3b, 0x7b, 0x9e, 0xf5, 0xe8, 0x14, 0x60, 0xf7, 0xbc, 0x4b, 0x56, 0x0f, 0xdb, 0xcd, 0xe3,
	0xfd, 0xcc, 0xbb, 0x29, 0xe0, 0x72, 0x21, 0x44, 0xc2, 0xf1, 0x85, 0x7e, 0x72, 0xa1, 0x3d, 0x91,
	0x4e, 0xed, 0x2c, 0xcf, 0x3f, 0x78, 0x7d, 0xd5, 0x71, 0x3d, 0x84, 0xf0, 0x82, 0x7e, 0x85, 0x69,
	0x52, 0x71, 0x85, 0x5b, 0x5e, 0x1d, 0xa4, 0x6b, 0x9a, 0xb4, 0x1f, 0xe0, 0xd6, 0xe3, 0x3f, 0x3d,
	0x11, 0x59, 0x7e, 0x86, 0xf7, 0xfb, 0xb2, 0x88, 0x55, 0x8a, 0xb1, 0x2a, 0x85, 0x8c, 0x4f, 0x26,
	0xcb, 0x31, 0xde, 0x17, 0xe2, 0x3e, 0x33, 0x18, 0xeb, 0x74, 0x39, 0x61, 0xf5, 0xda, 0xd5, 0x23,
	0x7d, 0x22, 0xe9, 0xd0, 0x6d, 0xdf, 0x97, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x52, 0x0c, 0x57,
	0x5d, 0x90, 0x02, 0x00, 0x00,
}
