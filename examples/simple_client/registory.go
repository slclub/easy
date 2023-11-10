package main

import (
	"simple_client/message/ID"
	"simple_client/message/json"
)

func InitWsRegister() {
	WsMgr.router.Register(ID.LOGIN_REQ, &json.LoginReq{}, nil)
	WsMgr.router.Register(ID.LOGIN_RES, &json.LoginRes{}, LoginRes)
}
