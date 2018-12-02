package nat

import (
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/storage/network/shareport"
	"github.com/golang/protobuf/proto"
	"net"
	"strconv"
	"time"
)

func (tunnel *KATunnel) DigHoeInPubNet(lAddr, rAddr *nbsnet.NbsUdpAddr,
	sid string, toPort int, task *ConnTask) {

	conn, err := shareport.DialUDP("udp4", "", rAddr.NatServer)
	if err != nil {
		task.err <- err
		return
	}

	connReq := &net_pb.NatConnect{
		NatServer:  lAddr.NatServer,
		TargetId:   rAddr.NetworkId,
		TargetPort: int32(toPort),
		SessionId:  sid,
		FromId:     lAddr.NetworkId,
	}
	data, _, _ := nbsnet.PackNatData(connReq, nbsnet.NatConnect)
	if _, err := conn.Write(data); err != nil {
		task.err <- err
		return
	}

	_, port, _ := net.SplitHostPort(conn.LocalAddr().String())
	task.locPort = port
	task.portCapConn = conn
	tunnel.digTask[sid] = task

	logger.Debug("hole punch step2-1 tell peer's nat server to dig out:->", conn.LocalAddr().String())

	return
}

/************************************************************************
*
*			server side
*
*************************************************************************/
//TIPS:: the server forward the connection invite to peer
func (nat *Manager) forwardDigOutReq(data []byte, peerAddr *net.UDPAddr) error {
	req := &net_pb.NatConnect{}
	if err := proto.Unmarshal(data, req); err != nil {
		return err
	}

	nat.cacheLock.Lock()
	defer nat.cacheLock.Unlock()

	toItem, ok := nat.cache[req.TargetId]
	if !ok {
		return NotFundErr
	}

	req.Public = peerAddr.String()
	forwardData, _ := proto.Marshal(req)
	if _, err := nat.sysNatServer.WriteToUDP(forwardData, toItem.pubAddr); err != nil {
		return err
	}
	logger.Debug("hole punch step2-2 forward dig out request:->", req.Public)

	return nil
}

func (tunnel *KATunnel) digOut(data []byte) {

	req := &net_pb.NatConnect{}
	if err := proto.Unmarshal(data, req); err != nil {
		logger.Warning("dig out unmarshal err :->", err)
		return
	}

	digTask := &ConnTask{
		err: make(chan error),
	}
	defer digTask.Close()

	go tunnel.notifyCaller(req, digTask)

	lPort := strconv.Itoa(int(req.TargetPort))
	conn, err := shareport.DialUDP("udp4", "0.0.0.0:"+lPort, req.Public)
	if err != nil {
		return
	}
	defer conn.Close()

	go tunnel.waitDigOutRes(digTask, conn)

	go tunnel.digDig(conn, digTask)

	select {
	case err := <-digTask.err:
		logger.Info("dig out finished:->", err)
	case <-time.After(HolePunchTimeOut):
		logger.Warning("dig out time out")
	}
}

func (tunnel *KATunnel) makeAHole(data []byte) {

	ack := &net_pb.NatConnectAck{}
	if err := proto.Unmarshal(data, ack); err != nil {
		logger.Error("failed to punch a hole:->", err)
		return
	}

	sid := ack.SessionId

	task, ok := tunnel.digTask[sid]
	if !ok {
		logger.Error("can't find the dig task")
		return
	}
	defer delete(tunnel.digTask, sid)

	conn, err := shareport.DialUDP("udp4", "0.0.0.0:"+task.locPort, ack.Public)
	if err != nil {
		task.err <- err
		return
	}

	task.udpConn = conn
	task.err <- nil

	logger.Debug("hole punch step2-7 create hole channel:->",
		conn.LocalAddr().String(), ack.Public)
}
