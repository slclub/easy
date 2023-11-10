package callback

import (
	"github.com/slclub/easy/nets/agent"
	"github.com/slclub/easy/servers"
	"github.com/slclub/log8q"
	"simple/lservers"
)

func RegisterCallerToLservers() {
	lservers.Server1().Hook().Append(servers.CONST_AGENT_NEW, handleOnServerNew)
	lservers.Server1().Hook().Append(servers.CONST_AGENT_CLOSE, handleOnServerClose)
}

func handleOnServerNew(ag agent.Agent) {
	log8q.Info("[CONNECTION.NEW] server1 create an new connection")
}

func handleOnServerClose(ag agent.Agent) {
	log8q.Info("[CONNECTION.CLOSE] server1 closed an old connection")
}
