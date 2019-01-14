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

	if err := gossip.GetGossipInstance().Publish(request.Topics, request.Message); err != nil {
		return nil, err
	}

	return &pb.PublishResponse{
		Result: "subscribe success!",
	}, nil
}
func (service *pubSubService) Subscribe(request *pb.SubscribeRequest, stream pb.PubSubTask_SubscribeServer) error {

	queue, err := gossip.GetGossipInstance().Subscribe(request.Topics)
	if err != nil {
		return err
	}

	defer gossip.GetGossipInstance().Unsubscribe(request.Topics)

	for {
		select {
		case msg := <-queue:
			res := &pb.SubscribeResponse{
				From:    msg.From,
				MsgType: msg.Typ,
				MsgData: msg.PayLoad,
			}

			if err := stream.Send(res); err != nil {
				return err
			}
		}
	}
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
