package rpcService

import (
	"errors"
	"fmt"
	"github.com/NBSChain/go-nbs/utils/cmdKits/pb"
	"github.com/NBSChain/go-nbs/utils/crypto"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
	"io"
	"os"
	"sync"
)

const DefaultStreamTaskNo = 30

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
			file, err := os.Create("server-" + request.FileName)
			if err != nil {
				return nil, err
			}

			file.Write(request.FileData)

			file.Close()

			return &pb.AddResponse{Message: "success"}, nil
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
	if !ok || len(header["sessionId"]) == 0 {
		return errors.New("unknown stream without session info")
	}

	sessionId := header["sessionId"][0]

	request := service.fileAddTask[sessionId]
	if request == nil {
		return errors.New(fmt.Sprintf("can't find the request of this session: %s", sessionId))
	}

	defer service.removeTask(sessionId)

	file, err := os.Create("test.mp4")
	if err != nil {
		return err
	}
	defer file.Close()

	var dataLen int64 = 0

	for {
		fileData, err := stream.Recv()
		if err != nil {

			if err == io.EOF {
				break
			} else {
				return err
			}
		}

		file.Write(fileData.Content)
	}

	stream.SendAndClose(&pb.AddResponse{
		DataReceive: dataLen,
	})

	return nil
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
