package agent

import (
	"github.com/slclub/easy/nets/conns"
	"time"
)

type Gate struct {
	MaxConnNum      int
	PendingWriteNum int
	MaxMsgLen       uint32
	Addr            string

	HTTPTimeout time.Duration
	CertFile    string
	KeyFile     string

	MsgDigit     int
	LittleEndian bool // 常用 默认值是 true
	Protocol     string
}

func (self *Gate) OnInit() {

}

func (self *Gate) Init(gate *Gate) {
	self.Addr = gate.Addr
	self.HTTPTimeout = gate.HTTPTimeout
	self.CertFile = gate.CertFile
	self.KeyFile = gate.KeyFile
	self.HTTPTimeout = gate.HTTPTimeout
	self.LittleEndian = gate.LittleEndian
	self.MsgDigit = gate.MsgDigit
	self.MaxConnNum = gate.MaxConnNum
	self.PendingWriteNum = gate.PendingWriteNum
	self.Protocol = gate.Protocol

	if self.MaxConnNum <= 0 {
		self.MaxConnNum = 100
	}
	if self.PendingWriteNum <= 0 {
		self.PendingWriteNum = 100
	}
	if self.MaxMsgLen <= 0 {
		self.MaxMsgLen = 4096
	}
	if self.HTTPTimeout <= 0 {
		self.HTTPTimeout = 10 * time.Second
	}
	if self.MsgDigit == 0 {
		self.MsgDigit = conns.CONST_MSG_DIGIT
	}
}
