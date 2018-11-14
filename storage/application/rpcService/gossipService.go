package rpcService

import (
	"context"
	"github.com/NBSChain/go-nbs/console/pb"
	"github.com/NBSChain/go-nbs/thirdParty/account"
)

type gossipService struct{}

func (service *gossipService) StartService(ctx context.Context, request *pb.StartRequest) (*pb.StartResponse, error) {

	peerId := account.GetAccountInstance().GetPeerID()
	if peerId == "" {
		return nil, account.ENoAvailableAccount
	}

	if err := gossipInst.StartUp(peerId); err != nil {
		return nil, err
	}

	return &pb.StartResponse{Result: "gossip service start success."}, nil
}
