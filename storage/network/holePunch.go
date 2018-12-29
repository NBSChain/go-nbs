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

func (network *nbsNetwork) noticePeerAndWait(lAddr, rAddr *nbsnet.NbsUdpAddr, toPort int, task *connTask) {

	natAddr := network.natClient.NatAddr
	if natAddr.NetType == nbsnet.MultiPort {
		task.finish(NTSportErr, nil)
		return
	}

	natHost := &net.UDPAddr{
		IP:   net.ParseIP(rAddr.NatServerIP),
		Port: utils.GetConfig().HolePuncherPort,
	}
	conn, err := net.ListenUDP("udp4", nil)

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
	if _, err := conn.WriteToUDP(data, natHost); err != nil {
		logger.Debug("write msg to peer's nat server err:->", err)
		task.finish(err, nil)
		return
	}

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
func (network *nbsNetwork) sendDigData(conn *net.UDPConn, peerAddr *net.UDPAddr, data []byte) error {

	for i := 0; i < DigTryTimesOnNat; i++ {
		if _, err := conn.WriteToUDP(data, peerAddr); err != nil {
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

	locPort := strconv.Itoa(int(req.TargetPort))
	conn, err := shareport.ListenUDP("udp4", "0.0.0.0:"+locPort)
	if err != nil {
		logger.Warning("dig out err:->", err)
		return err
	}
	defer conn.Close()

	network.notifyCaller(req, conn)

	msg := &net_pb.NatMsg{
		Typ:   nbsnet.NatDigOut,
		NetID: network.networkId,
	}
	data, _ := proto.Marshal(msg)
	pubIp, port, _ := nbsnet.SplitHostPort(req.Public)

	if req.NtType == nbsnet.SigIpSigPort {
		peerAddr := &net.UDPAddr{
			IP:   net.ParseIP(pubIp),
			Port: int(port),
		}
		return network.sendDigData(conn, peerAddr, data)
	}

	for _, ip := range req.PubIps {
		peerAddr := &net.UDPAddr{
			IP:   net.ParseIP(ip),
			Port: int(port),
		}
		if err := network.sendDigData(conn, peerAddr, data); err != nil {
			logger.Warning("fail to dig out when applied:->", err, peerAddr)
		}
	}

	return nil
}

func (network *nbsNetwork) notifyCaller(msg *net_pb.DigApply, conn *net.UDPConn) {

	serverAddr := &net.UDPAddr{
		IP:   net.ParseIP(msg.NatServer),
		Port: utils.GetConfig().HolePuncherPort,
	}

	natAddr := network.natClient.NatAddr
	confirmMsg := &net_pb.NatMsg{
		Typ:   nbsnet.NatDigConfirm,
		NetID: network.networkId,
		DigConfirm: &net_pb.DigConfirm{
			TargetId: msg.FromId,
			NtType:   natAddr.NetType,
			PubIps:   natAddr.AllPubIps,
		},
	}

	data, _ := proto.Marshal(confirmMsg)
	if _, err := conn.WriteToUDP(data, serverAddr); err != nil {
		logger.Warning("write data err:->", err)
		return
	}

	logger.Info("hole punch step2-3 notify caller:->", conn.LocalAddr().String())
}

func (network *nbsNetwork) makeAHole(params interface{}) error {
	msg, ok := params.(*net_pb.NatMsg)
	if !ok {
		return CmdTaskErr
	}
	ack := msg.DigConfirm
	sid := msg.NetID

	task, ok := network.digTask[sid]
	if !ok {
		return fmt.Errorf("can't find the dig taskQueue")
	}
	defer delete(network.digTask, sid)

	ip, port, _ := nbsnet.SplitHostPort(ack.Public)
	var conn *net.UDPConn
	var err error

	locAddr := task.portCapConn.LocalAddr().(*net.UDPAddr)

	if ack.NtType == nbsnet.SigIpSigPort {
		peerAddr := &net.UDPAddr{
			IP:   net.ParseIP(ip),
			Port: int(port),
		}
		task.portCapConn.Close()

		conn, err = net.DialUDP("udp4", locAddr, peerAddr)
		if err != nil {
			logger.Debug("dial up when make a hole:->", err)
			task.finish(err, nil)
			return err
		}
	} else {
		choseIp := &net_pb.NatMsg{
			Typ:   nbsnet.NatFindPubIpSyn,
			NetID: network.networkId,
		}
		data, _ := proto.Marshal(choseIp)

		go func() {

			for _, ip := range ack.PubIps {
				peerAddr := &net.UDPAddr{
					IP:   net.ParseIP(ip),
					Port: int(port),
				}
				if _, err := task.portCapConn.WriteToUDP(data, peerAddr); err != nil {
					logger.Warning("chose peer's nat pub ip:->", err)
					continue
				}
			}
		}()

		buffer := make([]byte, utils.NormalReadBuffer)
		task.portCapConn.SetReadDeadline(time.Now().Add(HolePunchTimeOut))

		_, peerAddr, err := task.portCapConn.ReadFromUDP(buffer)
		if err != nil {
			logger.Warning("didn't find one from multi ips:->", err)
			task.finish(err, nil)
			return err
		}
		task.portCapConn.Close()

		conn, err = net.DialUDP("udp4", locAddr, peerAddr)
		if err != nil {
			logger.Warning("setup hole connection from one of the ips err:->", err)
			task.finish(err, nil)
			return err
		}
	}

	task.finish(nil, conn)
	logger.Info("hole punch step2-7 create hole channel:->",
		conn.LocalAddr().String(), ack.Public)

	return nil
}
