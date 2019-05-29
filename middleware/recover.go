package middleware

import (
	"fmt"
	"os"
	"runtime"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"grpcwrapper/log"
)

func get_stack() string {
	const size = 64 << 10
	buf := make([]byte, size)
	rs := runtime.Stack(buf, false)
	if rs > size {
		rs = size
	}
	buf = buf[:size]
	return string(buf)
}

func UnaryServerRecoverInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if rerr := recover(); rerr != nil {
				sinfo := get_stack()
				linfo := fmt.Sprintf("server panic: %v %v %s", req, rerr, sinfo)
				fmt.Fprintf(os.Stderr, linfo)
				log.Errorf(linfo)
				err = status.Errorf(codes.Unknown, "server panic")
			}
		}()
		resp, err = handler(ctx, req)
		return
	}
}

func StreamServerRecoverInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if rerr := recover(); rerr != nil {
				sinfo := get_stack()
				log.Errorf("server panic: %v %v %s ", stream, rerr, sinfo)
				err = status.Errorf(codes.Unknown, "server panic")
			}
		}()
		err = handler(srv, stream)
		return
	}
}
