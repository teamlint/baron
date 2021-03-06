package service

import (
	"context"
	"fmt"
	"log"

	"github.com/goava/di"
	pb "github.com/teamlint/baron/_example/api/echo"
	"github.com/teamlint/baron/_example/echo-service/domain"
)

func NewService() *echoService {
	return &echoService{}
}

type echoService struct {
	di.Inject
	pb.UnimplementedEchoServer
	EchoService *domain.EchoService
}

func (s *echoService) Echo(ctx context.Context, in *pb.EchoRequest) (*pb.EchoResponse, error) {
	var resp pb.EchoResponse
	model, _ := s.EchoService.Get(in.In)
	result := model.Msg
	log.Printf("[Echo] EchoRequest=%+v\n", *in)
	if in.At != nil {
		result += "|" + fmt.Sprint(*in.At)
		log.Printf("[Echo] EchoRequest.Optional.Int64=%+v\n", *in.At)
	}
	if in.Desc != nil {
		result += "|" + fmt.Sprint(*in.Desc)
		log.Printf("[Echo] EchoRequest.Optional.String=%+v\n", *in.Desc)
	}
	if in.Debug != nil {
		result += "|" + fmt.Sprint(*in.Debug)
		log.Printf("[Echo] EchoRequest.Optional.Bool=%+v\n", *in.Debug)
	}
	resp.CreatedAt = in.CreatedAt
	resp.Out = result
	return &resp, nil
}

func (s *echoService) Louder(ctx context.Context, in *pb.LouderRequest) (*pb.EchoResponse, error) {
	var resp pb.EchoResponse
	result := fmt.Sprint(in.In) + "|" + fmt.Sprint(in.Loudness)
	resp.Out = result
	log.Printf("[Louder] result = %+v\n", resp.Out)
	return &resp, nil
}

func (s *echoService) LouderGet(ctx context.Context, in *pb.LouderRequest) (*pb.EchoResponse, error) {
	var resp pb.EchoResponse
	result := fmt.Sprint(in.In) + "|" + fmt.Sprint(in.Loudness)
	resp.Out = result
	log.Printf("[LouderGet] result = %+v\n", resp.Out)
	return &resp, nil
}
