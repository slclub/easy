package main

import (
	"context"
	"fmt"
	"github.com/slclub/easy/examples/rpc/helloworld"
	"github.com/slclub/easy/log"
	"github.com/slclub/easy/rpc/cgrpc"
	"github.com/slclub/easy/rpc/etcd"
	"google.golang.org/grpc"
	"time"
)

var (
	etcdAddr  string = "123.57.25.243:12379"
	namespace        = "easy"
)

func main() {
	eoption := &etcd.Option{}
	eoption.Conv(etcdAddr)
	etcd.NewWithOption(eoption)

	client := cgrpc.NewClient("server1", namespace, "")
	client.Start()

	// do your things
	handle(client.ClientConn)
	// just for test
	client.Wait()

	// close
	client.Close()
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
