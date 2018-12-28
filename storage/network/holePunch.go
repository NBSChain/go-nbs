package network

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/storage/network/shareport"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/NBSChain/go-nbs/utils/crypto"
	"github.com/golang/protobuf/proto"
	"strconv"
	"time"
)

func (network *nbsNetwork) noticePeerAndWait(lAddr, rAddr *nbsnet.NbsUdpAddr, toPort int, task *connTask) {

	natAddr := network.natClient.NatAddr
	if natAddr.NetType == nbsnet.MultiPort {
		task.finish(NTSportErr, nil)
		return
	}

	natHost := nbsnet.JoinHostPort(rAddr.NatServerIP, int32(utils.GetConfig().HolePuncherPort))
	conn, err := shareport.DialUDP("udp4", "", natHost)

	if err != nil {
		logger.Debug("dial up peer's nat server err:->", err)
		task.finish(err, nil)
		return
	}

	msg := &net_pb.NatMsg{
		Typ: nbsnet.NatDigApply,
		DigApply: &net_pb.DigApply{
			NatServer:  lAddr.NatServerIP,
			TargetId:   rAddr.NetworkId,
			TargetPort: int32(toPort),
			FromId:     lAddr.NetworkId,
			NtType:     natAddr.NetType,
			PubIps:     natAddr.AllPubIps,
		},
	}

	data, _ := proto.Marshal(msg)
	if _, err := conn.Write(data); err != nil {
		logger.Debug("write msg to peer's nat server err:->", err)
		task.finish(err, nil)
		return
	}

	task.lAddr = conn.LocalAddr().String()
	task.portCapConn = conn
	network.digTask[rAddr.NetworkId] = task

	logger.Info("hole punch step2-1 tell peer's nat server to dig out:->", conn.LocalAddr().String())
	return
}

/************************************************************************
*
*			server side
*
*************************************************************************/
//TIPS:: the server forward the connection invite to peer
func (network *nbsNetwork) sendDigData(locPort, peerAddr string) error {
	conn, err := shareport.DialUDP("udp4", "0.0.0.0:"+locPort, peerAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	msg := &net_pb.NatMsg{
		Typ: nbsnet.NatDigOut,
	}
	data, _ := proto.Marshal(msg)
	for i := 0; i < DigTryTimesOnNat; i++ {
		if _, err := conn.Write(data); err != nil {
			logger.Error(err)
			return err
		}

		logger.Info("hole punch step2-4  dig dig on peer's nat server:->", nbsnet.ConnString(conn))
	}
	return nil
}
func (network *nbsNetwork) digOut(params interface{}) error {
	req, ok := params.(*net_pb.DigApply)
	if !ok {
		return CmdTaskErr
	}

	natAddr := network.natClient.NatAddr
	if natAddr.NetType == nbsnet.MultiPort {
		return NTSportErr
	}

	go network.notifyCaller(req)
	lPort := strconv.Itoa(int(req.TargetPort))
	if req.NtType == nbsnet.SigIpSigPort {
		return network.sendDigData(lPort, req.Public)
	}

	_, port, _ := nbsnet.SplitHostPort(req.Public)
	for _, ip := range req.PubIps {
		peerAddr := nbsnet.JoinHostPort(ip, port)
		if err := network.sendDigData(lPort, peerAddr); err != nil {
			logger.Warning("fail to dig out when applied:->", err, peerAddr)
		}
	}

	return nil
}

func (network *nbsNetwork) notifyCaller(msg *net_pb.DigApply) {

	lPort := strconv.Itoa(int(msg.TargetPort))
	peersNatServer := nbsnet.JoinHostPort(msg.NatServer, int32(utils.GetConfig().HolePuncherPort))
	conn, err := shareport.DialUDP("udp4", "0.0.0.0:"+lPort, peersNatServer)
	if err != nil {
		logger.Warning("dial err:->", err)
		return
	}
	defer conn.Close()

	natAddr := network.natClient.NatAddr
	confirmMsg := &net_pb.NatMsg{
		Typ: nbsnet.NatDigConfirm,
		DigConfirm: &net_pb.DigConfirm{
			TargetId: msg.FromId,
			NtType:   natAddr.NetType,
			PubIps:   natAddr.AllPubIps,
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

	logger.Info("hole punch step2-3 notify caller:->", conn.LocalAddr().String())
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
	task.portCapConn.Close()

	conn, err := shareport.DialUDP("udp4", task.lAddr, ack.Public)
	if err != nil {
		logger.Debug("dial up when make a hole:->", err)
		task.finish(err, nil)
		return err
	}

	now := time.Now()
	msg := &net_pb.NatMsg{
		Typ: nbsnet.NatEnd,
		ConnKA: &net_pb.ConnKA{
			KA: crypto.MD5SS(now.String()),
		},
	}
	data, _ := proto.Marshal(msg)

	if _, err := conn.Write(data); err != nil {
		logger.Debug("make a hole write err:->", err)
		task.finish(err, nil)
		return err
	}

	conn.SetDeadline(time.Now().Add(HolePunchTimeOut))

	buffer := make([]byte, utils.NormalReadBuffer)
	if _, err := conn.Read(buffer); err != nil {
		logger.Debug("make a hole read err:->", err)
		task.finish(err, nil)
		return err
	}

	task.finish(nil, conn)
	logger.Info("hole punch step2-7 create hole channel:->",
		conn.LocalAddr().String(), ack.Public)

	return nil
}
