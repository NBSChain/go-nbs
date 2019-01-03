package nat

import (
	"context"
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/storage/network/shareport"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"net"
	"sync"
	"time"
)

var NoTimeOut = time.Time{}

const (
	KeepAliveTime    = time.Second * 100
	BootStrapTimeOut = time.Second * 4
	CmdTaskPoolSize  = 100
	CMDAnswerInvite  = 1
	CMDDigOut        = 2
	CMDDigSetup      = 3
)

type CmdProcess func(interface{}) error

type ClientCmd struct {
	CmdType int
	Params  interface{}
}

type Client struct {
	networkId  string
	CanServer  bool
	Ctx        context.Context
	closeCtx   context.CancelFunc
	NatAddr    *nbsnet.NbsUdpAddr
	CtrlConn   *net.TCPConn
	updateTime time.Time
	CmdTask    chan *ClientCmd
}

func NewNatClient(networkId string, canServer chan bool) (*Client, error) {

	ctx, cancel := context.WithCancel(context.Background())
	c := &Client{
		Ctx:        ctx,
		closeCtx:   cancel,
		networkId:  networkId,
		updateTime: time.Now(),
		CmdTask:    make(chan *ClientCmd, CmdTaskPoolSize),
	}

	if err := c.findWhoAmI(canServer); err != nil {
		return nil, err
	}

	if c.CanServer {
		return c, nil
	}

	c.checkMyNetType()

	go c.keepAlive()
	go c.readCmd()
	go c.listenInPrivate()
	return c, nil

}

func (c *Client) findWhoAmI(canSever chan bool) error {

	for _, serverIp := range utils.GetConfig().NatServerIP {

		natServerAddr := nbsnet.JoinHostPort(serverIp, int32(utils.GetConfig().NatServerPort))

		cc, err := net.DialTimeout("tcp4", natServerAddr, BootStrapTimeOut)
		if err != nil {
			logger.Warning("this nat server is done:->", serverIp, err)
			continue
		}
		conn := cc.(*net.TCPConn)
		localAddr := conn.LocalAddr().String()

		host, port, err := nbsnet.SplitHostPort(localAddr)
		request := &net_pb.NatMsg{
			Typ:   nbsnet.NatBootReg,
			NetID: c.networkId,
			BootReg: &net_pb.BootReg{
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
			NetworkId:   c.networkId,
			CanServe:    nbsnet.CanServe(ack.NatType),
			NatServerIP: serverIp,
			PubIp:       ack.PublicIp,
			NatIp:       ack.PublicIp,
			NatPort:     ack.PublicPort,
			PriIp:       host,
		}

		if ack.NatType == net_pb.NatType_ToBeChecked {
			select {
			case can := <-canSever:
				c.NatAddr.CanServe = can

			case <-time.After(BootStrapTimeOut):
				c.NatAddr.CanServe = false
			}
		}

		c.CanServer = c.NatAddr.CanServe
		c.CtrlConn = conn
		logger.Info("create client success for network:->\n", c.NatAddr.String())
		return nil
	}

	return fmt.Errorf("failed to find who am I")
}

/************************************************************************
*
*			client side
*
*************************************************************************/
func (c *Client) keepAlive() {
	request := &net_pb.NatMsg{
		Typ:   nbsnet.NatKeepAlive,
		NetID: c.networkId,
		KeepAlive: &net_pb.KeepAlive{
			PriAddr: c.CtrlConn.LocalAddr().String(),
		},
	}
	requestData, _ := proto.Marshal(request)
	for {
		select {
		case <-time.After(KeepAliveTime):

			if no, err := c.CtrlConn.Write(requestData); err != nil || no == 0 {
				logger.Warning("failed to send keep alive channel message:", err)
				c.closeCtx()
				return
			}
			logger.Debug("send  keep alive to nat server:->", request.KeepAlive.PriAddr)

		case <-c.Ctx.Done():
			logger.Info("exit sending thread cause's of context close")
			return
		}
	}
}

/************************************************************************
*
*			server side
*
*************************************************************************/

func (c *Client) refreshNatInfo(alive *net_pb.KeepAlive) {
	c.updateTime = time.Now()
	oriNatInfo := nbsnet.JoinHostPort(c.NatAddr.NatIp, c.NatAddr.NatPort)
	if oriNatInfo != alive.PubAddr {
		pubIp, pubPort, _ := nbsnet.SplitHostPort(alive.PubAddr)
		c.NatAddr.NatIp = pubIp
		c.NatAddr.NatPort = pubPort
		logger.Warning("node's nat info changed.", alive)
	}
}

func (c *Client) listenInPrivate() {

	lisConn, err := net.ListenTCP("tcp4", &net.TCPAddr{
		Port: utils.GetConfig().NatPrivatePingPort,
	})
	if err != nil {
		logger.Warning("start private nat ping listener err:->", err)
		return
	}
	defer logger.Warning("ping listening exit")

	res := &net_pb.NatMsg{
		Typ:   nbsnet.NatPriDigAck,
		NetID: c.networkId,
	}
	data, _ := proto.Marshal(res)

	logger.Info("start listen private connection request:->", lisConn.Addr().String())
	for {
		conn, err := lisConn.AcceptTCP()
		if err != nil {
			logger.Warning("private ping listener err:->", err)
			c.closeCtx()
			lisConn.Close()
			return
		}

		logger.Debug("accept a private connection:->", nbsnet.ConnString(conn))

		if _, err := conn.Write(data); err != nil {
			logger.Warning("write response err:->", err)
		}
		conn.Close()

		select {
		case <-c.Ctx.Done():
			lisConn.Close()
			logger.Info("exit sending thread cause's of context close")
			return
		default:
			logger.Info("Step 1-6:->answer dig in private:->", res)
		}
	}
}

func (c *Client) readCmd() {

	for {
		buffer := make([]byte, utils.NormalReadBuffer)

		n, err := c.CtrlConn.Read(buffer)
		if err != nil {
			logger.Warning("reading keep alive message failed:->", err)
			c.closeCtx()
			return
		}

		msg := &net_pb.NatMsg{}
		if err := proto.Unmarshal(buffer[:n], msg); err != nil {
			logger.Warning("keep alive msg Unmarshal failed:->", err)
			continue
		}

		logger.Debug("KA c receive message:->", msg)

		switch msg.Typ {
		case nbsnet.NatKeepAlive:
			c.refreshNatInfo(msg.KeepAlive)

		case nbsnet.NatReversInvite:
			c.CmdTask <- &ClientCmd{
				CmdType: CMDAnswerInvite,
				Params:  msg.ReverseInvite,
			}

		case nbsnet.NatDigApply:
			c.CmdTask <- &ClientCmd{
				CmdType: CMDDigOut,
				Params:  msg.DigApply,
			}

		case nbsnet.NatDigConfirm:
			c.CmdTask <- &ClientCmd{
				CmdType: CMDDigSetup,
				Params:  msg,
			}

		}

		select {
		case <-c.Ctx.Done():
			logger.Info("exit reading thread cause's of context close")
			return
		default:
		}
	}
}

func (c *Client) checkMyNetType() {

	logger.Debug("start to check nat type ")

	msg := net_pb.NatMsg{
		Typ: nbsnet.NatCheckNetType,
	}
	data, _ := proto.Marshal(&msg)

	ips, ports := make(map[string]struct{}), make(map[string]struct{})
	var waitGrp sync.WaitGroup
	var locker sync.Mutex

	checkAddr := &net.UDPAddr{
		Port: utils.GetConfig().NetTypeCheckPort,
	}
	for _, serverIp := range utils.GetConfig().NatServerIP {

		serverHost := nbsnet.JoinHostPort(serverIp, int32(utils.GetConfig().HolePuncherPort))
		waitGrp.Add(1)

		go func() {
			defer waitGrp.Done()
			conn, err := shareport.DialUDP("udp4", checkAddr.String(), serverHost)
			if err != nil {
				logger.Warning("this nat server is down:->", err, serverHost)
				return
			}
			defer conn.Close()

			logger.Debug("request from server:->", serverHost)

			if _, err := conn.Write(data); err != nil {
				logger.Warning("write check nat type msg err:->", err)
				return
			}

			conn.SetDeadline(time.Now().Add(BootStrapTimeOut))

			buffer := make([]byte, utils.NormalReadBuffer)
			n, err := conn.Read(buffer)
			if err != nil {
				logger.Warning("read nat type err:->", err)
				return
			}

			res := net_pb.NatMsg{}
			proto.Unmarshal(buffer[:n], &res)

			logger.Debug("get nat response:->", res)

			ack := res.NatTypeCheck
			locker.Lock()
			ips[ack.Ip] = struct{}{}
			ports[ack.Port] = struct{}{}
			locker.Unlock()
		}()
	}

	waitGrp.Wait()

	c.NatAddr.AllPubIps = make([]string, 0)
	var networkType = nbsnet.SigIpSigPort

	if len(ports) == 1 {
		if len(ips) > 1 {
			networkType = nbsnet.MulIpSigPort
			for ip := range ips {
				c.NatAddr.AllPubIps = append(c.NatAddr.AllPubIps, ip)
			}
		}

	} else {
		logger.Warning("we don't support multi port model to punch hole.")
		networkType = nbsnet.MultiPort
	}

	c.NatAddr.NetType = networkType
	logger.Debug("my network type:->", networkType, c.NatAddr.AllPubIps)
}
