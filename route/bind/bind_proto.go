package bind

import (
	"github.com/golang/protobuf/proto"
	"github.com/slclub/easy/log"
	"github.com/slclub/easy/route/element"
	"reflect"
)

type BindProto struct {
	BindBase
}

func NewBindProto(m *element.PathMap) *BindProto {
	bd := (&BindProto{})
	//fmt.Printf("type := %v %+v \n", bd, bd)
	bd.pathMap = m //element.NewPahtMap()
	return bd
}

// 绑定消息ID 和消息
func (self *BindProto) Register(id element.MID, msgany any) {
	msg, ok := msgany.(proto.Message)
	if !ok {
		log.Fatal("protobuf router bind message register error")
	}
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		log.Fatal("protobuf router bind message pointer required")
	}

	i := self.pathMap.GetNewByMID(id)

	i.Type = msgType
	i.MID = id
	self.pathMap.Add(i)
}
