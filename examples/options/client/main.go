package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	"grpcwrapper"
	"grpcwrapper/examples/echo"
	//"grpcwrapper/middleware"
)

func ClientHook(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	fmt.Println("客户端请求: 方法-> ", method, " 目标: ", cc.Target)
	return invoker(ctx, method, req, reply, cc, opts...)
}

func callUnaryEcho(ctx context.Context, client echo.EchoClient, message string) {
	resp, err := client.UnaryEcho(ctx, &echo.EchoRequest{Message: message})
	if err != nil {
		log.Fatalf("client.UnaryEcho(_) = _, %v: ", err)
	}
	fmt.Println("UnaryEcho: ", resp.Message)
}

func main() {
	fmt.Println("client")

	opts := make([]grpc.DialOption, 0)
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithUnaryInterceptor(ClientHook))
	opts = append(opts, grpcwrapper.DefaultDialKeepailveOpt())
	conn, err := grpc.Dial("0.0.0.0:8888", grpc.WithInsecure(), grpc.WithUnaryInterceptor(ClientHook))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	for i := 0; i < 1000; i++ {
		go func() {
			grpcwrapper.CallWithTimeOut(func(ctx context.Context) {
				rgc := echo.NewEchoClient(conn)
				callUnaryEcho(ctx, rgc, fmt.Sprintf("this is examples"))
			})
		}()
		time.Sleep(5 * time.Second)
	}

	select {}
}
