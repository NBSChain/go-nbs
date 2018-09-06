package rpcService

import (
	"errors"
	"github.com/NBSChain/go-nbs/utils/cmdKits/pb"
	"github.com/NBSChain/go-nbs/utils/crypto"
	"golang.org/x/net/context"
	"os"
)

type addService struct{}

const BigFileThreshold int64 = 50 << 20

func (s *addService) AddFile(ctx context.Context, request *pb.AddRequest) (*pb.AddResponse, error) {

	switch request.FileType {

	case pb.FileType_LARGEFILE:
		{
			sessionId := crypto.MD5(request.FullPath + string(request.FileSize))
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

func (s *addService) TransLargeFile(stream pb.AddTask_TransLargeFileServer) error {
	return nil
}
