package conns

const (
	// 用TCP 消息长度的 位数 byte 数量
	// 2 byte = 2的16次方  65535
	// 4 byte = 2 - 32 也就是4G
	CONST_MSG_DIGIT = 4

	// WebSocket 默认不需要这东西，已经内部完成了
)

type Encoder interface {
	Unmarshal(data []byte) (any, error)
	Marshal(msg any) ([]byte, error)
}

// 公用的引用状态 ； 例如：链接所以来的Server服务，解析器

type Option struct {
	MaxMsgLen uint32
	MinMsgLen uint32
	MsgDigit  int
	Encrypt   Encoder
	MsgParser FromConnReadWriter
}

func (self *Option) GetOption() *Option {
	return self
}
