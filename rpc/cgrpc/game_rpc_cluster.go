package cgrpc

import (
	"errors"
	"github.com/slclub/easy/log"
	"github.com/slclub/easy/vendors/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"net/url"
	"os"
	"strings"
)

type GameRpcCluster struct {
	gconns    []*GameGrpcConn
	watcher   WatchInterface
	pathName  string
	namespace string
	waiter
}

func NewGameGrpcCluster(assignment option.Assignment) *GameRpcCluster {
	g := &GameRpcCluster{
		gconns:  []*GameGrpcConn{},
		watcher: NewGameEtcdWatcher("easy"),
		waiter: waiter{
			wait: make(chan os.Signal),
		},
	}
	defer g.Start()

	if assignment == nil {
		return g
	}
	assignment.Target(g)
	assignment.Apply()

	g.watcher = NewGameEtcdWatcher(g.namespace)
	return g
}

func (self *GameRpcCluster) Register(whatcher WatchInterface) {
	self.watcher = whatcher
}

func (self *GameRpcCluster) Start() {
	go self.ResolveNow()
}

func (self *GameRpcCluster) ResolveNow() {
	self.watcher.Watch(
		resolver.Target{URL: url.URL{
			Scheme: self.namespace,
			Path:   strings.Replace("/"+self.pathName, "//", "/", 1)}},
		self, nil)
}

func (self *GameRpcCluster) UpdateState(servInfos []*ServInfo) error {
	if servInfos == nil {
		self.releaseConns()
		return errors.New("GAME.GRPC closed all grpc clients")
	}
	gconns := []*GameGrpcConn{}
	newInfos := []*ServInfo{}
	oldConnIndexs := make([]int, len(self.gconns))
	for k, info := range servInfos {

		j := -1
		for i, ggc := range self.gconns {
			if info.AddrToRpc == ggc.Info.AddrToRpc {
				gconns = append(gconns, self.gconns[i])
				j = i
				break
			}
		}
		if j == -1 {
			newInfos = append(newInfos, servInfos[k])
		} else {
			oldConnIndexs[j] = 1
		}
	}

	// release 旧的链接
	for i, v := range oldConnIndexs {
		if v == 1 {
			continue
		}
		defer self.releaseOne(self.gconns[i])
	}

	for _, info := range newInfos {
		obj := newGameRpcConn(option.OptionWith(info))
		err := obj.Dial(obj.Info.AddrToRpc, grpc.WithDefaultServiceConfig(`{}`),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatal("GAME.GRPC disconnection name:%v addr:%v client:%v err:%v", info.Name, info.AddrToRpc, info.AddrToClient, err)
			continue
		}
		gconns = append(gconns, obj)
	}

	self.gconns = gconns
	return nil
}

func (self *GameRpcCluster) GetGameClientConns() []*GameGrpcConn {
	return self.gconns
}

func (self *GameRpcCluster) Handle(fn func(gconns []*GameGrpcConn)) {
	fn(self.gconns)
}

func (self *GameRpcCluster) releaseOne(ggc *GameGrpcConn) {
	log.Error("GAME.GRPC disconnection name:%v addr:%v client:%v", ggc.Info.Name, ggc.Info.AddrToRpc, ggc.Info.AddrToClient)
	ggc.Close()
}

func (self *GameRpcCluster) releaseConns() {
	for _, cc := range self.gconns {
		cc.Close()
	}
	self.gconns = []*GameGrpcConn{}
}

func (self *GameRpcCluster) Close() {
	self.releaseConns()
}
