package gossip

import "github.com/NBSChain/go-nbs/thirdParty/gossip/message"

type BasicProtocol interface {
	Online(peerId string) error

	Offline() error

	IsOnline() bool

	Publish(channel string, message string) error

	Subscribe(channel string) (chan *message.MsgEntity, error)

	Unsubscribe(channel string)

	AllPeers(channel string, depth int) ([]string, []string)

	AllMyTopics() []string

	ShowInputViews() ([]string, error)

	ShowOutputViews() ([]string, error)

	ClearInputViews() int

	ClearOutputViews() int
}
