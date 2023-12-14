package lservers

import (
	"github.com/slclub/easy/nets/agent"
	"github.com/slclub/easy/servers"
	"github.com/slclub/easy/typehandle"
	"github.com/slclub/easy/vendors/option"
	"strconv"
)

var (
	server1 servers.ListenServer
	server2 servers.ListenServer
)

//
//type EasyServer struct {
//	Server servers.ListenServer
//}
//
//func (self *EasyServer) Init() {
//	self.Server = servers.NewWSServer()
//	self.Server.Init(&agent.Gate{
//		Addr:            "::18080",
//		Protocol:        typehandle.ENCRIPT_DATA_JSON,
//		PendingWriteNum: 2000,
//		LittleEndian:    true,
//		MaxConnNum:      2000,
//	})
//}

func Server1() servers.ListenServer {
	return server1
}

func Server2() servers.ListenServer {
	return server2
}

func InitListenServer(ListenPort int) {
	server1 = servers.NewWSServer()

	server1.Init(option.OptionWith(&agent.Gate{
		Addr:            ":" + strconv.Itoa(ListenPort),
		Protocol:        typehandle.ENCRIPT_DATA_JSON,
		PendingWriteNum: 2000,
		LittleEndian:    true,
		MaxConnNum:      20000,
	}).Default(option.DEFAULT_IGNORE_ZERO))

	server2 = servers.NewTCPServer()
	server2.Init(option.OptionWith(&agent.Gate{
		Addr:            ":" + strconv.Itoa(ListenPort+1),
		Protocol:        typehandle.ENCRIPT_DATA_JSON,
		PendingWriteNum: 2000,
		LittleEndian:    true,
		MaxConnNum:      20000,
	}).Default(option.DEFAULT_IGNORE_ZERO))
}
