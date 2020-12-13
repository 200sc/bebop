package generated_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/200sc/bebop"
	"github.com/200sc/bebop/testdata/generated"
	"github.com/google/go-cmp/cmp"
)

func TestMarshalCycleRecords(t *testing.T) {
	type testCase struct {
		name   string
		record bebop.Record
		// notable bug: unmarshalling to a non-empty record
		// causes random behavior based on the field types of the record.
		unmarshalTo  bebop.Record
		skipEquality bool
	}
	tcs := []testCase{{
		name:        "empty ArrayOfStrings",
		record:      &generated.ArrayOfStrings{},
		unmarshalTo: &generated.ArrayOfStrings{},
	}, {
		name: "ArrayOfStrings",
		record: &generated.ArrayOfStrings{
			Strings: []string{
				"hello",
				"world",
			},
		},
		unmarshalTo: &generated.ArrayOfStrings{},
	}, {
		name:        "empty BasicArrays",
		record:      &generated.BasicArrays{},
		unmarshalTo: &generated.BasicArrays{},
	}, {
		name: "BasicArrays",
		record: &generated.BasicArrays{
			A_bool:   []bool{true, false, true},
			A_uint16: []uint16{0, 2, 65535},
		},
		unmarshalTo: &generated.BasicArrays{},
	}, {
		name:        "empty TestInt32Array",
		record:      &generated.TestInt32Array{},
		unmarshalTo: &generated.TestInt32Array{},
	}, {
		name: "TestInt32Array",
		record: &generated.TestInt32Array{
			A: []int32{
				0, 2, 15412, 301523, 3441213,
			},
		},
		unmarshalTo: &generated.TestInt32Array{},
	}, {
		name:        "empty BasicTypes",
		record:      &generated.BasicTypes{},
		unmarshalTo: &generated.BasicTypes{},
	}, {
		name: "BasicTypes",
		record: &generated.BasicTypes{
			A_bool:    true,
			A_byte:    4,
			A_int16:   330,
			A_date:    time.Unix(444444, 0).UTC(),
			A_float64: 3.3333,
		},
		unmarshalTo: &generated.BasicTypes{},
	}, {
		name:        "empty DocS",
		record:      &generated.DocS{},
		unmarshalTo: &generated.DocS{},
	}, {
		name: "DocS",
		record: &generated.DocS{
			X: 203202003,
		},
		unmarshalTo: &generated.DocS{},
	}, {
		name:        "empty DepM",
		record:      &generated.DepM{},
		unmarshalTo: &generated.DepM{},
	}, {
		name: "DepM",
		record: &generated.DepM{
			X: int32p(444),
		},
		unmarshalTo: &generated.DepM{},
	}, {
		name:        "empty DocM",
		record:      &generated.DocM{},
		unmarshalTo: &generated.DocM{},
	}, {
		name: "DocM",
		record: &generated.DocM{
			X: int32p(14123),
			Y: int32p(12314502),
			Z: int32p(-2),
		},
		unmarshalTo: &generated.DocM{},
	}, {
		name:        "empty Foo",
		record:      &generated.Foo{},
		unmarshalTo: &generated.Foo{},
	}, {
		name: "Foo",
		record: &generated.Foo{
			Bar: generated.Bar{
				X: float64p(3.21312),
				Y: float64p(3.21333312),
				Z: float64p(3.21312421),
			},
		},
		unmarshalTo: &generated.Foo{},
	}, {
		name:        "empty Bar",
		record:      &generated.Bar{},
		unmarshalTo: &generated.Bar{},
	}, {
		name: "Bar",
		record: &generated.Bar{
			Y: float64p(19999999999999999.2),
		},
		unmarshalTo: &generated.Bar{},
	}, {
		name:         "empty Musician",
		record:       &generated.Musician{},
		unmarshalTo:  &generated.Musician{},
		skipEquality: true,
	}, {
		name: "empty Library",
		record: &generated.Library{
			Songs: map[[16]uint8]generated.Song{},
		},
		unmarshalTo: &generated.Library{},
	}, {
		name: "Library",
		record: &generated.Library{
			Songs: map[[16]byte]generated.Song{
				{0x35, 0x91, 0x8b, 0xc9, 0x19, 0x6d, 0x40, 0xea, 0x97, 0x79, 0x88, 0x9d, 0x79, 0xb7, 0x53, 0xf0}: {
					Title: stringp("song-title"),
					Year:  uint16p(2034),
				},
			},
		},
		unmarshalTo: &generated.Library{},
	}, {
		name:        "empty Song",
		record:      &generated.Song{},
		unmarshalTo: &generated.Song{},
	}, {name: "Song",
		record: &generated.Song{
			Title: stringp("song-title2"),
			Year:  uint16p(20342),
		},
		unmarshalTo: &generated.Song{},
	}, {
		name:        "empty VideoData",
		record:      &generated.VideoData{},
		unmarshalTo: &generated.VideoData{},
	}, {
		name: "VideoData",
		record: &generated.VideoData{
			Time:     -2042.122,
			Width:    9333,
			Height:   123,
			Fragment: []byte{0, 123, 5, 1, 3, 50, 123, 3, 3, 3, 3, 3},
		},
		unmarshalTo: &generated.VideoData{},
	}, {
		name:        "empty MediaMessage",
		record:      &generated.MediaMessage{},
		unmarshalTo: &generated.MediaMessage{},
	}, {
		name:        "empty SkipTestOld",
		record:      &generated.SkipTestOld{},
		unmarshalTo: &generated.SkipTestOld{},
	}, {
		name: "SkipTestOld",
		record: &generated.SkipTestOld{
			X: int32p(2222),
			Y: int32p(12315),
		},
		unmarshalTo: &generated.SkipTestOld{},
	}, {
		name:        "empty SkipTestNew",
		record:      &generated.SkipTestNew{},
		unmarshalTo: &generated.SkipTestNew{},
	}, {
		name: "SkipTestNew",
		record: &generated.SkipTestNew{
			X: int32p(222322),
			Y: int32p(123125),
			Z: int32p(-12344444),
		},
		unmarshalTo: &generated.SkipTestNew{},
	}, {
		name:        "empty SkipTestOldContainer",
		record:      &generated.SkipTestOldContainer{},
		unmarshalTo: &generated.SkipTestOldContainer{},
	}, {
		name:        "empty SkipTestNewContainer",
		record:      &generated.SkipTestNewContainer{},
		unmarshalTo: &generated.SkipTestNewContainer{},
	}, {
		name:         "empty S",
		record:       &generated.S{},
		unmarshalTo:  &generated.S{},
		skipEquality: true,
	}, {
		name: "empty SomeMaps",
		record: &generated.SomeMaps{
			M1: map[bool]bool{},
			M2: map[string]map[string]string{},
			M5: map[[16]byte]generated.M{},
		},
		unmarshalTo: &generated.SomeMaps{},
	}, {
		name:        "empty M",
		record:      &generated.M{},
		unmarshalTo: &generated.M{},
	}, {
		name:        "empty MsgpackComparison",
		record:      &generated.MsgpackComparison{},
		unmarshalTo: &generated.MsgpackComparison{},
	}, {
		name:         "empty Furniture",
		record:       &generated.Furniture{},
		unmarshalTo:  &generated.Furniture{},
		skipEquality: true,
	}, {
		name:         "empty RequestResponse",
		record:       &generated.RequestResponse{},
		unmarshalTo:  &generated.RequestResponse{},
		skipEquality: true,
	}, {
		name:        "empty RequestCatalog",
		record:      &generated.RequestCatalog{},
		unmarshalTo: &generated.RequestCatalog{},
	}}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			marshalData, err := bebop.Marshal(tc.record)
			if err != nil {
				t.Fatalf("initial marshal failed: %v", err)
			}
			err = bebop.Unmarshal(marshalData, tc.unmarshalTo)
			if err != nil {
				t.Fatalf("initial unmarshal failed: %v", err)
			}
			marshalData2, err := bebop.Marshal(tc.unmarshalTo)
			if err != nil {
				t.Fatalf("second marshal failed: %v", err)
			}
			// casting to string for easy equality
			if string(marshalData) != string(marshalData2) {
				fmt.Println(marshalData)
				fmt.Println(marshalData2)
				t.Fatal("second marshal did not have same bytes as first")
			}
			if !tc.skipEquality {
				if diff := cmp.Diff(tc.record, tc.unmarshalTo); diff != "" {
					fmt.Println(diff)
					t.Fatal("original did not match unmarshaled")
				}
			}
		})
	}
}

func int32p(i int32) *int32 {
	return &i
}

func float64p(f float64) *float64 {
	return &f
}

func stringp(s string) *string {
	return &s
}

func uint16p(i uint16) *uint16 {
	return &i
}
