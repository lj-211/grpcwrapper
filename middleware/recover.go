package middleware

import (
	"fmt"
	"os"
	"runtime"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	//"google.golang.org/grpc/codes"
	//"google.golang.org/grpc/status"

	"github.com/lj-211/common/ecode"
	"github.com/lj-211/grpcwrapper/log"
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
				err = ecode.Errorf(ecode.ServerErr, "server panic")
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
				linfo := fmt.Sprintf("server panic: %v %v %s", stream, rerr, sinfo)
				fmt.Fprintln(os.Stderr, linfo)
				log.Errorf(linfo)
				err = ecode.Errorf(ecode.ServerErr, "server panic")
			}
		}()
		err = handler(srv, stream)
		return
	}
}

func UnaryClientRecoverInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		defer func() {
			if rerr := recover(); rerr != nil {
				sinfo := get_stack()
				linfo := fmt.Sprintf("server panic: %v %v %s", req, rerr, sinfo)
				fmt.Fprintln(os.Stderr, linfo)
				err = ecode.Errorf(ecode.ServerErr, "server panic")
			}
		}()
		err = invoker(ctx, method, req, reply, cc, opts...)
		return
	}
}
