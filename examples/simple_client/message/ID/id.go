package ID

import "github.com/slclub/easy/route/element"

/**
 * 一个消息ID 在一个easy 监听的服务里面 只能有一个绑定
 * 只对应一个  消息体 和 Handle(可选的用nil 代替)
 */
const (
	LOGIN_REQ element.MID = 1001 // 这里可以不用定义类型 现用现转就行
	LOGIN_RES             = 1002
)
