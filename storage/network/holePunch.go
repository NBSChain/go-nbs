package network

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/storage/network/shareport"
	"github.com/golang/protobuf/proto"
	"net"
	"strconv"
	"time"
)

func (network *nbsNetwork) noticePeerAndWait(lAddr, rAddr *nbsnet.NbsUdpAddr,
	sid string, toPort int, task *connTask) {

	conn, err := shareport.DialUDP("udp4", "", rAddr.NatServer)
	if err != nil {
		task.err = err
		task.udpConn <- nil
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
		task.err = err
		task.udpConn <- nil
		return
	}

	_, port, _ := net.SplitHostPort(conn.LocalAddr().String())
	task.locPort = port
	task.portCapConn = conn
	network.digTask[sid] = task
	logger.Info("hole punch step2-1 tell peer's nat server to dig out:->", conn.LocalAddr().String())
	return
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
	go network.notifyCaller(req)

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
	if _, err := conn.Write(data); err != nil {
		logger.Error(err)
		return err
	}

	logger.Debug("hole punch step2-4  dig dig:->", nbsnet.ConnString(conn))
	return nil
}

func (network *nbsNetwork) notifyCaller(msg *net_pb.DigApply) {

	lPort := strconv.Itoa(int(msg.TargetPort))
	conn, err := shareport.DialUDP("udp4", "0.0.0.0:"+lPort, msg.NatServer)
	if err != nil {
		logger.Warning("dial err:->", err)
		return
	}
	defer conn.Close()

	confirmMsg := &net_pb.NatMsg{
		Typ: nbsnet.NatDigConfirm,
		DigConfirm: &net_pb.DigConfirm{
			SessionId: msg.SessionId,
			TargetId:  msg.FromId,
		},
	}

	data, err := proto.Marshal(confirmMsg)
	if err != nil {
		logger.Warning("pack data err:->", err)
		return
	}

	if _, err := conn.Write(data); err != nil {
		logger.Warning("write data err:->", err)
		return
	}

	logger.Debug("hole punch step2-3 notify caller:->", conn.LocalAddr().String())
}

func (network *nbsNetwork) makeAHole(params interface{}) error {
	ack, ok := params.(*net_pb.DigConfirm)
	if !ok {
		return CmdTaskErr
	}
	sid := ack.SessionId

	task, ok := network.digTask[sid]
	if !ok {
		return fmt.Errorf("can't find the dig taskQueue")
	}
	defer delete(network.digTask, sid)

	conn, err := shareport.DialUDP("udp4", "0.0.0.0:"+task.locPort, ack.Public)
	if err != nil {
		task.err = err
		task.udpConn <- nil
		return err
	}

	task.udpConn <- conn
	task.err = nil

	logger.Debug("hole punch step2-7 create hole channel:->",
		conn.LocalAddr().String(), ack.Public)
	return nil
}
