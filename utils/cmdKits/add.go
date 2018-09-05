package cmdKits

import (
	"github.com/NBSChain/go-nbs/utils"
	"github.com/NBSChain/go-nbs/utils/cmdKits/pb"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"log"
	"path/filepath"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add file to nbs network",
	Long:  `Add file to cache and find the peers to store it.`,
	Run:   addFile,
}

func addFile(cmd *cobra.Command, args []string) {

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

	request := &pb.AddRequest{}

	response := AddFile(request)

	logger.Info("Reading success......", response.Message)
}

func AddFile(request *pb.AddRequest) *pb.AddResponse {

	conn := DialToCmdService()

	defer conn.Close()
	client := pb.NewAddTaskClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.AddFile(ctx, request)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	return response
}
