package console

import (
	"encoding/hex"
	"fmt"
	"github.com/NBSChain/go-nbs/console/pb"
	"github.com/NBSChain/go-nbs/utils/crypto"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var accountCmd = &cobra.Command{
	Use:   	"account",
	Short: 	"nbs account unlock/list",
	Long:  	`nbs account unlock/list`,
	Run:	accountAction,
}

var accountUnlockCmd = &cobra.Command{
	Use:   	"unlock [# password]",
	Short: 	"unlock current account",
	Long:  	`unlock current account so you can make some security actions.`,
	Run:	unlockAccount,
	Args:	cobra.MinimumNArgs(1),
}


func init() {
	rootCmd.AddCommand(accountCmd)
	accountCmd.AddCommand(accountUnlockCmd)
}

func accountAction(cmd *cobra.Command, args []string) {
	fmt.Println("I'm ok!")
}

func unlockAccount(cmd *cobra.Command, args []string) {
	fmt.Println("unlock this:", args[0])

	password := args[0]

	if err := crypto.CheckPassword(password); err != nil {
		panic(err)
	}

	encodeKey := hex.EncodeToString([]byte(args[0]))
	request := &pb.AccountUnlockRequest{
		Password: encodeKey,
	}

	conn := DialToCmdService()
	defer conn.Close()

	client := pb.NewAccountTaskClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.AccountUnlock(ctx, request)
	if err != nil{
		logger.Fatalf("failed to unlock account:", err.Error())
	}
	logger.Info(response)
}