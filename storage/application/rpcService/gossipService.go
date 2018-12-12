package rpcService

import (
	"context"
	"github.com/NBSChain/go-nbs/console/pb"
	"github.com/NBSChain/go-nbs/thirdParty/account"
	"github.com/NBSChain/go-nbs/thirdParty/gossip"
)

type gossipService struct{}

func (service *gossipService) StartService(ctx context.Context, request *pb.StartRequest) (*pb.StartResponse, error) {

	peerId := account.GetAccountInstance().GetPeerID()
	if peerId == "" {
		return nil, account.ENoAvailableAccount
	}

	if err := gossip.GetGossipInstance().Online(peerId); err != nil {
		return nil, err
	}

	return &pb.StartResponse{Result: "gossip service start success."}, nil
}

func (service *gossipService) StopService(ctx context.Context, request *pb.StopRequest) (*pb.StopResponse, error) {

	if err := gossip.GetGossipInstance().Offline(); err != nil {
		return nil, err
	}

	return &pb.StopResponse{Result: "gossip service stop success."}, nil
}
func (service *gossipService) ShowViews(ctx context.Context, request *pb.ShowGossipView) (*pb.AllNodeInView, error) {

	if !gossip.GetGossipInstance().IsOnline() {
		return nil, gossip.ServiceNotValid
	}

	ins, err := gossip.GetGossipInstance().ShowInputViews()
	if err != nil {
		return nil, err
	}

	outs, err := gossip.GetGossipInstance().ShowOutputViews()
	if err != nil {
		return nil, err
	}

	return &pb.AllNodeInView{
		InputView:  ins,
		OutputView: outs,
	}, nil
}
