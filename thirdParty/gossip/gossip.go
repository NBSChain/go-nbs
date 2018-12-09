package gossip

type BasicProtocol interface {
	StartUp(peerId string) error

	Offline() error

	Publish(channel string, message []byte) error

	Subscribe(channel string) error

	Unsubscribe(channel string) error

	AllPeers(channel string, depth int) ([]string, []string)

	AllMyTopics() []string
}
