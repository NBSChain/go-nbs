package rpcService

import (
	"context"
	"fmt"
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
func (service *gossipService) Debug(ctx context.Context, request *pb.DebugCmd) (*pb.DebugResult, error) {
	switch request.Cmd {
	case "showIV":
		return service.showViews(1)
	case "showOV":
		return service.showViews(2)
	case "showAV":
		return service.showViews(3)
	case "clearIV":
		return service.clearViews(1)
	case "clearOV":
		return service.clearViews(2)
	case "clearAV":
		return service.clearViews(3)
	default:
		return nil, fmt.Errorf("unkown command")
	}
}

func (service *gossipService) showViews(typ int) (msg *pb.DebugResult, err error) {

	if !gossip.GetGossipInstance().IsOnline() {
		return nil, gossip.ServiceNotValid
	}
	var ins, outs []string
	if typ == 1 || typ == 3 {
		ins, err = gossip.GetGossipInstance().ShowInputViews()
		if err != nil {
			return nil, err
		}
	}
	if typ == 2 || typ == 3 {
		outs, err = gossip.GetGossipInstance().ShowOutputViews()
		if err != nil {
			return nil, err
		}
	}

	return &pb.DebugResult{
		InputViews: &pb.ViewInfos{
			Views: ins,
		},
		OutputViews: &pb.ViewInfos{
			Views: outs,
		},
	}, nil
}

func (service *gossipService) clearViews(typ int) (msg *pb.DebugResult, err error) {
	var noI, noO int
	if typ == 1 || typ == 3 {
		noI = gossip.GetGossipInstance().ClearInputViews()
	}
	if typ == 2 || typ == 3 {
		noO = gossip.GetGossipInstance().ClearOutputViews()
	}
	return &pb.DebugResult{
		Result: fmt.Sprintf("clear success totaly input=%d, output=%d:->", noI, noO),
	}, nil
}
