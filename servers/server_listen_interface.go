package servers

import (
	"github.com/slclub/easy/nets/agent"
	"github.com/slclub/easy/route"
)

type ListenServer interface {
	Init(*agent.Gate)
	OnInit()
	Router() route.Router
	Start()
	Hook() *hookAgent // agent 链接 回调
	//OnClose(func())
	Close()
}
