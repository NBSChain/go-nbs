package cmdKits

import (
	"errors"
	"github.com/W-B-S/nbs-node/utils"
	"github.com/W-B-S/nbs-node/utils/cmdKits/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"
)

type cmdService struct{}

func (s *cmdService) RpcTask(ctx context.Context, req *pb.CmdRequest) (*pb.CmdResponse, error) {
	return handleInputCmd(ctx, req)
}

func DialToCmdService(request *pb.CmdRequest) *pb.CmdResponse {

	var address = "127.0.0.1:" + utils.GetConfig().CmdServicePort

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		logger.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewCmdTaskClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.RpcTask(ctx, request)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	return response
}

func StartCmdService() {

	var address = "127.0.0.1:" + utils.GetConfig().CmdServicePort

	listener, err := net.Listen("tcp", address)
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}
	theServer := grpc.NewServer()

	pb.RegisterCmdTaskServer(theServer, &cmdService{})

	reflection.Register(theServer)
	if err := theServer.Serve(listener); err != nil {
		logger.Fatalf("failed to serve: %v", err)
	}
}

func handleInputCmd(ctx context.Context, req *pb.CmdRequest) (*pb.CmdResponse, error) {

	service := NbsCommandSet[req.CmdName]
	if nil == service {
		return nil, errors.New("can't find service to process this command")
	}

	return service(ctx, req)
}
