package cmdKits

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(getCmd)
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get file by hash value",
	Long:  `Get file by hash value.`,
	Run:   getFile,
}

func getFile(cmd *cobra.Command, args []string) {
	logger.Info("get command args:(", args, ")-->", cmd.CommandPath())
}
