package etcd

// register servers into the ETCD.

import (
	"github.com/slclub/easy/log"
	"github.com/slclub/easy/vendors/option"
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

func NewWithOption(assignment option.Assignment) {
	ecli = NewClientWithOption(assignment)
}

func NewClientWithOption(assignment option.Assignment) *clientv3.Client {

	v3config := clientv3.Config{}
	assignment.Target(&v3config)
	// set target default value
	// please running it before Apply() method.
	assignment.Default(option.OptionFunc(func() (string, any) {
		return "DialTimeout", CONNECTION_ETCD_TIMEOUT_DEFAULT * time.Second
	}))

	assignment.Apply()

	cli, err := clientv3.New(v3config)
	if err != nil {
		log.Fatal("[ETCD] client created error")
		panic("[ETCD] client created error")
		return nil
	}
	return cli
}
