module github.com/teamlint/baron/_example

go 1.15

require (
	github.com/go-kit/kit v0.10.0
	github.com/goava/di v1.2.1
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.4.3
	github.com/gorilla/mux v1.8.0
	github.com/nats-io/nats.go v1.9.1
	github.com/pkg/errors v0.9.1
	github.com/rs/xid v1.2.1 // indirect
	github.com/teamlint/baron v0.2.1
	golang.org/x/net v0.0.0-20201021035429-f5854403a974
	google.golang.org/grpc v1.33.0
	google.golang.org/protobuf v1.25.0
)

replace github.com/teamlint/baron v0.2.1 => ../
