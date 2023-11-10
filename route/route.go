package route

import (
	"github.com/slclub/easy/route/bind"
	"github.com/slclub/easy/route/element"
	"github.com/slclub/easy/route/encode"
	"github.com/slclub/easy/typehandle"
)

/**
 * r :=  NewRouter()
 *
 * r.Binding(bind.NewBindJson(r.PathMap()), encode.NewJson(r.PathMap()))
 */
type Route struct {
	DistributeCollection
	pathMap *element.PathMap
}

func NewRouter() Router {
	r := &Route{
		DistributeCollection: DistributeCollection{},
		pathMap:              element.NewPahtMap(),
	}
	binder := bind.NewBindProto(r.PathMap())
	r.Binding(encode.NewProtobuf(r.PathMap()), binder)
	return r
}

/** 注册路由
 *  The r is a Router
 *  r.Register(MSG_LOGIN_REQ, &JsonLoginReq{}, loginRequest) // with handle registory
 * 		if it is no handle needed can replace with nil
 *  Just register messages and MessageID
 *  	r.Binder().Register(MSG_LOGIN_REQ, &JsonLoginReq{}) //
 *  Just register pathMap with handler and message id.
 *  	r.Binder().RegisterHandle(id element.MID, handle typehandle.HandleMessage)
 */
func (self *Route) Register(ID element.MID, msg any, handle typehandle.HandleMessage) {
	self.Binder().Register(ID, msg)
	self.Binder().RegisterHandle(ID, handle)
}

// Automatic excuted it in your Server
// The binder achieve the element.RouteExecuter
func (self *Route) Route(msg any, ag any) error {
	return self.Binder().Route(msg, ag)
}

// 获取路由表
// 将消息ID 和消息结构体 定义在此表中
// 用 Binder.Register 等方法注册
func (self *Route) PathMap() *element.PathMap {
	return self.pathMap
}
