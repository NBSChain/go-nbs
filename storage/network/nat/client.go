package nat

import (
	"context"
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/storage/network/shareport"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"net"
	"strconv"
	"time"
)

var NoTimeOut = time.Time{}

const (
	KeepAliveTime    = time.Second * 100
	KeepAliveTimeOut = KeepAliveTime * 3
	HolePunchTimeOut = 4 * time.Second
	BootStrapTimeOut = time.Second * 2
	ErrNoBeforeRetry = 3
)

type ConnTask struct {
	err         error
	locPort     string
	udpConn     chan *net.UDPConn
	portCapConn *net.UDPConn
}

type Client struct {
	networkId  string
	CanServer  bool
	ctx        context.Context
	closeCtx   context.CancelFunc
	errNo      int
	natChanged chan struct{}
	NatAddr    *nbsnet.NbsUdpAddr
	conn       *net.UDPConn
	updateTime time.Time
	digTask    map[string]*ConnTask
}

func NewNatClient(networkId string, canServer chan bool) (*Client, error) {

	ctx, cancel := context.WithCancel(context.Background())
	c := &Client{
		ctx:        ctx,
		closeCtx:   cancel,
		networkId:  networkId,
		natChanged: make(chan struct{}),
		updateTime: time.Now(),
		digTask:    make(map[string]*ConnTask),
	}

	if err := c.findWhoAmI(canServer); err != nil {
		return nil, err
	}

	if c.CanServer {
		return c, nil
	}

	go c.keepAlive()

	go c.readCmd()

	return c, nil

}

func (c *Client) findWhoAmI(canSever chan bool) error {

	for _, serverIp := range utils.GetConfig().NatServerIP {

		natServerAddr := &net.UDPAddr{
			IP:   net.ParseIP(serverIp),
			Port: utils.GetConfig().NatServerPort,
		}

		conn, err := net.DialUDP("udp4", &net.UDPAddr{
			Port: utils.GetConfig().NatClientPort,
		}, natServerAddr)
		if err != nil {
			logger.Warning("this nat server is done:->", serverIp, err)
			continue
		}

		err = conn.SetDeadline(time.Now().Add(BootStrapTimeOut))
		if err != nil {
			logger.Warning("set dead line err:->", serverIp, err)
			conn.Close()
			continue
		}

		localAddr := conn.LocalAddr().String()
		host, port, err := nbsnet.SplitHostPort(localAddr)
		request := &net_pb.NatMsg{
			Typ: nbsnet.NatBootReg,
			Seq: time.Now().Unix(),
			BootReg: &net_pb.BootReg{
				NodeId:      c.networkId,
				PrivateIp:   host,
				PrivatePort: port,
			},
		}

		requestData, _ := proto.Marshal(request)
		if no, err := conn.Write(requestData); err != nil || no == 0 {
			logger.Error("failed to send nat request to selfNatServer:->", err, no)
			conn.Close()
			continue
		}

		responseData := make([]byte, utils.NormalReadBuffer)
		hasRead, err := conn.Read(responseData)
		if err != nil {
			logger.Error("reading failed from nat server:->", err)
			conn.Close()
			continue
		}

		response := &net_pb.NatMsg{}
		if err := proto.Unmarshal(responseData[:hasRead], response); err != nil {
			logger.Error("unmarshal err:->", err)
			conn.Close()
			continue
		}

		ack := response.BootAnswer
		c.NatAddr = &nbsnet.NbsUdpAddr{
			NetworkId: c.networkId,
			CanServe:  nbsnet.CanServe(ack.NatType),
			NatServer: natServerAddr.String(),
			PubIp:     ack.PublicIp,
			NatIp:     ack.PublicIp,
			NatPort:   ack.PublicPort,
			PriIp:     host,
		}

		if ack.NatType == net_pb.NatType_ToBeChecked {

			select {
			case can := <-canSever:
				c.NatAddr.CanServe = can

			case <-time.After(BootStrapTimeOut / 2):
				c.NatAddr.CanServe = false
			}
		}

		c.CanServer = c.NatAddr.CanServe
		conn.SetDeadline(NoTimeOut)
		c.conn = conn

		logger.Info("create client success for network")
		return nil
	}

	return fmt.Errorf("failed to find who am I")
}

//TODO:: data race
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
func (tunnel *Client) keepAlive() {

	for {
		select {
		case <-time.After(KeepAliveTime):
			if err := tunnel.sendKeepAlive(); err != nil {
				logger.Warning("failed to send nat keep alive message")
				if tunnel.errNo++; tunnel.errNo > ErrNoBeforeRetry {
					logger.Warning("too many times send errors")
					tunnel.reSetupChannel()
					return
				}
			}
		case <-tunnel.ctx.Done():
			logger.Info("exit sending thread cause's of context close")
			return
		}
	}
}

func (tunnel *Client) sendKeepAlive() error {

	request := &net_pb.NatMsg{
		Typ: nbsnet.NatKeepAlive,
		KeepAlive: &net_pb.KeepAlive{
			NodeId: tunnel.networkId,
			LAddr:  tunnel.sharedAddr,
		},
	}

	requestData, err := proto.Marshal(request)
	if err != nil {
		logger.Warning("failed to marshal keep alive message:", err)
		return err
	}

	if no, err := tunnel.kaConn.Write(requestData); err != nil || no == 0 {
		logger.Warning("failed to send keep alive channel message:", err)
		return err
	}

	return nil
}

//TODO::
func (tunnel *Client) reSetupChannel() {
	tunnel.errNo = 0
	tunnel.closeCtx()
}

/************************************************************************
*
*			server side
*
*************************************************************************/

func (tunnel *Client) answerInvite(invite *net_pb.ReverseInvite) {
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

func (tunnel *Client) refreshNatInfo(alive *net_pb.KeepAlive) {

	tunnel.updateTime = time.Now()

	if tunnel.NatAddr.NatIp != alive.PubIP &&
		tunnel.NatAddr.NatPort != alive.PubPort {

		tunnel.NatAddr.NatIp = alive.PubIP
		tunnel.NatAddr.NatPort = alive.PubPort
		tunnel.natChanged <- struct{}{}
		logger.Info("node's nat info changed.", alive)
	}
}

func (tunnel *Client) directDialInPriNet(lAddr, rAddr *nbsnet.NbsUdpAddr, task *ConnTask, toPort int, sessionID string) {

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
