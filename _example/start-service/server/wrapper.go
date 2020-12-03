package server

import (
	grpctransport "github.com/go-kit/kit/transport/grpc"
	httptransport "github.com/go-kit/kit/transport/http"
	natstransport "github.com/go-kit/kit/transport/nats"
)

var (
	_ grpctransport.Server
	_ httptransport.Server
	_ natstransport.Subscriber
)

// GRPCServerOptions GRPC 服务传输选项
func GRPCServerOptions() []grpctransport.ServerOption {
	// default
	// serverOptions := []grpctransport.ServerOption{
	// 	grpctransport.ServerBefore(GRPCMetadataToContext),
	// }
	return nil
}

// HTTPServerOptions GRPC 服务传输选项
func HTTPServerOptions() []httptransport.ServerOption {
	// default
	// serverOptions := []httptransport.ServerOption{
	// 	httptransport.ServerBefore(pb.HTTPHeadersToContext),
	// 	httptransport.ServerErrorEncoder(errorEncoder),
	// 	httptransport.ServerAfter(httptransport.SetContentType(contentType)),
	// }
	return nil
}

// NATSSubscriberOptions NATS 订阅器选项
func NATSSubscriberOptions() []natstransport.SubscriberOption {
	// default none
	// options := []natstransport.SubscriberOption{}
	return nil
}
