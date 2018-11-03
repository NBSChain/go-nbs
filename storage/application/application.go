package application

type Application interface {
	Start() error

	GetNodeId() string

	ReloadForNewAccount() error
}
