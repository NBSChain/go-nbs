package rpcService

import (
	"github.com/NBSChain/go-nbs/storage/application/rpcServiceImpl"
	"github.com/NBSChain/go-nbs/storage/merkledag/cid"
	"github.com/NBSChain/go-nbs/utils/cmdKits/pb"
	"io"
)

type getService struct{}

func (service *getService) Get(request *pb.GetRequest, stream pb.GetTask_GetServer) error {

	dataHash := request.DataUri

	//TODO:: didn't consider the naming system request right now.
	cidKey, err := cid.IsValidPath(dataHash)
	if err != nil{
		return err
	}

	resolver, err := rpcServiceImpl.ReadStreamData(cidKey)
	if err != nil{
		return nil
	}

	for {
		data, err := resolver.Next()

		if err != nil{
			if err != io.EOF{
				return err
			}else{
				break
			}
		}

		response := &pb.GetResponse{
			Content:data,
		}

		stream.Send(response)
	}

	return nil
}