package iohelp

import (
	"bytes"
	"encoding/binary"
	"math"
	"math/rand"
	"testing"
	"unsafe"
)

func randomUints() [1000]uint64 {
	cases := [1000]uint64{}
	for i := 0; i < len(cases); i++ {
		cases[i] = uint64(rand.Intn(math.MaxInt64))
	}
	return cases
}

func BenchmarkWriteBytes_BinaryLittleEndianCopy(b *testing.B) {
	b.StopTimer()
	cases := randomUints()
	writeTo := make([]byte, 8)
	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		// arbitrarily making this longer so the benchmark runs longer
		for j := 0; j < 100; j++ {
			tc := cases[rand.Intn(1000)]
			writeBinaryLittleEndianCopy(tc, writeTo)
			globalGot = readBinaryLittleEndianCopy(writeTo)
		}
	}
}

func BenchmarkWriteBytes_Unsafe(b *testing.B) {
	b.StopTimer()
	cases := randomUints()
	writeTo := make([]byte, 8)
	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			tc := cases[rand.Intn(1000)]
			writeUnsafe(tc, writeTo)
			globalGot = readUnsafe(writeTo)
		}
	}
}

func BenchmarkWriteBytes_BinaryLibrary(b *testing.B) {
	b.StopTimer()
	cases := randomUints()
	writeTo := make([]byte, 8)
	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			tc := cases[rand.Intn(1000)]
			binary.LittleEndian.PutUint64(writeTo, tc)
			globalGot = binary.LittleEndian.Uint64(writeTo)
		}
	}
}

var globalGot uint64

func TestWriteBytesEquivalent(t *testing.T) {
	cases := make([]uint64, 1000)
	for i := 0; i < len(cases); i++ {
		cases[i] = uint64(rand.Intn(math.MaxInt64))
	}
	writeToA := make([]byte, 8)
	writeToB := make([]byte, 8)
	for i := 0; i < 1000; i++ {
		tc := cases[i]
		writeBinaryLittleEndianCopy(tc, writeToA)
		writeUnsafe(tc, writeToB)
		if !bytes.Equal(writeToA, writeToB) {
			t.Fail()
		}
	}
}

func writeBinaryLittleEndianCopy(in uint64, data []byte) {
	_ = data[7]
	data[0] = byte(in)
	data[1] = byte(in >> 8)
	data[2] = byte(in >> 16)
	data[3] = byte(in >> 24)
	data[4] = byte(in >> 32)
	data[5] = byte(in >> 40)
	data[6] = byte(in >> 48)
	data[7] = byte(in >> 56)
}

func readBinaryLittleEndianCopy(data []byte) uint64 {
	_ = data[7]
	return uint64(data[0]) | uint64(data[1])<<8 | uint64(data[2])<<16 | uint64(data[3])<<24 |
		uint64(data[4])<<32 | uint64(data[5])<<40 | uint64(data[6])<<48 | uint64(data[7])<<56
}

func writeUnsafe(in uint64, data []byte) {
	_ = data[7]
	*(*uint64)(unsafe.Pointer(&data[0])) = in
}

func readUnsafe(data []byte) uint64 {
	_ = data[7]
	b := *(*uint64)(unsafe.Pointer(&data[0]))
	return b
}

func randomUints_16() [1000]uint16 {
	cases := [1000]uint16{}
	for i := 0; i < len(cases); i++ {
		cases[i] = uint16(rand.Intn(math.MaxInt16))
	}
	return cases
}

func BenchmarkWriteBytes_BinaryLittleEndianCopy_16(b *testing.B) {
	b.StopTimer()
	cases := randomUints_16()
	writeTo := make([]byte, 2)
	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		// arbitrarily making this longer so the benchmark runs longer
		for j := 0; j < 100; j++ {
			tc := cases[rand.Intn(1000)]
			writeBinaryLittleEndianCopy_16(tc, writeTo)
			globalGot_16 = readBinaryLittleEndianCopy_16(writeTo)
		}
	}
}

func BenchmarkWriteBytes_Unsafe_16(b *testing.B) {
	b.StopTimer()
	cases := randomUints_16()
	writeTo := make([]byte, 8)
	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			tc := cases[rand.Intn(1000)]
			writeUnsafe_16(tc, writeTo)
			globalGot_16 = readUnsafe_16(writeTo)
		}
	}
}

func BenchmarkWriteBytes_BinaryLibrary_16(b *testing.B) {
	b.StopTimer()
	cases := randomUints_16()
	writeTo := make([]byte, 2)
	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			tc := cases[rand.Intn(1000)]
			binary.LittleEndian.PutUint16(writeTo, tc)
			globalGot_16 = binary.LittleEndian.Uint16(writeTo)
		}
	}
}

var globalGot_16 uint16

func writeBinaryLittleEndianCopy_16(in uint16, data []byte) {
	_ = data[1]
	data[0] = byte(in)
	data[1] = byte(in >> 8)
}

func readBinaryLittleEndianCopy_16(data []byte) uint16 {
	_ = data[1]
	return uint16(data[0]) | uint16(data[1])<<8
}

func writeUnsafe_16(in uint16, data []byte) {
	_ = data[1]
	*(*uint16)(unsafe.Pointer(&data[0])) = in
}

func readUnsafe_16(data []byte) uint16 {
	_ = data[1]
	b := *(*uint16)(unsafe.Pointer(&data[0]))
	return b
}
