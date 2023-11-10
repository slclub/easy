package encode

// stream 解析器 encode decode操作
type Encoder interface {
	Unmarshal(data []byte) (any, error)
	Marshal(msg any) ([]byte, error)
	LittleEndian(...bool) bool
}
