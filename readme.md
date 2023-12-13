# EASY  FRAMEWORK
[![license](https://img.shields.io/badge/License-MIT_2023-blue)](https://github.com/slclub/easy/blob/master/LICENSE)

## Summary

A light net framework writen with golang.  You can use it to write games,web,or chat and so on.

Easy framework biggest advantage is light, flexibility and customeize. you can change some component of Easy to adapt to your needs.

It includes various Network Transport Protocol such as TCP,websocket,HTTP. The TCP transmission rule can be customed too.
so,in addition to using default rules, you can also choose your own rules, but you need to implement the interface of Rule.

## Directory
- [Installation](#Installation)
- [Quick Start](#quick-start)
- [Router for keepalive](#Router-for-keepalive)
- [ListenServers](#ListenServers)
- [Customize](#Customize)
- [Open Packages](#Open-packages)
- [NETS](#NETS)
- [RPC and discovery](#RPC)
- [Examples](#Examples)
- [Building With Docker](#Building-With-Docker)
- [Contribution](#Contribution)
- [License](#License)

## Installation
```go
go get github.com/slclub/easy
```

If you want the more detaild tutorial. Please refer to the sample testing project in the examples directory .

## Quick Start

>### Start Listen Servers

- import
```go

import (
    "github.com/slclub/easy/nets/agent"
    "github.com/slclub/easy/servers"
    "github.com/slclub/easy/typehandle"
)
```

- Set

Registe server with your configruation.Here we use the websocket server(servers.WSServer) as an example.

```go
server1 = servers.NewWSServer()

server1.Init(&agent.Gate{
    Addr:            ":18080",
    Protocol:        typehandle.ENCRIPT_DATA_JSON,
    PendingWriteNum: 2000,
    LittleEndian:    true,
    MaxConnNum:      2000,
})
```

- Start

Please use the easy.Serv to start all of your listening servers. 
you can start more than one listening server.Each server has its owned network protocol.

```go

func Start() {
    easy.Serv(
        lservers.Server1(), // websocket 监听服务 可以有多个
        //lservers.Server2(), // tcp 服务
    )
}

```

>### registe
Registe message id and message struct and handle to the router of servers

Each server has a corresponding router.

```go
// "github.com/slclub/easy/typehandle"
r1 := lservers.Server1().Router()
r1.Register(ID.LOGIN_RES, &json.LoginRes{}, nil)
r1.Register(ID.LOGIN_REQ, &json.LoginReq{}, typehandle.HandleMessage(login.HandleLogin))
```

>### handle 

Handle is a listening service  and processing function.
In the other word it's a controller(The C of MVC).

It is an entry function of logical business.

- First Argument

It is a Agent object. You just need to understand  like a link.

- Second Argument

It is a message that you defined. 

Please used the Pointer type of Golang.

- example

```go
import (
    "github.com/slclub/easy/nets/agent"
    "reflect"
    "simple/vendors/log8q"
)

func HandleLogin(agent1 agent.Agent, arg any) {
    log8q.Log().Info("WS controller.Handle.Login info: ", reflect.TypeOf(arg).Elem().Name())
    // agent1.WriteMsg(nil)
}

func HandleLoginTcp(agent2 agent.Agent, arg any) {
    log8q.Log().Info("TCP controller.Handle.Login info: ", reflect.TypeOf(arg).Elem().Name())
}
```


## Router for keepalive
Router is the link between servers and handle. All parts of the router are implemented using interfaces.
so,you can custome it by your self. especially for Binder and Encoder that they are related to rules and data tranmited .

>### Router

- definition

```go
type Router interface {
    element.Distributer
    Register(ID element.MID, msg any, handle typehandle.HandleMessage)
    Route(msg any, ag any) error
    PathMap() *element.PathMap
}
```

```go

// 为route 绑定插件
type Distributer interface {
    DistributePlug
    // 绑定 解码器=typehandle.Encoder ; 绑定器 = element.Binder
    // Binding(encoder typehandle.Encoder, binder Binder)
    Binding(plugins ...any)
}

type DistributePlug interface {
    Binder() Binder
    Encoder() encode.Encoder
}
```

- Binding method

The ```Binding(plugins ...any)``` method  can bind binder and encoder to Router. 
So use this method to replace Binder and Encoder if you want to custom yourself.

example:

```go
r := Server1().Router()
r.Binding(bind.NewBindJson(r.PathMap()), encode.NewJson(r.PathMap()))
```

- PathMap

This method return the storage for routing. Binder and Encoder should use it together. 


> ### Encoder

Encode/Decode transferred data. you can also customize this component.
By default, we support two types : json and protobuf.

- Interface

```go
// stream 解析器 encode decode操作
type Encoder interface {
	Unmarshal(data []byte) (any, error)
	Marshal(msg any) ([]byte, error)
	LittleEndian(...bool) bool
}
```

- JSON

- Protobuf

>### Binder

It has three basic functions: binding Encoder, binding handle and route functions.
It corresponds one-to-one with the encoder. the methods of Register, RegisterHandle and Route will be called by Router

- Interface

```go
type Binder interface {
    // 绑定消息ID 和消息
    Register(id MID, msg any)
    // 绑定 handle 到 路由
    RegisterHandle(id MID, handle typehandle.HandleMessage)
    // 继承执行器
    RouteExecuter
}


type RouteExecuter interface {
    // 路由分发消息  给 对应的handle
    Route(msg any, ag any) error
}
```



## ListenServers
I had defined various servers followed by network protocol.

Any of them can be New. so It is very convenient to use more than one servers in your project.
Any instance servers can use different components. Also they can be customed.

The Listen Server like a container of components, 
enable smooth collaboration among various components to complete tasks.

>### interface
All of the servers should implement the interface of ```ListenServer```.
They will all be used uniformly in the easy.Serv().

```go

import (
    "github.com/slclub/easy/nets/agent"
    "github.com/slclub/easy/route"
)

type ListenServer interface {
    Init(*agent.Gate)
    OnInit()
    Router() route.Router
    Start()
    Hook() *hookAgent // agent 链接 回调
    //OnClose(func())
    Close()
}
```

>### websocket

- Name

```
WSServer
```

- Create

```go
// import github.com/slclub/easy/servers
servers.NewWSServer()
```


>### TCP

- Name

```
TCPServer
```

- Create

```go
// import github.com/slclub/easy/servers
servers.NewTCPServer()
```

>### WEB

## Customize 

### components:

- Listen Server
- Router
- route.Binder
- route.Encoder
- Agent Of Net
- FromConnReadWriter of conns
- Conn of conns


## Open Packages

[option package](https://github.com/slclub/easy/blob/master/docs/option.md)

[aoi package](https://github.com/slclub/easy/blob/master/docs/aoi.md)

[events package](https://github.com/slclub/easy/blob/master/docs/events.md)


---- 

## NETS

>### agent.Agent

When you send or recvive messages from client,you will need to use this.

The first argment of the handle that registed in your Router is agent.Agent type.

You will bind it to your own entity.

```go
type Agent interface {
    WriteMsg(msg any)
    LocalAddr() net.Addr
    RemoteAddr() net.Addr
    Close()
    Destroy()
    UserData() any
    SetUserData(data any)
    LoopRecv(handle AgentHandle)
}
```

- Send Message

```Agent.WriteMsg(msg any)```

- Recive Message

Will be called by Router Excuter. So we do not need to care about it.

The handle functions is the processor that receives messages.


>### conns.Conn

It is does not matter with the logical business. Just open for framework customed.

```go
type Conn interface {
    ReadMsg() ([]byte, error)
    WriteMsg(args []byte) error
    LocalAddr() net.Addr
    RemoteAddr() net.Addr
    Close()
    Destroy()
    Done() chan struct{}
    GetOption() *Option
}
```

## RPC

The easy framework program integrates GRPC and ETCD. constitutes a complete service discovery and distributed RPC communication server architecture.

easy

[tutorials and example](https://github.com/slclub/easy/tree/master/examples#RPC)

[ETCD with docker](https://github.com/slclub/easy/blob/master/docs/etcd.md)

[grpc-go 官网](https://grpc.io/docs/languages/go/)

[etcd 官网](https://etcd.io/docs/)

docker image (```quay.io/coreos/etcd```). Ofcourse, you can used any image according to your preferences.  

## Examples

[Detail](https://github.com/slclub/easy/tree/master/examples/)

## Building With Docker
Under Construction.
## Contribution
## License

Copyright (c) 2023 许亚军

[MIT](https://github.com/slclub/easy/blob/master/LICENSE)
