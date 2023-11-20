package main

import (
	"context"
	"github.com/slclub/easy/examples/rpc/helloworld"
	"github.com/slclub/easy/log"
	cgrpc "github.com/slclub/easy/rpc/cgrpc"
	"github.com/slclub/easy/rpc/etcd"
	"google.golang.org/grpc"
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
	eoption := &etcd.Option{}
	eoption.Conv(etcdAddr)
	etcd.NewWithOption(eoption)

	server := cgrpc.NewServer(&cgrpc.Config{
		Name:      "server1",
		Addr:      serverAddr,
		Namespace: namespace,
	})

	server.RegisterService(
		func(server *grpc.Server) {
			helloworld.RegisterGreeterServer(server, &hello{})
		},
	)

	server.Serv()
}
