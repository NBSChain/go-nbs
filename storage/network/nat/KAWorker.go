//+build !windows

package nat

import (
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"time"
)

/************************************************************************
*
*			for linux unix darwin and so on
*
*************************************************************************/
func (tunnel *KATunnel) readKeepAlive() {

	for {
		buffer := make([]byte, utils.NormalReadBuffer)

		n, err := tunnel.kaConn.Read(buffer)
		if err != nil {
			logger.Warning("reading keep alive message failed:", err)
			continue
		}

		if err := tunnel.process(buffer[:n]); err != nil {
			logger.Warning("process nat response message failed")
			continue
		}
	}
}

func (tunnel *KATunnel) process(buffer []byte) error {

	response := &net_pb.Response{}
	if err := proto.Unmarshal(buffer, response); err != nil {
		logger.Warning("keep alive response Unmarshal failed:", err)
		return err
	}

	switch response.MsgType {
	case net_pb.NatMsgType_KeepAlive:
		tunnel.updateTime = time.Now()

	case net_pb.NatMsgType_Connect:
		tunnel.punchAHole(response.ConnRes)
	}

	return nil
}

func (tunnel *KATunnel) listening() {

	for {
		buffer := make([]byte, utils.NormalReadBuffer)
		n, err := tunnel.serverHub.Read(buffer)
		if err != nil {
			logger.Warning("receiving port:", err)
			continue
		}

		if err := tunnel.process(buffer[:n]); err != nil {
			logger.Warning("process nat response message failed")
			continue
		}
	}
}

//TODO::
func (tunnel *KATunnel) restoreNatChannel() {

}
