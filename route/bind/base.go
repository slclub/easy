package bind

import (
	"errors"
	"github.com/slclub/easy/nets/agent"
	"github.com/slclub/easy/route/element"
	"github.com/slclub/easy/typehandle"
	"github.com/slclub/easy/vendors/ants"
	"reflect"
)

/**
 * 绑定 控制器到，路由器中
 * 或许叫它  粘合器 更贴切一些；粘合 解码器 和路由解析能力
 * 它可以被自定义
 * 执行器与注册器 可以说是直接影响的；所以我们这里都定义到了一起，其实可以分开；我们也是为了 简化
 */
type BindBase struct {
	pathMap *element.PathMap
}

// 绑定消息ID 和消息
func (self *BindBase) Register(id element.MID, msgany any) {

}

// 绑定 handle 到 路由
func (self *BindBase) RegisterHandle(id element.MID, handle typehandle.HandleMessage) {
	i := self.pathMap.GetNewByMID(id)
	i.MID = id
	i.Handle = handle
	self.pathMap.Add(i)
}

func (self *BindBase) Route(msg any, ag any) error {
	msgType := reflect.TypeOf(msg)
	info := self.pathMap.GetByType(msgType)
	if info == nil {
		return errors.New("ROUTER NOT FOUND " + msgType.Kind().String())
	}
	if info.Handle == nil {
		return errors.New("ROUTER NOT FOUND Handle" + msgType.Kind().String())
	}
	switch val := ag.(type) {
	case agent.Agent:
		ants.Pool().Submit(func() {
			info.Handle(val, msg)
		})
	}
	return nil
}

var _ element.Binder = &BindBase{}
