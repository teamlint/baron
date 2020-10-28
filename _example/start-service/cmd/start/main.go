package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	// This Service Server
	"github.com/teamlint/baron/_example/start-service/server"
)

var config server.Config

func init() {
	flag.StringVar(&config.DebugAddr, "debug.addr", ":5060", "Debug and metrics listen address")
	flag.StringVar(&config.HTTPAddr, "http.addr", ":5050", "HTTP listen address")
	flag.StringVar(&config.GRPCAddr, "grpc.addr", ":5040", "gRPC (HTTP) listen address")
	flag.StringVar(&config.NATSAddr, "nats.addr", ":4222", "NATS listen address")
	flag.BoolVar(&config.Debug, "debug", false, "Debug mode")

	// Use environment variables, if set. Flags have priority over Env vars.
	if debug := os.Getenv("DEBUG"); debug != "" {
		config.Debug = "true" == strings.ToLower(debug)
	}
	if addr := os.Getenv("DEBUG_ADDR"); addr != "" {
		config.DebugAddr = addr
	}
	if port := os.Getenv("PORT"); port != "" {
		config.HTTPAddr = fmt.Sprintf(":%s", port)
	}
	if addr := os.Getenv("HTTP_ADDR"); addr != "" {
		config.HTTPAddr = addr
	}
	if addr := os.Getenv("GRPC_ADDR"); addr != "" {
		config.GRPCAddr = addr
	}
	if addr := os.Getenv("NATS_ADDR"); addr != "" {
		config.NATSAddr = addr
	}
}

func main() {
	// Update addresses if they have been overwritten by flags
	flag.Parse()

	server.Run(config)
}
