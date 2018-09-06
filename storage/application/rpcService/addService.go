package rpcService

import (
	"bytes"
	//"bytes"
	"errors"
	"fmt"
	"github.com/NBSChain/go-nbs/storage/core"
	"github.com/NBSChain/go-nbs/utils/cmdKits/pb"
	"github.com/NBSChain/go-nbs/utils/crypto"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
	"io"
	"sync"
)

const DefaultStreamTaskNo = 30
const StreamSessionIDKey = "sessionid"
const DefaultChunkSize = 1 << 18
const BigFileThreshold int64 = 50 << 20 //50M
const BigFileChunkSize int64 = 1 << 20  //1M

type addService struct {
	taskLock    sync.RWMutex
	fileAddTask map[string]*pb.AddRequest
}

func (service *addService) AddFile(ctx context.Context, request *pb.AddRequest) (*pb.AddResponse, error) {

	switch request.FileType {

	case pb.FileType_LARGEFILE:
		{
			if len(service.fileAddTask) >= DefaultStreamTaskNo {
				return nil, errors.New("too many add task running now. ")
			}

			sessionId := crypto.MD5(request.FullPath + string(request.FileSize))

			service.appendTask(sessionId, request)

			return &pb.AddResponse{Message: "continue", SessionId: sessionId}, nil
		}

	case pb.FileType_FILE:
		{
			if cs := request.ChunkSize; cs == 0 {
				request.ChunkSize = DefaultChunkSize
			}

			reader := &fileReader{
				reader: bytes.NewReader(request.FileData),
			}

			importer := &RpcFileImporter{
				chunkSize:   request.ChunkSize,
				reader:      reader,
				fileName:    request.FileName,
				fullPath:    request.FullPath,
				isDirectory: false,
			}

			err := core.ImportFile(importer)

			if err != nil {
				return nil, err
			} else {

				return &pb.AddResponse{Message: "success"}, nil
			}
		}
	case pb.FileType_DIRECTORY:
		{
			return &pb.AddResponse{Message: "success"}, nil
		}

	case pb.FileType_SYSTEMLINK:
		{
			return &pb.AddResponse{Message: "success"}, nil
		}
	default:
		return nil, errors.New("Unsupported file type! ")
	}

	return &pb.AddResponse{Message: "I want to add " + request.FileName}, nil
}

func (service *addService) TransLargeFile(stream pb.AddTask_TransLargeFileServer) error {
	header, ok := metadata.FromIncomingContext(stream.Context())
	if !ok || len(header[StreamSessionIDKey]) == 0 {
		return errors.New("unknown stream without session info")
	}

	sessionId := header[StreamSessionIDKey][0]

	request := service.fileAddTask[sessionId]
	if request == nil {
		return errors.New(fmt.Sprintf("can't find the "+
			"request of this session: %s", sessionId))
	}

	if cs := request.ChunkSize; cs == 0 {
		request.ChunkSize = DefaultChunkSize
	}

	reader := &streamReader{
		stream:    stream,
		service:   service,
		sessionId: sessionId,
	}

	importer := &RpcFileImporter{
		chunkSize:   request.ChunkSize,
		reader:      reader,
		fileName:    request.FileName,
		fullPath:    request.FullPath,
		isDirectory: false,
	}

	err := core.ImportFile(importer)

	importer.reader.Close()

	return err
}

func (service *addService) appendTask(sessionId string, request *pb.AddRequest) {

	service.taskLock.Lock()
	defer service.taskLock.Unlock()
	service.fileAddTask[sessionId] = request
}

func (service *addService) removeTask(sessionId string) {

	service.taskLock.Lock()
	defer service.taskLock.Unlock()
	delete(service.fileAddTask, sessionId)
}

type fileReader struct {
	reader io.Reader
}

func (r *fileReader) Read(p []byte) (n int, err error) {
	return r.reader.Read(p)
}
func (r *fileReader) Close() error {
	return nil
}

type streamReader struct {
	stream     pb.AddTask_TransLargeFileServer
	sessionId  string
	service    *addService
	cacheBytes []byte
}

func (s *streamReader) Read(p []byte) (n int, err error) {

	pLen := len(p)
	if pLen == 0 {
		return 0, nil
	}

	if s.cacheBytes == nil || len(s.cacheBytes) == 0 {

		fileChunk, err := s.stream.Recv()
		if err != nil {
			return 0, err
		}

		s.cacheBytes = make([]byte, BigFileChunkSize)
		copy(s.cacheBytes, fileChunk.Content)
	}

	cacheLen := len(s.cacheBytes)
	if pLen >= cacheLen {
		copy(p, s.cacheBytes)
		s.cacheBytes = nil
		return cacheLen, nil
	} else {
		copy(p, s.cacheBytes[:pLen])
		s.cacheBytes = s.cacheBytes[pLen:]
		return pLen, nil
	}
}

func (s *streamReader) Close() error {

	s.stream.SendAndClose(&pb.AddResponse{})

	s.service.removeTask(s.sessionId)

	return nil
}

type RpcFileImporter struct {
	reader      io.ReadCloser
	fileName    string
	fullPath    string
	isDirectory bool
	chunkSize   int32
}

func (importer *RpcFileImporter) Read(p []byte) (n int, err error) {
	return importer.reader.Read(p)
}
func (importer *RpcFileImporter) FileName() string {
	return importer.fileName
}

func (importer *RpcFileImporter) FullPath() string {
	return importer.fullPath
}

func (importer *RpcFileImporter) IsDirectory() bool {
	return importer.isDirectory
}

func (importer *RpcFileImporter) NextFile() (core.FileImporter, error) {
	return nil, io.EOF
}
func (importer *RpcFileImporter) Close() error {
	return importer.reader.Close()
}
