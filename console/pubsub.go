package console

import (
	"fmt"
	"github.com/NBSChain/go-nbs/console/pb"
	"github.com/spf13/cobra"
)

var pubSubCmd = &cobra.Command{
	Use:   "pubsub",
	Short: "gossip protocol publish and subscribe actions",
	Long:  `gossip protocol publish and subscribe actions`,
	Run:   pubAndSub,
}

var pubSubPubCmd = &cobra.Command{
	Use:   "pub [# topic] [# message]",
	Short: "publish message to topic.",
	Long:  `publish message to topic.`,
	Run:   publish,
	Args:  cobra.MinimumNArgs(2),
}

var pubSubSubCmd = &cobra.Command{
	Use:   "sub [# topic]",
	Short: "subscribe a topic",
	Long:  `subscribe a topic.`,
	Run:   subscribe,
	Args:  cobra.MinimumNArgs(1),
}

var pubSubPeersCmd = &cobra.Command{
	Use:   "peers",
	Short: "list all peers in my current node",
	Long:  `list all peers in my current node`,
	Run:   listAllPeers,
	Args:  cobra.MinimumNArgs(1),
}

var pubSubLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list all topics has subscribed",
	Long:  `list all topics has subscribed.`,
	Run:   listAllTopics,
}
var peerSearchDepth *int32

func init() {
	rootCmd.AddCommand(pubSubCmd)
	pubSubCmd.AddCommand(pubSubPubCmd)
	pubSubCmd.AddCommand(pubSubSubCmd)
	pubSubCmd.AddCommand(pubSubPeersCmd)
	pubSubCmd.AddCommand(pubSubLsCmd)
	peerSearchDepth = pubSubPeersCmd.Flags().Int32P("depth", "d", 4, "the depth you want to search")
}

func pubAndSub(cmd *cobra.Command, args []string) {
	fmt.Println("I'm ok!")
}

func publish(cmd *cobra.Command, args []string) {

	conn := DialToCmdService()
	defer conn.Close()

	client := pb.NewPubSubTaskClient(conn.c)

	request := &pb.PublishRequest{
		Topics:  args[0],
		Message: args[1],
	}

	response, err := client.Publish(conn.ctx, request)

	logger.Info(response, err)
}

func subscribe(cmd *cobra.Command, args []string) {

	conn := DialToCmdService()
	defer conn.Close()

	client := pb.NewPubSubTaskClient(conn.c)

	request := &pb.SubscribeRequest{
		Topics: args[0],
	}

	stream, err := client.Subscribe(conn.ctx, request)
	if err != nil {
		logger.Warning("failed to sub the topics:->", args[0])
		return
	}

	for {
		response, err := stream.Recv()
		if err != nil {
			logger.Warning("receive rpc response err:->", err)
			return
		}

		logger.Info("get msg:->", response.MsgData)
	}
}

func listAllPeers(cmd *cobra.Command, args []string) {

	conn := DialToCmdService()
	defer conn.Close()

	client := pb.NewPubSubTaskClient(conn.c)

	request := &pb.PeersRequest{
		Topics: args[0],
		Depth:  *peerSearchDepth,
	}

	response, err := client.Peers(conn.ctx, request)

	logger.Info(response, err)
}

func listAllTopics(cmd *cobra.Command, args []string) {

	conn := DialToCmdService()
	defer conn.Close()

	client := pb.NewPubSubTaskClient(conn.c)

	request := &pb.TopicsRequest{
		Message: "", //TODO:: I don't know what we need right now.
	}

	response, err := client.Topics(conn.ctx, request)

	logger.Info(response, err)
}
