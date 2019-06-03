package grpcwrapper

import (
	"context"

	"google.golang.org/grpc"

	"github.com/lj-211/grpcwrapper/middleware"
)

type ClientOpt struct {
	Interceptors []grpc.UnaryClientInterceptor
}

var defaultClientOpt *ClientOpt = &ClientOpt{
	Interceptors: make([]grpc.UnaryClientInterceptor, 0),
}

func DefaultClient() *ClientOpt {
	return defaultClientOpt
}

func (co *ClientOpt) DialContext(ctx context.Context, target string,
	opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	co.Interceptors = append(co.Interceptors, middleware.ClientErrorHook)
	co.Interceptors = append(co.Interceptors, middleware.UnaryClientRecoverInterceptor())

	dops := make([]grpc.DialOption, 0)
	// interceptors
	dops = append(dops, grpc.WithUnaryInterceptor(ChainUnaryClient(co.Interceptors)))

	dops = append(dops, grpc.WithInsecure())
	dops = append(dops, opts...)

	conn, err = grpc.DialContext(ctx, target, dops...)

	return
}
