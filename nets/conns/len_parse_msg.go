package conns

import (
	"errors"
	"github.com/slclub/easy/vendors/encode"
	"io"
	"math"
)

// --------------
//
//	Message with lenght
//
// --------------
type msgParser struct {
	encode.Code
	*Option
}

func NewMsgParser() *msgParser {
	p := new(msgParser)
	p.MsgDigit = 4
	p.MinMsgLen = 1
	p.MaxMsgLen = 4096

	return p
}

// It's dangerous to call the method on reading or writing
func (self *msgParser) SetMsgLen(option *Option) {
	if option.MsgDigit > 0 {
		self.MsgDigit = option.MsgDigit
	}
	if option.MinMsgLen != 0 {
		self.MinMsgLen = option.MinMsgLen
	}
	if option.MaxMsgLen != 0 {
		self.MaxMsgLen = option.MaxMsgLen
	}

	var max uint32
	switch self.MsgDigit {
	case 1:
		max = math.MaxUint8
	case 2:
		max = math.MaxUint16
	case 4:
		max = math.MaxUint32
	}
	if self.MinMsgLen > max {
		self.MinMsgLen = max
	}
	if self.MaxMsgLen > max {
		self.MaxMsgLen = max
	}
}

// goroutine safe
func (self *msgParser) Read(conn *TCPConn) ([]byte, error) {
	var b [4]byte
	bufMsgLen := b[:self.MsgDigit]

	// read len
	if _, err := io.ReadFull(conn, bufMsgLen); err != nil {
		return nil, err
	}

	// parse len
	var msgLen uint32 = self.Bytes2Uint32(bufMsgLen)

	// check len
	if msgLen > self.MaxMsgLen {
		return nil, errors.New("message too long")
	} else if msgLen < self.MinMsgLen {
		return nil, errors.New("message too short")
	}

	// data
	msgData := make([]byte, msgLen)
	if _, err := io.ReadFull(conn, msgData); err != nil {
		return nil, err
	}

	return msgData, nil
}

// goroutine safe
func (self *msgParser) Write(conn *TCPConn, args []byte) error {
	// get len
	msgLen := uint32(len(args))

	// check len
	if msgLen > self.MaxMsgLen {
		return errors.New("message too long")
	} else if msgLen < self.MinMsgLen {
		return errors.New("message too short")
	}

	msg := make([]byte, uint32(self.MsgDigit)+msgLen)

	// write len
	idbyte := self.Uint322Bytes(msgLen)
	copy(msg, idbyte)

	// write data
	l := self.MsgDigit
	copy(msg[l:], args)

	conn.Write(msg)

	return nil
}
