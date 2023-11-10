package controller

import (
	"github.com/slclub/easy/typehandle"
	"simple/controller/login"
	"simple/lservers"
	"simple/message/ID"
	"simple/message/json"
)

func InitBindingRoute() {
	r1 := lservers.Server1().Router()
	r1.Register(ID.LOGIN_REQ, &json.LoginReq{}, typehandle.HandleMessage(login.HandleLogin))
}

func InitBindingRouteServer2() {
	r2 := lservers.Server2().Router()
	r2.Register(ID.LOGIN_REQ, &json.LoginReq{}, typehandle.HandleMessage(login.HandleLoginTcp))
}
