package main

import (
	"simple_client/message/ID"
	"simple_client/message/json"
)

// websocket 服务消息注册
func InitWsRegister() {
	WsMgr.router.Register(ID.LOGIN_REQ, &json.LoginReq{}, nil)
	WsMgr.router.Register(ID.LOGIN_RES, &json.LoginRes{}, LoginRes)
}

// tcp 的服务 消息注册
func InitTCPRegister() {
	TCPMgr.router.Register(ID.LOGIN_REQ, &json.LoginReq{}, nil)
	TCPMgr.router.Register(ID.LOGIN_RES, &json.LoginRes{}, LoginResTCP)
}
