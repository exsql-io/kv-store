package util

import (
	"encoding/binary"
)

func Join(size int, bytes ...[]byte) []byte {
	b, i := make([]byte, size), 0
	for _, v := range bytes {
		i += copy(b[i:], v)
	}

	return b
}

func UInt32FromBytes(bytes []byte) uint32 {
	return binary.LittleEndian.Uint32(bytes)
}

func UInt32ToBytes(value uint32) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, value)

	return bytes
}
