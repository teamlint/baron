package middleware

import (
	"context"
	"os"
	"reflect"
	"runtime"
	"time"

	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	pb "github.com/teamlint/baron/_example/api/echo"
)

func Logging(in pb.Endpoints) pb.Endpoints {
	logger := kitlog.NewLogfmtLogger(os.Stderr)
	loggingMiddleware := logging(logger)
	in.EchoEndpoint = loggingMiddleware(in.EchoEndpoint)
	in.LouderEndpoint = loggingMiddleware(in.LouderEndpoint)
	in.LouderGetEndpoint = loggingMiddleware(in.LouderGetEndpoint)
	return in
}

func logging(logger kitlog.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				// logger.Log("endpoint", runtime.FuncForPC(reflect.ValueOf(next).Pointer()).Name(), "transport_error", err, "took", time.Since(begin))
				logger.Log("endpoint", InspectFunc(next).Name, "transport_error", err, "took", time.Since(begin))
			}(time.Now())
			return next(ctx, request)
		}
	}
}

// Func is a function description.
type Func struct {
	Name string
	reflect.Type
	reflect.Value
}

// InspectFunc inspects function.
func InspectFunc(fn interface{}) Func {
	if reflect.ValueOf(fn).Kind() != reflect.Func {
		return Func{}
	}
	val := reflect.ValueOf(fn)
	typ := val.Type()
	funcForPC := runtime.FuncForPC(val.Pointer())
	return Func{
		Name:  funcForPC.Name(),
		Type:  typ,
		Value: val,
	}
}
