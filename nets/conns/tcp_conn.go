package conns

import (
	"github.com/slclub/easy/log"
	"net"
)

type TCPConn struct {
	//sync.Mutex
	conn net.Conn
	connChan
	*Option
}

func NewTCPConn(conn net.Conn, pendingWriteNum int, option *Option) *TCPConn {
	tcpConn := new(TCPConn)
	tcpConn.conn = conn
	tcpConn.writeChan = make(chan []byte, pendingWriteNum)
	tcpConn.stopChan = make(chan struct{})
	tcpConn.Option = option

	if option.MsgParser == nil {
		option.MsgParser = NewMsgParser()
	}
	tcpConn.Option.MsgParser.SetWithOption(option)
	go tcpConn.loopSend()

	return tcpConn
}

func (self *TCPConn) loopSend() {
	defer self.Destroy()
	for {
		select {
		case <-self.Done():
			//self.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			self.conn.Write([]byte("TCP CLOSE"))
			return
		case b := <-self.writeChan:
			if b == nil {
				return
			}

			_, err := self.conn.Write(b)
			if err != nil {
				self.Done()
			}
		}
	}
}

func (self *TCPConn) doDestroy() {
	self.conn.(*net.TCPConn).SetLinger(0)
	self.conn.Close()

	if !self.closeFlag {
		close(self.writeChan)
		self.closeFlag = true
	}
}

func (self *TCPConn) Destroy() {
	//self.Lock()
	//defer self.Unlock()

	self.doDestroy()
}

func (self *TCPConn) Close() {
	//self.Lock()
	//defer self.Unlock()
	if self.closeFlag {
		return
	}
	self.doWrite(nil)
	self.closeFlag = true
	self.closeChan()
}

func (self *TCPConn) doWrite(b []byte) {
	if len(self.writeChan) == cap(self.writeChan) {
		log.Debug("close conn: channel full")
		self.doDestroy()
		return
	}

	self.writeChan <- b
}

// b must not be modified by the others goroutines
func (self *TCPConn) Write(b []byte) {
	//self.Lock()
	//defer self.Unlock()
	if self.closeFlag || b == nil {
		return
	}

	self.doWrite(b)
}

func (self *TCPConn) Read(b []byte) (int, error) {
	return self.conn.Read(b)
}

func (self *TCPConn) LocalAddr() net.Addr {
	return self.conn.LocalAddr()
}

func (self *TCPConn) RemoteAddr() net.Addr {
	return self.conn.RemoteAddr()
}

func (self *TCPConn) ReadMsg() ([]byte, error) {
	return self.GetOption().MsgParser.Read(self)
}

func (self *TCPConn) WriteMsg(args []byte) error {
	return self.GetOption().MsgParser.Write(self, args)
}

var _ Conn = &TCPConn{}
