package console

import (
	"context"
	"github.com/NBSChain/go-nbs/console/pb"
	"github.com/spf13/cobra"
	"io"
	"log"
	"os"
)

var targetFileName string
var getCmd = &cobra.Command{
	Use:   	"get",
	Short: 	"Get file by hash value",
	Long:  	`Get file by hash value.`,
	Run:	getFile,
	Args:	cobra.MinimumNArgs(1),
}

func init() {

	rootCmd.AddCommand(getCmd)

	getCmd.Flags().StringVarP(&targetFileName, "object", "o", "", "object file name")

}

func getFile(cmd *cobra.Command, args []string) {
	logger.Info("get command args:(", args, ")-->", targetFileName)
	dataHash := args[0]

	if targetFileName == ""{
		targetFileName = dataHash
	}

	request := &pb.GetRequest{
		DataUri:dataHash,
	}

	conn := DialToCmdService()
	defer conn.Close()

	client := pb.NewGetTaskClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := client.Get(ctx, request)
	if err != nil{
		log.Fatalf("could not add file : %v", err)
	}

	writeDataToFile(targetFileName, stream)
}

func writeDataToFile(filePath string, stream pb.GetTask_GetClient)  {

	file, err := os.Create(filePath)
	if err != nil{
		log.Fatal(err)
	}

	for {
		response, err := stream.Recv()
		if err !=  nil{
			if err == io.EOF{
				break
			}else{
				log.Fatal("failed to read data from node.")
			}
		}else{
			if _, err := file.Write(response.Content); err != nil{
				log.Fatal("failed to write data to local file:"+err.Error())
			}
		}
	}

	file.Close()

	logger.Info("Get file success......")
}
