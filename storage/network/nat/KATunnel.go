package nat

import (
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/golang/protobuf/proto"
	"net"
	"time"
)

const (
	KeepAlive        = 15
	KeepAliveTimeOut = 45
)

type ConnTask struct {
	Err       error
	sessionId string
	ConnCh    chan *net.UDPConn
}

type KATunnel struct {
	closed      chan bool
	networkId   string
	privateIP   string
	privatePort string
	publicIp    string
	publicPort  string
	receiveHub  *net.UDPConn
	kaConn      *net.UDPConn
	updateTime  time.Time
	natTask     map[string]*ConnTask
}

func (tunnel *KATunnel) InitNatChannel() error {

	request := &net_pb.NatRequest{
		MsgType: net_pb.NatMsgType_BootStrapReg,
		BootRegReq: &net_pb.BootNatRegReq{
			NodeId:      tunnel.networkId,
			PrivateIp:   tunnel.privateIP,
			PrivatePort: tunnel.privatePort,
		},
	}

	reqData, err := proto.Marshal(request)
	if err != nil {
		logger.Warning("failed to marshal nat keep alive message", err)
		return err
	}

	if no, err := tunnel.kaConn.Write(reqData); err != nil || no == 0 {
		logger.Warning("nat channel keep alive message failed", err, no)
		return err
	}

	if err := tunnel.readRegResponse(); err != nil {
		logger.Warning("failed to read channel initialize response.")
		return err
	}

	return nil
}

func (tunnel *KATunnel) Close() {

	tunnel.receiveHub.Close()
	tunnel.kaConn.Close()

	tunnel.closed <- true
}

func (tunnel *KATunnel) MakeANatConn(fromId, toId, connId string, port int) *ConnTask {

	sessionId := fromId + toId
	task := &ConnTask{
		sessionId: sessionId,
	}

	payload := &net_pb.NatConReq{
		FromPeerId: fromId,
		ToPeerId:   toId,
		SessionId:  sessionId,
	}
	request := &net_pb.NatRequest{
		MsgType: net_pb.NatMsgType_Connect,
		ConnReq: payload,
	}

	reqData, err := proto.Marshal(request)
	if err != nil {
		logger.Error("failed to marshal the nat connect request", err)
		task.Err = err
		return task
	}

	if no, err := tunnel.kaConn.Write(reqData); err != nil || no == 0 {
		logger.Warning("nat channel keep alive message failed", err, no)
		task.Err = err
		return task
	}

	task.ConnCh = make(chan *net.UDPConn)

	tunnel.natTask[sessionId] = task

	return task
}

func (tunnel *KATunnel) runLoop() {

	for {
		select {
		case <-time.After(time.Second * KeepAlive):
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
