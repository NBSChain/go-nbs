package nat

import (
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/golang/protobuf/proto"
	"net"
	"time"
)

const (
	KeepAliveTime       = time.Second * 15
	KeepAliveTimeOut    = 45
	HolePunchingTimeOut = 4
)

type ProxyTask struct {
	sessionID string
	toAddr    *net.UDPAddr
	digResult chan error
}

type ConnTask struct {
	UdpConn *net.UDPConn
	Err     chan error
}

type KATunnel struct {
	natAddr    *nbsnet.NbsUdpAddr
	networkId  string
	serverHub  *net.UDPConn
	kaConn     *net.UDPConn
	sharedAddr string
	updateTime time.Time
	workLoad   map[string]*ProxyTask
	inviteTask map[string]*ConnTask
}

/************************************************************************
*
*			client side
*
*************************************************************************/
func (tunnel *KATunnel) Close() {

	tunnel.serverHub.Close()
	tunnel.kaConn.Close()
}

func (tunnel *KATunnel) runLoop() {

	for {
		tunnel.sendKeepAlive()

		time.Sleep(KeepAliveTime)
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

//TODO::
func (tunnel *KATunnel) connManage() {
}

//TODO::
func (tunnel *KATunnel) restoreNatChannel() {
}
