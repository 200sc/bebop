package generated_test

import (
	"testing"

	"github.com/200sc/bebop/testdata/generated"
)

var benchMaps = &generated.SomeMaps{
	M1: map[bool]bool{
		true:  true,
		false: false,
	},
	M2: map[string]map[string]string{
		"hello": {
			"world": "!",
		},
		"foo": {
			"bar": "bizz",
		},
		"ursula": {
			"k": "leguin",
		},
		"mario": {
			"mario":    "",
			"luigi":    "",
			"brothers": "",
		},
	},
	M3: []map[int32][]map[bool]generated.S{
		{
			0: []map[bool]generated.S{{
				true:  generated.S{},
				false: generated.S{},
			}},
		}, {
			1: []map[bool]generated.S{{
				true:  generated.S{},
				false: generated.S{},
			}},
			2: []map[bool]generated.S{{
				true:  generated.S{},
				false: generated.S{},
			}},
			3: []map[bool]generated.S{{
				true:  generated.S{},
				false: generated.S{},
			}},
		}, {
			41111: []map[bool]generated.S{{
				true:  generated.S{},
				false: generated.S{},
			}},
		},
	},
	M4: []map[string][]float32{
		{
			"a": []float32{1321, 1423, 1423, 540, 12314, 1231, 4123, 1412, 1230, 4123, 123},
		},
	},
	M5: map[[16]byte]generated.M{
		{5: 3}: {B: float64p(0.0000002)},
	},
}

var benchMapsBytes = []byte{2, 0, 0, 0, 1, 1, 0, 0, 4, 0, 0, 0, 5, 0, 0, 0, 104, 101, 108, 108, 111, 1, 0, 0, 0, 5, 0,
	0, 0, 119, 111, 114, 108, 100, 1, 0, 0, 0, 33, 3, 0, 0, 0, 102, 111, 111, 1, 0, 0, 0, 3, 0, 0, 0, 98, 97, 114, 4, 0,
	0, 0, 98, 105, 122, 122, 6, 0, 0, 0, 117, 114, 115, 117, 108, 97, 1, 0, 0, 0, 1, 0, 0, 0, 107, 6, 0, 0, 0, 108, 101,
	103, 117, 105, 110, 5, 0, 0, 0, 109, 97, 114, 105, 111, 3, 0, 0, 0, 5, 0, 0, 0, 109, 97, 114, 105, 111, 0, 0, 0, 0,
	5, 0, 0, 0, 108, 117, 105, 103, 105, 0, 0, 0, 0, 8, 0, 0, 0, 98, 114, 111, 116, 104, 101, 114, 115, 0, 0, 0, 0, 3,
	0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0,
	0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 2, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 1, 0, 0,
	0, 2, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 1, 0, 0, 0, 2, 0, 0, 0, 1, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 151, 160, 0, 0, 1, 0, 0, 0, 2, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 97, 11, 0, 0, 0, 0, 32, 165, 68, 0, 224, 177,
	68, 0, 224, 177, 68, 0, 0, 7, 68, 0, 104, 64, 70, 0, 224, 153, 68, 0, 216, 128, 69, 0, 128, 176, 68, 0, 192, 153,
	68, 0, 216, 128, 69, 0, 0, 246, 66, 1, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 10, 0, 0, 0, 2, 72,
	175, 188, 154, 242, 215, 138, 62, 0}

var benchMaps2 *generated.SomeMaps

func BenchmarkMarshalSomeMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		out = benchMaps.MarshalBebop()
	}
}

func BenchmarkUnmarshalSomeMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchMaps2 = &generated.SomeMaps{}
		benchMaps2.MustUnmarshalBebop(benchMapsBytes)
	}
}

func BenchmarkUnmarshalSafeSomeMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchMaps2 = &generated.SomeMaps{}
		benchMaps2.UnmarshalBebop(benchMapsBytes)
	}
}

// SomeMap cannot be parsed by json as its structure is unsupported
// Most of SomeMap's fields cannot be encoded in protobuf
