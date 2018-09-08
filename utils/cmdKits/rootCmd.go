package cmdKits

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/application"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"os"
)

var nbsUsage = `TODO::......`
var logger = utils.GetLogInstance()

var rootCmd = &cobra.Command{
	Use: "nbs",

	Short: "Nbs is a new blockChain system with distributed storage and smart contract.",

	Long: nbsUsage,

	Run: mainRun,
}

func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func mainRun(cmd *cobra.Command, args []string) {

	logger.Info("root command args:", args)

	application.GetInstance().Start()

	logger.Info("Nbs daemon is ready......")
}

func DialToCmdService() *grpc.ClientConn {
	var address = "127.0.0.1:" + utils.GetConfig().CmdServicePort

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		logger.Fatalf("did not connect: %v", err)
		return nil
	}
	return conn
}
