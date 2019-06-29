package middleware

import (
	"context"

	"google.golang.org/grpc"
	gstatus "google.golang.org/grpc/status"
)

func ClientErrorHook(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	err := invoker(ctx, method, req, reply, cc, opts...)

	if err != nil {
		gst, ok := gstatus.FromError(err)
		if ok {
			ec := ToEcode(gst)
			err = ec
		}
	}

	return err
}

func UnaryServerErrorHook(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {
	resp, err = handler(ctx, req)
	if err != nil {
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
