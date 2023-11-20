package easy

import (
	"github.com/slclub/easy/log"
	"github.com/slclub/easy/servers"
	"github.com/slclub/easy/vendors/ants"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	log.Init()
}

func Serv(servs ...servers.ListenServer) {

	//
	for _, s := range servs {
		s.OnInit()
		s.Start()
	}

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	<-c

	for _, s := range servs {
		s.Close()
	}
	ants.Pool().Release()
}
