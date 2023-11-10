package login

import (
	"github.com/slclub/easy/nets/agent"
	"reflect"
	"simple/vendors/log8q"
)

func HandleLogin(agent agent.Agent, arg any) {

	log8q.Log().Info("WS controller.Handle.Login info: ", reflect.TypeOf(arg).Elem().Name())
}

func HandleLoginTcp(agent2 agent.Agent, arg any) {
	log8q.Log().Info("TCP controller.Handle.Login info: ", reflect.TypeOf(arg).Elem().Name())
}
