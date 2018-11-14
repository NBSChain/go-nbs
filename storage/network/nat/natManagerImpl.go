package nat

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/utils"
)

func (nat *nbsNatManager) FindWhoAmI() error {

	config := utils.GetConfig()

	for _, serverIP := range config.NatServerIP {

		//TIPS:: no need to bind local host and local port right now
		connection, err := nat.connectToNatServer(serverIP)
		if err != nil {
			logger.Error("can't know who am I", err)
			goto CloseConn
		}

		if err := nat.sendNatRequest(connection); err != nil {
			goto CloseConn
		}

		_, err = nat.parseNatResponse(connection)
		if err == nil {
			break
		}

	CloseConn:
		connection.Close()
	}

	return nil
}

func (nat *nbsNatManager) GetStatus() string {

	status := fmt.Sprintf("\n=========================================================================\n"+
		"\tnetworkId:\t%s\n"+
		"\tnatType:\t%v\n"+
		"\tpublicAddress:\t%s\n"+
		"\tprivateIP:\t%s\n"+
		"=========================================================================",
		nat.networkId,
		nat.natType,
		nat.publicAddress.String(),
		nat.privateIP)

	return status
}

func (nat *nbsNatManager) NatType() nat_pb.NatType {
	nat.Lock()
	defer nat.Unlock()
	return nat.natType
}
