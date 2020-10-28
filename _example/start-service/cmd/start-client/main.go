package main

import (
	"context"
	"flag"
	"log"

	"github.com/nats-io/nats.go"

	pb "github.com/teamlint/baron/_example/api/start"

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
	log.Println("[GRPC.Start]")
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

	// GRPC.Start.Status
	{
		ctx := context.Background()
		var in pb.StatusRequest
		out, err := baronGRPCClient.Status(ctx, &in)
		if err != nil {
			log.Fatalf("[Baron.GRPCClient] Start.Status err=%v\n", err)
		}
		log.Printf("[Baron.GRPCClient] Start.Status result=%+v\n", *out)
		pln()
	}

	// Baron HTTP Client
	log.Println("[HTTP][Start]")
	pln()

	// HTTP.Start.Status
	{
		ctx := context.Background()
		var in pb.StatusRequest
		baronHTTPClient, err := pb.NewHTTPClient(httpAddr)
		if err != nil {
			log.Fatal(err)
		}
		out, err := baronHTTPClient.Status(ctx, &in)
		if err != nil {
			log.Fatalf("[Baron.HTTPClient] Start.Status err=%v\n", err)
		}
		log.Printf("[Baron.HTTPClient] Start.Status result=%+v\n", *out)
		pln()
	}

	// Baron NATS Client
	log.Println("[NATS.Start]")
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

	// NATS.Start.Status
	{
		ctx := context.Background()
		var in pb.StatusRequest
		out, err := baronNATSClient.Status(ctx, &in)
		if err != nil {
			log.Fatalf("[Baron.NATSClient] Start.Status err=%v\n", err)
		}
		log.Printf("[Baron.NATSClient] Start.Status result=%+v\n", *out)
		pln()
	}

}
func pln() {
	log.Println("---------------------------------------------------------------------------------")
}
