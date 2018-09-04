package cmdKits

import "github.com/NBSChain/go-nbs/utils/cmdKits/pb"
import "golang.org/x/net/context"
import "google.golang.org/grpc"
import "google.golang.org/grpc/reflection"
import "net"
import "time"
import "io"
import "errors"
import "log"
import "github.com/NBSChain/go-nbs/utils"

type cmdService struct{}

type File interface {
	io.ReadCloser

	FileName() string

	FullPath() string

	IsDirectory() bool

	NextFile() (File, error)
}

const cmdNameAdd = "add"
const cmdVersion = "version"

type ServiceAction func(ctx context.Context, req *pb.CmdRequest) (*pb.CmdResponse, error)

var NbsCommandSet = map[string]ServiceAction{
	cmdNameAdd: ServiceTaskAddFile,
	cmdVersion: ServiceTaskVersion,
}

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
		logger.Fatalf("Failed to listen: %v", err)
	}
	theServer := grpc.NewServer()

	pb.RegisterCmdTaskServer(theServer, &cmdService{})

	reflection.Register(theServer)
	if err := theServer.Serve(listener); err != nil {
		logger.Fatalf("Failed to serve: %v", err)
	}
}

func handleInputCmd(ctx context.Context, req *pb.CmdRequest) (*pb.CmdResponse, error) {

	service := NbsCommandSet[req.CmdName]
	if nil == service {
		return nil, errors.New("can't find service to process this command")
	}

	return service(ctx, req)
}
