package rpcService

import (
	"context"
	"fmt"
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

	networkId, addr := network.GetInstance().GetNatAddr()

	info := fmt.Sprintf("\n=========================================================================\n"+
		"\tnetworkId:\t%s\n"+
		"\tcanServe:\t%v\n"+
		"\tpubIp:\t%s\n"+
		"\tpriIp:\t%s\n"+
		"=========================================================================",
		networkId,
		addr.CanServe,
		addr.PubIp,
		addr.PriIp)

	return &pb.NatStatusResponse{
		Message: info,
	}, nil
}
