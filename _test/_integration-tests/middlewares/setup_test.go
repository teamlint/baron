package test

import (
	"os"
	"testing"

	"github.com/teamlint/baron/cmd/_integration-tests/middlewares/middlewarestest-service/handlers"
	svc "github.com/teamlint/baron/cmd/_integration-tests/middlewares/middlewarestest-service/svc"
	pb "github.com/teamlint/baron/cmd/_integration-tests/middlewares/proto"
)

var middlewareEndpoints svc.Endpoints

func TestMain(m *testing.M) {

	var service pb.MiddlewaresTestServer
	{
		service = handlers.NewService()
	}

	// Endpoint domain.
	alwaysWrapped := svc.MakeAlwaysWrappedEndpoint(service)
	sometimesWrapped := svc.MakeSometimesWrappedEndpoint(service)
	labeledTestHandler := svc.MakeLabeledTestHandlerEndpoint(service)

	middlewareEndpoints = svc.Endpoints{
		AlwaysWrappedEndpoint:      alwaysWrapped,
		SometimesWrappedEndpoint:   sometimesWrapped,
		LabeledTestHandlerEndpoint: labeledTestHandler,
	}

	middlewareEndpoints = handlers.WrapEndpoints(middlewareEndpoints)

	os.Exit(m.Run())
}
