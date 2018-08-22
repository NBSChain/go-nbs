package cmdKits

import (
	"fmt"
	"github.com/W-B-S/nbs-node/utils"
	"github.com/W-B-S/nbs-node/utils/cmdKits/pb"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"os"
)

var nbsUsage = `
The nbs will start listening on ports on the network, which are
documented in (and can be modified through) 'ipfs config Addresses'.
For example, to change the 'Gateway' port:

  ipfs config Addresses.Gateway /ip4/127.0.0.1/tcp/8082

The API address can be changed the same way:

  ipfs config Addresses.API /ip4/127.0.0.1/tcp/5002

Make sure to restart the daemon after changing addresses.

By default, the gateway is only accessible locally. To expose it to
other computers in the network, use 0.0.0.0 as the ip address:

  ipfs config Addresses.Gateway /ip4/0.0.0.0/tcp/8080

Be careful if you expose the API. It is a security risk, as anyone could
control your node remotely. If you need to control the node remotely,
make sure to protect the port as you would other services or database
(firewall, authenticated proxy, etc).

HTTP Headers

ipfs supports passing arbitrary headers to the API and Gateway. You can
do this by setting headers on the API.HTTPHeaders and Gateway.HTTPHeaders
keys:

  ipfs config --json API.HTTPHeaders.X-Special-Header '["so special :)"]'
  ipfs config --json Gateway.HTTPHeaders.X-Special-Header '["so special :)"]'

Note that the value of the keys is an _array_ of strings. This is because
headers can have more than one value, and it is convenient to pass through
to other libraries.

CORS Headers (for API)

You can setup CORS headers the same way:

  ipfs config --json API.HTTPHeaders.Access-Control-Allow-Origin '["example.com"]'
  ipfs config --json API.HTTPHeaders.Access-Control-Allow-Methods '["PUT", "GET", "POST"]'
  ipfs config --json API.HTTPHeaders.Access-Control-Allow-Credentials '["true"]'

Shutdown

To shutdown the daemon, send a SIGINT signal to it (e.g. by pressing 'Ctrl-C')
or send a SIGTERM signal to it (e.g. with 'kill'). It may take a while for the
daemon to shutdown gracefully, but it can be killed forcibly by sending a
second signal.

IPFS_PATH environment variable

ipfs uses a repository in the local file system. By default, the repo is
located at ~/.ipfs. To change the repo location, set the $IPFS_PATH
environment variable:

  export IPFS_PATH=/path/to/ipfsrepo

Routing

IPFS by default will use a DHT for content routing. There is a highly
experimental alternative that operates the DHT in a 'client only' mode that
can be enabled by running the daemon as:

  ipfs daemon --routing=dhtclient

This will later be transitioned into a config option once it gets out of the
'experimental' stage.

DEPRECATION NOTICE

Previously, ipfs used an environment variable as seen below:

  export API_ORIGIN="http://localhost:8888/"

This is deprecated. It is still honored in this version, but will be removed
in a future version, along with this notice. Please move to setting the HTTP
Headers.
`
var logger = utils.GetLogInstance()

const cmdNameAdd = "add"
const cmdVersion = "version"

type ServiceAction func(ctx context.Context, req *pb.CmdRequest) (*pb.CmdResponse, error)

var NbsCommandSet = map[string]ServiceAction{
	cmdNameAdd: ServiceTaskVersionAddFile,
	cmdVersion: ServiceTaskVersion,
}

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
