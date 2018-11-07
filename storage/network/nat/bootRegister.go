package nat

import (
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"net"
	"strconv"
	"time"
)

func (nat *nbsNatManager) connectToNatServer(serverIP string) (*net.UDPConn, error) {

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

func (nat *nbsNatManager) sendNatRequest(connection *net.UDPConn) error {

	localAddr := connection.LocalAddr().String()

	host, port, err := net.SplitHostPort(localAddr)
	nat.privateIP = host
	bootRequest := &nat_pb.BootNatRegReq{
		NodeId:      nat.networkId,
		PrivateIp:   host,
		PrivatePort: port,
	}

	request := &nat_pb.NatRequest{
		MsgType:    nat_pb.NatMsgType_BootStrapReg,
		BootRegReq: bootRequest,
	}

	requestData, err := proto.Marshal(request)
	if err != nil {
		logger.Error("failed to marshal nat request", err)
		return err
	}

	if no, err := connection.Write(requestData); err != nil || no == 0 {
		logger.Error("failed to send nat request to natServer ", err, no)
		return err
	}

	return nil
}

func (nat *nbsNatManager) parseNatResponse(connection *net.UDPConn) (*nat_pb.BootNatRegRes, error) {

	responseData := make([]byte, NetIoBufferSize)
	hasRead, _, err := connection.ReadFromUDP(responseData)
	if err != nil {
		logger.Error("reading failed from nat server", err)
		return nil, err
	}

	response := &nat_pb.BootNatRegRes{}
	if err := proto.Unmarshal(responseData[:hasRead], response); err != nil {
		logger.Error("unmarshal err:", err)
		return nil, err
	}

	logger.Debug("response:", response)

	port, _ := strconv.Atoi(response.PublicPort)

	nat.Lock()
	nat.natType = response.NatType
	nat.Unlock()

	if response.NatType == nat_pb.NatType_BehindNat ||
		response.NatType == nat_pb.NatType_ToBeChecked {
		nat.publicAddress = &net.UDPAddr{
			IP:   net.ParseIP(response.PublicIp),
			Port: port,
			Zone: response.Zone,
		}
	}

	return response, nil
}

//TODO:: set multiple servers to make it stronger.
func (nat *nbsNatManager) FetchNatInfo() error {

	config := utils.GetConfig()

	for _, serverIP := range config.NatServerIP {

		//TIPS:: no need to bind local host and local port right now
		connection, err := nat.connectToNatServer(serverIP)
		if err != nil {
			logger.Error("can't know who am I", err)
			goto CloseConn
		}

		if err := nat.sendNatRequest(connection); err != nil {
			goto CloseConn
		}

		_, err = nat.parseNatResponse(connection)
		if err == nil {
			break
		}

	CloseConn:
		connection.Close()
	}

	return nil
}
