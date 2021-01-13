package server

import (
	pb "github.com/teamlint/baron/_example/api/start"
	"github.com/teamlint/baron/_example/start-service/service"
)

// Endpoints application domain
func Endpoints() pb.Endpoints {
	svc := service.NewService()
	endpoints := NewEndpoints(svc)

	return endpoints
}
