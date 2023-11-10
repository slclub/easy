package encode

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/slclub/easy/route/element"
	"github.com/slclub/easy/vendors/encode"
	"reflect"
)

type Protobuf struct {
	encode.Code
	pathMap *element.PathMap
}

func NewProtobuf(m *element.PathMap) *Protobuf {
	return &Protobuf{
		pathMap: m,
	}
}

func (self *Protobuf) Unmarshal(data []byte) (any, error) {
	if len(data) < 2 {
		return nil, errors.New("protobuf data too short")
	}

	// id
	mid := self.Bytes2Uint(data)

	//logs.Debug("protobuf.go :: Unmarshal:: msgID [%v]", msgID)

	// msg
	i := self.pathMap.GetByMID(element.MID(mid))
	if i == nil {
		return nil, fmt.Errorf("message id %v not registered", mid)
	}

	msg := reflect.New(i.Type.Elem()).Interface()
	return msg, proto.UnmarshalMerge(data[2:], msg.(proto.Message))
}

func (self *Protobuf) Marshal(msg any) ([]byte, error) {
	msgType := reflect.TypeOf(msg)

	info := self.pathMap.GetByType(msgType)
	if nil == info {
		err := fmt.Errorf("message %s not registered", msgType)
		return nil, err
	}

	id := self.Uint2Bytes(uint16(info.MID))

	// data
	data, err := proto.Marshal(msg.(proto.Message))

	//logs.Debug("protobuf.go :: Unmarshal:: protobuf-len [%v] msgID [%v]", len(data), len(id))

	var byteTemp []byte
	byteTemp = append(id[:], data...)
	//logs.Debug("protobuf.go :: Unmarshal:: byteTemp [%v] bytetemp-len[%v]", byteTemp, len(byteTemp))
	return byteTemp, err
}
