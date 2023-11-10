package route

import (
	"github.com/slclub/easy/log"
	"github.com/slclub/easy/nets/agent"
	"github.com/slclub/easy/route/element"
)

const (
	MSG_LOGIN_REQ = 1001
	MSG_LOGIN_RES = 1002
)

// JSON message define
type JsonLoginReq struct {
	MID  element.MID
	UID  string
	Name string
	Sex  int16
}

type JsonLoginReS struct {
	MID element.MID
	UID string
}

// Protobuf message define

// message handlers

func loginRequest(ag agent.Agent, msg any) {
	data, ok := msg.(*JsonLoginReq)
	if !ok {
		log.Error("loginRequest handle get an error message")
		return
	}
	log.Debug("Router.Route got handler is (loginRequest) MID:%v Name:%v", data.MID, data.Name)
}

func loginResponse(ag agent.Agent, msg any) {
	log.Debug("Router.Route got handler is (loginResponse)")
}
