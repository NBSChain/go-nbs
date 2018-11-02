package rpcService

import (
	"github.com/NBSChain/go-nbs/console/pb"
	"golang.org/x/net/context"
)
import "google.golang.org/grpc"
import "google.golang.org/grpc/reflection"
import "net"
import "github.com/NBSChain/go-nbs/utils"

const SplitterSize 		= 1 << 18	   	//256K
const BigFileThreshold 	int64 	= 50 << 20  	   	//50M
const BigFileChunkSize 	int64 	= SplitterSize << 2	//1M
const MaxFIleSize 	int 	= 10 << 30		//10G

var   logger 			= utils.GetLogInstance()

type cmdService struct{}

func StartCmdService() {

	var address = "127.0.0.1:" + utils.GetConfig().CmdServicePort

	listener, err := net.Listen("tcp", address)
	if err != nil {
		logger.Fatalf("Failed to listen: %v", err)
	}

	theServer := grpc.NewServer(grpc.MaxRecvMsgSize(MaxFIleSize), grpc.MaxSendMsgSize(MaxFIleSize))

	pb.RegisterAddTaskServer(theServer, &addService{
		fileAddTask: make(map[string]*pb.AddRequest),
	})

	pb.RegisterVersionTaskServer(theServer, &cmdService{})

	pb.RegisterGetTaskServer(theServer, &getService{})

	pb.RegisterAccountTaskServer(theServer, &accountService{})

	reflection.Register(theServer)
	if err := theServer.Serve(listener); err != nil {
		logger.Fatalf("Failed to serve: %v", err)
	}
}

func (s *cmdService) SystemVersion(ctx context.Context,
	request *pb.VersionRequest) (*pb.VersionResponse, error) {

	return &pb.VersionResponse{Result: "Current version is  " +
		utils.GetConfig().CurrentVersion}, nil
}
