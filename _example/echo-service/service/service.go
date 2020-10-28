package service

import (
	"context"
	"fmt"
	"log"

	pb "github.com/teamlint/baron/_example/api/echo"
	"github.com/teamlint/baron/_example/echo-service/global"
)

type DB struct {
	Conn string
}

func init() {
	global.Container.Provide(func() *DB {
		return &DB{Conn: "db 连接字符串"}
	})
}

// NewService returns a naïve, stateless implementation of Service.
func NewService() pb.EchoServer {
	var db *DB
	global.Container.Resolve(&db)
	return &echoService{DB: db}
}

type echoService struct {
	*DB
	pb.UnimplementedEchoServer
}

func (s *echoService) Echo(ctx context.Context, in *pb.EchoRequest) (*pb.EchoResponse, error) {
	var resp pb.EchoResponse
	result := s.DB.Conn + " " + in.In
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
