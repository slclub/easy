package conns

// 每个链接routine 的独自状态
type connChan struct {
	writeChan chan []byte
	closeFlag bool
	stopChan  chan struct{}
}

func (self *connChan) Done() chan struct{} {
	return self.stopChan
}

func (self *connChan) closeChan() {
	close(self.stopChan)
	close(self.writeChan)
}

func (self *connChan) release() {
	self.stopChan = nil
	self.writeChan = nil
}
