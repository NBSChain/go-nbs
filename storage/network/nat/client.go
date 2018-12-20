package nat

import (
	"context"
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"net"
	"time"
)

var NoTimeOut = time.Time{}

const (
	KeepAliveTime    = time.Second * 100
	KeepAliveTimeOut = KeepAliveTime * 3
	BootStrapTimeOut = time.Second * 2
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
	natConn    *net.UDPConn
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
	go c.keepAlive()
	go c.readCmd()
	go c.listenPing()
	return c, nil

}

func (c *Client) findWhoAmI(canSever chan bool) error {

	for _, serverIp := range utils.GetConfig().NatServerIP {

		natServerAddr := &net.UDPAddr{
			IP:   net.ParseIP(serverIp),
			Port: utils.GetConfig().NatServerPort,
		}
		conn, err := net.DialUDP("udp4", nil, natServerAddr)
		if err != nil {
			logger.Warning("this nat server is done:->", serverIp, err)
			continue
		}
		conn.SetDeadline(time.Now().Add(BootStrapTimeOut))

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
		c.natConn = conn

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
		Typ: nbsnet.NatKeepAlive,
		KeepAlive: &net_pb.KeepAlive{
			NodeId: c.networkId,
			LAddr:  c.natConn.LocalAddr().String(),
		},
	}
	requestData, _ := proto.Marshal(request)
	for {
		select {
		case <-time.After(KeepAliveTime):

			if no, err := c.natConn.Write(requestData); err != nil || no == 0 {
				logger.Warning("failed to send keep alive channel message:", err)
				c.closeCtx()
				return
			}
			logger.Debug("send  keep alive to nat server:->", request.KeepAlive.LAddr)

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
	//TODO::
	c.updateTime = time.Now()

	if c.NatAddr.NatIp != alive.PubIP &&
		c.NatAddr.NatPort != alive.PubPort {

		c.NatAddr.NatIp = alive.PubIP
		c.NatAddr.NatPort = alive.PubPort
		logger.Warning("node's nat info changed.", alive)
	}
}

func (c *Client) listenPing() {

	conn, err := net.ListenUDP("udp4", &net.UDPAddr{
		Port: utils.GetConfig().NatPrivatePingPort,
	})
	if err != nil {
		logger.Warning("start private nat ping listener err:->", err)
		return
	}
	defer conn.Close()
	defer logger.Warning("ping listening exit")

	for {
		buffer := make([]byte, utils.NormalReadBuffer)
		n, peerAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			logger.Warning("private ping listener err:->", err)
			c.closeCtx()
			return
		}
		msg := &net_pb.NatMsg{}
		if err := proto.Unmarshal(buffer[:n], msg); err != nil {
			logger.Warning("private ping listener Unmarshal failed:->", err)
			continue
		}
		logger.Debug("is it a private nat dig apply?:->", msg, peerAddr)
		if msg.Typ != nbsnet.NatPriDigSyn {
			logger.Warning("I can't this message:->", msg.Typ)
			continue
		}

		res := &net_pb.NatMsg{
			Typ: nbsnet.NatPriDigAck,
			Seq: msg.Seq + 1,
		}
		data, _ := proto.Marshal(res)
		if _, err := conn.WriteToUDP(data, peerAddr); err != nil {
			logger.Warning("answer NatPriDigAck err:->", err)
			c.closeCtx()
			return
		}

		select {
		case <-c.Ctx.Done():
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

		n, err := c.natConn.Read(buffer)
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
				Params:  msg.DigConfirm,
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
