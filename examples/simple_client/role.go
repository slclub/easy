package main

import (
	"github.com/slclub/easy/client"
)

type Role struct {
	client *client.WSClient
	Uid    int
}

func NewRole() *Role {
	return &Role{}
}

func (self *Role) Client() *client.WSClient {
	return self.client
}
