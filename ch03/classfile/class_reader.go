package classfile

import "encoding/binary"

type ClassReader struct {
	data []byte
}

// u1
func (cr *ClassReader) readUint8() uint8 {
	val := cr.data[0]
	cr.data = cr.data[1:]
	return val
}

// u2
// Go标准库encoding/binary包中定义了一个变量BigEndian，正好 可以从[]byte中解码多字节数据
func (cr *ClassReader) readUint16() uint16 {
	val := binary.BigEndian.Uint16(cr.data)
	cr.data = cr.data[2:]
	return val
}

// u4
func (cr *ClassReader) readUint32() uint32 {
	val := binary.BigEndian.Uint32(cr.data)
	cr.data = cr.data[4:]
	return val
}

func (cr *ClassReader) readUint64() uint64 {
	val := binary.BigEndian.Uint64(cr.data)
	cr.data = cr.data[8:]
	return val
}

// 读取uint16表，表的大小由开头的uint16数据指 出
func (cr *ClassReader) readUint16s() []uint16 {
	// 开头的uint16说明了后面表的大小
	n := cr.readUint16()
	// 创建指定长度的slice, 类型为uint16, 用于存储后面的uint16表
	s := make([]uint16, n)
	// 遍历[0,n), 依次读取uint16数据, 填充到slice
	for i := range s {
		s[i] = cr.readUint16()
	}
	return s
}

// 最后一个方法是readBytes()，用于读取指定数量的字节
func (cr *ClassReader) readBytes(n uint32) []byte {
	bytes := cr.data[:n]
	cr.data = cr.data[n:]
	return bytes
}
