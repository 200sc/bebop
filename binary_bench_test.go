package bebop

import (
	"encoding/binary"
	"testing"
)

type nullWriter struct{}

func (nullWriter) Write(data []byte) (int, error) {
	return len(data), nil
}

var err error

func BenchmarkBinaryWrite(b *testing.B) {
	w := nullWriter{}
	for i := 0; i < b.N; i++ {
		err = binary.Write(w, binary.LittleEndian, int64(i))
	}
}

func BenchmarkCustomBinaryWrite(b *testing.B) {
	w := nullWriter{}
	for i := 0; i < b.N; i++ {
		v := int64(i)
		b := make([]byte, 8)
		b[0] = byte(v)
		b[1] = byte(v >> 8)
		b[2] = byte(v >> 16)
		b[3] = byte(v >> 24)
		b[4] = byte(v >> 32)
		b[5] = byte(v >> 40)
		b[6] = byte(v >> 48)
		b[7] = byte(v >> 56)
		_, err = w.Write(b)
	}
}
