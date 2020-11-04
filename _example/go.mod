module github.com/teamlint/baron/_example

go 1.15

require (
	github.com/go-kit/kit v0.10.0
	github.com/goava/di v1.2.1
	github.com/gogo/protobuf v1.2.1
	github.com/golang/protobuf v1.4.3
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.0.1
	github.com/nats-io/nats.go v1.9.1
	github.com/pkg/errors v0.9.1
	github.com/rs/xid v1.2.1
	golang.org/x/net v0.0.0-20201021035429-f5854403a974
	golang.org/x/sys v0.0.0-20201020230747-6e5568b54d1a // indirect
	google.golang.org/genproto v0.0.0-20201019141844-1ed22bb0c154
	google.golang.org/grpc v1.33.1
	google.golang.org/protobuf v1.25.0
	gopkg.in/yaml.v2 v2.3.0 // indirect
)

replace github.com/teamlint/baron v0.2.4 => ../
