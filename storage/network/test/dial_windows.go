package main

func (peer *NatPeer) dialMultiTarget(ack *net_pb.DigConfirm, port int32, localAddr string) (*net.UDPConn, error) {
	lis, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		return nil, err
	}
	msg := net_pb.NatMsg{
		Typ: nbsnet.NatBlankKA,
		Seq: time.Now().Unix(),
	}
	data, _ := proto.Marshal(&msg)

	for _, ips := range ack.PubIps {
		tarAddr := &net.UDPAddr{
			IP:   net.ParseIP(ips),
			Port: int(port),
		}
		if _, err := lis.WriteToUDP(data, tarAddr); err != nil {
			return nil, err
		}
	}
	buffer := make([]byte, utils.NormalReadBuffer)
	n, p, e := lis.ReadFromUDP(buffer)
	if e != nil {
		return nil, err

	}
	logger.Info("get answer:->", n, p)
	lis.Close()
	conn, err := net.DialUDP("udp", digAddr, p)
	if err != nil {
		return nil, err
	}
	return conn, err
}
