package conns

import (
	"errors"
	"github.com/gorilla/websocket"
	"github.com/slclub/easy/log"
	"net"
)

type WSConn struct {
	//sync.Mutex
	conn *websocket.Conn
	connChan
	*Option
	write_over_num int
}

func NewWSConn(conn *websocket.Conn, option *Option, pendingWriteNum int, maxMsgLen uint32) *WSConn {
	wsConn := new(WSConn)
	wsConn.conn = conn
	wsConn.writeChan = make(chan []byte, pendingWriteNum)
	wsConn.stopChan = make(chan struct{})
	wsConn.Option = option

	go wsConn.loopSend()
	//go wsConn.loopRecv(handle)

	return wsConn
}

func (self *WSConn) loopSend() {
	defer self.Destroy()
	for {
		select {
		case <-self.Done():
			self.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			log.Debug("WSConn.loopSend STOP")
			return
		case b := <-self.writeChan:
			if b == nil {
				return
			}

			err := self.conn.WriteMessage(websocket.BinaryMessage, b)
			if err != nil {
				return
			}
		}
	}
}

func (self *WSConn) LoopRecv(handle Handle) {
	defer self.Close()
	for {
		select {
		case <-self.Done():
			//a.conn.WriteMsg()
			return
		default:
			data, err := self.ReadMsg()
			if err != nil {
				log.Debug("ws conn read connection [%v] error message: %v", self, err)
				return
			}
			handle(data)
		}
	}
}

func (wsConn *WSConn) doDestroy() {
	wsConn.conn.UnderlyingConn().(*net.TCPConn).SetLinger(0)
	wsConn.conn.Close()

	if !wsConn.closeFlag {
		//close(wsConn.writeChan)
		wsConn.closeFlag = true
	}
	wsConn.release()
}

func (wsConn *WSConn) Destroy() {
	//wsConn.Lock()
	//defer wsConn.Unlock()

	wsConn.doDestroy()
}

func (wsConn *WSConn) Close() {
	//wsConn.Lock()
	//defer wsConn.Unlock()
	if wsConn.closeFlag {
		return
	}
	//wsConn.doWrite(nil)
	wsConn.closeChan()
	//wsConn.conn.Close()
	wsConn.closeFlag = true
	//wsConn.conn.UnderlyingConn().(*net.TCPConn).SetLinger(0)
}

// 写消息到客户端
func (wsConn *WSConn) doWrite(b []byte) {
	if len(wsConn.writeChan) == cap(wsConn.writeChan) {
		wsConn.write_over_num++
		if wsConn.write_over_num%1000 == 0 {
			log.Debug("close conn: channel full closeFlag:%v", wsConn.closeFlag)
		}
		wsConn.doDestroy()
		return
	}

	wsConn.writeChan <- b
}

func (wsConn *WSConn) LocalAddr() net.Addr {
	return wsConn.conn.LocalAddr()
}

func (wsConn *WSConn) RemoteAddr() net.Addr {
	return wsConn.conn.RemoteAddr()
}

// goroutine not safe
func (wsConn *WSConn) ReadMsg() ([]byte, error) {
	if wsConn.closeFlag {
		return nil, errors.New("closed wsconn")
	}
	_, b, err := wsConn.conn.ReadMessage()
	return b, err
}

// args must not be modified by the others goroutines
func (wsConn *WSConn) WriteMsg(args []byte) error {
	//wsConn.Lock()
	//defer wsConn.Unlock()
	if wsConn.closeFlag {
		return nil
	}

	// get len
	msgLen := uint32(len(args))

	// check len
	if msgLen > wsConn.Option.MaxMsgLen {
		return errors.New("message too long")
	} else if msgLen < wsConn.Option.MinMsgLen {
		return errors.New("message too short")
	}

	// don't copy
	wsConn.doWrite(args)
	return nil
}

func (WSConn *WSConn) release() {
	WSConn.connChan.release()
	WSConn.Option = nil
	WSConn.conn = nil
}

var _ Conn = &WSConn{}
