package generated_test

import (
	"bytes"
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
		unmarshalTo2 bebop.Record
		skipEquality bool
	}
	tcs := []testCase{{
		name: "empty ArrayOfStrings",
		record: &generated.ArrayOfStrings{
			Strings: []string{},
		},
		unmarshalTo:  &generated.ArrayOfStrings{},
		unmarshalTo2: &generated.ArrayOfStrings{},
	}, {
		name: "ArrayOfStrings",
		record: &generated.ArrayOfStrings{
			Strings: []string{
				"hello",
				"world",
			},
		},
		unmarshalTo:  &generated.ArrayOfStrings{},
		unmarshalTo2: &generated.ArrayOfStrings{},
	}, {
		name: "empty BasicArrays",
		record: &generated.BasicArrays{
			A_bool:    []bool{},
			A_byte:    []byte{},
			A_int16:   []int16{},
			A_uint16:  []uint16{},
			A_int32:   []int32{},
			A_uint32:  []uint32{},
			A_int64:   []int64{},
			A_uint64:  []uint64{},
			A_float32: []float32{},
			A_float64: []float64{},
			A_string:  []string{},
			A_guid:    [][16]byte{},
		},
		unmarshalTo:  &generated.BasicArrays{},
		unmarshalTo2: &generated.BasicArrays{},
	}, {
		name: "BasicArrays",
		record: &generated.BasicArrays{
			A_bool:    []bool{true, false, true},
			A_uint16:  []uint16{0, 2, 65535},
			A_byte:    []byte{},
			A_int16:   []int16{},
			A_int32:   []int32{},
			A_uint32:  []uint32{},
			A_int64:   []int64{},
			A_uint64:  []uint64{},
			A_float32: []float32{},
			A_float64: []float64{},
			A_string:  []string{},
			A_guid:    [][16]byte{},
		},
		unmarshalTo:  &generated.BasicArrays{},
		unmarshalTo2: &generated.BasicArrays{},
	}, {
		name: "empty TestInt32Array",
		record: &generated.TestInt32Array{
			A: []int32{},
		},
		unmarshalTo:  &generated.TestInt32Array{},
		unmarshalTo2: &generated.TestInt32Array{},
	}, {
		name: "TestInt32Array",
		record: &generated.TestInt32Array{
			A: []int32{
				0, 2, 15412, 301523, 3441213,
			},
		},
		unmarshalTo:  &generated.TestInt32Array{},
		unmarshalTo2: &generated.TestInt32Array{},
	}, {
		name:         "empty BasicTypes",
		record:       &generated.BasicTypes{},
		unmarshalTo:  &generated.BasicTypes{},
		unmarshalTo2: &generated.BasicTypes{},
	}, {
		name: "BasicTypes",
		record: &generated.BasicTypes{
			A_bool:    true,
			A_byte:    4,
			A_int16:   330,
			A_date:    time.Unix(444444, 0).UTC(),
			A_float64: 3.3333,
		},
		unmarshalTo:  &generated.BasicTypes{},
		unmarshalTo2: &generated.BasicTypes{},
	}, {
		name:         "empty DocS",
		record:       &generated.DocS{},
		unmarshalTo:  &generated.DocS{},
		unmarshalTo2: &generated.DocS{},
	}, {
		name: "DocS",
		record: &generated.DocS{
			X: 203202003,
		},
		unmarshalTo:  &generated.DocS{},
		unmarshalTo2: &generated.DocS{},
	}, {
		name:         "empty DepM",
		record:       &generated.DepM{},
		unmarshalTo:  &generated.DepM{},
		unmarshalTo2: &generated.DepM{},
	}, {
		name: "DepM",
		record: &generated.DepM{
			X: int32p(444),
		},
		unmarshalTo:  &generated.DepM{},
		unmarshalTo2: &generated.DepM{},
	}, {
		name:         "empty DocM",
		record:       &generated.DocM{},
		unmarshalTo:  &generated.DocM{},
		unmarshalTo2: &generated.DocM{},
	}, {
		name: "DocM",
		record: &generated.DocM{
			X: int32p(14123),
			Y: int32p(12314502),
			Z: int32p(-2),
		},
		unmarshalTo:  &generated.DocM{},
		unmarshalTo2: &generated.DocM{},
	}, {
		name:         "empty Foo",
		record:       &generated.Foo{},
		unmarshalTo:  &generated.Foo{},
		unmarshalTo2: &generated.Foo{},
	}, {
		name: "Foo",
		record: &generated.Foo{
			Bar: generated.Bar{
				X: float64p(3.21312),
				Y: float64p(3.21333312),
				Z: float64p(3.21312421),
			},
		},
		unmarshalTo:  &generated.Foo{},
		unmarshalTo2: &generated.Foo{},
	}, {
		name:         "empty Bar",
		record:       &generated.Bar{},
		unmarshalTo:  &generated.Bar{},
		unmarshalTo2: &generated.Bar{},
	}, {
		name: "Bar",
		record: &generated.Bar{
			Y: float64p(19999999999999999.2),
		},
		unmarshalTo:  &generated.Bar{},
		unmarshalTo2: &generated.Bar{},
	}, {
		name:         "empty Musician",
		record:       &generated.Musician{},
		unmarshalTo:  &generated.Musician{},
		unmarshalTo2: &generated.Musician{},
		skipEquality: true,
	}, {
		name: "empty Library",
		record: &generated.Library{
			Songs: map[[16]uint8]generated.Song{},
		},
		unmarshalTo:  &generated.Library{},
		unmarshalTo2: &generated.Library{},
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
		unmarshalTo:  &generated.Library{},
		unmarshalTo2: &generated.Library{},
	}, {
		name:         "empty Song",
		record:       &generated.Song{},
		unmarshalTo:  &generated.Song{},
		unmarshalTo2: &generated.Song{},
	}, {name: "Song",
		record: &generated.Song{
			Title: stringp("song-title2"),
			Year:  uint16p(20342),
		},
		unmarshalTo:  &generated.Song{},
		unmarshalTo2: &generated.Song{},
	}, {
		name: "empty VideoData",
		record: &generated.VideoData{
			Fragment: []byte{},
		},
		unmarshalTo:  &generated.VideoData{},
		unmarshalTo2: &generated.VideoData{},
	}, {
		name: "VideoData",
		record: &generated.VideoData{
			Time:     -2042.122,
			Width:    9333,
			Height:   123,
			Fragment: []byte{0, 123, 5, 1, 3, 50, 123, 3, 3, 3, 3, 3},
		},
		unmarshalTo:  &generated.VideoData{},
		unmarshalTo2: &generated.VideoData{},
	}, {
		name:         "empty MediaMessage",
		record:       &generated.MediaMessage{},
		unmarshalTo:  &generated.MediaMessage{},
		unmarshalTo2: &generated.MediaMessage{},
	}, {
		name:         "empty SkipTestOld",
		record:       &generated.SkipTestOld{},
		unmarshalTo:  &generated.SkipTestOld{},
		unmarshalTo2: &generated.SkipTestOld{},
	}, {
		name: "SkipTestOld",
		record: &generated.SkipTestOld{
			X: int32p(2222),
			Y: int32p(12315),
		},
		unmarshalTo:  &generated.SkipTestOld{},
		unmarshalTo2: &generated.SkipTestOld{},
	}, {
		name:         "empty SkipTestNew",
		record:       &generated.SkipTestNew{},
		unmarshalTo:  &generated.SkipTestNew{},
		unmarshalTo2: &generated.SkipTestNew{},
	}, {
		name: "SkipTestNew",
		record: &generated.SkipTestNew{
			X: int32p(222322),
			Y: int32p(123125),
			Z: int32p(-12344444),
		},
		unmarshalTo:  &generated.SkipTestNew{},
		unmarshalTo2: &generated.SkipTestNew{},
	}, {
		name:         "empty SkipTestOldContainer",
		record:       &generated.SkipTestOldContainer{},
		unmarshalTo:  &generated.SkipTestOldContainer{},
		unmarshalTo2: &generated.SkipTestOldContainer{},
	}, {
		name:         "empty SkipTestNewContainer",
		record:       &generated.SkipTestNewContainer{},
		unmarshalTo:  &generated.SkipTestNewContainer{},
		unmarshalTo2: &generated.SkipTestNewContainer{},
	}, {
		name:         "empty S",
		record:       &generated.S{},
		unmarshalTo:  &generated.S{},
		unmarshalTo2: &generated.S{},
		skipEquality: true,
	}, {
		name: "empty SomeMaps",
		record: &generated.SomeMaps{
			M1: map[bool]bool{},
			M2: map[string]map[string]string{},
			M3: []map[int32][]map[bool]generated.S{},
			M4: []map[string][]float32{},
			M5: map[[16]byte]generated.M{},
		},
		unmarshalTo:  &generated.SomeMaps{},
		unmarshalTo2: &generated.SomeMaps{},
	}, {
		name: "SomeMaps1",
		record: &generated.SomeMaps{
			M3: []map[int32][]map[bool]generated.S{
				{
					0: []map[bool]generated.S{{
						true: generated.S{},
					}},
				},
			},
		},
		unmarshalTo:  &generated.SomeMaps{},
		unmarshalTo2: &generated.SomeMaps{},
		skipEquality: true,
	},
		// we can't do some maps 2 because it contains maps with more than one element, whose order is marshalled randomly.
		// {
		// 	name: "SomeMaps2",
		// 	record: &generated.SomeMaps{
		// 		M1: map[bool]bool{
		// 			true: true,
		// 		},
		// 		M2: map[string]map[string]string{
		// 			"mario": map[string]string{
		// 				"mario":    "",
		// 				"luigi":    "",
		// 				"brothers": "",
		// 			},
		// 		},
		// 		M3: []map[int32][]map[bool]generated.S{
		// 			{
		// 				0: []map[bool]generated.S{{
		// 					true: generated.S{},
		// 				}},
		// 			}, {
		// 				2: []map[bool]generated.S{{
		// 					true: generated.S{},
		// 				}},
		// 			}, {
		// 				41111: []map[bool]generated.S{{
		// 					false: generated.S{},
		// 				}},
		// 			},
		// 		},
		// 		M4: []map[string][]float32{
		// 			{
		// 				"a": []float32{1321, 1423, 1423, 540, 12314, 1231, 4123, 1412, 1230, 4123, 123},
		// 			},
		// 		},
		// 		M5: map[[16]byte]generated.M{
		// 			[16]byte{5: 3}: generated.M{B: float64p(0.0000002)},
		// 		},
		// 	},
		// 	unmarshalTo:  &generated.SomeMaps{},
		//  unmarshalTo2:  &generated.SomeMaps{},
		// 	skipEquality: true,
		// },
		{
			name:         "empty M",
			record:       &generated.M{},
			unmarshalTo:  &generated.M{},
			unmarshalTo2: &generated.M{},
		}, {
			name: "empty MsgpackComparison",
			record: &generated.MsgpackComparison{
				ARRAY0: []int32{},
				ARRAY1: []string{},
				ARRAY8: []int32{},
			},
			unmarshalTo:  &generated.MsgpackComparison{},
			unmarshalTo2: &generated.MsgpackComparison{},
		}, {
			name:         "empty Furniture",
			record:       &generated.Furniture{},
			unmarshalTo:  &generated.Furniture{},
			unmarshalTo2: &generated.Furniture{},
			skipEquality: true,
		}, {
			name:         "empty RequestResponse",
			record:       &generated.RequestResponse{},
			unmarshalTo:  &generated.RequestResponse{},
			unmarshalTo2: &generated.RequestResponse{},
			skipEquality: true,
		}, {
			name:         "empty RequestCatalog",
			record:       &generated.RequestCatalog{},
			unmarshalTo:  &generated.RequestCatalog{},
			unmarshalTo2: &generated.RequestCatalog{},
		}, {
			name:         "empty ReadOnlyMap",
			record:       &generated.ReadOnlyMap{},
			unmarshalTo:  &generated.ReadOnlyMap{},
			unmarshalTo2: &generated.ReadOnlyMap{},
			skipEquality: true,
			// Empty unions are not supported
			// }, {
			// 	name:         "empty Union U",
			// 	record:       &generated.U{},
			// 	unmarshalTo:  &generated.U{},
			// 	unmarshalTo2: &generated.U{},
		}, {
			name: "Union U: A",
			record: &generated.U{
				A: &generated.A{
					A: uint32p(2),
				},
			},
			unmarshalTo:  &generated.U{},
			unmarshalTo2: &generated.U{},
		}, {
			name: "Union U: B",
			record: &generated.U{
				B: &generated.B{
					B: true,
				},
			},
			unmarshalTo:  &generated.U{},
			unmarshalTo2: &generated.U{},
		}, {
			name: "Union U: C",
			record: &generated.U{
				C: &generated.C{},
			},
			unmarshalTo:  &generated.U{},
			unmarshalTo2: &generated.U{},
		}, {
			name: "Union U: WD",
			record: &generated.U{
				W: &generated.W{
					D: &generated.D{
						S: "first",
					},
				},
			},
			unmarshalTo:  &generated.U{},
			unmarshalTo2: &generated.U{},
		}, {
			name: "Union U: WX",
			record: &generated.U{
				W: &generated.W{
					X: &generated.X{
						X: true,
					},
				},
			},
			unmarshalTo:  &generated.U{},
			unmarshalTo2: &generated.U{},
		}, {
			name: "Union List",
			record: &generated.List{
				Cons: &generated.Cons{
					Head: 1,
					Tail: generated.List{
						Cons: &generated.Cons{
							Head: 2,
							Tail: generated.List{
								Nil: &generated.Nil{},
							},
						},
					},
				},
			},
			unmarshalTo:  &generated.List{},
			unmarshalTo2: &generated.List{},
		}, {
			name: "Date MyObj",
			// compilation + empty value tests-- times are rounded down
			// so comparing them exactly is tricky
			record: &generated.MyObj{
				Start: nowP(),
				End:   nowP(),
			},
			unmarshalTo:  &generated.MyObj{},
			unmarshalTo2: &generated.MyObj{},
		}}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			err := tc.record.EncodeBebop(buf)
			marshalData := buf.Bytes()
			if err != nil {
				t.Fatalf("initial marshal failed: %v", err)
			}
			err = tc.unmarshalTo.DecodeBebop(bytes.NewBuffer(marshalData))
			if err != nil {
				t.Fatalf("initial unmarshal failed: %v", err)
			}
			buf = &bytes.Buffer{}
			err = tc.unmarshalTo.EncodeBebop(buf)
			marshalData2 := buf.Bytes()
			if err != nil {
				t.Fatalf("second marshal failed: %v", err)
			}
			// casting to string for easy equality
			if string(marshalData) != string(marshalData2) {
				fmt.Println(tc.unmarshalTo)
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
			noWriterMarshal := tc.record.MarshalBebop()
			if string(marshalData) != string(noWriterMarshal) {
				fmt.Println(marshalData)
				fmt.Println(noWriterMarshal)
				t.Fatal("no-writer marshal did not have same bytes as with-writer")
			}
			err = tc.unmarshalTo2.UnmarshalBebop(noWriterMarshal)
			if err != nil {
				t.Fatalf("second unmarshal failed: %v", err)
			}
			if !tc.skipEquality {
				if diff := cmp.Diff(tc.record, tc.unmarshalTo2); diff != "" {
					fmt.Println(diff)
					t.Fatal("original did not match unmarshaled")
				}
			}
			marshalData3 := tc.unmarshalTo2.MarshalBebop()
			if string(marshalData) != string(marshalData3) {
				fmt.Println(marshalData)
				fmt.Println(marshalData3)
				t.Fatal("no-writer unmarshal did not have same bytes as prior unmarshals")
			}

			type MustUnmarshaler interface {
				MustUnmarshalBebop([]byte)
			}

			if mu, ok := tc.unmarshalTo2.(MustUnmarshaler); ok {
				mu.MustUnmarshalBebop(marshalData3)
				marshalData4 := tc.unmarshalTo2.MarshalBebop()
				if string(marshalData) != string(marshalData4) {
					fmt.Println(marshalData)
					fmt.Println(marshalData4)
					t.Fatal("must unmarshal did not have same bytes as prior unmarshals")
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

func uint32p(i uint32) *uint32 {
	return &i
}

func nowP() *time.Time {
	t := time.Now()
	return &t
}
