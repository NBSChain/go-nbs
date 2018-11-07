package idService

type ID string

type Identity struct {
	PeerID     string
	Encrypted  bool
	PrivateKey string
}

func (id ID) Pretty() string {
	return IDB58Encode(id)
}
