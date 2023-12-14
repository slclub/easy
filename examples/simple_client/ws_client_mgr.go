package main

import (
	"github.com/slclub/easy/client"
	"github.com/slclub/easy/route"
	"github.com/slclub/easy/typehandle"
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
	// create  two roles
	for i := 0; i < nws; i++ {
		self.append(self.CreateRole())
	}
	//self.append(self.CreateRole())
}

func (self *WsTestMgr) CreateRole() *Role {
	role := NewRole()
	cln := client.NewWsClient(self.gate, self.router)
	role.client = cln
	return role
}

func (self *WsTestMgr) Run() {

	self.run()
}

func (self *WsTestMgr) run() {
	for _, role := range self.roles {
		role.ListenServ()
	}
}

func (self *WsTestMgr) append(role *Role) {
	self.roles = append(self.roles, role)
}

func StartWs() {
	// 初始化 客户端
	WsMgr = NewWsTestMgr(&client.Gate{
		Protocol:     typehandle.ENCRIPT_DATA_JSON,
		Addr:         ":15080",
		LittleEndian: true,
	})

	WsMgr.Init()
	InitWsRegister() // 注册消息ID 和路由
	WsMgr.Run()
}
