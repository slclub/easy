package client

import (
	"github.com/slclub/easy/nets/conns"
	"time"
)

/**
 * the gate is dedicated for client. it is very different from server gate.
 * Dose not confuse the two of gate.
 */

type Gate struct {
	// connection
	Addr        string
	HTTPTimeout time.Duration

	// certicate
	CertFile string
	KeyFile  string

	// mesaage configurtion
	MsgDigit        int
	MaxMsgLen       uint32
	LittleEndian    bool // 常用 默认值是 true
	PendingWriteNum int
	Protocol        string
}

func (self *Gate) Init(gate *Gate) {
	self.Addr = gate.Addr
	self.HTTPTimeout = gate.HTTPTimeout
	self.CertFile = gate.CertFile
	self.KeyFile = gate.KeyFile
	self.HTTPTimeout = gate.HTTPTimeout
	self.LittleEndian = gate.LittleEndian
	self.MsgDigit = gate.MsgDigit
	self.PendingWriteNum = gate.PendingWriteNum
	self.Protocol = gate.Protocol

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
