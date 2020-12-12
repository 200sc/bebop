package generated_test

import (
	"fmt"
	"testing"

	"github.com/200sc/bebop"
	"github.com/200sc/bebop/testdata/generated"
)

func TestMarshalCycleRecords(t *testing.T) {
	type testCase struct {
		name   string
		record bebop.Record
		// notable bug: unmarshalling to a non-empty record
		// causes random behavior based on the field types of the record.
		unmarshalTo bebop.Record
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
		name:        "empty BasicTypes",
		record:      &generated.BasicTypes{},
		unmarshalTo: &generated.BasicTypes{},
	}, {
		name:        "empty DocS",
		record:      &generated.DocS{},
		unmarshalTo: &generated.DocS{},
	}, {
		name:        "empty DepM",
		record:      &generated.DepM{},
		unmarshalTo: &generated.DepM{},
	}, {
		name:        "empty DocM",
		record:      &generated.DocM{},
		unmarshalTo: &generated.DocM{},
	}, {
		name:        "empty Foo",
		record:      &generated.Foo{},
		unmarshalTo: &generated.Foo{},
	}, {
		name:        "empty Bar",
		record:      &generated.Bar{},
		unmarshalTo: &generated.Bar{},
	}, {
		name:        "empty Musician",
		record:      &generated.Musician{},
		unmarshalTo: &generated.Musician{},
	}, {
		name:        "empty Library",
		record:      &generated.Library{},
		unmarshalTo: &generated.Library{},
	}, {
		name:        "empty Song",
		record:      &generated.Song{},
		unmarshalTo: &generated.Song{},
	}, {
		name:        "empty VideoData",
		record:      &generated.VideoData{},
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
		name:        "empty SkipTestNew",
		record:      &generated.SkipTestNew{},
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
		name:        "empty S",
		record:      &generated.S{},
		unmarshalTo: &generated.S{},
	}, {
		name:        "empty SomeMaps",
		record:      &generated.SomeMaps{},
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
		name:        "empty Furniture",
		record:      &generated.Furniture{},
		unmarshalTo: &generated.Furniture{},
	}, {
		name:        "empty RequestResponse",
		record:      &generated.RequestResponse{},
		unmarshalTo: &generated.RequestResponse{},
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
		})
	}
}
