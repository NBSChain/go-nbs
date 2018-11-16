package nat

import (
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"net"
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

func (nat *Manager) sendNatRequest(connection *net.UDPConn) (string, error) {

	localAddr := connection.LocalAddr().String()

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

	if no, err := connection.Write(requestData); err != nil || no == 0 {
		logger.Error("failed to send nat request to natServer ", err, no)
		return "", err
	}

	return host, nil
}

func (nat *Manager) parseNatResponse(connection *net.UDPConn) (*net_pb.BootNatRegRes, error) {

	responseData := make([]byte, NetIoBufferSize)
	hasRead, _, err := connection.ReadFromUDP(responseData)
	if err != nil {
		logger.Error("reading failed from nat server", err)
		return nil, err
	}

	response := &net_pb.BootNatRegRes{}
	if err := proto.Unmarshal(responseData[:hasRead], response); err != nil {
		logger.Error("unmarshal err:", err)
		return nil, err
	}

	logger.Debug("response:", response)

	return response, nil
}
