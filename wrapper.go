package grpcwrapper

import (
	"fmt"
	"net"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"

	"grpcwrapper/middleware"
)

var (
	NotInitError error = errors.New("server not init")
)

// 服务实例
type Engine struct {
	Server       *grpc.Server                   // 原始grpc服务实例
	UnaryChains  []grpc.UnaryServerInterceptor  // unary拦截器
	StreamChains []grpc.StreamServerInterceptor // stream拦截器
	opts         []grpc.ServerOption            // 附加选项
}

func (e *Engine) lazyInit() {
	if e.UnaryChains == nil {
		e.UnaryChains = make([]grpc.UnaryServerInterceptor, 0)
	}
	if e.StreamChains == nil {
		e.StreamChains = make([]grpc.StreamServerInterceptor, 0)
	}
}

// NOTE：
//	1. 只有在InitServer之前调用才会生效
func (e *Engine) Use(umiddleware grpc.UnaryServerInterceptor, smiddleware grpc.StreamServerInterceptor) {
	e.lazyInit()
	if umiddleware != nil {
		e.UnaryChains = append(e.UnaryChains, umiddleware)
	}
	if smiddleware != nil {
		e.StreamChains = append(e.StreamChains, smiddleware)
	}
}

func (e *Engine) InitServer() error {
	e.opts = append(e.opts, grpc.UnaryInterceptor(ChainUnary(e.UnaryChains)))
	e.opts = append(e.opts, grpc.StreamInterceptor(ChainStream(e.StreamChains)))
	e.Server = grpc.NewServer(e.opts...)

	return nil
}

func (e *Engine) InitTLSServer(cert string, key string) error {
	creds, cerr := credentials.NewServerTLSFromFile(cert, key)
	if cerr != nil {
		return errors.Wrap(cerr, "创建credentials失败")
	}
	e.opts = append(e.opts, grpc.Creds(creds))
	e.opts = append(e.opts, grpc.UnaryInterceptor(ChainUnary(e.UnaryChains)))
	e.opts = append(e.opts, grpc.StreamInterceptor(ChainStream(e.StreamChains)))
	e.Server = grpc.NewServer(e.opts...)

	return nil
}

func (e *Engine) Run(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("监听地址[%s]失败", addr))
	}

	if serr := e.Server.Serve(lis); serr != nil {
		return errors.Wrap(serr, "启动服务异常")
	}
	return nil
}

func (e *Engine) GetRawServer() *grpc.Server {
	return e.Server
}

func (e *Engine) RunTLS(addr string, cert string, key string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("监听地址[%s]失败", addr))
	}
	if serr := e.Server.Serve(lis); serr != nil {
		return errors.Wrap(serr, "启动服务异常")
	}
	return nil
}

func (e *Engine) UseDefaultKeepaliveOpt() {
	var kaep = keepalive.EnforcementPolicy{
		MinTime:             5 * time.Second,
		PermitWithoutStream: true,
	}

	var kasp = keepalive.ServerParameters{
		MaxConnectionIdle: 15 * time.Minute,
		//MaxConnectionAge:      10 * time.Second,
		MaxConnectionAgeGrace: 5 * time.Second,
		Time:                  5 * time.Second,
		Timeout:               1 * time.Second,
	}

	ep := grpc.KeepaliveEnforcementPolicy(kaep)
	kp := grpc.KeepaliveParams(kasp)

	e.opts = append(e.opts, ep)
	e.opts = append(e.opts, kp)
}

func Default() *Engine {
	e := &Engine{}
	e.Use(middleware.UnaryServerRecoverInterceptor(),
		middleware.StreamServerRecoverInterceptor())
	return e
}
