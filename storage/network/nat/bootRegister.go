package nat

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"net"
	"strconv"
	"time"
)

func (nat *Manager) connectToNatServer(serverIP string) (*net.UDPConn, error) {

	config := utils.GetConfig()
	natServerAddr := &net.UDPAddr{
		IP:   net.ParseIP(serverIP),
		Port: config.NatServerPort,
	}

	conn, err := net.DialUDP("udp", nil, natServerAddr)
	if err != nil {
		return nil, err
	}

	conn.SetDeadline(time.Now().Add(BootStrapNatServerTimeOutInSec * time.Second))

	return conn, nil
}

func (nat *Manager) sendNatRequest(conn *net.UDPConn) (string, error) {

	localAddr := conn.LocalAddr().String()

	host, port, err := net.SplitHostPort(localAddr)
	bootRequest := &net_pb.BootNatRegReq{
		NodeId:      nat.networkId,
		PrivateIp:   host,
		PrivatePort: port,
	}

	request := &net_pb.NatRequest{
		MsgType:    net_pb.NatMsgType_BootStrapReg,
		BootRegReq: bootRequest,
	}

	requestData, err := proto.Marshal(request)
	if err != nil {
		logger.Error("failed to marshal nat request", err)
		return "", err
	}

	if no, err := conn.Write(requestData); err != nil || no == 0 {
		logger.Error("failed to send nat request to selfNatServer ", err, no)
		return "", err
	}

	return host, nil
}

func (nat *Manager) parseNatResponse(conn *net.UDPConn) (*net_pb.BootNatRegRes, error) {

	responseData := make([]byte, NetIoBufferSize)
	hasRead, _, err := conn.ReadFromUDP(responseData)
	if err != nil {
		logger.Error("reading failed from nat server", err)
		return nil, err
	}

	response := &net_pb.Response{}
	if err := proto.Unmarshal(responseData[:hasRead], response); err != nil {
		logger.Error("unmarshal Err:", err)
		return nil, err
	}

	logger.Debug("response:", response)

	return response.BootRegRes, nil
}

/************************************************************************
*
*			server side
*
*************************************************************************/

func (nat *Manager) checkWhoIsHe(request *net_pb.BootNatRegReq, peerAddr *net.UDPAddr) error {

	response := &net_pb.BootNatRegRes{}
	response.PublicIp = peerAddr.IP.String()
	response.PublicPort = fmt.Sprintf("%d", peerAddr.Port)
	if peerAddr.IP.Equal(net.ParseIP(request.PrivateIp)) {

		response.NatType = net_pb.NatType_NoNatDevice
	} else if strconv.Itoa(peerAddr.Port) == request.PrivatePort {

		response.NatType = net_pb.NatType_ToBeChecked
		go nat.ping(peerAddr)

	} else {

		response.NatType = net_pb.NatType_BehindNat
	}

	pbRes := &net_pb.Response{
		MsgType:    net_pb.NatMsgType_BootStrapReg,
		BootRegRes: response,
	}

	pbResData, err := proto.Marshal(pbRes)
	if err != nil {
		logger.Warning("failed to marshal nat response data", err)
		return err
	}

	if _, err := nat.selfNatServer.WriteToUDP(pbResData, peerAddr); err != nil {
		logger.Warning("failed to send nat response", err)
		return err
	}

	if !IsPublic(response.NatType) {
		item := &ClientItem{
			nodeId:     request.NodeId,
			pubIp:      response.PublicIp,
			pubPort:    response.PublicPort,
			priIp:      request.PrivateIp,
			priPort:    request.PrivatePort,
			updateTIme: time.Now(),
		}
		nat.cache[request.NodeId] = item
	}

	return nil
}
