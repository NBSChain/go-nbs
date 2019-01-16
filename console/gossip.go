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
	gossipCmd.AddCommand(gossipDebugCmd)
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

var gossipDebugCmd = &cobra.Command{
	Use:   "debug",
	Short: "gossip debug utils.",
	Long:  `gossip debug utils.`,
	Run:   gossipDebug,
	Args:  cobra.MinimumNArgs(1),
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
		Cmd: "",
	}

	response, err := client.StopService(conn.ctx, request)

	logger.Info(response, err)
}

func gossipDebug(cmd *cobra.Command, args []string) {

	logger.Debug(args)
	conn := DialToCmdService()
	defer conn.Close()

	client := pb.NewGossipTaskClient(conn.c)

	request := &pb.DebugCmd{
		Cmd: args[0],
	}

	response, err := client.Debug(conn.ctx, request)
	if err != nil {
		logger.Info("debug err:->", err)
	}
	processResult(args, response)
}
func processResult(args []string, res *pb.DebugResult) {

	switch args[0] {
	case "showIV":
		if res.InputViews != nil {
			fmt.Println(res.InputViews.Views)
		} else {
			fmt.Println("no input views")
		}
	case "showOV":
		if res.InputViews != nil {
			fmt.Println(res.OutputViews.Views)

		} else {
			fmt.Println("no output views")
		}
	case "showAV":
		fmt.Println("===>IN<===")
		if res.InputViews != nil {
			fmt.Println(res.InputViews.Views)
		} else {
			fmt.Println("no input views")
		}
		fmt.Println("===>Out<===")
		if res.InputViews != nil {
			fmt.Println(res.OutputViews.Views)

		} else {
			fmt.Println("no output views")
		}
	}
}
