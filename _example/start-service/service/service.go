package service

import (
	"context"

	pb "github.com/teamlint/baron/_example/api/start"
)

func NewService() *startService {
	return &startService{}
}

type startService struct {
	pb.UnimplementedStartServer
}

func (s *startService) Status(ctx context.Context, in *pb.StatusRequest) (*pb.StatusResponse, error) {
	var resp pb.StatusResponse
	return &resp, nil
}
