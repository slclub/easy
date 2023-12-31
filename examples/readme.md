# Program examples
用easy 实现的 服务端 和客户端样例

## Simple

### simple

比较简单的源码样例；

这是一个简单的服务端，你可以直接拿它做项目，扩展开发即可。

- 最基本的easy框架使用
- 简单的游戏架构，不包含数据层；
- simple代码以极简化为主，项目扩展要 结构化一些
- 目录多是空的

#### simple 目录介绍

```go
--conf // 配置
--controller 控制器也就是解读消息的入口
    --callback 放置一些基本的回调函数，如链接创建，服务平滑关闭等
    --login 登陆模块
    --player 用户玩家
    --store 商铺
    --world 大世界相关
--initialize 初始化，工程启动执行一次；与运行时无关
--lservers 接入easy监听服务 l 是 ```listen``` 当让也可以接入其他的监听服务
--message 消息定义
--models 数据模型，尽量只有数据结构的定义，和基本验证
--services 游戏逻辑存放区域，主要的逻辑都可以放在这里
--vendors 您项目的一些必要基础功能性的包，或者接入第三方包（且这个包需要配置等）；// 并非是替代 go mod
    注意：
    go mod 中也有一个verdor 且会产生vendor文件夹
    我们这里的vendors 仅仅是common 通用，基础，标准等的意思
    这里的包之间互相依赖也少，或者说机会是无
    比较大（功能性）的包引入后，总需要配置一些东西，甚至和自己的配置参数相关，那么放在这里改造一下（符合工程写法，结构要求等）就比较合适了
```


### simple_client

明显是simple 对应的客户端测试代码

- run

nws  = the number of websocket connections

ntcp = the number of tcp connections

```go
go build && ./simple_client -nws=1000 -ntcp=0
```

---

## RPC

我们使用的grpc和etcd，构成了一个完整的服务发现模式。可以轻松实现稳定分层，分布式架构的服务。

代码中出现的namespace 与 scheme 是同一个概念。

这介绍的简单教程，详细运行代码可以看example/rpc 的源码。

[ETCD with docker](https://github.com/slclub/easy/blob/master/docs/etcd.md)

### helloworld
 
这个子package 是官方的一个 接口定义的例子。也是最简单最easy的一个例子。
这里对接的grpc服务与你自己的业务代码结合的通道接口。与MVC 中 C 是一个位置，
它属于业务应用类代码。

### server

grpc的服务端，使用easy.rpc 只需要简短的代码就可以构筑，rpc应用服务端。

- run 

运行前先修改下etcd 的地址；集群多个etcd地址用分号隔开即可。为了运行命令简便，这里并没有使用flag等。

先跳转到 ```cd examples/rpc/server```

运行 ```go build && ./server ```

- 服务端配置ETCD

```go
    // 配置ETCD服务
    eoption := &etcd.Option{}
    eoption.Conv(etcdAddr)
    etcd.NewWithOption(eoption)
	
```

- 服务端配置grpc

业务接口注册需要在Serv()监听之前去注册。

```go	
    // New 一个rpc 监听服务
    server := cgrpc.NewServer(&cgrpc.Config{
        Name:      "server1",
        Addr:      serverAddr,
        Namespace: namespace,
    })
    
    // 绑定业务接口到 rpc服务
    // 可以被多次使用RegisterService，我们用的append
    server.RegisterService(
        func(server *grpc.Server) {
            helloworld.RegisterGreeterServer(server, &hello{})
        },
    )
    
    // 监听；如果您有主监听接口，那么可以用go 并发运行
    server.Serv()
```

- 业务接口

```go
// server is used to implement helloworld.GreeterServer.
type hello struct {
    helloworld.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *hello) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
    log.Info("Received: %v", in.GetName())
    return &helloworld.HelloReply{Message: "Hello " + in.GetName()}, nil
}
```

### client

grpc的客户端；load balance 是在客户端实现的.

- run

这里与server 是一样的，修改下etcd 地址即可

跳转到项目目录：
```cd examples/rpc/client```

运行：
```go build && ./client```

- 变量

```go
var (
	etcdAddr  string = "127.0.0.1:12379"
	namespace        = "easy"
)
```

- ETCD 初始化

与服务端一样的

```go
    // plan 1
    //eoption := &etcd.Option{}
    //eoption.Conv(etcdAddr)
    //etcd.NewWithOption(eoption)
	
    // plan2
    etcd.NewWithOption(option.OptionWith(nil).Default(
        option.OptionFunc(func() (string, any) {
        return "Endpoints", strings.Split(etcdAddr, ";")
        })),
    )
```

- grpc 客户端配置

代码中client对象是需要全局暴漏（public）。handle调用部分可以写在任意地方，是与原生grpc调用一样的，只要在ClientConn的有效生命周期内通信调用即可。

使用匿名函数封装可以避免暴漏全局变量，但为了兼容grpc使用习惯，就没有去封装。

```go
    client := cgrpc.NewClient(option.OptionWith(&struct {
        Name      string
        Namespace string
    }{"server1", namespace}))
	
    client.Start()
    
    // do your things
    handle(client.ClientConn)
    // just for test
    client.Wait()
    
    // close
    client.Close()
```

- handle 业务处理

handle 写这十几行，其实为了压力测试，持续性测试等。


```go
func handle(clientConn grpc.ClientConnInterface) {
    c := helloworld.NewGreeterClient(clientConn)
    ticker := time.NewTicker(2 * time.Second)
    i := 1
    for range ticker.C {
    
        resp1, err := c.SayHello(
            context.Background(),
            &helloworld.HelloRequest{Name: fmt.Sprintf("xiaoming-%d", i)},
        )
        if err != nil {
            log.Fatal("SayHello call error：%v", err)
            continue
        }
        log.Info("SayHello Response：%s\n", resp1.Message)
    
        i++
    }

}
```

仅仅是这2句是用的grpc，也是核心调用，其他的是多余的
```go
    c := helloworld.NewGreeterClient(clientConn)
    resp1, err := c.SayHello(
        context.Background(),
        &helloworld.HelloRequest{Name: fmt.Sprintf("xiaoming-%d", i)},
    )
```

---
