package nat

type Manager interface {
	FindWhoAmI() error

	GetStatus() string
}
