package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/rs/xid"

	// timestamppb "github.com/golang/protobuf/ptypes/timestamp"
	pb "github.com/teamlint/baron/_example/api/echo"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"

	// "github.com/teamlint/baron/protobuf/types/known/timestamppb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"

	"google.golang.org/grpc"
)

var (
	grpcAddr string
	httpAddr string
	natsAddr string
)

func init() {
	flag.StringVar(&grpcAddr, "grpc.addr", ":5040", "gRPC (HTTP) listen address")
	flag.StringVar(&httpAddr, "http.addr", ":5050", "HTTP listen address")
	flag.StringVar(&natsAddr, "nats.addr", ":4222", "NATS listen address")
}

func main() {
	flag.Parse()

	// Baron GRPC Client
	log.Println("[GRPC.Echo]")
	pln()
	conn, err := grpc.Dial(
		grpcAddr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	baronGRPCClient, err := pb.NewGRPCClient(conn)
	if err != nil {
		log.Fatal(err)
	}

	// GRPC.Echo.Echo
	{
		ctx := context.Background()
		var in pb.EchoRequest
		in.In = "grpc->" + xid.New().String()
		now := time.Now().Unix()
		in.At = &now
		in.CreatedAt = timestamppb.Now()
		out, err := baronGRPCClient.Echo(ctx, &in)
		if err != nil {
			log.Fatalf("[Baron.GRPCClient] Echo.Echo err=%v\n", err)
		}
		log.Printf("[Baron.GRPCClient] Echo.Echo result=%+v\n", *out)
		pln()
	}

	// GRPC.Echo.Louder
	{
		ctx := context.Background()
		var in pb.LouderRequest
		out, err := baronGRPCClient.Louder(ctx, &in)
		if err != nil {
			log.Fatalf("[Baron.GRPCClient] Echo.Louder err=%v\n", err)
		}
		log.Printf("[Baron.GRPCClient] Echo.Louder result=%+v\n", *out)
		pln()
	}

	// GRPC.Echo.LouderGet
	{
		ctx := context.Background()
		var in pb.LouderRequest
		out, err := baronGRPCClient.LouderGet(ctx, &in)
		if err != nil {
			log.Fatalf("[Baron.GRPCClient] Echo.LouderGet err=%v\n", err)
		}
		log.Printf("[Baron.GRPCClient] Echo.LouderGet result=%+v\n", *out)
		pln()
	}

	// Baron HTTP Client
	log.Println("[HTTP][Echo]")
	pln()

	// HTTP.Echo.Echo
	{
		ctx := context.Background()
		var in pb.EchoRequest
		in.In = "http->" + xid.New().String()
		now := time.Now().Unix()
		in.At = &now
		desc := "HTTP测试可选字段"
		in.Desc = &desc
		in.CreatedAt = timestamppb.Now()
		in.JsonStr = wrapperspb.String("wrappers字段")
		baronHTTPClient, err := pb.NewHTTPClient(httpAddr)
		if err != nil {
			log.Fatal(err)
		}
		out, err := baronHTTPClient.Echo(ctx, &in)
		if err != nil {
			log.Fatalf("[Baron.HTTPClient] Echo.Echo err=%v\n", err)
		}
		log.Printf("[Baron.HTTPClient] Echo.Echo result=%+v\n", *out)
		pln()
	}

	// HTTP.Echo.Louder
	{
		ctx := context.Background()
		var in pb.LouderRequest
		baronHTTPClient, err := pb.NewHTTPClient(httpAddr)
		if err != nil {
			log.Fatal(err)
		}
		out, err := baronHTTPClient.Louder(ctx, &in)
		if err != nil {
			log.Fatalf("[Baron.HTTPClient] Echo.Louder err=%v\n", err)
		}
		log.Printf("[Baron.HTTPClient] Echo.Louder result=%+v\n", *out)
		pln()
	}

	// HTTP.Echo.LouderGet
	{
		ctx := context.Background()
		var in pb.LouderRequest
		baronHTTPClient, err := pb.NewHTTPClient(httpAddr)
		if err != nil {
			log.Fatal(err)
		}
		out, err := baronHTTPClient.LouderGet(ctx, &in)
		if err != nil {
			log.Fatalf("[Baron.HTTPClient] Echo.LouderGet err=%v\n", err)
		}
		log.Printf("[Baron.HTTPClient] Echo.LouderGet result=%+v\n", *out)
		pln()
	}

	// Baron NATS Client
	log.Println("[NATS.Echo]")
	pln()
	nc, err := nats.Connect(natsAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	baronNATSClient, err := pb.NewNATSClient(nc)
	if err != nil {
		log.Fatal(err)
	}

	// NATS.Echo.Echo
	{
		ctx := context.Background()
		var in pb.EchoRequest
		in.In = "nats->" + xid.New().String()
		now := time.Now().Unix()
		in.At = &now
		in.CreatedAt = timestamppb.Now()
		out, err := baronNATSClient.Echo(ctx, &in)
		if err != nil {
			log.Fatalf("[Baron.NATSClient] Echo.Echo err=%v\n", err)
		}
		log.Printf("[Baron.NATSClient] Echo.Echo result=%+v\n", *out)
		pln()
	}

	// NATS.Echo.Louder
	{
		ctx := context.Background()
		var in pb.LouderRequest
		out, err := baronNATSClient.Louder(ctx, &in)
		if err != nil {
			log.Fatalf("[Baron.NATSClient] Echo.Louder err=%v\n", err)
		}
		log.Printf("[Baron.NATSClient] Echo.Louder result=%+v\n", *out)
		pln()
	}

	// NATS.Echo.LouderGet
	{
		ctx := context.Background()
		var in pb.LouderRequest
		out, err := baronNATSClient.LouderGet(ctx, &in)
		if err != nil {
			log.Fatalf("[Baron.NATSClient] Echo.LouderGet err=%v\n", err)
		}
		log.Printf("[Baron.NATSClient] Echo.LouderGet result=%+v\n", *out)
		pln()
	}

}
func pln() {
	log.Println("---------------------------------------------------------------------------------")
}
