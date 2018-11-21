package nat

import (
	"github.com/NBSChain/go-nbs/storage/network/denat"
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
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
	}

	return nil
}

/************************************************************************
*
*			server side
*
*************************************************************************/
//TODO::Find peers from nat gossip protocol
func (nat *Manager) notifyConnInvite(req *net_pb.NatConInvite, peerAddr *net.UDPAddr) error {
	nat.cacheLock.Lock()
	defer nat.cacheLock.Unlock()

	connType := nbsnet.ConnType(req.CType)

	if connType == nbsnet.CTypeNatReverseDirect ||
		connType == nbsnet.CTypeNatReverseWithProxy {

		toItem, ok := nat.cache[req.FromAddr.NetworkId]

		if !ok {
			if err := denat.GetDNSInstance().ProxyConnInvite(req); err != nil {
				return err
			}
		} else {
			response := &net_pb.Response{
				MsgType: net_pb.NatMsgType_Connect,
				ConnRes: req,
			}
			responseData, err := proto.Marshal(response)
			if err != nil {
				return err
			}

			nat.sysNatServer.WriteToUDP(responseData, toItem.pubAddr)
		}
	}

	return nil
}
