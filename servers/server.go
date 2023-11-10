package servers

import (
	"crypto/tls"
	"errors"
	"github.com/slclub/easy/log"
	"github.com/slclub/easy/nets/agent"
	"github.com/slclub/easy/nets/conns"
	"github.com/slclub/easy/route"
	"github.com/slclub/easy/route/bind"
	"github.com/slclub/easy/route/encode"
	"github.com/slclub/easy/typehandle"
	"io"
	"net"
	"sync"
)

/**
 * 基类 监听服务
 */
type Server struct {
	agent.Gate
	ln         net.Listener
	box        *ConnBox
	router     route.Router
	hook       *hookAgent
	connOption *conns.Option
}

func (self *Server) Init(gate *agent.Gate) {
	self.Gate.Init(gate)
	if self.Protocol == "" {
		self.Protocol = typehandle.EnCriPT_DATA_PROTOBUF
	}

	self.defaultRouteEncripty()
	self.Router().Encoder().LittleEndian(self.LittleEndian)
}

func (self *Server) startBefore() {

	self.connOption = &conns.Option{
		Encrypt:   self.Router().Encoder(),
		MaxMsgLen: self.MaxMsgLen,
		MinMsgLen: 2,
		MsgDigit:  self.MsgDigit,
	}
	ln, err := net.Listen("tcp", self.Addr)
	if err != nil {
		log.Fatal("Listen error: %v", err)
	}
	if self.CertFile != "" || self.KeyFile != "" {
		config := &tls.Config{}
		config.NextProtos = []string{"http/1.1"}

		var err error
		config.Certificates = make([]tls.Certificate, 1)
		config.Certificates[0], err = tls.LoadX509KeyPair(self.CertFile, self.KeyFile)
		if err != nil {
			log.Fatal("%v", err)
		}
		ln = tls.NewListener(ln, config)
	}
	self.ln = ln

	// 链接容器
	self.box = &ConnBox{
		server: self,
		conns:  make(ConnSet),
		mutex:  sync.Mutex{},
	}
}

func (self *Server) defaultRouteEncripty() {
	switch self.Protocol {
	case typehandle.ENCRIPT_DATA_JSON:
		self.router.Binding(bind.NewBindJson(self.router.PathMap()), encode.NewJson(self.router.PathMap()))
	default:
		self.router.Binding(bind.NewBindProto(self.router.PathMap()), encode.NewProtobuf(self.router.PathMap()))
	}
}

func (self *Server) Close() {
	// 利用钩子优雅关闭
	self.hook.EmitWithKey(CONST_SERVER_CLOSE, nil)

	// 关闭监听
	self.ln.Close()

	// 关闭哦链接池
	self.box.Close()

}

func (self *Server) Router() route.Router {
	return self.router
}

func (self *Server) Hook() *hookAgent {
	return self.hook
}

// -----------------------------------------------------------
// connection box
type ConnSet map[any]struct{}

// 链接池容器
type ConnBox struct {
	server *Server
	conns  ConnSet
	mutex  sync.Mutex
}

func (self *ConnBox) Add(key any) error {
	if self.conns == nil {
		return errors.New("WSConnBox.conn is nil")
	}
	self.mutex.Lock()
	defer self.mutex.Unlock()
	if len(self.conns) >= self.server.MaxConnNum {
		return errors.New("WSConnBox connection fail too many conns")
	}
	self.conns[key] = struct{}{}
	return nil
}

func (self *ConnBox) Del(key any) error {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	delete(self.conns, key)
	return nil
}

func (self *ConnBox) Close() {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	for conn, _ := range self.conns {
		switch val := conn.(type) {
		case io.Closer:
			val.Close()
		}
	}
	self.conns = nil

}

func (self *ConnBox) Len() int {
	return len(self.conns)
}
