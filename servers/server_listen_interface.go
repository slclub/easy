package servers

import (
	"github.com/slclub/easy/route"
	"github.com/slclub/easy/vendors/option"
)

type ListenServer interface {
	Init(assignment option.Assignment)
	OnInit()
	Router() route.Router
	Start()
	Hook() *hookAgent // agent 链接 回调
	//OnClose(func())
	Close()
}
