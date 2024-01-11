package callback

import (
	"github.com/slclub/easy/nets/agent"
	"github.com/slclub/easy/servers"
	"simple/lservers"
	"simple/vendors/log8q"
)

func RegisterCallerToLservers() {
	lservers.Server1().Hook().Append(servers.CONST_AGENT_NEW, handleOnConnNew)
	lservers.Server1().Hook().Append(servers.CONST_AGENT_CLOSE, handleOnConnClose)
	lservers.Server1().Hook().Append(servers.CONST_SERVER_CLOSE, handleOnServerClose)

	lservers.Server2().Hook().Append(servers.CONST_AGENT_NEW, handleOnConnNew)
	lservers.Server2().Hook().Append(servers.CONST_AGENT_CLOSE, handleOnConnClose)
}

func handleOnConnNew(ag agent.Agent) {
	log8q.Info("[CONNECTION.NEW] server create an new connection")
}

func handleOnConnClose(ag agent.Agent) {
	log8q.Info("[CONNECTION.CLOSE] server closed an old connection")
}

// the current listening server is closing
// smoothly shutdown the server
func handleOnServerClose(ag agent.Agent) {
	// ag == nil
	// 执行一些 平滑停服务的逻辑
}
