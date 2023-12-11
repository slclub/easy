package cgrpc

import (
	"github.com/slclub/easy/vendors/option"
	"google.golang.org/grpc/resolver"
)

type WatchInterface interface {
	Scheme() string
	Watch(target resolver.Target, cc GameClientClusterInterface, assignment option.Assignment) (Resolver, error)
}

type Resolver interface {
	GameClientClusterInterface
}

type GameClientClusterInterface interface {
	UpdateState([]*ServInfo) error
}
