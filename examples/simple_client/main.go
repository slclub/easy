package main

import (
	"os"
	"os/signal"
)

func main() {

	StartWs()
	StartTCP()

	// 业务逻辑入口
	RunBusiniess()

	wait()
}

func RunBusiniess() {
	Do(WsMgr.roles)
	Do(TCPMgr.roles)
}

func wait() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}
