
## 队列最小单元
RabbitMQ

## 同一个微服务 实现的工作模式
```json
大多是为了训练，demo
实际业务中，用消息队列，多是 不同的微服务代表一个角色
```

## Work 模式 可以用RabbitMQ 组织成 工作模式
```json
RabbitMQ
  相同 Exchange
  相同 Queue

一个微服务 为Product
其他N个微服务 为Worker Consume 接收者
```

## 发布订阅 模式
```json
RabbitMQ
  相同 Exchange
  不同 Queue

一个微服务微 Product
其他N个微服务 的Queue Name 各自不相同， 且为Consume
```



