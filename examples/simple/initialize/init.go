package initialize

import (
	"simple/controller"
	"simple/controller/callback"
	"simple/lservers"
	"simple/message"
)

func Init() {
	// 读取配置
	// do configurition

	// listen servers initialization
	forServerInitialize()

	callback.RegisterCallerToLservers()
}

func forServerInitialize() {
	// Init configure data to listening server.
	lservers.InitListenServer()

	// registing messages to the Router of listening server
	message.Init()
	controller.InitBindingRoute()
}
