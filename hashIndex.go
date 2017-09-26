package main

import "fmt"
import "io"

// hash index is a number, uint64, representing a deterministic non-redundant []byte.
// all index values have a unique byte sequence, but not all 8-byte sequences can be produced, because variable length means some states are encoded by the length.
// the un=indexable []byte's are above {0xFE,0xFE,0xFE,0xFE,0xFE,0xFE,0xFE,0xFE}
type hashIndexType uint64

// reader implemenation returning the []byte represented by hashindex
func (hi hashIndexType) Read(b []byte) (n int, err error) {
	if hi == 0 {
		return 0, io.EOF
	}
	if hi < 0x0101 {
		hi -= 1
		b[0] = uint8(hi)
		return 1, io.EOF
	}
	if hi < 0x010101 {
		hi -= 0x0101
		b[1] = uint8(hi)
		b[0] = uint8(hi >> 8)
		return 2, io.EOF
	}
	if hi < 0x01010101 {
		hi -= 0x010101
		b[2] = uint8(hi)
		b[1] = uint8(hi >> 8)
		b[0] = uint8(hi >> 16)
		return 3, io.EOF
	}
	if hi < 0x0101010101 {
		hi -= 0x01010101
		b[3] = uint8(hi)
		b[2] = uint8(hi >> 8)
		b[1] = uint8(hi >> 16)
		b[0] = uint8(hi >> 24)
		return 4, io.EOF
	}
	if hi < 0x010101010101 {
		hi -= 0x0101010101
		b[4] = uint8(hi)
		b[3] = uint8(hi >> 8)
		b[2] = uint8(hi >> 16)
		b[1] = uint8(hi >> 24)
		b[0] = uint8(hi >> 32)
		return 5, io.EOF
	}
	if hi < 0x01010101010101 {
		hi -= 0x010101010101
		b[5] = uint8(hi)
		b[4] = uint8(hi >> 8)
		b[3] = uint8(hi >> 16)
		b[2] = uint8(hi >> 24)
		b[1] = uint8(hi >> 32)
		b[0] = uint8(hi >> 40)
		return 6, io.EOF
	}
	if hi < 0x0101010101010101 {
		hi -= 0x01010101010101
		b[6] = uint8(hi)
		b[5] = uint8(hi >> 8)
		b[4] = uint8(hi >> 16)
		b[3] = uint8(hi >> 24)
		b[2] = uint8(hi >> 32)
		b[1] = uint8(hi >> 40)
		b[0] = uint8(hi >> 48)
		return 7, io.EOF
	}
	hi -= 0x0101010101010101
	b[7] = uint8(hi)
	b[6] = uint8(hi >> 8)
	b[5] = uint8(hi >> 16)
	b[4] = uint8(hi >> 24)
	b[3] = uint8(hi >> 32)
	b[2] = uint8(hi >> 40)
	b[1] = uint8(hi >> 48)
	b[0] = uint8(hi >> 56)
	return 8, io.EOF
}

// new hash index from a []byte
func NewHashIndexType(b []byte) (hi hashIndexType) {
	switch len(b) {
	case 0:
		return hashIndexType(0)
	case 1:
		return hashIndexType(1 + uint64(b[0]))
	case 2:
		return hashIndexType(0x0101 + uint64(b[1]) + uint64(b[0])<<8)
	case 3:
		return hashIndexType(0x010101 + uint64(b[2]) + uint64(b[1])<<8 + uint64(b[0])<<16)
	case 4:
		return hashIndexType(0x01010101 + uint64(b[3]) + uint64(b[2])<<8 + uint64(b[1])<<16 + uint64(b[0])<<24)
	case 5:
		return hashIndexType(0x0101010101 + uint64(b[4]) + uint64(b[3])<<8 + uint64(b[2])<<16 + uint64(b[1])<<24 + uint64(b[0])<<32)
	case 6:
		return hashIndexType(0x010101010101 + uint64(b[5]) + uint64(b[4])<<8 + uint64(b[3])<<16 + uint64(b[2])<<24 + uint64(b[1])<<32 + uint64(b[0])<<40)
	case 7:
		return hashIndexType(0x01010101010101 + uint64(b[6]) + uint64(b[5])<<8 + uint64(b[4])<<16 + uint64(b[3])<<24 + uint64(b[2])<<32 + uint64(b[1])<<40 + uint64(b[0])<<48)
	default:
		return hashIndexType(0x0101010101010101 + uint64(b[7]) + uint64(b[6])<<8 + uint64(b[5])<<16 + uint64(b[4])<<24 + uint64(b[3])<<32 + uint64(b[2])<<40 + uint64(b[1])<<48 + uint64(b[0])<<56)
	}
}

// string rep as hexadecimal
func (hi hashIndexType) String() string {
	b := make([]byte, 8, 8)
	n, _ := hi.Read(b)
	return fmt.Sprintf("% x", b[:n])
}

// return new hashindex whoses rep is as the source heashindex but with added byte(s)
func (hi hashIndexType) Append(b byte) hashIndexType{
	buf := make([]byte, 8, 8)
	n, _ := hi.Read(buf)
	buf[n]=b
	return NewHashIndexType(buf[:n+1])
}

// return new hashindex whoses rep is as the source heashindex but with removed byte(s)
func (hi hashIndexType) Truncate(c int) hashIndexType{
	buf := make([]byte, 8, 8)
	n, _ := hi.Read(buf)
	return NewHashIndexType(buf[:n-c])
}


