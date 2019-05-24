package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	"grpcwrapper/examples/echo"
	//"grpcwrapper/middleware"
)

func ClientHook(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	fmt.Println("客户端请求: 方法-> ", method, " 目标: ", cc.Target)
	return invoker(ctx, method, req, reply, cc, opts...)
}

func callUnaryEcho(client echo.EchoClient, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.UnaryEcho(ctx, &echo.EchoRequest{Message: message})
	if err != nil {
		log.Fatalf("client.UnaryEcho(_) = _, %v: ", err)
	}
	fmt.Println("UnaryEcho: ", resp.Message)
}

func main() {
	fmt.Println("client")
	conn, err := grpc.Dial("0.0.0.0:8888", grpc.WithInsecure(), grpc.WithUnaryInterceptor(ClientHook))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	for i := 0; i < 1000; i++ {
		go func() {
			tmp := i
			rgc := echo.NewEchoClient(conn)
			for j := 0; j < 1000; j++ {
				callUnaryEcho(rgc, fmt.Sprintf("this is examples %d", tmp*10000+j))
			}
		}()
	}

	select {}
}
