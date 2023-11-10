package bind

import (
	"github.com/slclub/easy/log"
	"github.com/slclub/easy/route/element"
	"reflect"
)

type BindJson struct {
	BindBase
}

func NewBindJson(m *element.PathMap) *BindJson {
	bd := &BindJson{}
	bd.pathMap = m
	return bd
}

// 绑定消息ID 和消息
func (self *BindJson) Register(id element.MID, msgany any) {
	if msgany == nil {
		log.Fatal("json router bind message register error")
	}
	msgType := reflect.TypeOf(msgany)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		log.Fatal("json router bind message pointer required")
	}

	i := self.pathMap.GetNewByMID(id)

	i.Type = msgType
	i.MID = id
	self.pathMap.Add(i)
}
