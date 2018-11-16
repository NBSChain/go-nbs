//+build !windows

package nat

import (
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"time"
)

func (ch *KAChannel) readRegResponse() error{

	responseData := make([]byte, NetIoBufferSize)
	hasRead, err := ch.kaConn.Read(responseData)
	if err != nil {
		logger.Warning("reading failed from nat server", err)
		return err
	}

	response := &net_pb.Response{}
	if err := proto.Unmarshal(responseData[:hasRead], response); err != nil {
		logger.Warning("unmarshal err:", err)
		return err
	}

	logger.Info("response:", response)

	resValue := response.BootRegRes

	ch.publicIp = resValue.PublicIp
	ch.publicPort = resValue.PublicPort

	go ch.runLoop()
	go ch.listening()
	go ch.readKeepAlive()
	return nil
}

func (ch *KAChannel) runLoop() {

	for {
		select {
			case <- time.After(time.Second * NatKeepAlive):
				ch.sendKeepAlive()
			case <- ch.closed:
				logger.Info("keep alive channel is closed.")
				return
		}
	}
}

func (ch *KAChannel) sendKeepAlive() error{

	request := &net_pb.NatRequest{
		MsgType: net_pb.NatMsgType_KeepAlive,
		KeepAlive:&net_pb.NatKeepAlive{
			NodeId: ch.networkId,
		},
	}

	requestData, err := proto.Marshal(request)
	if err != nil {
		logger.Warning("failed to marshal keep alive request:", err)
		return err
	}

	logger.Debug("keep alive channel start")

	if no, err := ch.kaConn.Write(requestData); err != nil || no == 0 {
		logger.Warning("failed to send keep alive channel message:", err)
		return err
	}

	return nil
}

func (ch *KAChannel) readKeepAlive()  {

	for{
		buffer := make([]byte, utils.NormalReadBuffer)

		n, err := ch.kaConn.Read(buffer)
		if err != nil{
			logger.Warning("reading keep alive message failed:", err)
			continue
		}

		if err := ch.processResponse(buffer[:n]); err != nil{
			logger.Warning("process nat response message failed")
			continue
		}
	}
}

func (ch *KAChannel) processResponse(buffer []byte) error{

	response :=  &net_pb.Response{}
	if err := proto.Unmarshal(buffer, response); err != nil{
		logger.Warning("keep alive response Unmarshal failed:", err)
		return err
	}

	switch response.MsgType {
	case net_pb.NatMsgType_KeepAlive:
		ch.updateTime = time.Now()
		
	}

	return nil
}


func (ch *KAChannel) listening() {

	for{
		buffer := make([]byte, utils.NormalReadBuffer)
		n, err := ch.kaConn.Read(buffer)
		if err != nil{
			logger.Warning("reading keep alive message failed:", err)
			continue
		}

		if err := ch.processResponse(buffer[:n]); err != nil{
			logger.Warning("process nat response message failed")
			continue
		}
	}


}

func (ch *KAChannel) restoreNatChannel(){

}