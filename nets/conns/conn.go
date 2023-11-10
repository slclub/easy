package conns

import (
	"net"
)

type Conn interface {
	ReadMsg() ([]byte, error)
	WriteMsg(args []byte) error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
	Destroy()
	Done() chan struct{}
	GetOption() *Option
}

type Handle func([]byte)
