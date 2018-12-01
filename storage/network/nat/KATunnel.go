package nat

import (
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/storage/network/shareport"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"net"
	"strconv"
	"time"
)

const (
	KeepAliveTime    = time.Second * 15
	KeepAliveTimeOut = 45
	HolePunchTimeOut = 4 * time.Second
	TryDigHoleTimes  = 3
	BootStrapTimeOut = time.Second * 4
)

const (
	_ int32 = iota
	FromPriNet
	FromPubNet
	ToPubNet
)

type ProxyTask struct {
	sessionID string
	toAddr    *net.UDPAddr
	err       chan error
}

type ConnTask struct {
	err     chan error
	udpConn *net.UDPConn
}

type KATunnel struct {
	natChanged chan struct{}
	natAddr    *nbsnet.NbsUdpAddr
	networkId  string
	serverHub  *net.UDPConn
	kaConn     *net.UDPConn
	sharedAddr string
	updateTime time.Time
	workLoad   map[string]*ProxyTask
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
		if err := tunnel.sendKeepAlive(); err != nil {
			logger.Warning("failed to send nat keep alive message")
			//TODO::if failed more than 3 times, wo need to find new nat server
		}

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
		tunnel.natChanged <- struct{}{}
		logger.Info("node's nat info changed.", alive)
	}
}

func (tunnel *KATunnel) directDialInPriNet(lAddr, rAddr *nbsnet.NbsUdpAddr, task *ConnTask, toPort int, sessionID string) {

	conn, err := net.DialUDP("udp4", &net.UDPAddr{
		IP:   net.ParseIP(lAddr.PriIp),
		Port: int(lAddr.PriPort),
	}, &net.UDPAddr{
		IP:   net.ParseIP(rAddr.PriIp),
		Port: toPort,
	})

	if err != nil {
		logger.Warning("Step 2-1:can't dial by private network.", err)
		task.err <- err
		return
	}

	holeMsg := &net_pb.NatRequest{
		MsgType: net_pb.NatMsgType_DigIn,
		HoleMsg: &net_pb.HoleDig{
			SessionId:   sessionID,
			NetworkType: FromPriNet,
		},
	}
	data, _ := proto.Marshal(holeMsg)

	go tunnel.waitDigResponse(task, conn)

	logger.Info("Step 2-4:->dig in private network:->")
	tunnel.digDig(data, conn, task)
}

func (tunnel *KATunnel) digDig(data []byte, conn *net.UDPConn, task *ConnTask) {

	for i := 0; i < TryDigHoleTimes; i++ {

		if _, err := conn.Write(data); err != nil {
			logger.Error(err)
		}
		select {
		case <-task.err:
			logger.Debug("dig in action finished.")
			return
		case <-time.After(time.Second):
			logger.Debug("dig again")
		}
	}
}

func (tunnel *KATunnel) waitDigResponse(task *ConnTask, conn *net.UDPConn) {

	conStr := "[" + conn.LocalAddr().String() + "]-->[" + conn.RemoteAddr().String() + "]"

	buffer := make([]byte, utils.NormalReadBuffer)
	n, err := conn.Read(buffer)
	if err != nil {
		logger.Error("reading dig result failed:->", err, conStr)
		task.err <- err
		return
	}
	response := &net_pb.NatResponse{}
	if err = proto.Unmarshal(buffer[:n], response); err != nil {
		logger.Warning("reading dig result Unmarshal failed:", err, conStr)
		task.err <- err
		return
	}

	logger.Debug("get dig response:->", response, conStr)

	switch response.MsgType {
	case net_pb.NatMsgType_DigIn, net_pb.NatMsgType_DigOut:
		res := &net_pb.NatResponse{
			MsgType: net_pb.NatMsgType_DigSuccess,
			HoleMsg: response.HoleMsg,
		}

		data, _ := proto.Marshal(res)

		if _, err := conn.Write(data); err != nil {
			logger.Warning("failed to confirm the dig :->", conStr)
		}
	case net_pb.NatMsgType_DigSuccess:
		logger.Info("dig dig success:->", conStr)
	}

	task.err <- nil
	task.udpConn = conn
	return
}
