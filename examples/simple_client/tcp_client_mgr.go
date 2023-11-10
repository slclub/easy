package main

import (
	"github.com/slclub/easy/client"
	"github.com/slclub/easy/route"
	"github.com/slclub/easy/typehandle"
)

type TCPTestMgr struct {
	router route.Router
	roles  []*Role
	gate   *client.Gate
}

var TCPMgr *TCPTestMgr

func NewTCPTestMgr(gate *client.Gate) *TCPTestMgr {
	mgr := &TCPTestMgr{
		gate:   gate,
		router: route.NewRouter(),
	}
	client.RouterWithProtocol(mgr.router, mgr.gate.Protocol)
	mgr.router.Encoder().LittleEndian(mgr.gate.LittleEndian)
	return mgr
}

func (self *TCPTestMgr) Init() {
	// create three roles
	self.append(self.CreateRole())
	self.append(self.CreateRole())
	self.append(self.CreateRole())
}

func (self *TCPTestMgr) CreateRole() *Role {
	role := NewRole()
	cln := client.NewTCPClient(self.gate, self.router)
	role.client = cln
	return role
}

func (self *TCPTestMgr) Run() {

	self.run()
}

func (self *TCPTestMgr) run() {
	for _, role := range self.roles {
		role.ListenServ()
	}
}

func (self *TCPTestMgr) append(role *Role) {
	self.roles = append(self.roles, role)
}

func StartTCP() {
	// 初始化 客户端
	TCPMgr = NewTCPTestMgr(&client.Gate{
		Protocol:     typehandle.ENCRIPT_DATA_JSON,
		Addr:         ":18081",
		LittleEndian: true,
	})

	TCPMgr.Init()
	InitTCPRegister() // 注册消息ID 和路由
	TCPMgr.Run()
}
