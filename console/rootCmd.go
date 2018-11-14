package console

import (
	"context"
	"fmt"
	"github.com/NBSChain/go-nbs/storage/application"
	"github.com/NBSChain/go-nbs/storage/application/rpcService"
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
var natServiceConf *bool

func init() {
	natServiceConf = rootCmd.Flags().BoolP("without-nat",
		"n", false,
		"Without nat service, this is used when you're sure in private network")

}

func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type CmdConnection struct {
	c      *grpc.ClientConn
	ctx    context.Context
	cancel context.CancelFunc
}

func mainRun(cmd *cobra.Command, args []string) {

	logger.Info("get command args:(", args, ")-->", *natServiceConf)

	utils.GetConfig().NatServiceOff = *natServiceConf

	application.GetInstance().Start()

	rpcService.StartCmdService()
}

func DialToCmdService() *CmdConnection {
	var address = "127.0.0.1:" + utils.GetConfig().CmdServicePort

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		logger.Fatalf("did not connect: %v", err)
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &CmdConnection{
		c:      conn,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (conn *CmdConnection) Close() {
	conn.c.Close()
	conn.cancel()
}
