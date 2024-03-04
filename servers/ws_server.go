package servers

import (
	"github.com/gorilla/websocket"
	"github.com/slclub/easy/log"
	"github.com/slclub/easy/nets/agent"
	"github.com/slclub/easy/nets/conns"
	"github.com/slclub/easy/route"
	"net/http"
)

/**
 * web socket 监听服务
 */
type WSServer struct {
	Server
	upgrader websocket.Upgrader
}

func NewWSServer() *WSServer {
	ser := Server{
		router: route.NewRouter(),
		hook:   newHookAgent(),
	}
	return &WSServer{
		Server: ser,
	}
}

func (self *WSServer) Start() {
	self.startBefore()

	handleHttp := &WebSocketHandle{
		server: self,
		handle: dealHandle(&self.Server),
	}
	httpServer := &http.Server{
		Addr:           self.Addr,
		Handler:        handleHttp,
		ReadTimeout:    self.HTTPTimeout,
		WriteTimeout:   self.HTTPTimeout,
		MaxHeaderBytes: 1024,
	}

	go httpServer.Serve(self.ln)
}

var _ ListenServer = &WSServer{}

// ------------------------------------------------------------
// serverHttpHandle
type WebSocketHandle struct {
	server *WSServer
	handle agent.AgentHandle
}

func (self *WebSocketHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	conn, err := self.server.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Debug("upgrade error: %v", err)
		return
	}
	conn.SetReadLimit(int64(self.server.MaxMsgLen))

	if err = self.server.box.Add(conn); err != nil {
		conn.Close()
		log.Error("", err)
		return
	}
	wsConn := conns.NewWSConn(conn, self.server.connOption, self.server.PendingWriteNum, self.server.MaxMsgLen)
	ag := agent.NewAgent(wsConn)
	self.server.hook.EmitWithKey(CONST_AGENT_NEW, ag)
	ag.LoopRecv(self.handle)

	//ag.Close()
	//ag.OnClose()
	self.server.hook.EmitWithKey(CONST_AGENT_CLOSE, ag)
	self.server.box.Del(conn)
}

// function handle 路由分发
func dealHandle(serv *Server) agent.AgentHandle {
	return func(data []byte, ag agent.Agent) {

		msg, err := serv.Router().Encoder().Unmarshal(data)
		if err != nil {
			log.Debug("unmarshal message error: %v", err)
			return
		}
		err = serv.Router().Route(msg, ag)
		if err != nil {
			log.Debug("route message error: %v", err)
			return
		}
	}
}
