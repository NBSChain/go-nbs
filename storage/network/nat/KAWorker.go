//+build !windows

package nat

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/storage/network/shareport"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"net"
	"strconv"
	"time"
)

func (tunnel *KATunnel) readRegResponse() error {

	responseData := make([]byte, NetIoBufferSize)
	hasRead, err := tunnel.kaConn.Read(responseData)
	if err != nil {
		logger.Warning("reading failed from nat server", err)
		return err
	}

	response := &net_pb.Response{}
	if err := proto.Unmarshal(responseData[:hasRead], response); err != nil {
		logger.Warning("unmarshal Err:", err)
		return err
	}

	logger.Info("response:", response)

	resValue := response.BootRegRes

	tunnel.publicIp = resValue.PublicIp
	tunnel.publicPort = resValue.PublicPort

	go tunnel.runLoop()
	go tunnel.listening()
	go tunnel.readKeepAlive()
	return nil
}

func (tunnel *KATunnel) readKeepAlive() {

	for {
		buffer := make([]byte, utils.NormalReadBuffer)

		n, err := tunnel.kaConn.Read(buffer)
		if err != nil {
			logger.Warning("reading keep alive message failed:", err)
			continue
		}

		if err := tunnel.process(buffer[:n]); err != nil {
			logger.Warning("process nat response message failed")
			continue
		}
	}
}

func (tunnel *KATunnel) process(buffer []byte) error {

	response := &net_pb.Response{}
	if err := proto.Unmarshal(buffer, response); err != nil {
		logger.Warning("keep alive response Unmarshal failed:", err)
		return err
	}

	switch response.MsgType {
	case net_pb.NatMsgType_KeepAlive:
		tunnel.updateTime = time.Now()

	case net_pb.NatMsgType_Connect:
		tunnel.punchAHole(response.ConnRes)
	}

	return nil
}

func (tunnel *KATunnel) listening() {

	for {
		buffer := make([]byte, utils.NormalReadBuffer)
		n, err := tunnel.receiveHub.Read(buffer)
		if err != nil {
			logger.Warning("receiving port:", err)
			continue
		}

		if err := tunnel.process(buffer[:n]); err != nil {
			logger.Warning("process nat response message failed")
			continue
		}
	}
}

//TODO::
func (tunnel *KATunnel) restoreNatChannel() {

}

//TIPS::get peer's addr info and make a connection.
func (tunnel *KATunnel) punchAHole(response *net_pb.NatConRes) {
	sessionId := response.SessionId
	task, ok := tunnel.natTask[sessionId]
	if !ok {
		logger.Error("can't find the nat connection task")
		return
	}

	task.ProxyAddr = &net.UDPAddr{
		IP:   net.ParseIP(tunnel.privateIP),
		Port: int(response.TargetPort),
	}
	if response.IsCaller {
		task.CType = ConnTypeNat
	} else {
		task.CType = ConnTypeNatInverse
	}

	if response.CanServe {

		port, _ := strconv.Atoi(response.PublicPort)
		conn, err := net.DialUDP("udp4", nil, &net.UDPAddr{
			IP:   net.ParseIP(response.PublicIp),
			Port: port,
		})
		if err != nil {
			logger.Warning("failed to make a nat connection while peer can be a server.", err)
			task.Err = err
			task.ConnCh <- nil
			return
		}
		task.ConnCh <- conn

	} else {

		priConn, priErr := shareport.DialUDP("udp4", tunnel.privateIP+":"+tunnel.privatePort,
			response.PrivateIp+":"+response.PrivatePort)

		if priErr != nil {
			logger.Warning("failed to make a nat connection from private network while peer is behind nat")
		} else {
			task.ConnCh <- priConn
			return
		}

		pubConn, pubErr := shareport.DialUDP("udp4", tunnel.privateIP+":"+tunnel.privatePort,
			response.PublicIp+":"+response.PublicPort)
		if pubErr != nil {
			logger.Error("failed to make a nat connection from public network while peer is behind nat")
			err := fmt.Errorf("failed to make a nat connection while peer is behind nat.%s-%s", pubErr, priErr)
			task.Err = err
			task.ConnCh <- nil
			return
		}

		task.ConnCh <- pubConn
	}

	delete(tunnel.natTask, task.sessionId)

}
