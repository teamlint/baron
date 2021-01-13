package domain

// EchoRepository domain repository
type EchoRepository interface {
	Get(in string) (*Echo, error)
}
