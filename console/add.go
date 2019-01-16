package console

import (
	"github.com/NBSChain/go-nbs/console/pb"
	"github.com/NBSChain/go-nbs/storage/application/rpcService"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/metadata"
	"io"
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
	Run:   addFileCmd,
	Args:  cobra.MinimumNArgs(1),
}

func addFileCmd(cmd *cobra.Command, args []string) {

	logger.Info("Add command args:(", args, ")-->", cmd.CommandPath())

	fileName := args[0]

	fileInfo, ok := utils.FileExists(fileName)
	if !ok {
		log.Fatal("File is not available.")
	}

	fullPath, err := filepath.Abs(fileName)
	if err != nil {
		logger.Fatal(err)
	}

	request := &pb.AddRequest{
		FileName: fileName,
		FullPath: fullPath,
		FileSize: fileInfo.Size(),
	}

	if fileInfo.IsDir() {

		request.FileType = pb.FileType_DIRECTORY

	} else if fileInfo.Mode() == os.ModeSymlink {

		request.FileType = pb.FileType_SYSTEMLINK

	} else {

		if fileInfo.Size() >= rpcService.BigFileThreshold {

			request.FileType = pb.FileType_LARGEFILE

			response := addFile(request)

			logger.Info("Send file metadata success......", response)

			sendFileStream(response.SessionId, fullPath)

		} else {

			request.FileType = pb.FileType_FILE

			request.FileData = make([]byte, fileInfo.Size())

			file, err := os.Open(fullPath)
			if err != nil {
				log.Fatal(err)
			}

			file.Read(request.FileData)

			response := addFile(request)

			defer file.Close()

			logger.Info("Add file success......", response)
		}
	}
}

func sendFileStream(sessionId, fileName string) {

	conn := DialToCmdService()
	defer conn.Close()

	client := pb.NewAddTaskClient(conn.c)

	header := metadata.Pairs(rpcService.StreamSessionIDKey, sessionId)
	ctx := metadata.NewOutgoingContext(conn.ctx, header)

	stream, err := client.TransLargeFile(ctx)
	if err != nil {
		log.Fatalf("send file stream failed : %v", err)
	}

	file, err := os.Open(fileName)
	if err != nil {
		logger.Fatal("failed to open the file:", err)
	}
	defer file.Close()

	//buffer := make([]byte, 123)
	buffer := make([]byte, rpcService.BigFileChunkSize)

	for {

		no, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				logger.Fatal("failed to read file chunck:", err)
			}
		}

		stream.Send(&pb.FileChunk{
			Content: buffer[:no],
		})
	}

	response, err := stream.CloseAndRecv()
	logger.Info(response, err)

}

func addFile(request *pb.AddRequest) *pb.AddResponse {

	conn := DialToCmdService()
	defer conn.Close()

	client := pb.NewAddTaskClient(conn.c)

	response, err := client.AddFile(conn.ctx, request)
	if err != nil {
		logger.Fatalf("could not add file : %v", err)
	}

	return response
}
