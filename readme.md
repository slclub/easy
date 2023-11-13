# EASY  FRAMEWORK

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

Handle is a service listening and processing function.

```go
import (
    "github.com/slclub/easy/nets/agent"
    "reflect"
    "simple/vendors/log8q"
)

func HandleLogin(agent agent.Agent, arg any) {
    log8q.Log().Info("WS controller.Handle.Login info: ", reflect.TypeOf(arg).Elem().Name())
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

The ```Binding(plugins ...any)``` method  can bind binder and encodersto Router. 
So use this method to replace Binder and Encoder if you want to custom yourself.

example:

```go
r := Server1().Router()
r.Binding(bind.NewBindJson(r.PathMap()), encode.NewJson(r.PathMap()))
```

> ### Encoder

Encode/Decode transferred data. you can alose customize this component.
By default, we support two typs : json and protobuf.

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

It has two basic functions: binding Encoder and route functions.
It corresponds one-to-one with the encoder.

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
## Customize 
## Open Packages
## NETS
## Examples
## Building With Docker
## Contribution
## License
