package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
)

var (
	nws  int
	ntcp int
)

func init() {
	flag.IntVar(&nws, "nws", 1, " ws client numbers")
	flag.IntVar(&nws, "ntcp", 1, " tcp client numbers")
}

func main() {
	flag.Parse()

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
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	<-c
}
