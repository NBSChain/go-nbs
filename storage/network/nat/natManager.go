package nat

import "github.com/NBSChain/go-nbs/storage/network/pb"

type Manager interface {
	FindWhoAmI() error

	GetStatus() string

	NatType() nat_pb.NatType
}
