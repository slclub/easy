package main

import (
	"fmt"
	"reflect"
	msgjson "simple_client/message/json"
	"time"
)

// 去写role 的动作，是发消息请求，还是其他
func Do(roles []*Role) {
	sendMsgWithPeriod(roles)
}

// 多少周期发一次消息
func sendMsgWithPeriod(roles []*Role) {
	go func() {
		ticker := time.NewTicker(time.Second * 5)
		for {
			select {
			case <-ticker.C:
				sendMsgRoles(roles, &msgjson.LoginReq{
					UID:  10,
					Name: "axigl",
				})
			}
		}
	}()
}

func sendMsgRoles(roles []*Role, msg any) {
	if len(roles) == 0 {
		return
	}
	for _, role := range roles {
		if role.Agent() == nil {
			continue
			//panic(any("agent of role is nil"))
		}
		role.Agent().WriteMsg(msg)
	}
	fmt.Println("------ DO roles ", len(roles), " client: ", reflect.TypeOf(roles[0].client).Elem().Name())
}
