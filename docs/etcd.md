
## ETCD build

### ETCD config

```go
name: single
data-dir: /etcd-data
listen-client-urls: http://0.0.0.0:12379
advertise-client-urls: http://0.0.0.0:12379
listen-peer-urls: http://0.0.0.0:12380
initial-advertise-peer-urls: http://0.0.0.0:12380
initial-cluster: single=http://0.0.0.0:12380
initial-cluster-token: etcd-cluster-token
initial-cluster-state: new
log-level: info
logger: zap
log-outputs: stderr
```

### Docker-compose 

- Single

```go
version: '3.8'

services:
  single:
    image: quay.io/coreos/etcd
    restart: on-failure
    entrypoint: ["/usr/local/bin/etcd", "--config-file", "/tmp/etcd/conf/etcd.yml"]
    ports:
      - "12379:12379"
      - "12380:12380"
    environment:
      ETCDCTL_API: 3
    volumes:
      - type: bind
        source: /dbstore/etcd/
        target: /tmp/etcd
```

- Cluster

如果使用需要先创建对应的目录和配置文件，参考 single 版本的加即可。

```go
version: '3.8'

services:
  etcd-1:
    image: gcr.io/etcd-development/etcd:v3.4.25
    entrypoint: [ "/usr/local/bin/etcd", "--config-file", "/tmp/etcd/conf/etcd.yml" ]
    ports:
      - "23791:2379"
    environment:
      ETCDCTL_API: 3
    volumes:
      - type: bind
        source: /tmp/etcd/cluster/etcd1
        target: /tmp/etcd
    networks:
      etcd-net:
        ipv4_address: 172.25.0.101

  etcd-2:
    image: gcr.io/etcd-development/etcd:v3.4.25
    entrypoint: [ "/usr/local/bin/etcd", "--config-file", "/tmp/etcd/conf/etcd.yml" ]
    ports:
      - "23792:2379"
    environment:
      ETCDCTL_API: 3
    volumes:
      - type: bind
        source: /tmp/etcd/cluster/etcd2
        target: /tmp/etcd
    networks:
      etcd-net:
        ipv4_address: 172.25.0.102

  etcd-3:
    image: gcr.io/etcd-development/etcd:v3.4.25
    entrypoint: [ "/usr/local/bin/etcd", "--config-file", "/tmp/etcd/conf/etcd.yml" ]
    ports:
      - "23793:2379"
    environment:
      ETCDCTL_API: 3
    volumes:
      - type: bind
        source: /tmp/etcd/cluster/etcd3
        target: /tmp/etcd
    networks:
      etcd-net:
        ipv4_address: 172.25.0.103

networks:
  etcd-net:
    driver: bridge
    ipam:
      config:
        - subnet: 172.25.0.0/16
          gateway: 172.25.0.1
```

[CSDN原文参考](https://www.cnblogs.com/NezhaYoung/p/17347450.html)


### DockerImage


```go
docker pull quay.io/coreos/etcd
```
