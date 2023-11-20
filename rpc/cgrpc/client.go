package cgrpc

import (
	"context"
	"github.com/slclub/easy/log"
	"github.com/slclub/easy/rpc/etcd"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type Client struct {
	Name        string
	Namespace   string
	DialBalance string // default:round_robin
	ClientConn  *grpc.ClientConn
	wait        chan os.Signal
	url         string
}

// resolver
type etcdResolver struct {
	Namespace string
	conn      resolver.ClientConn
}

// --------------------------------------------------
// client
// --------------------------------------------------
/**
example :
	c := proto.NewGreeterClient(Client.ClientConn)
	c.SyaHello()
*/

func NewClient(name, space, balance string) *Client {
	if balance == "" {
		balance = "round_robin"
	}
	return &Client{
		Name:        name,
		Namespace:   space,
		DialBalance: balance,
		wait:        make(chan os.Signal),
	}
}

// you should run it with go command.
func (self *Client) Start() {
	// new a resolver customed.
	rs := newResolver(self.Namespace)
	resolver.Register(rs)
	var err error

	// init default dial url
	self.DialUrl(rs.Scheme() + "://author/" + self.Name)

	// get conn
	self.ClientConn, err = self.Dial(
		self.DialUrl(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"`+self.DialBalance+`"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatal("GRPC connection error: %v", err)
		panic(any("err"))
		return
	}
}

// as same as grpc.Dial.
func (self *Client) Dial(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	var err error
	self.ClientConn, err = grpc.Dial(target, opts...)
	return self.ClientConn, err
}

func (self *Client) Close() {
	self.ClientConn.Close()
	defer signal.Stop(self.wait)
}

func (self *Client) Wait() {
	signal.Notify(self.wait, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	<-self.wait
}

func (self *Client) DialUrl(urls ...string) string {
	if len(urls) == 0 {
		return self.url
	}
	self.url = urls[0]
	return self.url
}

// --------------------------------------------------
// resolver
// --------------------------------------------------

func newResolver(scheme string) resolver.Builder {
	return &etcdResolver{Namespace: scheme}
}

func (self *etcdResolver) Scheme() string {
	return self.Namespace
}

func (self *etcdResolver) ResolveNow(rn resolver.ResolveNowOptions) {
	log.Info("resolver.ResolveNow")
}

func (self *etcdResolver) Close() {
	log.Info("resolver.Close")
}

func (self *etcdResolver) Build(target resolver.Target, clientConn resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	self.conn = clientConn
	go self.watch("/" + target.URL.Scheme + target.URL.Path + "/")
	return self, nil
}

func (self *etcdResolver) watch(pathPrefix string) {

	var addrList []resolver.Address
	resp, err := etcd.EClient().Get(context.Background(), pathPrefix, clientv3.WithPrefix())

	if err != nil {
		log.Fatal("GRPC client get server list error:", err)
		return
	}
	for i := range resp.Kvs {
		o := resp.Kvs[i]
		addrList = append(addrList, resolver.Address{Addr: strings.TrimPrefix(string(o.Key), pathPrefix)})
	}

	state := resolver.State{Addresses: addrList}
	self.conn.UpdateState(state)

	// 监听服务列表
	watchChan := etcd.EClient().Watch(context.Background(), pathPrefix, clientv3.WithPrefix())
	for respWathc := range watchChan {
		for _, ev := range respWathc.Events {
			addr := strings.TrimPrefix(string(ev.Kv.Key), pathPrefix)
			switch ev.Type {
			case 0: // PUT
				if exists(addrList, addr) {
					continue
				}
				addrList = append(addrList, resolver.Address{Addr: addr})
				state1 := resolver.State{Addresses: addrList}
				self.conn.UpdateState(state1)
				log.Info("ETCD has joined in a new server")
			case 1: // DELETE
				if list, ok := remove(addrList, addr); ok {
					addrList = list
					state1 := resolver.State{Addresses: addrList}
					self.conn.UpdateState(state1)
				}
				log.Info("ETCD logout a server addr:", addr)
			}
		}
	}
}

// --------------------------------------------------
// functions
// --------------------------------------------------

func exists(l []resolver.Address, addr string) bool {
	for i := range l {
		if l[i].Addr == addr {
			return true
		}
	}
	return false
}

func remove(s []resolver.Address, addr string) ([]resolver.Address, bool) {
	for i := range s {
		if s[i].Addr == addr {
			s[i] = s[len(s)-1]
			return s[:len(s)-1], true
		}
	}
	return nil, false
}
