package domain

// EchoService domain service
type EchoService struct {
	Repo EchoRepository
}

func NewEchoService(repo EchoRepository) *EchoService {
	return &EchoService{Repo: repo}
}

func (s *EchoService) Get(in string) (*Echo, error) {
	return s.Repo.Get(in)
}
