package console

import (
	"fmt"
	"github.com/NBSChain/go-nbs/console/pb"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(natCmd)
	natCmd.AddCommand(natStatusCmd)
}

var natCmd = &cobra.Command{
	Use:   "nat",
	Short: "NAT information and operations for current network",
	Long:  `NAT information and operations for current network`,
	Run:   natAction,
	Args:  cobra.MinimumNArgs(1),
}

var natStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "check the NAT status of current network.",
	Long:  `check the NAT status of current network.`,
	Run:   natStatus,
}

func natAction(cmd *cobra.Command, args []string) {
	fmt.Println("please input your option!")
}

func natStatus(cmd *cobra.Command, args []string) {

	conn := DialToCmdService()
	defer conn.Close()

	client := pb.NewNatTaskClient(conn.c)

	request := &pb.NatStatusResQuest{
		Password: "", //TODO:: Need to check password?
	}

	response, err := client.NatStatus(conn.ctx, request)
	if err != nil {
		logger.Fatalf("failed to create account:", err.Error())
	}
	logger.Info(response.Message)
}
