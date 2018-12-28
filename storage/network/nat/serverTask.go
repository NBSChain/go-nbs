package nat

import (
	"github.com/NBSChain/go-nbs/storage/network/nbsnet"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"net"
	"time"
)

func (nat *Server) checkWhoIsHe(conn *net.TCPConn, request *net_pb.BootReg, nodeId string) error {

	peerAddr := conn.LocalAddr().(*net.TCPAddr)
	response := &net_pb.BootAnswer{}
	response.PublicIp = peerAddr.IP.String()
	response.PublicPort = int32(peerAddr.Port)

	if peerAddr.IP.Equal(net.ParseIP(request.PrivateIp)) {
		response.NatType = net_pb.NatType_NoNatDevice
	} else if peerAddr.Port == int(request.PrivatePort) {
		if nat.networkId == nodeId {
			logger.Info("yeah, I'm the bootstrap node:->")
			response.NatType = net_pb.NatType_CanBeNatServer
		} else {
			response.NatType = net_pb.NatType_ToBeChecked
			go nat.ping(peerAddr, nodeId)
		}
	} else {
		response.NatType = net_pb.NatType_BehindNat
	}

	pbRes := &net_pb.NatMsg{
		Typ:        nbsnet.NatBootAnswer,
		BootAnswer: response,
	}

	resData, _ := proto.Marshal(pbRes)
	if _, err := conn.Write(resData); err != nil {
		return err
	}

	if response.NatType == net_pb.NatType_BehindNat ||
		response.NatType == net_pb.NatType_ToBeChecked {

		priAddr := nbsnet.JoinHostPort(request.PrivateIp, request.PrivatePort)
		item := &HostBehindNat{
			UpdateTIme: time.Now(),
			Conn:       conn,
			PriAddr:    priAddr,
		}

		nat.cacheLock.Lock()
		defer nat.cacheLock.Unlock()
		nat.cache[nodeId] = item
	}
	return nil
}

func (nat *Server) updateKATime(req *net_pb.KeepAlive, conn *net.TCPConn, nodeId string) error {

	nat.cacheLock.Lock()
	defer nat.cacheLock.Unlock()

	if item, ok := nat.cache[nodeId]; ok {
		item.UpdateTIme = time.Now()
		item.Conn = conn
		item.PriAddr = req.PriAddr
	} else {
		item := &HostBehindNat{
			UpdateTIme: time.Now(),
			Conn:       conn,
			PriAddr:    req.PriAddr,
		}
		nat.cache[nodeId] = item
	}

	res := &net_pb.NatMsg{
		Typ:   nbsnet.NatKeepAlive,
		NetID: nodeId,
		KeepAlive: &net_pb.KeepAlive{
			PubAddr: conn.RemoteAddr().String(),
		},
	}

	resData, _ := proto.Marshal(res)
	if _, err := conn.Write(resData); err != nil {
		delete(nat.cache, nodeId)
		return err
	}

	return nil
}

func (nat *Server) ping(peerAddr *net.TCPAddr, reqNodeId string) {

	peerServer := &net.TCPAddr{
		IP:   peerAddr.IP,
		Port: utils.GetConfig().NatServerPort,
	}
	conn, err := net.DialTimeout("udp4", peerServer.String(), BootStrapTimeOut)

	if err != nil {
		logger.Warning("ping DialUDP err:->", err)
		return
	}

	defer conn.Close()

	request := &net_pb.NatMsg{
		Typ: nbsnet.NatPingPong,
		PingPong: &net_pb.PingPong{
			Ping:  nat.networkId,
			Pong:  reqNodeId,
			Nonce: "", //TODO::security nonce
			TTL:   1,
		},
	}

	reqData, _ := proto.Marshal(request)
	if _, err := conn.Write(reqData); err != nil {
		logger.Warning("ping Write err:->", err)
		return
	}
}
