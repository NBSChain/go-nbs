package nat

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/utils"
	"time"
)

func (nat *NbsNatManager) FindWhoAmI() (address *NbsAddress, err error) {

	config := utils.GetConfig()

	for _, serverIP := range config.NatServerIP {

		conn, err := nat.connectToNatServer(serverIP)
		if err != nil {
			logger.Error("can't know who am I", err)
			conn.Close()
			continue
		}
		conn.SetDeadline(time.Now().Add(time.Second * 3))

		localHost, err := nat.sendNatRequest(conn)
		if err != nil {
			logger.Error("failed to read nat response:", err)
			conn.Close()
			continue
		}

		response, err := nat.parseNatResponse(conn)
		if err != nil {
			logger.Debug("get NAT server info success.")
			conn.Close()
			continue
		}

		address = &NbsAddress{
			PublicIP:  response.PublicIp,
			PrivateIp: localHost,
			IsInPub:   IsPublic(response.NatType),
		}

		if response.NatType == net_pb.NatType_ToBeChecked {

			select {
			case canServer := <-nat.canServe:
				address.IsInPub = canServer
			case <-time.After(time.Second * BootStrapNatServerTimeOutInSec / 2):
				address.IsInPub = false
			}
			close(nat.canServe)
		}

		return address, nil
	}

	return nil, fmt.Errorf("can't find available NAT server")
}
