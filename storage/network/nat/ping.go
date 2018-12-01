package nat

import (
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"net"
	"time"
)

func (nat *Manager) pong(data []byte, peerAddr *net.UDPAddr) error {

	ping := &net_pb.NatPing{}
	if err := proto.Unmarshal(data, ping); err != nil {
		return err
	}

	if ping.TTL <= 0 {
		nat.canServe <- true
		return nil
	}

	pong := &net_pb.NatPing{
		Ping:  ping.Ping,
		Pong:  nat.networkId,
		TTL:   ping.TTL - 1,
		Nonce: "", //TODO::
	}

	pongData, err := proto.Marshal(pong)
	if err != nil {
		logger.Warning("failed to marshal pong data", err)
		return err
	}

	if _, err := nat.sysNatServer.WriteToUDP(pongData, peerAddr); err != nil {
		logger.Warning("failed to send pong", err)
		return err
	}

	return nil
}

func (nat *Manager) ping(peerAddr *net.UDPAddr) {

	conn, ping, err := nat.createPingConn(peerAddr)
	defer conn.Close()
	if err != nil {
		logger.Warning("create ping message failed:", err)
		return
	}

	if err := nat.sendPing(ping, conn); err != nil {
		logger.Warning("send ping message failed:", err)
		return
	}

	pong, err := nat.readPong(conn)
	if err != nil {
		logger.Warning("read pong message failed:", err)
		return
	}

	pong.TTL = pong.TTL - 1

	pData, _ := proto.Marshal(pong)
	if err := nat.sendPing(pData, conn); err != nil {
		logger.Warning("send ping again message failed:", err)
		return
	}

}

func (nat *Manager) readPong(conn *net.UDPConn) (*net_pb.NatPing, error) {

	responseData := make([]byte, utils.NormalReadBuffer)
	hasRead, err := conn.Read(responseData)
	if err != nil {
		logger.Warning("get pong failed", err)
		return nil, err
	}

	pong := &net_pb.NatPing{}
	if err := proto.Unmarshal(responseData[:hasRead], pong); err != nil {
		logger.Warning("Unmarshal pong failed", err)
		return nil, err
	}

	logger.Debug("get pong", pong)

	return pong, nil
}

func (nat *Manager) createPingConn(peerAddr *net.UDPAddr) (*net.UDPConn, []byte, error) {

	conn, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   peerAddr.IP,
		Port: utils.GetConfig().NatServerPort,
	})

	if err != nil {
		return nil, nil, err
	}

	conn.SetDeadline(time.Now().Add(BootStrapTimeOut / 2))

	ping := &net_pb.NatPing{
		Ping:  nat.networkId,
		Nonce: "", //TODO::security nonce
		TTL:   2,  //time to live
	}
	data, _ := proto.Marshal(ping)
	return conn, data, nil
}

func (nat *Manager) sendPing(data []byte, conn *net.UDPConn) error {

	request := &net_pb.NatMsg{
		T: nbsnet.NatPingPong,
		L: int32(len(data)),
		V: data,
	}

	reqData, _ := proto.Marshal(request)
	if _, err := conn.Write(reqData); err != nil {
		return err
	}

	return nil
}
