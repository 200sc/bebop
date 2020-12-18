// Package iohelp provides common io utilities for bebop generated code.
package iohelp

import (
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
		Reader: r,
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

func ReadString(r ErrorReader) string {
	ln := uint32(0)
	binary.Read(r, binary.LittleEndian, &ln)
	data := make([]byte, ln)
	r.Read(data)
	return string(data)
}

func ReadTime(r ErrorReader) time.Time {
	tm := int64(0)
	binary.Read(r, binary.LittleEndian, (&tm))
	tm *= 100
	t := time.Time{}
	if tm == 0 {
		return t
	}
	return time.Unix(0, tm).UTC()
}

func ReadGUID(r ErrorReader) [16]byte {
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
	_, err := io.ReadFull(r, r.Buffer[:1])
	if err != nil {
		r.Err = err
		return 0
	}
	return r.Buffer[0]
}

func ReadUint8(r ErrorReader) uint8 {
	_, err := io.ReadFull(r, r.Buffer[:1])
	if err != nil {
		r.Err = err
	}
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

func WriteGUIDBytes(b []byte, guid [16]byte) {
	_ = b[15]
	b[0] = guid[3]
	b[1] = guid[2]
	b[2] = guid[1]
	b[3] = guid[0]
	b[4] = guid[5]
	b[5] = guid[4]
	b[6] = guid[7]
	b[7] = guid[6]
	b[8] = guid[8]
	b[9] = guid[9]
	b[10] = guid[10]
	b[11] = guid[11]
	b[12] = guid[12]
	b[13] = guid[13]
	b[14] = guid[14]
	b[15] = guid[15]
}

func WriteInt64(w ErrorWriter, i int64) {
	WriteInt64Bytes(w.Buffer, i)
	w.Write(w.Buffer)
}

func WriteInt64Bytes(b []byte, i int64) {
	_ = b[7]
	b[0] = byte(i)
	b[1] = byte(i >> 8)
	b[2] = byte(i >> 16)
	b[3] = byte(i >> 24)
	b[4] = byte(i >> 32)
	b[5] = byte(i >> 40)
	b[6] = byte(i >> 48)
	b[7] = byte(i >> 56)
}

func WriteUint64(w ErrorWriter, i uint64) {
	WriteUint64Bytes(w.Buffer, i)
	w.Write(w.Buffer)
}

func WriteUint64Bytes(b []byte, i uint64) {
	_ = b[7]
	b[0] = byte(i)
	b[1] = byte(i >> 8)
	b[2] = byte(i >> 16)
	b[3] = byte(i >> 24)
	b[4] = byte(i >> 32)
	b[5] = byte(i >> 40)
	b[6] = byte(i >> 48)
	b[7] = byte(i >> 56)
}

func WriteInt32(w ErrorWriter, i int32) {
	WriteInt32Bytes(w.Buffer, i)
	w.Write(w.Buffer[:4])
}

func WriteInt32Bytes(b []byte, i int32) {
	_ = b[3]
	b[0] = byte(i)
	b[1] = byte(i >> 8)
	b[2] = byte(i >> 16)
	b[3] = byte(i >> 24)
}

func WriteUint32(w ErrorWriter, i uint32) {
	WriteUint32Bytes(w.Buffer, i)
	w.Write(w.Buffer[:4])
}

func WriteUint32Bytes(b []byte, i uint32) {
	_ = b[3]
	b[0] = byte(i)
	b[1] = byte(i >> 8)
	b[2] = byte(i >> 16)
	b[3] = byte(i >> 24)
}

func WriteInt16(w ErrorWriter, i int16) {
	WriteInt16Bytes(w.Buffer, i)
	w.Write(w.Buffer[:2])
}

func WriteInt16Bytes(b []byte, i int16) {
	_ = b[1]
	b[0] = byte(i)
	b[1] = byte(i >> 8)
}

func WriteUint16(w ErrorWriter, i uint16) {
	WriteUint16Bytes(w.Buffer, i)
	w.Write(w.Buffer[:2])
}

func WriteUint16Bytes(b []byte, i uint16) {
	_ = b[1]
	b[0] = byte(i)
	b[1] = byte(i >> 8)
}

func WriteByte(w ErrorWriter, b byte) {
	w.Write([]byte{b})
}

func WriteByteBytes(b []byte, by byte) {
	b[0] = by
}

func WriteUint8(w ErrorWriter, b uint8) {
	w.Write([]byte{b})
}

func WriteUint8Bytes(b []byte, by uint8) {
	b[0] = by
}

func WriteBool(w ErrorWriter, b bool) {
	if b {
		w.Write([]byte{1})
	} else {
		w.Write([]byte{0})
	}
}

func WriteBoolBytes(b []byte, bl bool) {
	if bl {
		b[0] = 1
	} else {
		b[0] = 0
	}
}

func WriteFloat32(w ErrorWriter, f float32) {
	WriteUint32(w, math.Float32bits(f))
}

func WriteFloat32Bytes(b []byte, f float32) {
	WriteUint32Bytes(b, math.Float32bits(f))
}

func WriteFloat64(w ErrorWriter, f float64) {
	WriteUint64(w, math.Float64bits(f))
}

func WriteFloat64Bytes(b []byte, f float64) {
	WriteUint64Bytes(b, math.Float64bits(f))
}
