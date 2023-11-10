package element

type RouteExecuter interface {
	// 路由分发消息  给 对应的handle
	Route(msg any, ag any) error
}
