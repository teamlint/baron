package service

import (
	"context"
	"log"

	pb "github.com/teamlint/baron/_example/api/start"
)

func NewService() pb.StartServer {
	return &startService{}
}

type startService struct {
	pb.UnimplementedStartServer
}

func (s *startService) Status(ctx context.Context, in *pb.StatusRequest) (*pb.StatusResponse, error) {
	var resp pb.StatusResponse
	if in.Full {
		resp.Status = pb.ServiceStatus_OK
	} else {

		resp.Status = pb.ServiceStatus_FAIL
	}
	if in.Msg != nil {
		log.Println(*in.Msg)
	}
	return &resp, nil
}
