package servers

import (
	"github.com/slclub/easy/log"
	"github.com/slclub/easy/nets/agent"
	"github.com/slclub/easy/typehandle"
)

const (
	CONST_AGENT_CLOSE  = "CloseAgent"
	CONST_AGENT_NEW    = "NewAgent"
	CONST_SERVER_CLOSE = "CloseServer" // running your defined functions where you closed the listening server
)

type hookAgent struct {
	hook map[string][]typehandle.AgentHandle
}

func newHookAgent() *hookAgent {
	return &hookAgent{
		hook: make(map[string][]typehandle.AgentHandle),
	}
}

func (self *hookAgent) Append(hkey string, handle typehandle.AgentHandle) {
	if self.hook[hkey] == nil {
		self.hook[hkey] = []typehandle.AgentHandle{}
	}
	self.hook[hkey] = append(self.hook[hkey], handle)
}

func (self *hookAgent) EmitWithKey(hkey string, ag agent.Agent) {
	if nil == self.hook[hkey] || len(self.hook[hkey]) == 0 {
		return
	}
	for _, fn := range self.hook[hkey] {
		fn(ag)
	}
}

// defualt functions

func TcpNewAgent(a agent.Agent) {
	log.Debug("live connection NewAgent Addr[%v]", a.RemoteAddr().String())
}

func TcpCloseAgent(a agent.Agent) {
	log.Debug("live connection CloseAgent Addr[%v]", a.RemoteAddr().String())
}
