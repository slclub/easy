package main

import (
	"context"
	"github.com/slclub/easy/examples/rpc/helloworld"
	"github.com/slclub/easy/log"
	cgrpc "github.com/slclub/easy/rpc/cgrpc"
	"github.com/slclub/easy/rpc/etcd"
	"github.com/slclub/easy/vendors/option"
	"google.golang.org/grpc"
	"strings"
)

var (
	etcdAddr   string = "123.57.25.243:12379"
	namespace         = "easy"
	serverAddr        = "127.0.0.1:13001"
)

// server is used to implement helloworld.GreeterServer.
type hello struct {
	helloworld.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *hello) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	log.Info("Received: %v", in.GetName())
	return &helloworld.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	// setting ETCD service
	// eoption := &etcd.Option{}
	// eoption.Conv(etcdAddr)
	etcd.NewWithOption(option.OptionWith(nil).Default(
		option.OptionFunc(func() (string, any) {
			return "Endpoints", strings.Split(etcdAddr, ";")
		})),
	)

	// New 一个rpc 监听服务
	server := cgrpc.NewServer(option.OptionWith(&cgrpc.Config{
		PathName:  "server1",
		Addr:      serverAddr,
		Namespace: namespace,
		//TTL:       15,
	}).Default(
		option.DEFAULT_IGNORE_ZERO,
		option.OptionFunc(func() (string, any) {
			return "ID", "abluo"
		}),
		option.OptionFunc(func() (string, any) {
			return "AddrToClient", "127.0.0.1:18080"
		}),
	))

	// 绑定业务接口到 rpc服务
	// 可以被多次使用RegisterService，我们用的append
	server.RegisterService(
		func(server *grpc.Server) {
			helloworld.RegisterGreeterServer(server, &hello{})
		},
	)

	// 监听；如果您有主监听，那么可以用go 并发运行
	server.Serv()
}
