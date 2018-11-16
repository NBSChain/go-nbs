package nat

import (
	"github.com/NBSChain/go-nbs/utils"
	"strconv"
)

type decentralizedNatServer struct {
	hosts []string
}

func newDecentralizedNatServer() *decentralizedNatServer {

	officerServer := utils.GetConfig().NatServerIP
	server := &decentralizedNatServer{
		hosts: make([]string, len(officerServer)),
	}

	port := strconv.Itoa(utils.GetConfig().NatServerPort)

	for _, host := range officerServer{
		server.hosts = append(server.hosts, host + ":" + port)
	}

	return server
}

//TODO:: use gossip protocol to manager all nat servers. we use official nat servers right now.
func (s *decentralizedNatServer) GossipNatServer() string{
	return s.hosts[0]//TIPS:: simply use the first server.
}