package gossip

type BasicProtocol interface {
	Online(peerId string) error

	Offline() error

	IsOnline() bool

	Publish(channel string, message []byte) error

	Subscribe(channel string) error

	Unsubscribe(channel string) error

	AllPeers(channel string, depth int) ([]string, []string)

	AllMyTopics() []string

	ShowInputViews() ([]string, error)

	ShowOutputViews() ([]string, error)
}
