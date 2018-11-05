package nat

type Manager interface {
	FetchNatInfo() error
}
