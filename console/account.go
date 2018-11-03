package console

import (
	"fmt"
	"github.com/NBSChain/go-nbs/console/pb"
	"github.com/NBSChain/go-nbs/utils/crypto"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "nbs account unlock/list",
	Long:  `nbs account unlock/list`,
	Run:   accountAction,
}

var accountUnlockCmd = &cobra.Command{
	Use:   "unlock [# password]",
	Short: "unlock current account",
	Long:  `unlock current account so you can make some security actions.`,
	Run:   unlockAccount,
	Args:  cobra.MinimumNArgs(1),
}

var accountCreateCmd = &cobra.Command{
	Use:   "create [# password]",
	Short: "create an account with the password",
	Long:  `create an account with the password.`,
	Run:   createAccount,
	Args:  cobra.MinimumNArgs(1),
}

var accountPassword string

func init() {
	rootCmd.AddCommand(accountCmd)
	accountCmd.AddCommand(accountUnlockCmd)
	accountCmd.AddCommand(accountCreateCmd)
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

	request := &pb.AccountUnlockRequest{
		Password: crypto.MD5SS(password),
	}

	conn := DialToCmdService()
	defer conn.Close()

	client := pb.NewAccountTaskClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.AccountUnlock(ctx, request)
	if err != nil {
		logger.Fatalf("failed to unlock account:", err.Error())
	}
	logger.Info(response)
}

func createAccount(cmd *cobra.Command, args []string) {

	fmt.Println("unlock this:", args[0])
	password := args[0]

	if err := crypto.CheckPassword(password); err != nil {
		panic(err)
	}

	request := &pb.CreateAccountRequest{
		Password: crypto.MD5SS(password),
	}

	conn := DialToCmdService()
	defer conn.Close()

	client := pb.NewAccountTaskClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.CreateAccount(ctx, request)
	if err != nil {
		logger.Fatalf("failed to create account:", err.Error())
	}
	logger.Info(response)
}
