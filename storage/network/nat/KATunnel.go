package nat

import (
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/golang/protobuf/proto"
	"net"
	"time"
)

const (
	KeepAliveTime    = time.Second * 15
	KeepAliveTimeOut = 45
)

type KATunnel struct {
	networkId  string
	closed     chan bool
	serverHub  *net.UDPConn
	kaConn     *net.UDPConn
	sharedAddr string
	updateTime time.Time
	natTask    map[string]*nbsnet.ConnTask
}

/************************************************************************
*
*			client side
*
*************************************************************************/
func (tunnel *KATunnel) Close() {

	tunnel.serverHub.Close()
	tunnel.kaConn.Close()
	tunnel.closed <- true
}

func (tunnel *KATunnel) runLoop() {

	for {
		select {
		case <-time.After(KeepAliveTime):
			tunnel.sendKeepAlive()
		case <-tunnel.closed:
			logger.Info("keep alive channel is closed.")
			return
		}
	}
}

func (tunnel *KATunnel) sendKeepAlive() error {

	request := &net_pb.NatRequest{
		MsgType: net_pb.NatMsgType_KeepAlive,
		KeepAlive: &net_pb.NatKeepAlive{
			NodeId: tunnel.networkId,
			LAddr:  tunnel.sharedAddr,
		},
	}

	requestData, err := proto.Marshal(request)
	if err != nil {
		logger.Warning("failed to marshal keep alive request:", err)
		return err
	}

	logger.Debug("keep alive channel start")

	if no, err := tunnel.kaConn.Write(requestData); err != nil || no == 0 {
		logger.Warning("failed to send keep alive channel message:", err)
		return err
	}

	return nil
}

func (tunnel *KATunnel) invitePeer(task *nbsnet.ConnTask, lAddr, rAddr *nbsnet.NbsUdpAddr, connId string) error {

	connRes := &net_pb.NatConReq{
		FromPeerId: lAddr.NetworkId,
		PublicIp:   lAddr.PriIp,
		PrivateIp:  lAddr.PriIp,
		SessionId:  connId,
		CType:      int32(task.CType),
	}

	response := &net_pb.Response{
		MsgType: net_pb.NatMsgType_Connect,
		ConnRes: connRes,
	}

	toItemData, err := proto.Marshal(response)
	if err != nil {
		return err
	}

	if _, err := tunnel.kaConn.Write(toItemData); err != nil {
		return err
	}
	return nil
}
