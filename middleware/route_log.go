package middleware

import (
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

/*
 * 这是一个路由打印中间件，目前是测试代码；
 * 后续会进行扩展
 */

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		fmt.Println("拦截器1: ", req)
		return handler(ctx, req)
	}
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return handler(srv, stream)
	}
}
