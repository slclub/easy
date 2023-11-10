package message

import (
	"simple/lservers"
	"simple/message/ID"
	"simple/message/json"
)

// 将不需要handle 处理的消息 尽量放在这里注册
// 可以将所有注册消息都放在这里也可以
func Init() {
	InitJson()
	InitProtobuf()
}

func InitJson() {
	r1 := lservers.Server1().Router()
	r1.Register(ID.LOGIN_RES, &json.LoginRes{}, nil)

	r2 := lservers.Server2().Router()
	r2.Register(ID.LOGIN_RES, &json.LoginRes{}, nil)
}

func InitProtobuf() {
	//r2 := lservers.SimpleServ1.Router()
}
