package grpcwrapper

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"github.com/lj-211/grpcwrapper/middleware"
)

var defaultDialTimeout = time.Second * 3

type ClientOpt struct {
	Interceptors []grpc.UnaryClientInterceptor
	CallBlock    bool
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
	if co.CallBlock {
		dops = append(dops, grpc.WithBlock())
	}
	dops = append(dops, opts...)

	new_ctx, cancel := context.WithTimeout(ctx, defaultDialTimeout)
	conn, err = grpc.DialContext(new_ctx, target, dops...)
	defer cancel()

	return
}
