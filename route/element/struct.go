package element

import (
	"github.com/slclub/easy/typehandle"
	"reflect"
)

type PathMap struct {
	idTypes   map[MID]*MessageRouterType
	typeTypes map[reflect.Type]*MessageRouterType
}

// 消息类型，消息ID，消息处理器
type MessageRouterType struct {
	MID    MID
	Type   reflect.Type
	Handle typehandle.HandleMessage
}

type MID int // message ID

// PathMap
func NewPahtMap() *PathMap {
	return &PathMap{
		idTypes:   make(map[MID]*MessageRouterType),
		typeTypes: make(map[reflect.Type]*MessageRouterType),
	}
}
func (self *PathMap) Add(mt *MessageRouterType) {
	self.idTypes[mt.MID] = mt
	self.typeTypes[mt.Type] = mt
}

func (self *PathMap) GetByMID(mid MID) *MessageRouterType {
	return self.idTypes[mid]
}

func (self *PathMap) GetNewByMID(mid MID) *MessageRouterType {
	va := self.idTypes[mid]
	if va != nil {
		return va
	}
	return new(MessageRouterType)
}

func (self *PathMap) GetByType(t reflect.Type) *MessageRouterType {
	return self.typeTypes[t]
}
