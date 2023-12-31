package main

import (
	"github.com/slclub/easy/nets/agent"
	"github.com/slclub/log8q"
	msgjson "simple_client/message/json"
)

func LoginRes(a agent.Agent, msg any) {
	rs, _ := msg.(*msgjson.LoginRes)
	log8q.Info("WS client.LoginReq ", *rs)

}

func LoginResTCP(a agent.Agent, msg any) {
	rs, _ := msg.(*msgjson.LoginRes)
	log8q.Info("TCP client.LoginReq ", *rs)
}
