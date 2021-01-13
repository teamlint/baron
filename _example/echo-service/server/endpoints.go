package server

import (
	"log"

	"github.com/goava/di"
	pb "github.com/teamlint/baron/_example/api/echo"
	"github.com/teamlint/baron/_example/echo-service/adapter"
	"github.com/teamlint/baron/_example/echo-service/domain"
	"github.com/teamlint/baron/_example/echo-service/service"
)

// Endpoints application domain
func Endpoints() pb.Endpoints {
	// return goway()
	return ioc()
}

func goway() pb.Endpoints {
	echoRepo := adapter.NewEchoRepository()
	echoSvc := domain.NewEchoService(echoRepo)
	svc := service.NewService()
	svc.EchoService = echoSvc
	endpoints := NewEndpoints(svc)

	return endpoints
}

func ioc() pb.Endpoints {
	di.SetTracer(&di.StdTracer{})
	c, _ := di.New(
		di.Provide(adapter.NewEchoRepository, di.As(new(domain.EchoRepository))),
		di.Provide(domain.NewEchoService),
		di.Provide(service.NewService, di.As(new(pb.EchoServer))),
		di.Provide(NewEndpoints),
	)

	var endpoints pb.Endpoints
	if err := c.Resolve(&endpoints); err != nil {
		log.Fatal(err)
	}

	return endpoints
}
