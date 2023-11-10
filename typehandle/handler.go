package typehandle

import "github.com/slclub/easy/nets/agent"

// 消息处理控制器
type HandleMessage func(agent.Agent, any)

// 链接connection 控制器
type AgentHandle func(agent.Agent)
