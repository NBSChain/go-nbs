package rpcService

import (
	"context"
	"github.com/NBSChain/go-nbs/console/pb"
	"github.com/NBSChain/go-nbs/storage/network"
)

type natService struct{}

/*****************************************************************
*
*		service callback function.
*
*****************************************************************/
func (service *natService) NatStatus(ctx context.Context,
	request *pb.NatStatusResQuest) (*pb.NatStatusResponse, error) {

	info := network.GetInstance().GetNatInfo()
	return &pb.NatStatusResponse{
		Message: info,
	}, nil
}
