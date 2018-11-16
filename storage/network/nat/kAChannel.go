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
const(
	NatKeepAlive = 15
)

type KAChannel struct {
	closed		chan bool
	networkId	string
	privateIP	string
	privatePort	string
	publicIp	string
	publicPort	string
	receiveHub *net.UDPConn
	kaConn     *net.UDPConn
	updateTime	time.Time
}

func (ch *KAChannel) InitNatChannel() error{

	request := &net_pb.NatRequest{
		MsgType: net_pb.NatMsgType_BootStrapReg,
		BootRegReq: &net_pb.BootNatRegReq{
			NodeId:      ch.networkId,
			PrivateIp:   ch.privateIP,
			PrivatePort: ch.privatePort,
		},
	}

	requestData, err := proto.Marshal(request)
	if err != nil {
		logger.Warning("failed to marshal nat keep alive message", err)
		return err
	}

	if no, err := ch.kaConn.Write(requestData); err != nil || no == 0 {
		logger.Warning("nat channel keep alive message failed", err, no)
		return err
	}

	if err := ch.readRegResponse(); err != nil{
		logger.Warning("failed to read channel initialize response.")
		return err
	}

	return nil
}



func (nat *Manager) NewKAChannel() (*KAChannel, error) {

	port := strconv.Itoa(utils.GetConfig().NatClientPort)
	natServer := nat.dNatServer.GossipNatServer()

	listener, err := shareport.ListenUDP("udp4", port)
	if err != nil{
		logger.Warning("create share listening udp failed.")
		return nil, err
	}

	client, err := shareport.DialUDP("udp4", "0.0.0.0:" + port, natServer)
	if err != nil{
		logger.Warning("create share port dial udp connection failed.")
		return nil, err
	}

	localAddr := client.LocalAddr().String()
	priIP, priPort, err := net.SplitHostPort(localAddr)

	channel := &KAChannel{
		closed:		make(chan bool),
		networkId:  	nat.networkId,
		receiveHub: 	listener,
		kaConn:     	client,
		privateIP:	priIP,
		privatePort:	priPort,
		updateTime:	time.Now(),
	}

	return channel,nil
}

func (ch *KAChannel) Close(){

	ch.receiveHub.Close()
	ch.kaConn.Close()

	ch.closed <- true
}
