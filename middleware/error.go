package middleware

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	gstatus "google.golang.org/grpc/status"
)

func ClientErrorHook(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	err := invoker(ctx, method, req, reply, cc, opts...)

	gst, ok := gstatus.FromError(err)
	if ok {
		ec := ToEcode(gst)
		err = errors.WithMessage(ec, gst.Message())
	}

	return err
}

func UnaryServerErrorHook(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {
	resp, err = handler(ctx, req)
	if err != nil {
		fmt.Println("server: convert err to grpc.Status")
		err = FromError(err).Err()
	}
	return
}

func StreamServerErrorHook(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) (err error) {
	err = handler(srv, stream)
	if err != nil {
		err = FromError(err).Err()
	}

	return
}
