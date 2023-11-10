package login

import (
	"github.com/slclub/easy/nets/agent"
	"reflect"
	"simple/vendors/log8q"
)

func HandleLogin(agent agent.Agent, arg any) {

	log8q.Log().Info("controller.Handle.Login info: ", reflect.TypeOf(arg).Elem().Name())
}
