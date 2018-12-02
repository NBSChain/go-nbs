package nat

import (
	"fmt"
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
	KeepAliveTimeOut = time.Second * 45
	HolePunchTimeOut = 2 * time.Second
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

	KeepAlive := &net_pb.NatKeepAlive{
		NodeId: tunnel.networkId,
		LAddr:  tunnel.sharedAddr,
	}
	kaData, _ := proto.Marshal(KeepAlive)
	request := &net_pb.NatMsg{
		Typ:     nbsnet.NatKeepAlive,
		Len:     int32(len(kaData)),
		PayLoad: kaData,
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

func (tunnel *KATunnel) answerInvite(data []byte) {

	invite := &net_pb.ReverseInvite{}
	if err := proto.Unmarshal(data, invite); err != nil {
		logger.Warning("answer invite unmarshal err:->", err)
		return
	}
	myPort := strconv.Itoa(int(invite.ToPort))

	conn, err := shareport.DialUDP("udp4", "0.0.0.0:"+myPort,
		invite.PubIp+":"+invite.FromPort)
	if err != nil {
		logger.Errorf("failed to dial up peer to answer inviter:", err)
		return
	}
	defer conn.Close()
	InviteAck := &net_pb.ReverseInviteAck{
		SessionId: invite.SessionId,
	}
	ackData, _ := proto.Marshal(InviteAck)
	req := &net_pb.NatMsg{
		Typ:     nbsnet.NatReversDigAck,
		Len:     int32(len(ackData)),
		PayLoad: ackData,
	}

	reqData, _ := proto.Marshal(req)
	if _, err := conn.Write(reqData); err != nil {
		logger.Errorf("failed to write answer to inviter:", err)
		return
	}

	logger.Debug("Step4: answer the invite:->", conn.LocalAddr().String(), invite, req)
}

func (tunnel *KATunnel) refreshNatInfo(data []byte) {

	alive := &net_pb.NatKeepAlive{}

	if err := proto.Unmarshal(data, alive); err != nil {
		logger.Warning("unmarshal nat keep alive err:->", err)
		return
	}

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

	holeMsg := &net_pb.NatMsg{
		Typ: nbsnet.NatPriDigSyn,
		Seq: time.Now().Unix(),
	}
	data, _ := proto.Marshal(holeMsg)
	if _, err := conn.Write(data); err != nil {
		logger.Error("Step 2-2:private network dig dig failed:->", err)
		task.err <- err
		return
	}

	conStr := "[" + conn.LocalAddr().String() + "]-->[" + conn.RemoteAddr().String() + "]"
	logger.Info("Step 2-4:->dig in private network:->", conStr)

	if err := conn.SetReadDeadline(time.Now().Add(HolePunchTimeOut / 2)); err != nil {
		task.err <- err
		return
	}

	buffer := make([]byte, utils.NormalReadBuffer)
	_, err = conn.Read(buffer)
	if err != nil {
		logger.Error("Step 2-3:private network reading dig result failed:->", err, conStr)
		task.err <- err
		return
	}
	resMsg := &net_pb.NatMsg{}
	if err := proto.Unmarshal(buffer, resMsg); err != nil {
		task.err <- err
		return
	}

	if resMsg.Typ != nbsnet.NatPriDigAck || resMsg.Seq != holeMsg.Seq+1 {
		task.err <- fmt.Errorf("wrong ack package")
		return
	}

	task.err <- nil
	task.udpConn = conn
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
	response := &net_pb.NatMsg{}
	if err = proto.Unmarshal(buffer[:n], response); err != nil {
		logger.Warning("reading dig result Unmarshal failed:", err, conStr)
		task.err <- err
		return
	}

	logger.Debug("get dig response:->", response, conStr)

	switch response.Typ {
	case nbsnet.NatDigIn, nbsnet.NatDigOut:
		res := &net_pb.NatMsg{
			Typ:     nbsnet.NatDigSuccess,
			PayLoad: response.PayLoad,
			Len:     int32(len(response.PayLoad)),
		}

		data, _ := proto.Marshal(res)

		if _, err := conn.Write(data); err != nil {
			logger.Warning("failed to confirm the dig :->", conStr)
		}
	case nbsnet.NatDigSuccess:
		logger.Info("dig dig success:->", conStr)
	}

	task.err <- nil
	task.udpConn = conn
	return
}

func (tunnel *KATunnel) unpackMsg(payLoad *net_pb.HolePayLoad) {
	//sid := payLoad.SessionId
	//if pTask, ok := tunnel.workLoad[sid]; !ok{
	//	tunnel.
	//}
}
