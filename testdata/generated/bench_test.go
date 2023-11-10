package generated_test

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/200sc/bebop/testdata/generated"
)

func BenchmarkArraySamplesMarshalBebopTo(b *testing.B) {
	b.StopTimer()
	b1 := make([][][]byte, 100)
	for i := 0; i < 100; i++ {
		b1[i] = make([][]byte, 100)
		for j := 0; j < 100; j++ {
			b1[i][j] = make([]byte, 10)
			for k := 0; k < 10; k++ {
				b1[i][j][k] = byte(rand.Intn(255))
			}
		}
	}
	b2 := make([][][]byte, 100)
	for i := 0; i < 100; i++ {
		b2[i] = make([][]byte, 100)
		for j := 0; j < 100; j++ {
			b2[i][j] = make([]byte, 10)
			for k := 0; k < 10; k++ {
				b2[i][j][k] = byte(rand.Intn(255))
			}
		}
	}
	v := generated.ArraySamples{
		Bytes:  b1,
		Bytes2: b2,
	}
	out := make([]byte, v.Size())

	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		v.MarshalBebopTo(out)
	}
	benchOut = out
}

var benchOut []byte

func BenchmarkArraySamplesEncodeBebop(b *testing.B) {
	b.StopTimer()
	b1 := make([][][]byte, 100)
	for i := 0; i < 100; i++ {
		b1[i] = make([][]byte, 100)
		for j := 0; j < 100; j++ {
			b1[i][j] = make([]byte, 10)
			for k := 0; k < 10; k++ {
				b1[i][j][k] = byte(rand.Intn(255))
			}
		}
	}
	b2 := make([][][]byte, 100)
	for i := 0; i < 100; i++ {
		b2[i] = make([][]byte, 100)
		for j := 0; j < 100; j++ {
			b2[i][j] = make([]byte, 10)
			for k := 0; k < 10; k++ {
				b2[i][j][k] = byte(rand.Intn(255))
			}
		}
	}
	v := generated.ArraySamples{
		Bytes:  b1,
		Bytes2: b2,
	}

	w := bytes.NewBuffer(make([]byte, 0, v.Size()))

	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		err := v.EncodeBebop(w)
		if err != nil {
			b.Fatal(err)
		}
	}
	benchOut = w.Bytes()
}
