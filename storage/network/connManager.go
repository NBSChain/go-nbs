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
	close(task.udpConn)
	task.udpConn = nil
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

	lisConn, err := net.ListenUDP("udp4", nil)
	if err != nil {
		return nil, err
	}
	defer lisConn.Close()

	natAddr := &net.UDPAddr{
		IP:   net.ParseIP(rAddr.NatServerIP),
		Port: utils.GetConfig().HolePuncherPort,
	}

	localHost := lisConn.LocalAddr().(*net.UDPAddr)
	_, fromPort, _ := net.SplitHostPort(localHost.String())

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

	if _, err := lisConn.WriteToUDP(reqData, natAddr); err != nil {
		return nil, err
	}
	logger.Debug("Step1: notify applier's nat server and wait:->", req)

	lisConn.SetReadDeadline(time.Now().Add(HolePunchTimeOut))
	buffer := make([]byte, utils.NormalReadBuffer)

	_, peerAddr, err := lisConn.ReadFromUDP(buffer)
	if err != nil {
		logger.Warning("read the invite ack err:->", err)
		return nil, err
	}
	lisConn.Close()

	conn, err := net.DialUDP("udp4", localHost, peerAddr)
	if err != nil {
		return nil, err
	}

	logger.Debug("Step5: get answer and make a connection:->", localHost.String(), peerAddr.String())
	return conn, nil
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

func (network *nbsNetwork) directDialInPriNet(lAddr, rAddr *nbsnet.NbsUdpAddr, task *connTask, toPort int) {

	priHost := nbsnet.JoinHostPort(rAddr.PriIp, int32(utils.GetConfig().NatPrivatePingPort))
	pingConn, err := net.DialTimeout("tcp4", priHost, HolePunchTimeOut/2)
	if err != nil {
		logger.Info("Step 1-1:can't dial by private network.", err)
		task.finish(err, nil)
		return
	}
	defer pingConn.Close()

	logger.Debug("1-2 success in private network:->", nbsnet.ConnString(pingConn))

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
	logger.Debug("Step4: answer the invite:->", nbsnet.ConnString(conn))

	for i := 0; i < DigTryTimesOnNat; i++ {
		if _, err := conn.Write(reqData); err != nil {
			logger.Errorf("failed to write answer to inviter:", err)
			return err
		}
		logger.Debug(" reverse invite dig dig on peer's nat server:->", nbsnet.ConnString(conn))
	}
	return nil
}
