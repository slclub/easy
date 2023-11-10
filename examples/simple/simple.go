package main

import (
	"github.com/slclub/easy"
	"simple/initialize"
	"simple/lservers"
)

func main() {
	initialize.Init()
	Start()
}

func Start() {
	easy.Serv(
		lservers.Server1(), // websocket 监听服务 可以有多个
		lservers.Server2(), // tcp 服务
	)
}
