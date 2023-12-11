package main

import (
	"context"
	"fmt"
	"github.com/slclub/easy/examples/rpc/helloworld"
	"github.com/slclub/easy/log"
	"github.com/slclub/easy/rpc/cgrpc"
	"github.com/slclub/easy/rpc/etcd"
	"github.com/slclub/easy/vendors/option"
	"google.golang.org/grpc"
	"strings"
	"time"
)

var (
	etcdAddr  string = "123.57.25.243:12379"
	namespace        = "easy"
)

func main() {
	// plan1
	//eoption := &etcd.Option{}
	//eoption.Conv(etcdAddr)
	//etcd.NewWithOption(option.OptionWith(eoption).Default(option.DEFAULT_IGNORE_ZERO))

	// plan2 using the default value setting function.
	etcd.NewWithOption(option.OptionWith(nil).Default(
		option.OptionFunc(func() (string, any) {
			return "Endpoints", strings.Split(etcdAddr, ";")
		}),
	))

	client := cgrpc.NewClient(option.OptionWith(&struct {
		PathName  string
		Namespace string
	}{"server1", namespace}))

	client.Start()

	// do your things
	handle(client.ClientConn)
	// just for test
	client.Wait()

	// close
	client.Close()
}

func optionChoice() any {
	opt := struct {
		Endpoints []string
	}{}
	opt.Endpoints = strings.Split(etcdAddr, ";")
	return opt
}

func handle(clientConn grpc.ClientConnInterface) {
	c := helloworld.NewGreeterClient(clientConn)
	ticker := time.NewTicker(2 * time.Second)
	i := 1
	for range ticker.C {

		resp1, err := c.SayHello(
			context.Background(),
			&helloworld.HelloRequest{Name: fmt.Sprintf("xiaoming-%d", i)},
		)
		if err != nil {
			log.Fatal("SayHello call error：%v", err)
			continue
		}
		log.Info("SayHello Response：%s\n", resp1.Message)

		i++
	}

}
