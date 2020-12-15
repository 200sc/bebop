// Package iohelp provides common io utilities for bebop generated code.
package iohelp

import (
	"bufio"
	"encoding/binary"
	"io"
	"math"
	"time"
)

type ErrorReader struct {
	Reader io.Reader
	Err    error
	Buffer []byte
}

func NewErrorReader(r io.Reader) ErrorReader {
	return ErrorReader{
		Reader: bufio.NewReader(r),
		Buffer: make([]byte, 8),
	}
}

func (er ErrorReader) Read(b []byte) (n int, err error) {
	n, err = er.Reader.Read(b)
	if err != nil {
		er.Err = err
	}
	return n, err
}

func NewErrorWriter(w io.Writer) ErrorWriter {
	return ErrorWriter{
		Writer: w,
		Buffer: make([]byte, 8),
	}
}

type ErrorWriter struct {
	Writer io.Writer
	Err    error
	Buffer []byte
}

func (ew ErrorWriter) Write(b []byte) (n int, err error) {
	n, err = ew.Writer.Write(b)
	if err != nil {
		ew.Err = err
	}
	return n, err
}

func ReadString(r io.Reader) string {
	ln := uint32(0)
	binary.Read(r, binary.LittleEndian, &ln)
	data := make([]byte, ln)
	r.Read(data)
	return string(data)
}

func ReadTime(r io.Reader) time.Time {
	tm := int64(0)
	binary.Read(r, binary.LittleEndian, (&tm))
	tm *= 100
	t := time.Time{}
	if tm == 0 {
		return t
	}
	return time.Unix(0, tm).UTC()
}

func ReadGUID(r io.Reader) [16]byte {
	data := make([]byte, 16)
	r.Read(data)
	flipped := [16]byte{
		data[3], data[2], data[1], data[0],
		data[5], data[4],
		data[7], data[6],
		data[8], data[9], data[10], data[11], data[12], data[13], data[14], data[15],
	}
	return flipped
}

func ReadBool(r ErrorReader) bool {
	io.ReadFull(r, r.Buffer[:1])
	return r.Buffer[0] == 1
}

func ReadByte(r ErrorReader) byte {
	io.ReadFull(r, r.Buffer[:1])
	return r.Buffer[0]
}

func ReadUint16(r ErrorReader) uint16 {
	io.ReadFull(r, r.Buffer[:2])
	return uint16(r.Buffer[0]) | uint16(r.Buffer[1])<<8
}

func ReadInt16(r ErrorReader) int16 {
	io.ReadFull(r, r.Buffer[:2])
	return int16(r.Buffer[0]) | int16(r.Buffer[1])<<8
}

func ReadUint32(r ErrorReader) uint32 {
	io.ReadFull(r, r.Buffer[:4])
	return uint32(r.Buffer[0]) | uint32(r.Buffer[1])<<8 | uint32(r.Buffer[2])<<16 | uint32(r.Buffer[3])<<24
}

func ReadInt32(r ErrorReader) int32 {
	io.ReadFull(r, r.Buffer[:4])
	return int32(r.Buffer[0]) | int32(r.Buffer[1])<<8 | int32(r.Buffer[2])<<16 | int32(r.Buffer[3])<<24
}

func ReadUint64(r ErrorReader) uint64 {
	io.ReadFull(r, r.Buffer)
	return uint64(r.Buffer[0]) | uint64(r.Buffer[1])<<8 | uint64(r.Buffer[2])<<16 | uint64(r.Buffer[3])<<24 |
		uint64(r.Buffer[4])<<32 | uint64(r.Buffer[5])<<40 | uint64(r.Buffer[6])<<48 | uint64(r.Buffer[7])<<56
}

func ReadInt64(r ErrorReader) int64 {
	io.ReadFull(r, r.Buffer)
	return int64(r.Buffer[0]) | int64(r.Buffer[1])<<8 | int64(r.Buffer[2])<<16 | int64(r.Buffer[3])<<24 |
		int64(r.Buffer[4])<<32 | int64(r.Buffer[5])<<40 | int64(r.Buffer[6])<<48 | int64(r.Buffer[7])<<56
}

func ReadFloat32(r ErrorReader) float32 {
	return math.Float32frombits(ReadUint32(r))
}

func ReadFloat64(r ErrorReader) float64 {
	return math.Float64frombits(ReadUint64(r))
}

func WriteGUID(w ErrorWriter, guid [16]byte) {
	// 3 2 1 0
	// 5 4
	// 7 6
	// 8 9 10 11 12 13 14 15
	flipped := [16]byte{
		guid[3], guid[2], guid[1], guid[0],
		guid[5], guid[4],
		guid[7], guid[6],
		guid[8], guid[9], guid[10], guid[11], guid[12], guid[13], guid[14], guid[15],
	}
	w.Write(flipped[:])
}

func WriteInt64(w ErrorWriter, i int64) {
	w.Buffer[0] = byte(i)
	w.Buffer[1] = byte(i >> 8)
	w.Buffer[2] = byte(i >> 16)
	w.Buffer[3] = byte(i >> 24)
	w.Buffer[4] = byte(i >> 32)
	w.Buffer[5] = byte(i >> 40)
	w.Buffer[6] = byte(i >> 48)
	w.Buffer[7] = byte(i >> 56)
	w.Write(w.Buffer)
}

func WriteUint64(w ErrorWriter, i uint64) {
	w.Buffer[0] = byte(i)
	w.Buffer[1] = byte(i >> 8)
	w.Buffer[2] = byte(i >> 16)
	w.Buffer[3] = byte(i >> 24)
	w.Buffer[4] = byte(i >> 32)
	w.Buffer[5] = byte(i >> 40)
	w.Buffer[6] = byte(i >> 48)
	w.Buffer[7] = byte(i >> 56)
	w.Write(w.Buffer)
}

func WriteInt32(w ErrorWriter, i int32) {
	w.Buffer[0] = byte(i)
	w.Buffer[1] = byte(i >> 8)
	w.Buffer[2] = byte(i >> 16)
	w.Buffer[3] = byte(i >> 24)
	w.Write(w.Buffer[:4])
}

func WriteUint32(w ErrorWriter, i uint32) {
	w.Buffer[0] = byte(i)
	w.Buffer[1] = byte(i >> 8)
	w.Buffer[2] = byte(i >> 16)
	w.Buffer[3] = byte(i >> 24)
	w.Write(w.Buffer[:4])
}

func WriteInt16(w ErrorWriter, i int16) {
	w.Buffer[0] = byte(i)
	w.Buffer[1] = byte(i >> 8)
	w.Write(w.Buffer[:2])
}

func WriteUint16(w ErrorWriter, i uint16) {
	w.Buffer[0] = byte(i)
	w.Buffer[1] = byte(i >> 8)
	w.Write(w.Buffer[:2])
}

func WriteByte(w ErrorWriter, b byte) {
	w.Write([]byte{b})
}

func WriteUint8(w ErrorWriter, b uint8) {
	w.Write([]byte{b})
}

func WriteBool(w ErrorWriter, b bool) {
	if b {
		w.Write([]byte{1})
	} else {
		w.Write([]byte{0})
	}
}

func WriteFloat32(w ErrorWriter, f float32) {
	WriteUint32(w, math.Float32bits(f))
}

func WriteFloat64(w ErrorWriter, f float64) {
	WriteUint64(w, math.Float64bits(f))
}
