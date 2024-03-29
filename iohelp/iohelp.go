// Package iohelp provides common io utilities for bebop generated code.
package iohelp

import (
	"errors"
	"io"
	"math"
	"time"
	"unsafe"
)

var (

	// ErrUnpopulatedUnion indicates a union had no contents, when exactly
	// one member should be populated
	ErrUnpopulatedUnion error = errors.New("union has no populated member")
)

// An ErrorReader wraps an io.Reader with a reusable buffer for small allocations
// of structured data types and persistent error tracking.
type ErrorReader struct {
	Reader io.Reader
	Err    error
	buffer []byte
}

// NewErrorReader wraps a reader with error tracking. An ErrorReader not created
// via NewErrorReader is not safe to use.
func NewErrorReader(r io.Reader) *ErrorReader {
	if er, ok := r.(*ErrorReader); ok {
		return er
	}
	return &ErrorReader{
		Reader: r,
		buffer: make([]byte, 8),
	}
}

func (er *ErrorReader) Read(b []byte) (n int, err error) {
	n, err = io.ReadFull(er.Reader, b)
	if err != nil {
		er.Err = err
	}
	return n, err
}

func (er *ErrorReader) Drain() {
	_, _ = io.ReadAll(er.Reader)
}

// An ErrorWriter wraps an io.Writer with a reusable buffer for small allocations
// of structured data types and persistent error tracking.
type ErrorWriter struct {
	Writer io.Writer
	Err    error
	buffer []byte
}

// NewErrorWriter wraps a writer with error tracking. An ErrorWriter not created
// via NewErrorWriter is not safe to use.
func NewErrorWriter(w io.Writer) *ErrorWriter {
	if ew, ok := w.(*ErrorWriter); ok {
		return ew
	}
	return &ErrorWriter{
		Writer: w,
		buffer: make([]byte, 8),
	}
}

func (ew *ErrorWriter) Write(b []byte) (n int, err error) {
	n, err = ew.Writer.Write(b)
	if err != nil {
		ew.Err = err
	}
	return n, err
}

// SafeWrite will only write if ew.Err is not set, and will set it if
// a write fails.
func (ew *ErrorWriter) SafeWrite(b []byte) int {
	if ew.Err != nil {
		return 0
	}
	n, _ := ew.Write(b)
	return n
}

func ReadString(r *ErrorReader) string {
	data := make([]byte, ReadUint32(r))
	_, _ = r.Read(data)
	return string(data)
}

func MustReadStringBytes(buf []byte) string {
	sz := ReadUint32Bytes(buf)
	return string(buf[4 : 4+sz])
}

func MustReadStringBytesSharedMemory(buf []byte) string {
	sz := ReadUint32Bytes(buf)
	cut := buf[4 : 4+sz]
	return *(*string)(unsafe.Pointer(&cut))
}

func ReadStringBytes(buf []byte) (string, error) {
	if len(buf) < 4 {
		return "", io.ErrUnexpectedEOF
	}
	sz := ReadUint32Bytes(buf)
	if len(buf) < int(sz)+4 {
		return "", io.ErrUnexpectedEOF
	}
	return string(buf[4 : 4+sz]), nil
}

func ReadStringBytesSharedMemory(buf []byte) (string, error) {
	if len(buf) < 4 {
		return "", io.ErrUnexpectedEOF
	}
	sz := ReadUint32Bytes(buf)
	if len(buf) < int(sz)+4 {
		return "", io.ErrUnexpectedEOF
	}
	cut := buf[4 : 4+sz]
	return *(*string)(unsafe.Pointer(&cut)), nil
}

func ReadDate(r *ErrorReader) time.Time {
	_, _ = io.ReadFull(r, r.buffer)
	return ReadDateBytes(r.buffer)
}

func ReadDateBytes(buf []byte) time.Time {
	tm := ReadInt64Bytes(buf)
	tm *= 100
	t := time.Time{}
	if tm == 0 {
		return t
	}
	return time.Unix(0, tm).UTC()
}

func ReadGUID(r *ErrorReader) [16]byte {
	data := make([]byte, 16)
	_, _ = r.Read(data)
	return ReadGUIDBytes(data)
}

func ReadGUIDBytes(buf []byte) [16]byte {
	return [16]byte{
		buf[3], buf[2], buf[1], buf[0],
		buf[5], buf[4],
		buf[7], buf[6],
		buf[8], buf[9], buf[10], buf[11], buf[12], buf[13], buf[14], buf[15],
	}
}

func ReadBool(r *ErrorReader) bool {
	_, _ = io.ReadFull(r, r.buffer[:1])
	return r.buffer[0] == 1
}

func ReadBoolBytes(buf []byte) bool {
	// Technically a value other than 0 or 1 is invalid.
	return buf[0] == 1
}

func ReadByte(r *ErrorReader) byte {
	_, err := io.ReadFull(r, r.buffer[:1])
	if err != nil {
		r.Err = err
		return 0
	}
	return r.buffer[0]
}

func ReadByteBytes(buf []byte) byte {
	return buf[0]
}

func ReadUint8(r *ErrorReader) uint8 {
	_, err := io.ReadFull(r, r.buffer[:1])
	if err != nil {
		r.Err = err
	}
	return r.buffer[0]
}

func ReadUint8Bytes(buf []byte) uint8 {
	return buf[0]
}

func ReadUint16(r *ErrorReader) uint16 {
	_, _ = io.ReadFull(r, r.buffer[:2])
	return ReadUint16Bytes(r.buffer)
}

func ReadUint16Bytes(buf []byte) uint16 {
	_ = buf[1]
	return *(*uint16)(unsafe.Pointer(&buf[0]))
}

func ReadInt16(r *ErrorReader) int16 {
	_, _ = io.ReadFull(r, r.buffer[:2])
	return ReadInt16Bytes(r.buffer)
}

func ReadInt16Bytes(buf []byte) int16 {
	_ = buf[1]
	return *(*int16)(unsafe.Pointer(&buf[0]))
}

func ReadUint32(r *ErrorReader) uint32 {
	_, _ = io.ReadFull(r, r.buffer[:4])
	return ReadUint32Bytes(r.buffer)
}

func ReadUint32Bytes(buf []byte) uint32 {
	_ = buf[3]
	return *(*uint32)(unsafe.Pointer(&buf[0]))
}

func ReadInt32(r *ErrorReader) int32 {
	_, _ = io.ReadFull(r, r.buffer[:4])
	return ReadInt32Bytes(r.buffer)
}

func ReadInt32Bytes(buf []byte) int32 {
	_ = buf[3]
	return *(*int32)(unsafe.Pointer(&buf[0]))
}

func ReadUint64(r *ErrorReader) uint64 {
	_, _ = io.ReadFull(r, r.buffer)
	return ReadUint64Bytes(r.buffer)
}

func ReadUint64Bytes(buf []byte) uint64 {
	_ = buf[7]
	return *(*uint64)(unsafe.Pointer(&buf[0]))
}

func ReadInt64(r *ErrorReader) int64 {
	_, _ = io.ReadFull(r, r.buffer)
	return ReadInt64Bytes(r.buffer)
}

func ReadInt64Bytes(buf []byte) int64 {
	_ = buf[7]
	return *(*int64)(unsafe.Pointer(&buf[0]))
}

func ReadFloat32(r *ErrorReader) float32 {
	return math.Float32frombits(ReadUint32(r))
}

func ReadFloat32Bytes(buf []byte) float32 {
	return math.Float32frombits(ReadUint32Bytes(buf))
}

func ReadFloat64(r *ErrorReader) float64 {
	return math.Float64frombits(ReadUint64(r))
}

func ReadFloat64Bytes(buf []byte) float64 {
	return math.Float64frombits(ReadUint64Bytes(buf))
}

func WriteGUID(w *ErrorWriter, guid [16]byte) {
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
	_, _ = w.Write(flipped[:])
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

func WriteInt64(w *ErrorWriter, i int64) {
	WriteInt64Bytes(w.buffer, i)
	_, _ = w.Write(w.buffer)
}

func WriteInt64Bytes(b []byte, i int64) {
	_ = b[7]
	*(*int64)(unsafe.Pointer(&b[0])) = i
}

func WriteUint64(w *ErrorWriter, i uint64) {
	WriteUint64Bytes(w.buffer, i)
	_, _ = w.Write(w.buffer)
}

func WriteUint64Bytes(b []byte, i uint64) {
	// See bench_test.go
	_ = b[7]
	*(*uint64)(unsafe.Pointer(&b[0])) = i
}

func WriteInt32(w *ErrorWriter, i int32) {
	WriteInt32Bytes(w.buffer, i)
	_, _ = w.Write(w.buffer[:4])
}

func WriteInt32Bytes(b []byte, i int32) {
	_ = b[3]
	*(*int32)(unsafe.Pointer(&b[0])) = i
}

func WriteUint32(w *ErrorWriter, i uint32) {
	WriteUint32Bytes(w.buffer, i)
	_, _ = w.Write(w.buffer[:4])
}

func WriteUint32Bytes(b []byte, i uint32) {
	_ = b[3]
	*(*uint32)(unsafe.Pointer(&b[0])) = i
}

func WriteInt16(w *ErrorWriter, i int16) {
	WriteInt16Bytes(w.buffer, i)
	_, _ = w.Write(w.buffer[:2])
}

func WriteInt16Bytes(b []byte, i int16) {
	_ = b[1]
	*(*int16)(unsafe.Pointer(&b[0])) = i
}

func WriteUint16(w *ErrorWriter, i uint16) {
	WriteUint16Bytes(w.buffer, i)
	_, _ = w.Write(w.buffer[:2])
}

func WriteUint16Bytes(b []byte, i uint16) {
	_ = b[1]
	*(*uint16)(unsafe.Pointer(&b[0])) = i
}

func WriteByte(w *ErrorWriter, b byte) {
	_, _ = w.Write([]byte{b})
}

func WriteByteBytes(b []byte, by byte) {
	b[0] = by
}

func WriteUint8(w *ErrorWriter, b uint8) {
	_, _ = w.Write([]byte{b})
}

func WriteUint8Bytes(b []byte, by uint8) {
	b[0] = by
}

func WriteBool(w *ErrorWriter, b bool) {
	if b {
		_, _ = w.Write([]byte{1})
	} else {
		_, _ = w.Write([]byte{0})
	}
}

func WriteBoolBytes(b []byte, bl bool) {
	if bl {
		b[0] = 1
	} else {
		b[0] = 0
	}
}

func WriteFloat32(w *ErrorWriter, f float32) {
	WriteUint32(w, math.Float32bits(f))
}

func WriteFloat32Bytes(b []byte, f float32) {
	WriteUint32Bytes(b, math.Float32bits(f))
}

func WriteFloat64(w *ErrorWriter, f float64) {
	WriteUint64(w, math.Float64bits(f))
}

func WriteFloat64Bytes(b []byte, f float64) {
	WriteUint64Bytes(b, math.Float64bits(f))
}
