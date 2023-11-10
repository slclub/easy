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

// 链接的消息传输 规则 接口
// 注意与 传输用的数据格式区分开来，举例 web scoket 不需要这一步，内部已经有规则了。 数据格式还是需要我们自己定义，是json 还是protobuf 或者其他格式
type FromConnReadWriter interface {
	SetWithOption(option *Option)
	Read(conn *TCPConn) ([]byte, error)
	Write(conn *TCPConn, args []byte) error
}

type Handle func([]byte)
