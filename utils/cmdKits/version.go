package cmdKits

import (
	"github.com/W-B-S/nbs-node/utils"
	"github.com/W-B-S/nbs-node/utils/cmdKits/pb"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show the current software's version.",
	Long:  `show the current software's version.`,
	Run: func(cmd *cobra.Command, args []string) {

		request := &pb.CmdRequest{CmdName: cmdVersion}

		response := DialToCmdService(request)
		logger.Info(response)
	},
}

func ServiceTaskVersion(ctx context.Context, req *pb.CmdRequest) (*pb.CmdResponse, error) {
	return &pb.CmdResponse{Message: "Current version is  " +
		utils.GetConfig().CurrentVersion}, nil
}
