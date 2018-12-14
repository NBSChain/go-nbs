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
	KeepAliveTime    = time.Second * 20
	KeepAliveTimeOut = KeepAliveTime * 3
	HolePunchTimeOut = 4 * time.Second
	BootStrapTimeOut = time.Second * 2
)

type ConnTask struct {
	err         error
	locPort     string
	udpConn     chan *net.UDPConn
	portCapConn *net.UDPConn
}

type KATunnel struct {
	natChanged chan struct{}
	natAddr    *nbsnet.NbsUdpAddr
	networkId  string
	serverHub  *net.UDPConn
	kaConn     *net.UDPConn
	sharedAddr string
	updateTime time.Time
	digTask    map[string]*ConnTask
}

func (task *ConnTask) Close() {
	if task.portCapConn != nil {
		_ = task.portCapConn.Close()
	}
}

/************************************************************************
*
*			client side
*
*************************************************************************/
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

	request := &net_pb.NatMsg{
		Typ: nbsnet.NatKeepAlive,
		KeepAlive: &net_pb.KeepAlive{
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

	req := &net_pb.NatMsg{
		Typ: nbsnet.NatReversInviteAck,
		ReverseInviteAck: &net_pb.ReverseInviteAck{
			SessionId: invite.SessionId,
		},
	}

	reqData, _ := proto.Marshal(req)
	logger.Debug("Step4: answer the invite:->", conn.LocalAddr().String(), invite.SessionId)

	for i := 0; i < 3; i++ {
		if _, err := conn.Write(reqData); err != nil {
			logger.Errorf("failed to write answer to inviter:", err)
			return
		}
	}
}

func (tunnel *KATunnel) refreshNatInfo(alive *net_pb.KeepAlive) {

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
		logger.Warning("Step 1-1:can't dial by private network.", err)
		task.err = err
		task.udpConn <- nil
		return
	}
	conStr := "[" + conn.LocalAddr().String() + "]-->[" + conn.RemoteAddr().String() + "]"

	logger.Debug("hole punch step1-2 start in private network:->", conStr)

	holeMsg := &net_pb.NatMsg{
		Typ: nbsnet.NatPriDigSyn,
		Seq: time.Now().Unix(),
	}
	data, _ := proto.Marshal(holeMsg)
	if _, err := conn.Write(data); err != nil {
		logger.Error("Step 1-3:private network dig dig failed:->", err)
		task.err = err
		task.udpConn <- nil
		return
	}

	if err := conn.SetReadDeadline(time.Now().Add(HolePunchTimeOut / 2)); err != nil {
		task.err = err
		task.udpConn <- nil
		return
	}

	buffer := make([]byte, utils.NormalReadBuffer)
	n, err := conn.Read(buffer)
	if err != nil {
		logger.Error("Step 1-5:private network reading dig result err:->", err)
		task.err = err
		task.udpConn <- nil
		return
	}
	resMsg := &net_pb.NatMsg{}
	if err := proto.Unmarshal(buffer[:n], resMsg); err != nil {
		logger.Info("Step 1-4:->dig in private network Unmarshal err:->", err)
		task.err = err
		task.udpConn <- nil
		return
	}

	if resMsg.Typ != nbsnet.NatPriDigAck || resMsg.Seq != holeMsg.Seq+1 {
		task.err = fmt.Errorf("wrong ack package")
		task.udpConn <- nil
		return
	}

	if err := conn.SetDeadline(time.Time{}); err != nil {
		task.err = err
		task.udpConn <- nil
	}

	task.err = nil
	task.udpConn <- conn
}
