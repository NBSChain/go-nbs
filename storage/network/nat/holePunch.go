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
	defer conn.Close()

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
	tunnel.digTask[sid] = task

	return
}

/************************************************************************
*
*			server side
*
*************************************************************************/
//TIPS:: the server forward the connection invite to peer
func (nat *Manager) forwardDigRequest(data []byte, peerAddr *net.UDPAddr) error {
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
	simpleNat := &net_pb.SimpleNatInfo{
		NatInfo: req.Public,
	}
	resData, _, _ := nbsnet.PackNatData(simpleNat, nbsnet.NatConnect)
	if _, err := nat.sysNatServer.WriteToUDP(resData, peerAddr); err != nil {
		return err
	}

	logger.Info("Step 2:-> forward to peer:->", req)

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
	defer close(digTask.err)

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

func (tunnel *KATunnel) digSuccessRes(data []byte, peerAddr *net.UDPAddr) {

	res := &net_pb.NatMsg{
		Typ:     nbsnet.NatDigSuccess,
		Len:     int32(len(data)),
		PayLoad: data,
	}

	resData, _ := proto.Marshal(res)
	if _, err := tunnel.serverHub.WriteTo(resData, peerAddr); err != nil {
		logger.Warning("failed to response the dig confirm.")
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
}
