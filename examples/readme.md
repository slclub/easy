## Program examples
用easy 实现的 服务端 和客户端样例

>### simple

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


>### simple_client

明显是simple 对应的客户端测试代码