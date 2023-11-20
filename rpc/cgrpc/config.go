package cgrpc

// contain configration of grpc listened and etcd paramters.
type Config struct {
	Name      string
	Addr      string
	TTL       int64
	Namespace string
}
