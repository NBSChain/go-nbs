package rpcService

import (
	"github.com/NBSChain/go-nbs/console/pb"
	"github.com/NBSChain/go-nbs/thirdParty/account"
	"golang.org/x/net/context"
)

type accountService struct{}

func (service *accountService) AccountUnlock(ctx context.Context,
	request *pb.AccountUnlockRequest) (*pb.AccountUnlockResponse, error) {

	if err := account.GetAccountInstance().UnlockAccount(request.Password); err != nil {
		return nil, err
	}

	return &pb.AccountUnlockResponse{
		Message: "Unlock account success, It will be expire in 5 minutes",
	}, nil
}

func (service *accountService) CreateAccount(ctx context.Context,
	request *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {

	if err := account.CreateAccount(request.Password); err != nil {
		return nil, err
	}

	return &pb.CreateAccountResponse{
		Message: "Create account success!",
	}, nil
}
