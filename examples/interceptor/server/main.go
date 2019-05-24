package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"grpcwrapper"
	"grpcwrapper/examples/echo"
	"grpcwrapper/middleware"
)

type EchoServer struct{}

func (s *EchoServer) UnaryEcho(ctx context.Context, in *echo.EchoRequest) (*echo.EchoResponse, error) {
	fmt.Printf("收到消息: %q\n", in.Message)
	return &echo.EchoResponse{Message: in.Message}, nil
}

func (s *EchoServer) ServerStreamingEcho(in *echo.EchoRequest, stream echo.Echo_ServerStreamingEchoServer) error {
	return status.Error(codes.Unimplemented, "not implemented")
}

func (s *EchoServer) ClientStreamingEcho(stream echo.Echo_ClientStreamingEchoServer) error {
	return status.Error(codes.Unimplemented, "not implemented")
}

func (s *EchoServer) BidirectionalStreamingEcho(stream echo.Echo_BidirectionalStreamingEchoServer) error {
	return nil
}

func main() {
	engine := grpcwrapper.Default()
	engine.Use(middleware.UnaryServerInterceptor(), middleware.StreamServerInterceptor())
	engine.Use(middleware.UnaryServerInterceptor(), middleware.StreamServerInterceptor())
	engine.InitServer()

	echo.RegisterEchoServer(engine.GetRawServer(), &EchoServer{})
	fmt.Println("启动服务...")
	engine.Run("0.0.0.0:8888")
}
