package client

import (
	"github.com/slclub/easy/nets/agent"
	"github.com/slclub/easy/nets/conns"
	"github.com/slclub/easy/route"
	"net"
)

type TCPClient struct {
	Client
}

func NewTCPClient(gate *Gate, rs ...route.Router) *TCPClient {
	// route init
	self := &TCPClient{}
	self.router = autoRouter(rs)
	self.Gate.Init(gate)

	self.init()

	return self
}

func (self *TCPClient) init() {

	conn_origin, err := net.Dial("tcp", self.Addr)
	if err != nil {
		panic(any("tcp client dial " + err.Error()))
	}
	self.connOption = &conns.Option{
		Encrypt:   self.router.Encoder(),
		MsgDigit:  self.Gate.MsgDigit,
		MaxMsgLen: self.Gate.MaxMsgLen,
		MinMsgLen: 2,
	}
	conn_tcp := conns.NewTCPConn(conn_origin, self.PendingWriteNum, self.connOption)
	self.agent = agent.NewAgent(conn_tcp)
}

func (self *TCPClient) startBefore() {

}

func (self *TCPClient) Start() {
	self.startBefore()

	go self.Agent().LoopRecv(dealHandle(&self.Client))
}
