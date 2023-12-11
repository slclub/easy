package cgrpc

// contain configration of grpc listened and etcd paramters.
type Config struct {
	PathName  string
	Addr      string
	TTL       int64
	Namespace string
}
