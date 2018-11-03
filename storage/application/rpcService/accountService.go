package rpcService

import (
	"fmt"
	"github.com/NBSChain/go-nbs/console/pb"
	"github.com/NBSChain/go-nbs/storage/application"
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

	acc := account.GetAccountInstance()
	if acc.GetPeerID() != "" {
		return nil, fmt.Errorf("can't create another account, we support only one account right now")
	}

	if err := acc.CreateAccount(request.Password); err != nil {
		return nil, err
	}

	application.GetInstance().ReloadForNewAccount()

	return &pb.CreateAccountResponse{
		Message: "Create account success!",
	}, nil
}
