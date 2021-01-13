package adapter

import (
	"fmt"
	"time"

	"github.com/teamlint/baron/_example/echo-service/domain"
)

type EchoRepository struct{}

func NewEchoRepository() *EchoRepository {
	return &EchoRepository{}
}

func (r *EchoRepository) Get(in string) (*domain.Echo, error) {
	msg := fmt.Sprintf("%s@%v", in, time.Now().Unix())
	echo := domain.Echo{
		Msg: msg,
	}
	return &echo, nil
}
