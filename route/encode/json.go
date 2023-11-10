package encode

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/slclub/easy/route/element"
	"github.com/slclub/easy/vendors/encode"
	"reflect"
)

type Json struct {
	encode.Code
	pathMap *element.PathMap
}

func NewJson(m *element.PathMap) *Json {
	return &Json{
		pathMap: m,
	}
}

func (self *Json) Unmarshal(data []byte) (any, error) {
	if len(data) < 2 {
		return nil, errors.New("EASY DECODE JSON data too short")
	}
	mid := self.Bytes2Uint(data)
	i := self.pathMap.GetByMID(element.MID(mid))
	if i == nil {
		return nil, fmt.Errorf("message id %v not registered", mid)
	}

	msg := reflect.New(i.Type.Elem()).Interface()
	err := json.Unmarshal(data[2:], &msg)
	return msg, err
}

func (self *Json) Marshal(msg interface{}) ([]byte, error) {
	msgType := reflect.TypeOf(msg)
	info := self.pathMap.GetByType(msgType)
	if nil == info {
		err := fmt.Errorf("message %s not registered", msgType)
		return nil, err
	}

	id := self.Uint2Bytes(uint16(info.MID))

	data, err := json.Marshal(msg)
	var byteTemp []byte
	byteTemp = append(id[:], data...)
	//logs.Debug("protobuf.go :: Unmarshal:: byteTemp [%v] bytetemp-len[%v]", byteTemp, len(byteTemp))
	return byteTemp, err
}
