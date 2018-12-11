package console

import (
	"fmt"
	"github.com/NBSChain/go-nbs/console/pb"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(gossipCmd)
	gossipCmd.AddCommand(gossipStartCmd)
	gossipCmd.AddCommand(gossipStopCmd)
	gossipCmd.AddCommand(gossipViewCmd)
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
var gossipStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop the gossip service.",
	Long:  `stop the gossip service.`,
	Run:   gossipStop,
}

var gossipViewCmd = &cobra.Command{
	Use:   "showViews",
	Short: "show the nodes in partial view and input view.",
	Long:  `show the nodes in partial view and input view.`,
	Run:   gossipShowView,
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

func gossipStop(cmd *cobra.Command, args []string) {
	conn := DialToCmdService()
	defer conn.Close()

	client := pb.NewGossipTaskClient(conn.c)

	request := &pb.StopRequest{
		Cmd: "", //TODO:: Need to check password?
	}

	response, err := client.StopService(conn.ctx, request)

	logger.Info(response, err)
}

func gossipShowView(cmd *cobra.Command, args []string) {

	conn := DialToCmdService()
	defer conn.Close()

	client := pb.NewGossipTaskClient(conn.c)

	request := &pb.ShowGossipView{
		Cmd: "", //TODO:: Need to check password?
	}

	response, err := client.ShowViews(conn.ctx, request)

	logger.Info(response, err)
}
