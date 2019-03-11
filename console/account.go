package console

import (
	"fmt"
	"github.com/NBSChain/go-nbs/console/pb"
	"github.com/NBSChain/go-nbs/thirdParty/account"
	"github.com/NBSChain/go-nbs/utils/crypto"
	"github.com/spf13/cobra"
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
var offlineMode *bool

func init() {
	rootCmd.AddCommand(accountCmd)
	accountCmd.AddCommand(accountUnlockCmd)
	accountCmd.AddCommand(accountCreateCmd)
	offlineMode = accountCreateCmd.Flags().BoolP("offline",
		"o", false,
		"Create account in offline model")
}

func accountAction(cmd *cobra.Command, args []string) {
	request := &pb.AccountRequest{}

	conn := DialToCmdService()
	defer conn.Close()

	client := pb.NewAccountTaskClient(conn.c)

	response, err := client.Account(conn.ctx, request)

	logger.Info(response, err)
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

	client := pb.NewAccountTaskClient(conn.c)

	response, err := client.AccountUnlock(conn.ctx, request)

	logger.Info(response, err)
}

func createAccount(cmd *cobra.Command, args []string) {

	password := args[0]

	if err := crypto.CheckPassword(password); err != nil {
		panic(err)
	}

	request := &pb.CreateAccountRequest{
		Password: crypto.MD5SS(password),
	}

	if *offlineMode {
		createAccountOffline(crypto.MD5SS(password))
		return
	}

	conn := DialToCmdService()
	defer conn.Close()

	client := pb.NewAccountTaskClient(conn.c)

	response, err := client.CreateAccount(conn.ctx, request)

	logger.Info(response, err)
}

func createAccountOffline(password string) {

	acc := account.GetAccountInstance()
	if acc.GetPeerID() != "" {
		logger.Error("can't create another account, we support only one account right now!!")
		return
	}

	accId, err := acc.CreateAccount(password)
	if err != nil {
		logger.Info(err)
	}
	logger.Info("create account success:", accId)
}
