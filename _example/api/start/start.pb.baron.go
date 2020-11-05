// Code generated by baron. DO NOT EDIT.
// Rerunning baron will overwrite this file.
// Version: v0.2.4-4-g6637cde-dirty
// Version Date: 2020-11-05T15:49:40+08:00

package start

// This file contains methods to make individual endpoints from services,
// request and response types to serve those endpoints, as well as encoders and
// decoders for those types, for all of our supported transport serialization
// formats.

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-kit/kit/endpoint"
	any "google.golang.org/protobuf/types/known/anypb"
	duration "google.golang.org/protobuf/types/known/durationpb"
	empty "google.golang.org/protobuf/types/known/emptypb"
	_struct "google.golang.org/protobuf/types/known/structpb"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
	wrappers "google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	field_mask "google.golang.org/protobuf/types/known/fieldmaskpb"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	httptransport "github.com/go-kit/kit/transport/http"
	natstransport "github.com/go-kit/kit/transport/nats"
)

var (
	_ = fmt.Sprint
	_ = bytes.Compare
	_ = strconv.Atoi
	_ = httptransport.NewServer
	_ = ioutil.NopCloser
	_ = io.Copy
	_ = errors.Wrap
	// google.protobuf types
	_ any.Any
	_ duration.Duration
	_ empty.Empty
	_ _struct.Struct
	_ timestamp.Timestamp
	_ wrappers.StringValue
	_ field_mask.FieldMask
)

/************************************** Endpoints ******************************************/

// In a server, it's useful for functions that need to operate on a per-endpoint
// basis. For example, you might pass an Endpoints to a function that produces
// an http.Handler, with each method (endpoint) wired up to a specific path. (It
// is probably a mistake in design to invoke the Service methods on the
// Endpoints struct in a server.)
//
// In a client, it's useful to collect individually constructed endpoints into a
// single type that implements the Service interface. For example, you might
// construct individual endpoints using transport/http.NewClient, combine them into an Endpoints, and return it to the caller as a Service.
type Endpoints struct {
	UnimplementedStartServer

	StatusEndpoint endpoint.Endpoint
}

// Endpoints

func (e Endpoints) Status(ctx context.Context, in *StatusRequest) (*StatusResponse, error) {
	response, err := e.StatusEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return response.(*StatusResponse), nil
}

// Make Endpoints

func MakeStatusEndpoint(s StartServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*StatusRequest)
		v, err := s.Status(ctx, req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

// WrapAllExcept wraps each Endpoint field of struct Endpoints with a
// go-kit/kit/endpoint.Middleware.
// Use this for applying a set of middlewares to every endpoint in the service.
// Optionally, endpoints can be passed in by name to be excluded from being wrapped.
// WrapAllExcept(middleware, "Status", "Ping")
func (e *Endpoints) WrapAllExcept(middleware endpoint.Middleware, excluded ...string) {
	included := map[string]struct{}{
		"Status": {},
	}

	for _, ex := range excluded {
		if _, ok := included[ex]; !ok {
			panic(fmt.Sprintf("Excluded endpoint '%s' does not exist; see middlewares/endpoints.go", ex))
		}
		delete(included, ex)
	}

	for inc := range included {
		if inc == "Status" {
			e.StatusEndpoint = middleware(e.StatusEndpoint)
		}
	}
}

// LabeledMiddleware will get passed the endpoint name when passed to
// WrapAllLabeledExcept, this can be used to write a generic metrics
// middleware which can send the endpoint name to the metrics collector.
type LabeledMiddleware func(string, endpoint.Endpoint) endpoint.Endpoint

// WrapAllLabeledExcept wraps each Endpoint field of struct Endpoints with a
// LabeledMiddleware, which will receive the name of the endpoint. See
// LabeldMiddleware. See method WrapAllExept for details on excluded
// functionality.
func (e *Endpoints) WrapAllLabeledExcept(middleware func(string, endpoint.Endpoint) endpoint.Endpoint, excluded ...string) {
	included := map[string]struct{}{
		"Status": {},
	}

	for _, ex := range excluded {
		if _, ok := included[ex]; !ok {
			panic(fmt.Sprintf("Excluded endpoint '%s' does not exist; see middlewares/endpoints.go", ex))
		}
		delete(included, ex)
	}

	for inc := range included {
		if inc == "Status" {
			e.StatusEndpoint = middleware("Status", e.StatusEndpoint)
		}
	}
}

/************************************** GRPCServer ******************************************/

// MakeGRPCServer makes a set of endpoints available as a gRPC StartServer.
func MakeGRPCServer(endpoints Endpoints, options ...grpctransport.ServerOption) StartServer {
	serverOptions := []grpctransport.ServerOption{
		grpctransport.ServerBefore(GRPCMetadataToContext),
	}
	serverOptions = append(serverOptions, options...)
	// start
	return &grpcServer{

		status: grpctransport.NewServer(
			endpoints.StatusEndpoint,
			DecodeGRPCStatusRequest,
			EncodeGRPCStatusResponse,
			serverOptions...,
		),
	}
}

// grpcServer implements the StartServer interface
type grpcServer struct {
	UnimplementedStartServer

	status grpctransport.Handler
}

// Methods for grpcServer to implement StartServer interface

func (s *grpcServer) Status(ctx context.Context, req *StatusRequest) (*StatusResponse, error) {
	_, rep, err := s.status.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*StatusResponse), nil
}

// Server Decode

// DecodeGRPCStatusRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC status request to a user-domain status request. Primarily useful in a server.
func DecodeGRPCStatusRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*StatusRequest)
	return req, nil
}

// Server Encode

// EncodeGRPCStatusResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain status response to a gRPC status reply. Primarily useful in a server.
func EncodeGRPCStatusResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*StatusResponse)
	return resp, nil
}

// Helpers

func GRPCMetadataToContext(ctx context.Context, md metadata.MD) context.Context {
	for k, v := range md {
		if v != nil {
			// The key is added both in metadata format (k) which is all lower
			// and the http.CanonicalHeaderKey of the key so that it can be
			// accessed in either format
			ctx = context.WithValue(ctx, k, v[0])
			ctx = context.WithValue(ctx, http.CanonicalHeaderKey(k), v[0])
		}
	}

	return ctx
}

/************************************** GRPCClient ******************************************/

// NewGRPCClient returns an service backed by a gRPC client connection. It is the
// responsibility of the caller to dial, and later close, the connection.
func NewGRPCClient(conn *grpc.ClientConn, options ...grpctransport.ClientOption) (StartServer, error) {
	var statusEndpoint endpoint.Endpoint
	{
		statusEndpoint = grpctransport.NewClient(
			conn,
			"start.Start",
			"Status",
			EncodeGRPCStatusRequest,
			DecodeGRPCStatusResponse,
			StatusResponse{},
			options...,
		).Endpoint()
	}

	return Endpoints{
		StatusEndpoint: statusEndpoint,
	}, nil
}

// GRPC Client Decode

// DecodeGRPCStatusResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC status reply to a user-domain status response. Primarily useful in a client.
func DecodeGRPCStatusResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*StatusResponse)
	return reply, nil
}

// GRPC Client Encode

// EncodeGRPCStatusRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain status request to a gRPC status request. Primarily useful in a client.
func EncodeGRPCStatusRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*StatusRequest)
	return req, nil
}

// ContextValuesToGRPCMetadata is a grpctransport.ClientRequestFunc
func ContextValuesToGRPCMetadata(keys []string) grpctransport.ClientRequestFunc {
	return func(ctx context.Context, md *metadata.MD) context.Context {
		var pairs []string
		for _, k := range keys {
			if v, ok := ctx.Value(k).(string); ok {
				pairs = append(pairs, k, v)
			}
		}

		if pairs != nil {
			*md = metadata.Join(*md, metadata.Pairs(pairs...))
		}

		return ctx
	}
}

/************************************** HTTPServer ******************************************/

const contentType = "application/json; charset=utf-8"

// MakeHTTPHandler returns a handler that makes a set of endpoints available
// on predefined paths.
func MakeHTTPHandler(endpoints Endpoints, options ...httptransport.ServerOption) http.Handler {
	serverOptions := []httptransport.ServerOption{
		httptransport.ServerBefore(HTTPHeadersToContext),
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerAfter(httptransport.SetContentType(contentType)),
	}
	serverOptions = append(serverOptions, options...)
	m := mux.NewRouter()

	m.Methods("GET").Path("/status").Handler(httptransport.NewServer(
		endpoints.StatusEndpoint,
		DecodeHTTPStatusZeroRequest,
		EncodeHTTPGenericResponse,
		serverOptions...,
	))
	return m
}

// ErrorEncoder writes the error to the ResponseWriter, by default a content
// type of application/json, a body of json with key "error" and the value
// error.Error(), and a status code of 500. If the error implements Headerer,
// the provided headers will be applied to the response. If the error
// implements json.Marshaler, and the marshaling succeeds, the JSON encoded
// form of the error will be used. If the error implements StatusCoder, the
// provided StatusCode will be used instead of 500.
func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	body, _ := json.Marshal(errorWrapper{Error: err.Error()})
	if marshaler, ok := err.(json.Marshaler); ok {
		if jsonBody, marshalErr := marshaler.MarshalJSON(); marshalErr == nil {
			body = jsonBody
		}
	}
	w.Header().Set("Content-Type", contentType)
	if headerer, ok := err.(httptransport.Headerer); ok {
		for k := range headerer.Headers() {
			w.Header().Set(k, headerer.Headers().Get(k))
		}
	}
	code := http.StatusInternalServerError
	if sc, ok := err.(httptransport.StatusCoder); ok {
		code = sc.StatusCode()
	}
	w.WriteHeader(code)
	w.Write(body)
}

func errorDecoder(buf []byte) error {
	var w errorWrapper
	if err := json.Unmarshal(buf, &w); err != nil {
		const size = 8196
		if len(buf) > size {
			buf = buf[:size]
		}
		return fmt.Errorf("response body '%s': cannot parse non-json request body", buf)
	}

	return errors.New(w.Error)
}

type errorWrapper struct {
	Error string `json:"error"`
}

// httpError satisfies the Headerer and StatusCoder interfaces in
// package github.com/go-kit/kit/transport/http.
type httpError struct {
	error
	statusCode int
	headers    map[string][]string
}

func (h httpError) StatusCode() int {
	return h.statusCode
}

func (h httpError) Headers() http.Header {
	return h.headers
}

// Server Decode

// DecodeHTTPStatusZeroRequest is a transport/http.DecodeRequestFunc that
// decodes a JSON-encoded status request from the HTTP request
// body. Primarily useful in a server.
func DecodeHTTPStatusZeroRequest(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	var req StatusRequest
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read body of http request")
	}
	if len(buf) > 0 {
		unmarshaller := protojson.UnmarshalOptions{
			AllowPartial:   true,
			DiscardUnknown: true,
		}
		if err = unmarshaller.Unmarshal(buf, &req); err != nil {
			const size = 8196
			if len(buf) > size {
				buf = buf[:size]
			}
			return nil, httpError{errors.Wrapf(err, "request body '%s': cannot parse non-json request body", buf),
				http.StatusBadRequest,
				nil,
			}
		}
	}

	pathParams := mux.Vars(r)
	_ = pathParams

	queryParams := r.URL.Query()
	_ = queryParams

	if FullStatusStrArr, ok := queryParams["full"]; ok {
		FullStatusStr := FullStatusStrArr[0]
		FullStatus, err := strconv.ParseBool(FullStatusStr)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Error while extracting FullStatus from query, queryParams: %v", queryParams))
		}
		req.Full = FullStatus
	}

	if MsgStatusStrArr, ok := queryParams["msg"]; ok {
		MsgStatusStr := MsgStatusStrArr[0]
		MsgStatus := MsgStatusStr
		req.Msg = &MsgStatus
	}

	return &req, err
}

// EncodeHTTPGenericResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer. Primarily useful in a server.
func EncodeHTTPGenericResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	marshaller := protojson.MarshalOptions{
		AllowPartial:    true,
		UseProtoNames:   true,
		UseEnumNumbers:  true,
		EmitUnpopulated: true,
	}
	buf, err := marshaller.Marshal(response.(proto.Message))
	if err != nil {
		return err
	}
	_, err = w.Write(buf)
	return err
}

// Helper functions

func HTTPHeadersToContext(ctx context.Context, r *http.Request) context.Context {
	for k := range r.Header {
		// The key is added both in http format (k) which has had
		// http.CanonicalHeaderKey called on it in transport as well as the
		// strings.ToLower which is the grpc metadata format of the key so
		// that it can be accessed in either format
		ctx = context.WithValue(ctx, k, r.Header.Get(k))
		ctx = context.WithValue(ctx, strings.ToLower(k), r.Header.Get(k))
	}

	// Tune specific change.
	// also add the request url
	ctx = context.WithValue(ctx, "request-url", r.URL.Path)
	ctx = context.WithValue(ctx, "transport", "HTTPJSON")

	return ctx
}

/************************************** HTTPClient ******************************************/

// NewHTTPClient returns a service backed by an HTTP server living at the remote
// instance. We expect instance to come from a service discovery system, so
// likely of the form "host:port".
func NewHTTPClient(instance string, options ...httptransport.ClientOption) (StartServer, error) {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		return nil, err
	}
	_ = u

	var StatusZeroEndpoint endpoint.Endpoint
	{
		StatusZeroEndpoint = httptransport.NewClient(
			"GET",
			copyURL(u, "/status"),
			EncodeHTTPStatusZeroRequest,
			DecodeHTTPStatusResponse,
			options...,
		).Endpoint()
	}

	return Endpoints{
		StatusEndpoint: StatusZeroEndpoint,
	}, nil
}

func copyURL(base *url.URL, path string) *url.URL {
	next := *base
	next.Path = path
	return &next
}

// CtxValuesToSend configures the http client to pull the specified keys out of
// the context and add them to the http request as headers.  Note that keys
// will have net/http.CanonicalHeaderKey called on them before being send over
// the wire and that is the form they will be available in the server context.
func CtxValuesToSend(keys ...string) httptransport.ClientOption {
	return httptransport.ClientBefore(func(ctx context.Context, r *http.Request) context.Context {
		for _, k := range keys {
			if v, ok := ctx.Value(k).(string); ok {
				r.Header.Set(k, v)
			}
		}
		return ctx
	})
}

// HTTP Client Decode

// DecodeHTTPStatusResponse is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded StatusResponse response from the HTTP response body.
// If the response has a non-200 status code, we will interpret that as an
// error and attempt to decode the specific error message from the response
// body. Primarily useful in a client.
func DecodeHTTPStatusResponse(_ context.Context, r *http.Response) (interface{}, error) {
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err == io.EOF {
		return nil, errors.New("response http body empty")
	}
	if err != nil {
		return nil, errors.Wrap(err, "cannot read http body")
	}

	if r.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(errorDecoder(buf), "status code: '%d'", r.StatusCode)
	}

	var resp StatusResponse
	unmarshaller := protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}
	if err = unmarshaller.Unmarshal(buf, &resp); err != nil {
		return nil, errorDecoder(buf)
	}

	return &resp, nil
}

// HTTP Client Encode

// EncodeHTTPStatusZeroRequest is a transport/http.EncodeRequestFunc
// that encodes a status request into the various portions of
// the http request (path, query, and body).
func EncodeHTTPStatusZeroRequest(_ context.Context, r *http.Request, request interface{}) error {
	strval := ""
	_ = strval
	req := request.(*StatusRequest)
	_ = req

	r.Header.Set("transport", "HTTPJSON")
	r.Header.Set("request-url", r.URL.Path)

	// Set the path parameters
	path := strings.Join([]string{
		"",
		"status",
	}, "/")
	u, err := url.Parse(path)
	if err != nil {
		return errors.Wrapf(err, "couldn't unmarshal path %q", path)
	}
	r.URL.RawPath = u.RawPath
	r.URL.Path = u.Path

	// Set the query parameters
	values := r.URL.Query()
	var tmp []byte
	_ = tmp

	values.Add("full", fmt.Sprint(req.Full))

	if req.Msg != nil {
		values.Add("msg", fmt.Sprint(*req.Msg))
	}

	r.URL.RawQuery = values.Encode()
	return nil
}

/************************************** NATSServer ******************************************/

// ServeNATS
func ServeNATS(conn *nats.Conn, srv *natsServer) error {

	statusSub, err := conn.QueueSubscribe("Start.Status", "Start", srv.Status.ServeMsg(conn))
	if err != nil {
		return err
	}
	_ = statusSub

	return nil

}

// MakeNATSServer makes a set of endpoints available as a NATS StartServer.
func MakeNATSServer(endpoints Endpoints, options ...natstransport.SubscriberOption) *natsServer {
	serverOptions := []natstransport.SubscriberOption{
		// grpctransport.ServerBefore(GRPCMetadataToContext),
	}
	serverOptions = append(serverOptions, options...)
	// start
	return &natsServer{

		Status: natstransport.NewSubscriber(
			endpoints.StatusEndpoint,
			DecodeNATSStatusRequest,
			// EncodeNATSStatusResponse,
			EncodeNATSGenericResponse,
			serverOptions...,
		),
	}
}

// natsServer implements the StartServer interface
type natsServer struct {
	Status *natstransport.Subscriber
}

// Server Decode

func DecodeNATSStatusRequest(ctx context.Context, msg *nats.Msg) (interface{}, error) {
	var req StatusRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		return nil, err
	}
	return &req, nil
}

// Server Encode

// EncodeNATSGenericResponse is a transport/nats.EncodeResponseFunc that encodes
// JSON object to the subscriber reply. Many JSON-over services can use it as
// a sensible default.
func EncodeNATSGenericResponse(ctx context.Context, reply string, nc *nats.Conn, response interface{}) error {
	return natstransport.EncodeJSONResponse(ctx, reply, nc, response)
}

/************************************** NATSClient ******************************************/

// NewNATSClient returns an service backed by a nats client connection. It is the
// responsibility of the caller to dial, and later close, the connection.
func NewNATSClient(conn *nats.Conn, options ...natstransport.PublisherOption) (StartServer, error) {
	var statusEndpoint endpoint.Endpoint
	{
		statusEndpoint = natstransport.NewPublisher(
			conn,
			"Start.Status",
			EncodeNATSGenericRequest,
			DecodeNATSStatusResponse,
			options...,
		).Endpoint()
	}

	return Endpoints{
		StatusEndpoint: statusEndpoint,
	}, nil
}

// NATS Client Decode

// DecodeNATSStatusResponse is a transport/nats.DecodeResponseFunc that converts a
// nats.Msg to a user-domain status response. Primarily useful in a client.
func DecodeNATSStatusResponse(ctx context.Context, msg *nats.Msg) (interface{}, error) {
	// if err := json.Unmarshal(msg.Data, &errResponse); err != nil {
	// 	return nil, err
	// }
	// if errResponse.Error != nil && errResponse.Error.Code == "_internal_" {
	// 	return nil, errors.New(errResponse.Error.Message)
	// }
	var resp StatusResponse
	if err := json.Unmarshal(msg.Data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// NATS Client Encode
// EncodeNATSGenericRequest is a transport/nats.EncodeRequestFunc that serializes the request as a
// JSON object to the Data of the Msg. Many JSON-over-NATS services can use it as
// a sensible default.
func EncodeNATSGenericRequest(ctx context.Context, msg *nats.Msg, request interface{}) error {
	return natstransport.EncodeJSONRequest(ctx, msg, request)
}
