package main

import (
	"github.com/slclub/easy"
	"github.com/slclub/easy/log"
	"github.com/slclub/log8q"
	"simple/initialize"
	"simple/lservers"
)

func main() {
	log.LEVEL = log8q.ALL_LEVEL // 开放框架的全日志

	initialize.Init(func() {
		initialize.ListenPort = 15080
	})
	Start()
}

func Start() {
	easy.Serv(
		lservers.Server1(), // websocket 监听服务 可以有多个
		lservers.Server2(), // tcp 服务
	)
}
