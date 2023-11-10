package main

import (
	"fmt"
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
	for _, role := range roles {
		role.Client().Agent().WriteMsg(msg)
	}
	fmt.Println("------ DO roles ", len(roles))
}
