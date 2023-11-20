package etcd

// register servers into the ETCD.

import (
	"github.com/slclub/easy/log"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

const (
	CONNECTION_ETCD_TIMEOUT_DEFAULT = 15
)

var ecli *clientv3.Client

func EClient() *clientv3.Client {
	return ecli
}

func NewWithOption(option *Option) {
	var err error
	timeout := option.DialTimeout
	if timeout <= 0 {
		timeout = CONNECTION_ETCD_TIMEOUT_DEFAULT * time.Second
	}
	ecli, err = clientv3.New(clientv3.Config{
		Endpoints:   option.Endpoints,
		DialTimeout: timeout,
	})
	if err != nil {
		log.Fatal("[ETCD] client created error")
		return
	}
}
