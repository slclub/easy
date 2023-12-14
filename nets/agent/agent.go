package agent

import (
	"github.com/slclub/easy/log"
	"github.com/slclub/easy/nets/conns"
	"net"
	"reflect"
)

type agent struct {
	conn     conns.Conn
	userData any
}

type AgentHandle func([]byte, Agent)

func NewAgent(conn conns.Conn) Agent {
	return &agent{
		conn: conn,
	}
}

func (a *agent) LoopRecv(handle AgentHandle) {
	defer a.conn.Close()
	for {
		select {
		case <-a.conn.Done():
			//a.conn.WriteMsg()
			log.Debug("agent.LoopRecv STOP")
			return
		default:
			data, err := a.conn.ReadMsg()
			if err != nil {
				log.Debug("agent read connection error message: %v", err)
				return
			}
			if handle == nil {
				continue
			}
			handle(data, a)
		}
	}
}

func (a *agent) OnClose() {
	//if a.gate.AgentChanRPC != nil {
	//	err := a.gate.AgentChanRPC.Call0("CloseAgent", a)
	//	if err != nil {
	//		log.Error("chanrpc error: %v", err)
	//	}
	//}
}

func (a *agent) WriteMsg(msg any) {
	data, err := a.encrypt().Marshal(msg)
	if err != nil {
		log.Error("marshal message %v error: %v", reflect.TypeOf(msg), err)
		return
	}
	err = a.conn.WriteMsg(data)
	if err != nil {
		log.Error("write message %v error: %v", reflect.TypeOf(msg), err)
	}
}

func (a *agent) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *agent) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

func (a *agent) Close() {
	//a.Destroy()
	a.conn.Close()
	a.userData = nil
}

func (a *agent) Destroy() {
	a.conn.Destroy()
}

func (a *agent) UserData() any {
	return a.userData
}

func (a *agent) SetUserData(data any) {
	a.userData = data
}

func (a *agent) encrypt() conns.Encoder {
	return a.conn.GetOption().Encrypt
}

var _ Agent = &agent{}
