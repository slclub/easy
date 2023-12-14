package main

import (
	"flag"
	"github.com/slclub/easy/log"
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
	flag.IntVar(&ntcp, "ntcp", 0, " tcp client numbers")
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
	log.Debug("ntpc=%v nws:=%v  WebSocket.Roles=%v TCP.Roles=%v", ntcp, nws, len(WsMgr.roles), len(TCPMgr.roles))
}

func wait() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	<-c
}
