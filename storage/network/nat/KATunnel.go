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
	networkId   string
	closed      chan bool
	serverHub   *net.UDPConn
	kaConn      *net.UDPConn
	sharedAddr  string
	updateTime  time.Time
	natTask     map[string]*nbsnet.ConnTask
	connManager map[string]*nbsnet.HoleConn
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

func (tunnel *KATunnel) natHoleStep1(task *nbsnet.ConnTask, lAddr, rAddr *nbsnet.NbsUdpAddr, connId string) error {

	connReq := &net_pb.NatConInvite{

		FromAddr: &net_pb.NbsAddr{
			NetworkId: lAddr.NetworkId,
			CanServer: lAddr.CanServe,
			PubIp:     lAddr.PubIp,
			PubPort:   int32(lAddr.PubPort),
			PriIp:     lAddr.PriIp,
			PriPort:   int32(lAddr.PriPort),
		},
		ToAddr: &net_pb.NbsAddr{
			NetworkId: rAddr.NetworkId,
			CanServer: rAddr.CanServe,
			PubIp:     rAddr.PubIp,
			PubPort:   int32(rAddr.PubPort),
			PriIp:     rAddr.PriIp,
			PriPort:   int32(rAddr.PriPort),
		},
		SessionId: connId,
		CType:     int32(task.CType),
	}

	response := &net_pb.NatRequest{
		MsgType: net_pb.NatMsgType_Connect,
		ConnReq: connReq,
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
