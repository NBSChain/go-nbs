package nat

import (
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"net"
	"time"
)

func (nat *nbsNatManager) pong(ping *nat_pb.NatPing, peerAddr *net.UDPAddr) error {

	if ping.TTL == 2 {
		nat.confirmNatType()
		return nil
	}

	pong := &nat_pb.NatPing{
		Ping:  ping.Ping,
		Pong:  nat.networkId,
		TTL:   ping.TTL + 1,
		Nonce: "", //TODO::
	}

	pongData, err := proto.Marshal(pong)
	if err != nil {
		logger.Warning("failed to marshal pong data", err)
		return err
	}

	if _, err := nat.natServer.WriteToUDP(pongData, peerAddr); err != nil {
		logger.Warning("failed to send pong", err)
		return err
	}

	return nil
}

func (nat *nbsNatManager) ping(peerAddr *net.UDPAddr) {

	conn, ping, err := nat.createPingConn(peerAddr)
	if err != nil {
		return
	}

	if err := nat.sendPing(ping, conn); err != nil {
		return
	}

	pong, err := nat.readPong(conn)
	if err != nil {
		return
	}

	pong.TTL = pong.TTL + 1

	nat.sendPing(pong, conn)

	defer conn.Close()
}

func (nat *nbsNatManager) readPong(conn *net.UDPConn) (*nat_pb.NatPing, error) {

	responseData := make([]byte, NetIoBufferSize)
	hasRead, _, err := conn.ReadFromUDP(responseData)
	if err != nil {
		logger.Warning("get pong failed", err)
		return nil, err
	}

	pong := &nat_pb.NatPing{}
	if err := proto.Unmarshal(responseData[:hasRead], pong); err != nil {
		logger.Warning("Unmarshal pong failed", err)
		return nil, err
	}

	logger.Debug("get pong", pong)

	return pong, nil
}

func (nat *nbsNatManager) createPingConn(peerAddr *net.UDPAddr) (*net.UDPConn, *nat_pb.NatPing, error) {

	conn, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   peerAddr.IP,
		Port: utils.GetConfig().NatServerPort,
	})

	if err != nil {
		return nil, nil, err
	}

	conn.SetDeadline(time.Now().Add(time.Second * BootStrapNatServerTimeOutInSec / 2))

	ping := &nat_pb.NatPing{
		Ping:  nat.networkId,
		Nonce: "", //TODO::security nonce
		TTL:   0,
	}

	return conn, ping, nil
}

func (nat *nbsNatManager) sendPing(ping *nat_pb.NatPing, conn *net.UDPConn) error {

	request := &nat_pb.NatRequest{
		Ping:    ping,
		MsgType: nat_pb.NatMsgType_Ping,
	}

	data, _ := proto.Marshal(request)

	if _, err := conn.Write(data); err != nil {
		return err
	}

	return nil
}
