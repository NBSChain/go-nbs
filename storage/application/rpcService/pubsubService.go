package rpcService

import (
	"context"
	"github.com/NBSChain/go-nbs/console/pb"
	"github.com/NBSChain/go-nbs/thirdParty/gossip"
)

type pubSubService struct{}

/*****************************************************************
*
*		service callback function.
*
*****************************************************************/

func (service *pubSubService) Publish(ctx context.Context, request *pb.PublishRequest) (*pb.PublishResponse, error) {

	if err := gossip.GetGossipInstance().Publish(request.Topics, []byte(request.Message)); err != nil {
		return nil, err
	}

	return &pb.PublishResponse{
		Result: "subscribe success!",
	}, nil
}
func (service *pubSubService) Subscribe(ctx context.Context, request *pb.SubscribeRequest) (*pb.SubscribeResponse, error) {

	if err := gossip.GetGossipInstance().Subscribe(request.Topics); err != nil {
		return nil, err
	}

	return &pb.SubscribeResponse{
		Result: "subscribe success!",
	}, nil
}

func (service *pubSubService) Peers(ctx context.Context, request *pb.PeersRequest) (*pb.PeersResponse, error) {

	inputViews, outputView := gossip.GetGossipInstance().AllPeers(request.Topics, int(request.Depth))

	return &pb.PeersResponse{
		InPeers:  inputViews,
		OutPeers: outputView,
	}, nil
}
func (service *pubSubService) Topics(ctx context.Context, request *pb.TopicsRequest) (*pb.TopicsResponse, error) {

	topics := gossip.GetGossipInstance().AllMyTopics()

	return &pb.TopicsResponse{
		Topics: topics,
	}, nil
}
