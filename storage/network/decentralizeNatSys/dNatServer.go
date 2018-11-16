package decentralizeNatSys

import (
	"github.com/NBSChain/go-nbs/utils"
	"strconv"
)

type DecentralizedNatServer struct {
	hosts []string
}

func NewDecentralizedNatServer() *DecentralizedNatServer {

	officerServer := utils.GetConfig().NatServerIP
	server := &DecentralizedNatServer{
		hosts: make([]string, len(officerServer)),
	}

	port := strconv.Itoa(utils.GetConfig().NatServerPort)

	for _, host := range officerServer {
		server.hosts = append(server.hosts, host+":"+port)
	}

	return server
}

//TODO:: use gossip protocol to manager all nat servers. we use official nat servers right now.
func (s *DecentralizedNatServer) GossipNatServer() string {
	return s.hosts[0] //TIPS:: simply use the first server.
}
