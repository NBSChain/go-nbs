package rpcService

import (
	"errors"
	"github.com/NBSChain/go-nbs/utils/cmdKits/pb"
	"github.com/NBSChain/go-nbs/utils/crypto"
	"golang.org/x/net/context"
)

type addService struct{}

func (s *addService) AddFile(ctx context.Context, request *pb.AddRequest) (*pb.AddResponse, error) {

	logger.Info(request)

	switch request.FileType {

	case pb.FileType_LARGEFILE:
		{
			sessionId := crypto.MD5(request.FullPath + string(request.FileSize))
			return &pb.AddResponse{Message: "continue", SessionId: sessionId}, nil
		}
	case pb.FileType_FILE:
		{
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
