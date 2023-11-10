package client

import (
	"github.com/gorilla/websocket"
	"github.com/slclub/easy/log"
	"github.com/slclub/easy/nets/agent"
	"github.com/slclub/easy/nets/conns"
	"github.com/slclub/easy/route"
	"net/http"
)

type WSClient struct {
	Client
}

// you just need to using one Router enetity for all clients when you created more than one client
// the second param is Router allow you set existed one.
func NewWsClient(gate *Gate, rs ...route.Router) *WSClient {
	// route init
	var r route.Router
	if len(rs) == 0 {
		r = route.NewRouter()
	} else {
		r = rs[0]
	}

	self := &WSClient{}
	self.router = r
	self.Gate.Init(gate)

	self.init()

	return self
}

func (self *WSClient) init() {
	dialer := websocket.Dialer{}

	addr, err := GinWebSocketSchceme(self.Addr)
	if err != nil {
		panic(any(err.Error()))
	}
	conn_origin, _, err := dialer.Dial(addr, http.Header{})
	if err != nil {
		panic(any("client dial " + err.Error()))
	}
	self.connOption = &conns.Option{
		Encrypt:   self.router.Encoder(),
		MsgDigit:  self.Gate.MsgDigit,
		MaxMsgLen: self.Gate.MaxMsgLen,
		MinMsgLen: 2,
	}
	conn_ws := conns.NewWSConn(conn_origin, self.connOption, self.PendingWriteNum, self.MaxMsgLen)
	self.agent = agent.NewAgent(conn_ws)
}

func (self *WSClient) Agent() agent.Agent {
	return self.agent
}

func (self *WSClient) startBefore() {

}

func (self *WSClient) Start() {
	self.startBefore()

	go self.Agent().LoopRecv(dealHandle(&self.Client))
}

//--------------------------------------------
// functions
//--------------------------------------------

// function handle 路由分发
func dealHandle(client *Client) agent.AgentHandle {
	return func(data []byte, ag agent.Agent) {
		msg, err := client.router.Encoder().Unmarshal(data)
		if err != nil {
			log.Debug("unmarshal message error: %v", err)
			return
		}
		err = client.router.Route(msg, ag)
		if err != nil {
			log.Debug("route message error: %v", err)
			return
		}
	}
}
