package network

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/nat"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/storage/network/shareport"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"net"
	"strconv"
	"time"
)

func (network *nbsNetwork) noticePeerAndWait(lAddr, rAddr *nbsnet.NbsUdpAddr,
	sid string, toPort int, task *connTask) {

	conn, err := shareport.DialUDP("udp4", "", rAddr.NatServer)

	if err != nil {
		logger.Warning("dial peer's nat server err:->", err)
		task.finish(err, nil)
		return
	}

	msg := &net_pb.NatMsg{
		Typ: nbsnet.NatDigApply,
		DigApply: &net_pb.DigApply{
			NatServer:  lAddr.NatServer,
			TargetId:   rAddr.NetworkId,
			TargetPort: int32(toPort),
			SessionId:  sid,
			FromId:     lAddr.NetworkId,
		},
	}

	data, _ := proto.Marshal(msg)
	if _, err := conn.Write(data); err != nil {
		logger.Warning("write application to peer's nat server err:->", err)
		task.finish(err, nil)
		return
	}

	conn.SetDeadline(time.Now().Add(nat.BootStrapTimeOut))
	buffer := make([]byte, utils.NormalReadBuffer)
	n, err := conn.Read(buffer)
	if err != nil {
		logger.Warning("read ack from peer's nat server err:->", err)
		task.finish(err, nil)
		return
	}
	ack := &net_pb.NatMsg{}
	proto.Unmarshal(buffer[:n], ack)
	locAddr := conn.LocalAddr()
	logger.Info("hole punch step2-1 tell peer's nat server to dig out and wait the answer:->", locAddr.String(), ack.DigACK.PublicAddr)
	conn.Close()

	lisConn, err := shareport.ListenUDP("udp4", locAddr.String())
	if err != nil {
		logger.Warning("listen on waiting port err:->", err)
		task.finish(err, nil)
		return
	}

	lisConn.SetDeadline(time.Now().Add(HolePunchTimeOut))
	buffer = make([]byte, utils.NormalReadBuffer)
	_, peerAddr, err := lisConn.ReadFromUDP(buffer)
	if err != nil {
		logger.Warning("read data from waiting port err:->", err)
		task.finish(err, nil)
		return
	}
	ack = &net_pb.NatMsg{}
	proto.Unmarshal(buffer[:n], ack)
	if ack.Typ != nbsnet.NatDigOut {
		logger.Warning("this is not i am waiting:->", ack)
		task.finish(fmt.Errorf("this is not i am waiting"), nil)
		return
	}

	logger.Info("hole punch step2-1-2 peer dig out success:->", peerAddr)
	lisConn.Close()

	realConn, err := net.DialUDP("udp4", locAddr.(*net.UDPAddr), peerAddr)
	if err != nil {
		logger.Warning("the final dial up err:->", err)
		task.finish(err, nil)
		return
	}

	task.finish(nil, realConn)
}

/************************************************************************
*
*			server side
*
*************************************************************************/
//TIPS:: the server forward the connection invite to peer

func (network *nbsNetwork) digOut(params interface{}) error {
	req, ok := params.(*net_pb.DigApply)
	if !ok {
		return CmdTaskErr
	}

	lPort := strconv.Itoa(int(req.TargetPort))
	conn, err := shareport.DialUDP("udp4", "0.0.0.0:"+lPort, req.Public)
	if err != nil {
		return err
	}
	defer conn.Close()

	msg := &net_pb.NatMsg{
		Typ: nbsnet.NatDigOut,
		Seq: time.Now().Unix(),
	}
	data, _ := proto.Marshal(msg)
	for i := 0; i < DigTryTimesOnNat; i++ {
		if _, err := conn.Write(data); err != nil {
			logger.Error(err)
			return err
		}
		time.Sleep(time.Millisecond * 300)
		logger.Info("hole punch step2-4  dig dig on peer's nat server:->", nbsnet.ConnString(conn))
	}

	return nil
}
