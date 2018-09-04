package cmdKits

import (
	"fmt"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/spf13/cobra"
	"os"
)

var nbsUsage = `TODO::......`
var logger = utils.GetLogInstance()

var rootCmd = &cobra.Command{
	Use: "nbs",

	Short: "Nbs is a new blockChain system with distributed storage and smart contract.",

	Long: nbsUsage,

	Run: mainRun,
}

func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func mainRun(cmd *cobra.Command, args []string) {

	logger.Info("root command args:", args)

	StartCmdService()

	logger.Info("Nbs daemon is ready......")
}
