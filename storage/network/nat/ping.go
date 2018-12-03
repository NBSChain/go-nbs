package nat

import (
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"net"
	"time"
)

func (nat *Manager) pong(ping *net_pb.PingPong, peerAddr *net.UDPAddr) error {

	nat.canServe <- true

	logger.Debug("I can serve as in public network.")
	return nil
}

func (nat *Manager) ping(peerAddr *net.UDPAddr, reqNodeId string) {

	conn, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   peerAddr.IP,
		Port: utils.GetConfig().NatServerPort,
	})

	if err != nil {
		logger.Warning("ping DialUDP err:->", err)
		return
	}
	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(BootStrapTimeOut / 2)); err != nil {
		logger.Warning("ping SetDeadline err:->", err)
		return
	}

	request := &net_pb.NatMsg{
		Typ: nbsnet.NatPingPong,
		PingPong: &net_pb.PingPong{
			Ping:  nat.networkId,
			Pong:  reqNodeId,
			Nonce: "", //TODO::security nonce
			TTL:   1,  //time to live
		},
	}
	reqData, _ := proto.Marshal(request)
	if _, err := conn.Write(reqData); err != nil {
		logger.Warning("ping Write err:->", err)
		return
	}
}
