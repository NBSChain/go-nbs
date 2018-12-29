package network

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

type connTask struct {
	err         error
	udpConn     chan *net.UDPConn
	portCapConn *net.UDPConn
}

func (task *connTask) finish(err error, conn *net.UDPConn) {
	task.err = err
	task.udpConn <- conn
}

/************************************************************************
*
*			private functions
*
*************************************************************************/

func (network *nbsNetwork) makeDirectConn(lAddr, rAddr *nbsnet.NbsUdpAddr, toPort int) (*nbsnet.NbsUdpConn, error) {

	if rAddr == nil {
		return nil, fmt.Errorf("remote address can't be nil")
	}

	rUdpAddr := &net.UDPAddr{
		Port: toPort,
		IP:   net.ParseIP(rAddr.PubIp),
	}

	lUdpAddr := &net.UDPAddr{
		IP: net.ParseIP(lAddr.PriIp),
	}

	c, err := net.DialUDP("udp4", lUdpAddr, rUdpAddr)
	if err != nil {
		return nil, err
	}

	logger.Debug("Step6:make direct connection:->", c.LocalAddr().String(), c.RemoteAddr().String())
	return nbsnet.NewNbsConn(c, nbsnet.CTypeNormal), nil
}

func (network *nbsNetwork) invitePeerBehindNat(lAddr, rAddr *nbsnet.NbsUdpAddr, toPort int) (*net.UDPConn, error) {

	conn, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   net.ParseIP(rAddr.NatServerIP),
		Port: utils.GetConfig().HolePuncherPort,
	})

	if err != nil {
		return nil, err
	}

	connChan := &connTask{
		udpConn: make(chan *net.UDPConn),
	}
	go network.waitInviteAnswer(conn.LocalAddr().(*net.UDPAddr), connChan)

	localHost := conn.LocalAddr().String()
	_, fromPort, _ := net.SplitHostPort(localHost)

	req := &net_pb.NatMsg{
		Typ: nbsnet.NatReversInvite,
		ReverseInvite: &net_pb.ReverseInvite{
			PubIp:    lAddr.PubIp,
			ToPort:   int32(toPort),
			PeerId:   rAddr.NetworkId,
			FromPort: fromPort,
		},
	}

	reqData, _ := proto.Marshal(req)
	if _, err := conn.Write(reqData); err != nil {
		conn.Close()
		return nil, err
	}
	logger.Debug("Step1: notify applier's nat server:->", req)
	conn.Close()

	select {
	case conn := <-connChan.udpConn:
		if conn != nil {
			return conn, nil
		} else {
			return nil, connChan.err
		}
	case <-time.After(HolePunchTimeOut):
		return nil, fmt.Errorf("time out")
	}
}

func (network *nbsNetwork) punchANatHole(lAddr, rAddr *nbsnet.NbsUdpAddr, toPort int) (*net.UDPConn, nbsnet.ConnType, error) {

	priConnTask := &connTask{
		udpConn: make(chan *net.UDPConn),
	}

	go network.directDialInPriNet(lAddr, rAddr, priConnTask, toPort)
	pubConnTask := &connTask{
		udpConn: make(chan *net.UDPConn),
	}
	go network.noticePeerAndWait(lAddr, rAddr, toPort, pubConnTask)

	var pubFail, priFail bool
	for i := 2; i > 0; i-- {
		select {
		case conn := <-priConnTask.udpConn:
			logger.Info("hole punch step1-8 dig direct in private network finished:->", priConnTask.err)
			if conn != nil {
				return conn, nbsnet.CTypeNormal, nil
			} else {
				priFail = true
				if pubFail {
					return nil, 0, priConnTask.err
				}
			}
		case conn := <-pubConnTask.udpConn:
			logger.Info("hole punch step2-x dig in public network finished:->", pubConnTask.err)
			if conn != nil {
				return conn, nbsnet.CTypeNatDuplex, nil
			} else {
				pubFail = true
				if priFail {
					return nil, 0, pubConnTask.err
				}
			}
		case <-time.After(HolePunchTimeOut / 2):
			return nil, 0, fmt.Errorf("time out")
		}
	}
	return nil, 0, fmt.Errorf("time out")
}

func (network *nbsNetwork) waitInviteAnswer(host *net.UDPAddr, task *connTask) {

	lisConn, err := net.ListenUDP("udp4", host)
	if err != nil {
		task.finish(err, nil)
		return
	}

	logger.Debug("Step2: wait the answer:", host)

	buffer := make([]byte, utils.NormalReadBuffer)
	n, peerAddr, err := lisConn.ReadFromUDP(buffer)
	if err != nil {
		logger.Warning("read the invite ack err:->", err)
		lisConn.Close()
		task.finish(err, nil)
		return
	}

	res := &net_pb.NatMsg{}
	proto.Unmarshal(buffer[:n], res)
	lisConn.Close()

	logger.Debug("step5-1: get data:->", res)

	conn, err := net.DialUDP("udp4", host, peerAddr)
	if err != nil {
		task.finish(err, nil)
		return
	}
	task.finish(err, conn)

	logger.Debug("Step5: get answer and make a connection:->", host, peerAddr.String())
}

func (network *nbsNetwork) directDialInPriNet(lAddr, rAddr *nbsnet.NbsUdpAddr, task *connTask, toPort int) {

	priHost := nbsnet.JoinHostPort(rAddr.PriIp, int32(utils.GetConfig().NatPrivatePingPort))
	pingConn, err := net.DialTimeout("tcp4", priHost, HolePunchTimeOut/2)
	defer pingConn.Close()
	if err != nil {
		logger.Info("Step 1-1:can't dial by private network.", err)
		task.finish(err, nil)
		return
	}
	logger.Debug("create connection in private network success:->", nbsnet.ConnString(pingConn))

	conn, err := net.DialUDP("udp4", &net.UDPAddr{
		IP: net.ParseIP(lAddr.PriIp),
	}, &net.UDPAddr{
		IP:   net.ParseIP(rAddr.PriIp),
		Port: toPort,
	})

	task.finish(nil, conn)
}

func (network *nbsNetwork) answerInvite(params interface{}) error {
	invite, ok := params.(*net_pb.ReverseInvite)
	if !ok {
		return CmdTaskErr
	}
	myPort := strconv.Itoa(int(invite.ToPort))

	conn, err := shareport.DialUDP("udp4", "0.0.0.0:"+myPort,
		invite.PubIp+":"+invite.FromPort)
	if err != nil {
		logger.Errorf("failed to dial up peer to answer inviter:", err)
		return err
	}
	defer conn.Close()

	req := &net_pb.NatMsg{
		Typ:   nbsnet.NatReversInviteAck,
		NetID: network.networkId,
	}

	reqData, _ := proto.Marshal(req)
	logger.Debug("Step4: answer the invite:->", conn.LocalAddr().String())

	for i := 0; i < DigTryTimesOnNat; i++ {
		if _, err := conn.Write(reqData); err != nil {
			logger.Errorf("failed to write answer to inviter:", err)
			return err
		}
		logger.Debug(" reverse invite dig dig on peer's nat server:->", nbsnet.ConnString(conn))
	}
	return nil
}
