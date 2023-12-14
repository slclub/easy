package servers

import (
	"github.com/slclub/easy/log"
	"github.com/slclub/easy/nets/agent"
	"github.com/slclub/easy/nets/conns"
	"github.com/slclub/easy/route"
	"net"
	"time"
)

type TCPServer struct {
	Server
}

func NewTCPServer() *TCPServer {
	ser := Server{
		router: route.NewRouter(),
		hook:   newHookAgent(),
	}
	return &TCPServer{
		Server: ser,
	}
}

func (self *TCPServer) Start() {
	self.startBefore()
	self.run()
}

func (self *TCPServer) run() {
	var tempDelay time.Duration
	for {
		conn, err := self.ln.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.Release("accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return
		}
		tempDelay = 0

		if self.box.Len() >= self.MaxConnNum {
			conn.Close()
			log.Debug("EASY.TCP too many connections")
			continue
		}
		self.box.Add(conn)

		tcpConn := conns.NewTCPConn(conn, self.PendingWriteNum, self.connOption)
		ag := agent.NewAgent(tcpConn)
		go func() {
			self.hook.EmitWithKey(CONST_AGENT_NEW, ag)
			ag.LoopRecv(dealHandle(&self.Server))

			ag.Close()
			//ag.OnClose()
			self.hook.EmitWithKey(CONST_AGENT_CLOSE, ag)
			self.box.Del(conn)

			// cleanup
		}()
	}
}
