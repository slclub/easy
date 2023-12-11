package cgrpc

import (
	"os"
	"os/signal"
	"syscall"
)

type waiter struct {
	wait chan os.Signal
}

func (self *waiter) Wait() {
	signal.Notify(self.wait, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	<-self.wait
	defer signal.Stop(self.wait)
}

func (self *waiter) close() {
	close(self.wait)
}
