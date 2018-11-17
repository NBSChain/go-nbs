package nat

import (
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/storage/network/shareport"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"net"
	"strconv"
	"time"
)

const (
	KeepAlive = 15
)

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

func (nat *Manager) NewKAChannel() (*KATunnel, error) {

	port := strconv.Itoa(utils.GetConfig().NatClientPort)
	natServer := nat.dNatServer.GossipNatServer()

	listener, err := shareport.ListenUDP("udp4", port)
	if err != nil {
		logger.Warning("create share listening udp failed.")
		return nil, err
	}

	client, err := shareport.DialUDP("udp4", "0.0.0.0:"+port, natServer)
	if err != nil {
		logger.Warning("create share port dial udp connection failed.")
		return nil, err
	}

	localAddr := client.LocalAddr().String()
	priIP, priPort, err := net.SplitHostPort(localAddr)

	channel := &KATunnel{
		closed:      make(chan bool),
		networkId:   nat.networkId,
		receiveHub:  listener,
		kaConn:      client,
		privateIP:   priIP,
		privatePort: priPort,
		updateTime:  time.Now(),
	}

	return channel, nil
}

func (tunnel *KATunnel) MakeANatConn(fromId, toId, connId string, port int) (chan *net.UDPConn, error) {

	payload := &net_pb.NatConReq{
		FromPeerId: fromId,
		ToPeerId:   toId,
	}
	request := &net_pb.NatRequest{
		MsgType: net_pb.NatMsgType_Connect,
		ConnReq: payload,
	}

	reqData, err := proto.Marshal(request)
	if err != nil {
		logger.Error("failed to marshal the nat connect request", err)
		return nil, err
	}

	if no, err := tunnel.kaConn.Write(reqData); err != nil || no == 0 {
		logger.Warning("nat channel keep alive message failed", err, no)
		return nil, err
	}

	connChan := make(chan *net.UDPConn)

	return connChan, nil
}

func (nat *Manager) invitePeers(req *net_pb.NatConReq, peerAddr *net.UDPAddr) error {
	return nil
}
