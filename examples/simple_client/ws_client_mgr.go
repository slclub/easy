package main

import (
	"github.com/slclub/easy/client"
	"github.com/slclub/easy/route"
	"os"
	"os/signal"
)

type WsTestMgr struct {
	router route.Router
	roles  []*Role
	gate   *client.Gate
}

var WsMgr *WsTestMgr

func NewWsTestMgr(gate *client.Gate) *WsTestMgr {
	mgr := &WsTestMgr{
		gate:   gate,
		router: route.NewRouter(),
	}
	client.RouterWithProtocol(mgr.router, mgr.gate.Protocol)
	mgr.router.Encoder().LittleEndian(mgr.gate.LittleEndian)
	return mgr
}

func (self *WsTestMgr) Init() {
	// create roles
	role := NewRole()
	cln := client.NewWsClient(self.gate, self.router)
	role.client = cln

	self.append(role)
}

func (self *WsTestMgr) Run() {
	c := make(chan os.Signal)

	self.run()
	Do(self.roles)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}

func (self *WsTestMgr) run() {
	for _, role := range self.roles {
		role.Client().Start()
	}
}

func (self *WsTestMgr) append(role *Role) {
	self.roles = append(self.roles, role)
}
