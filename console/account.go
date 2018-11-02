package console

import (
	"fmt"
	"github.com/NBSChain/go-nbs/utils/crypto"
	"github.com/spf13/cobra"
)

var accountCmd = &cobra.Command{
	Use:   	"account",
	Short: 	"nbs account unlock/list",
	Long:  	`nbs account unlock/list`,
	Run:	accountAction,
}

var accountUnlockCmd = &cobra.Command{
	Use:   	"unlock",
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

	//encodeKey := hex.EncodeToString([]byte(args[0]))
}