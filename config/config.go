package config

import "io"

// Config defines the inputs to a baron service generation
type Config struct {
	// The first path in $GOPATH
	GoPath []string

	// The go package where .pb.go files protoc-gen-go creates will be written
	PBPackage string
	PBPath    string
	// The go package where the service code will be written
	ServicePackage string
	ServicePath    string

	// The paths to each of the .proto files baron is being run against
	DefPaths []string
	// The files of a previously generated service, may be nil
	PrevGen map[string]io.Reader
	// Generate service client CLI
	GenClient bool
	// Service transport protocol: all|grpc|http|nats
	Transport string
	// 服务定义
	Svcdef bool
}
