package application

import (
	"github.com/NBSChain/go-nbs/utils/cmdKits/pb"
	"github.com/NBSChain/go-nbs/utils/crypto"
	"github.com/W-B-S/nbs-node/storage/application"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"os"
)
import "google.golang.org/grpc"
import "google.golang.org/grpc/reflection"
import "net"
import "github.com/NBSChain/go-nbs/utils"

type cmdService struct{}

func startCmdService() {

	var address = "127.0.0.1:" + utils.GetConfig().CmdServicePort

	listener, err := net.Listen("tcp", address)
	if err != nil {
		logger.Fatalf("Failed to listen: %v", err)
	}
	theServer := grpc.NewServer()

	pb.RegisterAddTaskServer(theServer, &cmdService{})
	pb.RegisterVersionTaskServer(theServer, &cmdService{})

	reflection.Register(theServer)
	if err := theServer.Serve(listener); err != nil {
		logger.Fatalf("Failed to serve: %v", err)
	}
}

func (s *cmdService) AddFile(ctx context.Context, request *pb.AddRequest) (*pb.AddResponse, error) {

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
		return nil, errors.New("Unsupported file type!")
	}

	fileName := request.FileName

	app := application.GetInstance()

	file, err := os.Open(fileName)
	if err != nil {
		logger.Warning("Failed to open file.")
		return nil, err
	}

	app.AddFile(file)

	return &pb.AddResponse{Message: "I want to add " + request.FileName}, nil
}

func (s *cmdService) TransLargeFile(stream pb.AddTask_TransLargeFileServer) error {
	return nil
}

func (s *cmdService) SystemVersion(ctx context.Context,
	request *pb.VersionRequest) (*pb.VersionResponse, error) {

	return &pb.VersionResponse{Message: "Current version is  " +
		utils.GetConfig().CurrentVersion}, nil
}
