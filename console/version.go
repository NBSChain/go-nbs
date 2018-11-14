package console

import (
	"github.com/NBSChain/go-nbs/console/pb"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show the current software's version.",
	Long:  `show the current software's version.`,
	Run: func(cmd *cobra.Command, args []string) {

		request := &pb.VersionRequest{CmdName: "version"}

		response := versionReq(request)
		logger.Info(response)
	},
}

func versionReq(request *pb.VersionRequest) *pb.VersionResponse {

	conn := DialToCmdService()
	defer conn.Close()

	client := pb.NewVersionTaskClient(conn.c)

	response, err := client.SystemVersion(conn.ctx, request)
	if err != nil {
		logger.Fatalf("could not greet: %v", err)
	}

	return response
}
