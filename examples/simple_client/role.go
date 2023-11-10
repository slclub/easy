package main

import (
	"github.com/slclub/easy/client"
	"github.com/slclub/easy/nets/agent"
)

type Role struct {
	client any
	Uid    int
}

func NewRole() *Role {
	return &Role{}
}

func (self *Role) Agent() agent.Agent {
	switch val := self.client.(type) {
	case *client.WSClient:
		return val.Agent()
	case *client.TCPClient:
		return val.Agent()
	}
	return nil
}

func (self *Role) ListenServ() {
	switch val := self.client.(type) {
	case *client.WSClient:
		val.Start()
	case *client.TCPClient:
		val.Start()
	}
	return
}
