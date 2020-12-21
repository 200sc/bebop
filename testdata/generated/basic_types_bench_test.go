package generated_test

import (
	"encoding/json"
	"testing"

	"github.com/200sc/bebop/testdata/generated"
	"github.com/200sc/bebop/testdata/generated/protos"
	"google.golang.org/protobuf/proto"
)

var benchTypes = &generated.BasicTypes{
	A_bool:    true,
	A_byte:    8,
	A_int16:   8,
	A_uint16:  8,
	A_int32:   234436345,
	A_uint32:  33453566,
	A_int64:   34634566,
	A_uint64:  8,
	A_float32: 7,
	A_float64: 7,
	A_string:  "0123151234123123",
	A_guid:    [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
}
var benchTypesBytes = []byte{1, 8, 8, 0, 8, 0, 249, 54, 249, 13, 254, 117, 254, 1, 70, 123, 16, 2, 0, 0, 0, 0, 8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 224, 64, 0, 0, 0, 0, 0, 0, 28, 64, 16, 0, 0, 0, 48, 49, 50, 51, 49, 53, 49, 50, 51, 52, 49, 50, 51, 49, 50, 51, 3, 2, 1, 0, 5, 4, 7, 6, 8, 9, 10, 11, 12, 13, 14, 15, 0, 0, 0, 0, 0, 0, 0, 0}

var benchTypes2 *generated.BasicTypes

func BenchmarkMarshalBasicTypes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		out = benchTypes.MarshalBebop()
	}
}

func BenchmarkUnmarshalBasicTypes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchTypes2 = &generated.BasicTypes{}
		benchTypes2.MustUnmarshalBebop(benchTypesBytes)
	}
}

func BenchmarkUnmarshalSafeBasicTypes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchTypes2 = &generated.BasicTypes{}
		benchTypes2.UnmarshalBebop(benchTypesBytes)
	}
}

var basicTypesJSONBytes = []byte(`{"A_bool":true,"A_byte":8,"A_int16":8,"A_uint16":8,"A_int32":234436345,"A_uint32":33453566,"A_int64":34634566,"A_uint64":8,"A_float32":7,"A_float64":7,"A_string":"0123151234123123","A_guid":[0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15],"A_date":"0001-01-01T00:00:00Z"}`)

func BenchmarkMarshalBasicTypesJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		out, _ = json.Marshal(benchTypes)
	}
}

func BenchmarkUnmarshalBasicTypesJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		json.Unmarshal(basicTypesJSONBytes, benchTypes2)
	}
}

var basicTypesProto = &protos.BasicTypes{
	ABool:    true,
	AInt16:   8,
	AUint16:  8,
	AInt32:   234436345,
	AUint32:  33453566,
	AInt64:   5346345345,
	AUint64:  8,
	AFloat32: 7,
	AFloat64: 7,
	AString:  "0123151234123123",
	AGuid:    []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
}

func BenchmarkMarshalBasicTypesProto(b *testing.B) {
	for i := 0; i < b.N; i++ {
		out, _ = proto.Marshal(basicTypesProto)
	}
}

var basicTypesProtoBytes = []byte{8, 1, 24, 8, 32, 8, 40, 249, 237, 228, 111, 48, 254, 235, 249, 15, 56, 129, 131, 171, 245, 19, 64, 8, 77, 0, 0, 224, 64, 81, 0, 0, 0, 0, 0, 0, 28, 64, 90, 16, 48, 49, 50, 51, 49, 53, 49, 50, 51, 52, 49, 50, 51, 49, 50, 51, 98, 16, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
var basicTypesProto2 = &protos.BasicTypes{}

func BenchmarkUnmarshalBasicTypesProto(b *testing.B) {
	for i := 0; i < b.N; i++ {
		proto.Unmarshal(basicTypesProtoBytes, basicTypesProto2)
	}
}
