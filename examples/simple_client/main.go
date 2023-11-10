package main

import (
	"github.com/slclub/easy/client"
	"github.com/slclub/easy/typehandle"
)

func main() {

	// 初始化 客户端
	WsMgr = NewWsTestMgr(&client.Gate{
		Protocol:     typehandle.ENCRIPT_DATA_JSON,
		Addr:         ":18080",
		LittleEndian: true,
	})

	WsMgr.Init()
	InitWsRegister() // 注册消息ID 和路由
	WsMgr.Run()
}
