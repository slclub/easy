package element

import (
	"github.com/slclub/easy/typehandle"
)

type Binder interface {
	// 绑定消息ID 和消息
	Register(id MID, msg any)
	// 绑定 handle 到 路由
	RegisterHandle(id MID, handle typehandle.HandleMessage)
	// 继承执行器
	RouteExecuter
}
