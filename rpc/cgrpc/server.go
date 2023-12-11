package cgrpc

/**

 */
import (
	"context"
	"fmt"
	"github.com/slclub/easy/log"
	"github.com/slclub/easy/rpc/etcd"
	"github.com/slclub/easy/vendors/option"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"net"
	"os"
	"strings"
	"time"
)

type Server struct {
	ID           string
	pathName     string
	Addr         string
	AddrToClient string
	TTL          int64
	Namespace    string
	server       *grpc.Server
	waiter
	greeterHandle []func(*grpc.Server)
}

func NewServer(assignment option.Assignment) *Server {
	ser := &Server{
		waiter:        waiter{make(chan os.Signal)},
		greeterHandle: []func(*grpc.Server){},
	}
	assignment.Target(ser)
	assignment.Default(option.OptionFunc(func() (string, any) {
		return "TTL", 15
	}))
	assignment.Apply()
	return ser
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
	self.waiter.close()
	self.Delete()
	//if i, ok := r.(syscall.Signal); ok {
	//	os.Exit(i)
	//} else {
	//	os.Exit(0)
	//}
}

func (self *Server) register() {
	ticker := time.NewTicker(time.Second * time.Duration(self.TTL))
	defer ticker.Stop()

	// do while
	for {
		self.registerSoon()
		<-ticker.C
	}
}

func (self *Server) registerSoon() {
	resp, err := etcd.EClient().Get(context.Background(), self.pathKey(self.pathName, self.Addr))
	if err != nil {
		log.Fatal("GRPC get server info from etcd error:%v", err)
		return
	}
	// 已注册
	if resp.Count > 0 {
		return
	}
	err = self.keepAlive()
	if err != nil {
		log.Fatal("GRPC keep alive error:%v", err)
	}
}

func (self *Server) keepAlive() error {
	// create lease
	leaseResp, err := etcd.EClient().Grant(context.Background(), self.TTL)
	if err != nil {
		return err
	}

	// register the service to etcd.
	key := self.pathKey(self.pathName, self.Addr)

	_, err = etcd.EClient().Put(context.Background(), key, self.pathValue(self.Addr, self.AddrToClient, self.GetID()), clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return err
	}
	log.Info("GRPC havs successfully registered it to etcd %v", key)

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

func (self *Server) pathValue(args ...string) string {
	value := ""
	for i, v := range args {
		if i == 0 {
			value += v
			continue
		}
		value += ";" + v
	}
	//log.Fatal("-------------- %v", value)
	return value
}

func (self *Server) GetID() string {
	if self.ID == "" {
		return self.Addr
	}
	return self.ID
}

func (self *Server) pathKey(args ...string) string {
	if len(args) == 0 {
		return ""
	}

	return "/" + self.Namespace + "/" + strings.Join(args, "/")
}

func (self *Server) Delete() {
	etcd.EClient().Delete(context.Background(), self.pathKey(self.pathName, self.Addr))
}
