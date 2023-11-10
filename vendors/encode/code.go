package encode

import "encoding/binary"

// 于 Encoder 接口没有任何关系
type Code struct {
	littleEndian bool
}

func (self *Code) LittleEndian(b ...bool) bool {
	if len(b) >= 1 {
		self.littleEndian = b[0]
	}
	return self.littleEndian
}

func (self *Code) Uint2Bytes(mid uint16) []byte {
	id_byte := make([]byte, 2)
	if self.littleEndian {
		binary.LittleEndian.PutUint16(id_byte, mid)
	} else {
		binary.BigEndian.PutUint16(id_byte, mid)
	}
	return id_byte
}

func (self *Code) Bytes2Uint(data []byte) uint16 {
	var mid uint16
	if self.littleEndian {
		mid = binary.LittleEndian.Uint16(data)
	} else {
		mid = binary.BigEndian.Uint16(data)
	}
	return mid
}

func (self *Code) Uint322Bytes(mid uint32) []byte {
	id_byte := make([]byte, 2)
	if self.littleEndian {
		binary.LittleEndian.PutUint32(id_byte, mid)
	} else {
		binary.BigEndian.PutUint32(id_byte, mid)
	}
	return id_byte
}

func (self *Code) Bytes2Uint32(data []byte) uint32 {
	var mid uint32
	if self.littleEndian {
		mid = binary.LittleEndian.Uint32(data)
	} else {
		mid = binary.BigEndian.Uint32(data)
	}
	return mid
}
