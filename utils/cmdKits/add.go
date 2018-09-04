package cmdKits

import (
	"github.com/NBSChain/go-nbs/storage/application"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/NBSChain/go-nbs/utils/cmdKits/pb"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"log"
	"os"
	"path/filepath"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add file to nbs network",
	Long:  `Add file to cache and find the peers to store it.`,
	Run:   shellAddFile,
}

func shellAddFile(cmd *cobra.Command, args []string) {

	logger.Info("Add command args:(", args, ")-->", cmd.CommandPath())

	if len(args) == 0 {
		logger.Fatal("You should specify the file target to add.")
	}

	fileName := args[0]

	fileInfo, ok := utils.FileExists(fileName)
	if !ok || fileInfo.IsDir() {
		log.Fatal("File is not available.")
		//TODO::Going to support directory.
	}

	fileName, err := filepath.Abs(fileName)
	if err != nil {
		logger.Fatal(err)
	}

	request := &pb.CmdRequest{
		CmdName: cmdNameAdd,
		Args: []string{
			fileName,
		},
	}

	response := DialToCmdService(request)

	logger.Info("Reading success......", response.Message)
}

func ServiceTaskAddFile(ctx context.Context, req *pb.CmdRequest) (*pb.CmdResponse, error) {

	logger.Info(req)

	fileName := req.Args[0]

	app := application.GetInstance()

	file, err := os.Open(fileName)
	if err != nil {
		logger.Warning("Failed to open file.")
		return nil, err
	}

	app.AddFile(file)

	return &pb.CmdResponse{Message: "I want to  " + req.CmdName}, nil
}
