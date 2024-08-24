package cgrpc

import (
	"context"
	"github.com/slclub/easy/log"
	"github.com/slclub/easy/vendors/option"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"strings"
)

type etcdWatcher struct {
	scheme  string
	cluster *GameRpcCluster
	client  *clientv3.Client
}

func NewGameEtcdWatcher(scheme string, client *clientv3.Client) *etcdWatcher {
	return &etcdWatcher{
		scheme: scheme,
		client: client,
	}
}

func (self *etcdWatcher) Scheme() string {
	return self.scheme
}

func (self *etcdWatcher) Watch(target resolver.Target, ggcc GameClientClusterInterface, assignment option.Assignment) (Resolver, error) {
	prefix := self.pathPrefix(target)
	resp, err := self.client.Get(context.Background(), prefix, clientv3.WithPrefix())

	if err != nil {
		log.Error("GRPC client get server list error:", err)
		return ggcc, err
	}
	servInfos := []*ServInfo{}
	for i := range resp.Kvs {
		o := resp.Kvs[i]
		infos := strings.Split(string(o.Value), ";")
		servInfos = append(servInfos, NewServInfo(infos))
	}
	ggcc.UpdateState(servInfos)

	// 监听服务列表
	watchChan := self.client.Watch(context.Background(), prefix, clientv3.WithPrefix())
	for respWathc := range watchChan {
		for _, ev := range respWathc.Events {
			o := ev.Kv
			info := strings.Split(string(o.Value), ";")
			switch ev.Type {
			case 0: // PUT
				if existsServInfo(servInfos, info[0]) {
					continue
				}
				servInfos = append(servInfos, NewServInfo(info))
				ggcc.UpdateState(servInfos)
				log.Info("ETCD has joined in a new server RPC:%v", string(o.Value))
			case 1: // DELETE
				if list, ok := removeServInfo(servInfos, strings.TrimPrefix(string(o.Key), prefix)); ok {
					servInfos = list
					ggcc.UpdateState(servInfos)
				}
				log.Info("ETCD has deleted a server len:%v key:%v addr:%v ", len(servInfos), string(o.Key), string(o.Value))
			}
		}
	}
	return ggcc, nil
}

func (self *etcdWatcher) pathPrefix(target resolver.Target) string {
	return "/" + target.URL.Scheme + target.URL.Path + "/"
}

var _ WatchInterface = &etcdWatcher{}

func existsServInfo(l []*ServInfo, addr string) bool {
	for i := range l {
		if l[i].AddrToRpc == addr {
			return true
		}
	}
	return false
}

func removeServInfo(s []*ServInfo, addr string) ([]*ServInfo, bool) {
	for i := range s {
		if s[i].AddrToRpc == addr {
			s[i] = s[len(s)-1]
			return s[:len(s)-1], true
		}
	}
	return nil, false
}
