package nat

import (
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/storage/network/shareport"
	"github.com/golang/protobuf/proto"
	"net"
	"strconv"
	"time"
)

const (
	KeepAliveTime       = time.Second * 15
	KeepAliveTimeOut    = 45
	HolePunchingTimeOut = 6
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

/************************************************************************
*
*			server side
*
*************************************************************************/

func (tunnel *KATunnel) answerInvite(invite *net_pb.ReverseInvite) {

	myPort := strconv.Itoa(int(invite.ToPort))

	conn, err := shareport.DialUDP("udp4", "0.0.0.0:"+myPort,
		invite.PubIp+":"+invite.FromPort)
	if err != nil {
		logger.Errorf("failed to dial up peer to answer inviter:", err)
		return
	}
	defer conn.Close()

	req := &net_pb.NatRequest{
		MsgType: net_pb.NatMsgType_ReverseDigACK,
		InviteAck: &net_pb.ReverseInviteAck{
			SessionId: invite.SessionId,
		},
	}

	data, _ := proto.Marshal(req)
	if _, err := conn.Write(data); err != nil {
		logger.Errorf("failed to write answer to inviter:", err)
		return
	}

	logger.Debug("Step4: answer the invite:->", conn.LocalAddr().String(), invite, req)
}

func (tunnel *KATunnel) refreshNatInfo(alive *net_pb.NatKeepAlive) {
	tunnel.updateTime = time.Now()

	if tunnel.natAddr.NatIp != alive.PubIP &&
		tunnel.natAddr.NatPort != alive.PubPort {

		tunnel.natAddr.NatIp = alive.PubIP
		tunnel.natAddr.NatPort = alive.PubPort
		//TODO::tell other node ,this has changed.
		logger.Info("node's nat info changed.")
	}
}
