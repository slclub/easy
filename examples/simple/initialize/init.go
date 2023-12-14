package initialize

import (
	"simple/controller"
	"simple/controller/callback"
	"simple/lservers"
	"simple/message"
)

var ListenPort = 15080

func Init(fn func()) {
	// 读取配置
	// do configurition

	if fn != nil {
		fn()
	}

	// listen servers initialization
	forServerInitialize()

	callback.RegisterCallerToLservers()
}

func forServerInitialize() {
	// Init configure data to listening server.
	lservers.InitListenServer(ListenPort)

	// registing messages to the Router of listening server
	message.Init()
	controller.InitBindingRoute()
	controller.InitBindingRouteServer2()
}
