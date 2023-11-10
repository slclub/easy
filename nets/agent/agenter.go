package agent

import "net"

type Agent interface {
	WriteMsg(msg any)
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
	Destroy()
	UserData() any
	SetUserData(data any)
	LoopRecv(handle AgentHandle)
}
