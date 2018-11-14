package console

import (
	"fmt"
	"github.com/spf13/cobra"
)

var pubsubCmd = &cobra.Command{
	Use:   "pubsub",
	Short: "gossip protocol publish and subscribe actions",
	Long:  `gossip protocol publish and subscribe actions`,
	Run:   pubAndSub,
}

var pubsubPubCmd = &cobra.Command{
	Use:   "pub [# topic] [# message]",
	Short: "publish message to topic.",
	Long:  `publish message to topic.`,
	Run:   publish,
	Args:  cobra.MinimumNArgs(2),
}

var pubsubSubCmd = &cobra.Command{
	Use:   "sub [# topic]",
	Short: "subscribe a topic",
	Long:  `subscribe a topic.`,
	Run:   subscribe,
	Args:  cobra.MinimumNArgs(1),
}

var pubsubPeersCmd = &cobra.Command{
	Use:   "peers",
	Short: "list all peers in my current node",
	Long:  `list all peers in my current node`,
	Run:   listAllPeers,
}

var pubsubLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list all topics has subscribed",
	Long:  `list all topics has subscribed.`,
	Run:   listAllTopics,
}

func init() {
	rootCmd.AddCommand(pubsubCmd)
	pubsubCmd.AddCommand(pubsubPubCmd)
	pubsubCmd.AddCommand(pubsubSubCmd)
	pubsubCmd.AddCommand(pubsubPeersCmd)
	pubsubCmd.AddCommand(pubsubLsCmd)
}

func pubAndSub(cmd *cobra.Command, args []string) {
	fmt.Println("I'm ok!")
}
func publish(cmd *cobra.Command, args []string) {
	fmt.Println("I'm ok!")
}
func subscribe(cmd *cobra.Command, args []string) {

	conn := DialToCmdService()
	defer conn.Close()

	//client := pb.NewAddTaskClient(conn.c)

	fmt.Println("I'm ok!")
}
func listAllPeers(cmd *cobra.Command, args []string) {
	fmt.Println("I'm ok!")
}
func listAllTopics(cmd *cobra.Command, args []string) {
	fmt.Println("I'm ok!")
}
