package nat

import (
	"fmt"
	"time"
)

func (ch *KATunnel) readRegResponse() error {

	resChan := make(chan error)
	defer close(resChan)

	go ch.runLoop()

	select {
	case err := <-resChan:
		return err
	case <-time.After(time.Second * 2):
		return fmt.Errorf("timeout when reading channel register messae.")
	}
}

func (ch *KATunnel) runLoop(resChan chan error) {

	for {
		responseData := make([]byte, 2048)
		hasRead, peerAddr, err := peer.receivingHub.ReadFrom(responseData)
		if err != nil {
			fmt.Println("failed to read nat response from natServer", err)
			continue
		}

		response := &net_pb.Response{}
		if err := proto.Unmarshal(responseData[:hasRead], response); err != nil {
			fmt.Println("failed to unmarshal nat response data", err)
			continue
		}

		switch response.MsgType {
		case net_pb.NatMsgType_BootStrapReg:

			resValue := response.BootRegRes
			ch.publicIp = resValue.PublicIp
			ch.publicPort = resValue.PublicPort
			resChan <- nil
		}
	}
}
