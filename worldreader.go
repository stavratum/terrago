package main

import (
	"encoding/binary"
	"os"
)

type WorldReader struct {
	B8 []byte // buffer

	*os.File
}

func (r WorldReader) Uint64(b *uint64) {
	r.File.Read(r.B8)
	*b = binary.LittleEndian.Uint64(r.B8)
}

func (r WorldReader) Uint32(b *uint32) {
	r.File.Read(r.B8[:4])
	*b = binary.LittleEndian.Uint32(r.B8[:4])
}

func (r WorldReader) Uint16(b *uint16) {
	r.File.Read(r.B8[:2])
	*b = binary.LittleEndian.Uint16(r.B8[:2])
}

func (r WorldReader) Int32(b *int32) {
	r.File.Read(r.B8[:4])
	*b = int32(binary.LittleEndian.Uint32(r.B8[:4]))
}

func (r WorldReader) Int16(b *int16) {
	r.File.Read(r.B8[:2])
	*b = int16(binary.LittleEndian.Uint16(r.B8[:2]))
}

func (r WorldReader) Byte() byte {
	r.File.Read(r.B8[:1])
	return r.B8[0]
}

// read a capped string
func (r WorldReader) String(s *string, cap int) {
	b := make([]byte, cap)
	r.File.Read(b)

	*s = string(b)
}

func NewReader(f *os.File) WorldReader {
	return WorldReader{make([]byte, 8), f}
}

/*
func (r WorldReader) Read(data any) error {
	return binary.Read(r, binary.LittleEndian, data)
}*/
