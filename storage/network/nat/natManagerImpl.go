package nat

import (
	"fmt"
	"github.com/NBSChain/go-nbs/utils"
	"time"
)

func (nat *NbsNatManager) FindWhoAmI() error {

	config := utils.GetConfig()
	var success = false
	for _, serverIP := range config.NatServerIP {

		conn, err := nat.connectToNatServer(serverIP)
		if err != nil {
			logger.Error("can't know who am I", err)
			goto CloseConn
		}
		conn.SetDeadline(time.Now().Add(time.Second * 3))

		if err := nat.sendNatRequest(conn); err != nil {
			logger.Error("failed to read nat response:", err)
			goto CloseConn
		}

		_, err = nat.parseNatResponse(conn)
		if err == nil {
			logger.Debug("get NAT server info success.")
			success = true
			break
		}

	CloseConn:
		conn.Close()
	}

	if !success {
		return fmt.Errorf("failed to get nat information")
	}

	return nil
}

func (nat *NbsNatManager) GetStatus() string {

	status := fmt.Sprintf("\n=========================================================================\n"+
		"\tnetworkId:\t%s\n"+
		"\tnatType:\t%v\n"+
		"\tpublicAddress:\t%s\n"+
		"\tprivateIP:\t%s\n"+
		"=========================================================================",
		nat.networkId,
		nat.NatType,
		nat.PublicAddress.String(),
		nat.PrivateIP)

	return status
}
