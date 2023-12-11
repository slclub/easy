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

	client := cgrpc.NewGameGrpcCluster(option.OptionWith(&struct {
		PathName  string
		Namespace string
	}{"server1", namespace}))

	//client.Start()

	// do your things
	testHandle(client)

	// just for test
	client.Wait()

	// close
	client.Close()
}

func handleSay(clientConn grpc.ClientConnInterface) {
	c := helloworld.NewGreeterClient(clientConn)

	resp1, err := c.SayHello(
		context.Background(),
		&helloworld.HelloRequest{Name: fmt.Sprintf("xiaoming-%d", 9)},
	)
	if err != nil {
		log.Fatal("Say Hello error:%v", err)
		return
	}
	log.Info("SayHello Response：%s\n ", resp1.Message)

}

func testHandle(client *cgrpc.GameRpcCluster) {
	ticker := time.NewTicker(time.Second)
	for {
		<-ticker.C
		client.Handle(func(gconns []*cgrpc.GameGrpcConn) {
			log.Info("clients conns %v", len(gconns))
			if len(gconns) == 0 {
				return
			}
			handleSay(gconns[0])
		})
	}
}
