package nat

import (
	"github.com/NBSChain/go-nbs/storage/network/denat"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"net"
)

//TIPS::get peer's addr info and make a connection.
func (tunnel *KATunnel) punchAHole(response *net_pb.NatConInvite) error {

	connType := nbsnet.ConnType(response.CType)
	sessionId := response.SessionId
	if connType == nbsnet.CTypeNatReverseDirect ||
		connType == nbsnet.CTypeNatReverseWithProxy {

		holeConn := &nbsnet.HoleConn{}
		lAddr := response.FromAddr
		rAddr := response.ToAddr
		lfc, err := net.DialUDP("udp4", nil, &net.UDPAddr{
			IP:   net.ParseIP(rAddr.PriIp),
			Port: int(rAddr.PriPort),
		})
		if err != nil {
			return err
		}
		holeConn.LocalForwardConn = lfc

		rdc, err := net.DialUDP("udp4", nil, &net.UDPAddr{
			IP:   net.ParseIP(lAddr.PubIp),
			Port: int(lAddr.PubPort),
		})
		if err != nil {
			return err
		}
		holeConn.RemoteDataConn = rdc
		holeConn.SessionId = sessionId

		req := &net_pb.NatRequest{
			MsgType: net_pb.NatMsgType_holeAck,
			HoleAck: &net_pb.HoleAck{
				SessionId: sessionId,
			},
		}
		reqData, err := proto.Marshal(req)
		if err != nil {
			return err
		}

		if _, err := rdc.Write(reqData); err != nil {
			return err
		}

		tunnel.connManager[sessionId] = holeConn
		go tunnel.holeConnRead(holeConn)
		go tunnel.holeConnWrite(holeConn)
	}

	return nil
}

/************************************************************************
*
*			server side
*
*************************************************************************/
//TODO::Find peers from nat gossip protocol
func (nat *Manager) notifyConnInvite(request *net_pb.NatRequest, peerAddr *net.UDPAddr) error {

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

	if _, err := nat.sysNatServer.WriteTo(rawData, peerAddr); err != nil {
		return err
	}

	return nil
}

func (tunnel *KATunnel) holeConnRead(conn *nbsnet.HoleConn) {

	for {
		buffer := make([]byte, utils.NormalReadBuffer)
		n, err := conn.RemoteDataConn.Read(buffer)
		if err != nil {
			tunnel.closeHole(conn)
			logger.Info("read failed hole connection closed", conn.SessionId, err)
			return
		}

		conn.LocalForwardConn.Write(buffer[:n])
	}
}

func (tunnel *KATunnel) holeConnWrite(conn *nbsnet.HoleConn) {

	for {
		buffer := make([]byte, utils.NormalReadBuffer)
		n, err := conn.LocalForwardConn.Read(buffer)
		if err != nil {
			tunnel.closeHole(conn)
			logger.Info("write failed hole connection closed:", conn.SessionId, err)
			return
		}

		conn.RemoteDataConn.Write(buffer[:n])
	}
}

func (tunnel *KATunnel) closeHole(conn *nbsnet.HoleConn) {
	conn.Close()
	delete(tunnel.connManager, conn.SessionId)
}
