package cgrpc

/**

 */
import (
	"context"
	"fmt"
	"github.com/slclub/easy/log"
	"github.com/slclub/easy/rpc/etcd"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Server struct {
	Name          string
	Addr          string
	TTL           int64
	Namespace     string
	server        *grpc.Server
	wait          chan os.Signal
	greeterHandle []func(*grpc.Server)
}

func NewServer(conf *Config) *Server {
	if conf.TTL <= 0 {
		conf.TTL = 15
	}
	return &Server{
		Name:          conf.Name,
		Addr:          conf.Addr,
		TTL:           conf.TTL,
		Namespace:     conf.Namespace,
		wait:          make(chan os.Signal),
		greeterHandle: []func(*grpc.Server){},
	}
}

// suitable grpc option configration
func (self *Server) Serv(opts ...grpc.ServerOption) {
	listener, err := net.Listen("tcp", self.Addr)
	if err != nil {
		fmt.Println("GRPC Server start error:", err)
		return
	}
	defer listener.Close()

	self.server = grpc.NewServer(opts...)
	defer self.server.GracefulStop()

	// add service to the listen server.
	for _, handler := range self.greeterHandle {
		handler(self.server)
	}

	// register to etcd
	go self.register()

	//

	// start listening rpc server terminal.
	log.Info("GRPC server start the listening service of RPC.")
	err = self.server.Serve(listener)
	if err != nil {
		log.Fatal("GRPC clust Serv error:", err)
		return
	}
	log.Info("GRPC server start successful")
	self.exitHandle()

}

func (self *Server) RegisterService(caller func(server *grpc.Server)) {
	self.greeterHandle = append(self.greeterHandle, caller)
}

func (self *Server) exitHandle() {
	signal.Notify(self.wait, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	<-self.wait
	defer signal.Stop(self.wait)

	self.Delete()
	//if i, ok := r.(syscall.Signal); ok {
	//	os.Exit(i)
	//} else {
	//	os.Exit(0)
	//}
}

func (self *Server) register() {
	ticker := time.NewTicker(time.Second * time.Duration(self.TTL))
	key := self.pathKey(self.Name, self.Addr)
	defer ticker.Stop()
	for {
		<-ticker.C
		resp, err := etcd.EClient().Get(context.Background(), key)
		if err != nil {
			log.Fatal("GRPC get server info from etcd error:%v", err)
			continue
		}
		// 已注册
		if resp.Count > 0 {
			continue
		}
		err = self.keepAlive()
		if err != nil {
			log.Fatal("GRPC keep alive error:%v", err)
		}
	}
}

func (self *Server) keepAlive() error {
	// create lease
	leaseResp, err := etcd.EClient().Grant(context.Background(), self.TTL)
	if err != nil {
		return err
	}

	// register the service to etcd.
	key := self.pathKey(self.Name, self.Addr)
	_, err = etcd.EClient().Put(context.Background(), key, self.Addr, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return err
	}
	log.Info("GRPC havs successfully registered it to etcd")

	// keep alive with etcd
	channelLeaseAlive, err := etcd.EClient().KeepAlive(context.Background(), leaseResp.ID)
	if err != nil {
		return err
	}

	// clear keep alive channel
	// TODO read the response.
	go func() {
		for {
			<-channelLeaseAlive
		}
	}()
	return nil
}

func (self *Server) pathKey(args ...string) string {
	if len(args) == 0 {
		return ""
	}

	return "/" + self.Namespace + "/" + strings.Join(args, "/")
}

func (self *Server) Delete() {
	etcd.EClient().Delete(context.Background(), self.pathKey(self.Name, self.Addr))
}
