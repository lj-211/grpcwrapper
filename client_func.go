package grpcwrapper

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type RpcCall func(context.Context)

var defaultRpcCallTimeout = 5 * time.Second

func CallWithTimeOut(call RpcCall) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultRpcCallTimeout)
	defer cancel()
	call(ctx)
}

func DefaultDialKeepailveOpt() grpc.DialOption {
	var kacp = keepalive.ClientParameters{
		Time:                5 * time.Second,
		Timeout:             time.Second,
		PermitWithoutStream: true,
	}

	return grpc.WithKeepaliveParams(kacp)
}
