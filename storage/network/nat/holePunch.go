package nat

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/denat"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/storage/network/shareport"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"net"
	"strconv"
	"time"
)

//conn inviter call first
func (tunnel *KATunnel) natHoleStep1(task *nbsnet.ConnTask, lAddr, rAddr *nbsnet.NbsUdpAddr, connId string) error {

	connReq := &net_pb.NatConInvite{

		FromAddr: &net_pb.NbsAddr{
			NetworkId: lAddr.NetworkId,
			CanServer: lAddr.CanServe,
			PubIp:     lAddr.PubIp,
			PubPort:   int32(lAddr.PubPort),
			PriIp:     lAddr.PriIp,
			PriPort:   int32(lAddr.PriPort),
		},
		ToAddr: &net_pb.NbsAddr{
			NetworkId: rAddr.NetworkId,
			CanServer: rAddr.CanServe,
			PubIp:     rAddr.PubIp,
			PubPort:   int32(rAddr.PubPort),
			PriIp:     rAddr.PriIp,
			PriPort:   int32(rAddr.PriPort),
		},
		SessionId: connId,
	}

	response := &net_pb.NatRequest{
		MsgType: net_pb.NatMsgType_Connect,
		ConnReq: connReq,
	}

	toItemData, err := proto.Marshal(response)
	if err != nil {
		return err
	}

	if _, err := tunnel.kaConn.Write(toItemData); err != nil {
		return err
	}

	return nil
}

//caller make a direct connection to peer's public address
func (tunnel *KATunnel) natHoleStep2(task *nbsnet.ConnTask, rAddr *nbsnet.NbsUdpAddr) {

	port := strconv.Itoa(rAddr.PubPort)
	conn, err := shareport.DialUDP("udp4", tunnel.sharedAddr, rAddr.PubIp+":"+port)

	if err != nil {
		task.Conn = nil
		task.Err <- err
		return
	}

	if err := tunnel.sendDigData(task.SessionId, conn); err != nil {

		task.Conn = nil
		task.Err <- err
		conn.Close()
		return
	}

	task.Conn = conn
	task.Err <- nil
}

//TIPS::get peer's addr info and make a connection.
func (tunnel *KATunnel) natHoleStep4(response *net_pb.NatConInvite) error {

	port := strconv.Itoa(int(response.FromAddr.PubPort))
	remoteAddr := response.FromAddr.PubIp + ":" + port
	conn, err := shareport.DialUDP("udp4", tunnel.sharedAddr, remoteAddr)
	if err != nil {
		logger.Error("failed to setup hole connection")
		return err
	}
	sessionId := response.SessionId

	if err := tunnel.sendDigData(sessionId, conn); err != nil {
		logger.Error("failed to try to setup nat tunnel")
		conn.Close()
		return err
	}

	proxyAddr := &net.UDPAddr{
		IP:   net.ParseIP(response.ToAddr.PriIp),
		Port: int(response.ToAddr.PriPort),
	}

	item := &proxyConnItem{
		conn:       conn,
		sessionId:  sessionId,
		targetAddr: proxyAddr,
	}

	tunnel.proxyCache[sessionId] = item

	return nil
}

func (tunnel *KATunnel) sendDigData(sessionId string, conn *net.UDPConn) error {

	holeMsg := &net_pb.NatRequest{
		MsgType: net_pb.NatMsgType_DigDig,
		HoleMsg: &net_pb.HoleDig{
			SessionId: sessionId,
		},
	}

	data, _ := proto.Marshal(holeMsg)
	buffer := make([]byte, utils.NormalReadBuffer)

	for i := 0; i < HolePunchingTimeOut; i++ {

		if _, err := conn.Write(data); err != nil {
			logger.Debug("dig hole failed:", err)
		}
		conn.SetReadDeadline(time.Now().Add(time.Second))

		if _, err := conn.Read(buffer); err != nil {
			logger.Debug("read hole msg:", err)
			continue
		}

		return nil
	}
	return fmt.Errorf("time out")
}

/************************************************************************
*
*			server side
*
*************************************************************************/
//TIPS:: the server forward the connection invite to peer
func (nat *Manager) natHoleStep3(request *net_pb.NatRequest, peerAddr *net.UDPAddr) error {

	req := request.ConnReq
	rawData, _ := proto.Marshal(request)

	nat.cacheLock.Lock()
	defer nat.cacheLock.Unlock()

	toItem, ok := nat.cache[req.ToAddr.NetworkId]
	if !ok {
		if err := denat.GetDNSInstance().ProxyConnInvite(req); err != nil {
			return err
		}
	} else {
		if _, err := nat.sysNatServer.WriteToUDP(rawData, toItem.pubAddr); err != nil {
			return err
		}
	}

	return nil
}

/************************************************************************
*
*			proxy
*
*************************************************************************/
func (proxy *proxyConnItem) send(b []byte) (int, error) {
	return proxy.conn.Write(b)
}

func (proxy *proxyConnItem) receive(b []byte) (int, error) {
	return proxy.conn.Read(b)
}

func (proxy *proxyConnItem) close() {
	proxy.conn.Close()
	proxy.isClosed = true
}
