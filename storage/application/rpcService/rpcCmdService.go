package rpcService

import (
	"github.com/NBSChain/go-nbs/utils/cmdKits/pb"
	"golang.org/x/net/context"
)
import "google.golang.org/grpc"
import "google.golang.org/grpc/reflection"
import "net"
import "github.com/NBSChain/go-nbs/utils"

type cmdService struct{}

var logger = utils.GetLogInstance()

func StartCmdService() {

	var address = "127.0.0.1:" + utils.GetConfig().CmdServicePort

	listener, err := net.Listen("tcp", address)
	if err != nil {
		logger.Fatalf("Failed to listen: %v", err)
	}
	theServer := grpc.NewServer()

	pb.RegisterAddTaskServer(theServer, &addService{})
	pb.RegisterVersionTaskServer(theServer, &cmdService{})

	reflection.Register(theServer)
	if err := theServer.Serve(listener); err != nil {
		logger.Fatalf("Failed to serve: %v", err)
	}
}

func (s *cmdService) SystemVersion(ctx context.Context,
	request *pb.VersionRequest) (*pb.VersionResponse, error) {

	return &pb.VersionResponse{Message: "Current version is  " +
		utils.GetConfig().CurrentVersion}, nil
}
