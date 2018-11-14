package console

import (
	"fmt"
	"github.com/NBSChain/go-nbs/console/pb"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(gossipCmd)
	gossipCmd.AddCommand(gossipStartCmd)
}

var gossipCmd = &cobra.Command{
	Use:   "gossip",
	Short: "to detect the gossip service information.",
	Long:  `to detect the gossip service information.`,
	Run:   gossipAction,
	Args:  cobra.MinimumNArgs(1),
}

var gossipStartCmd = &cobra.Command{
	Use:   "start",
	Short: "start the gossip service.",
	Long:  `start the gossip service.`,
	Run:   gossipStart,
}

func gossipAction(cmd *cobra.Command, args []string) {
	fmt.Println("please input your option!")
}

func gossipStart(cmd *cobra.Command, args []string) {
	conn := DialToCmdService()
	defer conn.Close()

	client := pb.NewGossipTaskClient(conn.c)

	request := &pb.StartRequest{
		Cmd: "", //TODO:: Need to check password?
	}

	response, err := client.StartService(conn.ctx, request)

	logger.Info(response, err)
}
