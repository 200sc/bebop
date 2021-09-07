package generated_test

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
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
		unmarshalRecord func() bebop.Record
		skipEquality    bool
		tsName          string
	}
	tcs := []testCase{{
		name:   "empty ArrayOfStrings",
		tsName: "ArrayOfStrings",
		record: &generated.ArrayOfStrings{
			Strings: []string{},
		},
		unmarshalRecord: func() bebop.Record { return &generated.ArrayOfStrings{} },
	}, {
		name:   "ArrayOfStrings",
		tsName: "ArrayOfStrings",
		record: &generated.ArrayOfStrings{
			Strings: []string{
				"hello",
				"world",
			},
		},
		unmarshalRecord: func() bebop.Record { return &generated.ArrayOfStrings{} },
	}, {
		name:   "empty BasicArrays",
		tsName: "BasicArrays",
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
		unmarshalRecord: func() bebop.Record { return &generated.BasicArrays{} },
	}, {
		name: "BasicArrays",
		// fails, js's b64 tells us to create a huge paylaod
		tsName: "BasicArrays",
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
		unmarshalRecord: func() bebop.Record { return &generated.BasicArrays{} },
	}, {
		name:   "empty TestInt32Array",
		tsName: "TestInt32Array",
		record: &generated.TestInt32Array{
			A: []int32{},
		},
		unmarshalRecord: func() bebop.Record { return &generated.TestInt32Array{} },
	}, {
		name: "TestInt32Array",
		// fails js side before during decode
		tsName: "TestInt32Array",
		record: &generated.TestInt32Array{
			A: []int32{
				0, 2, 15412, 301523, 3441213,
			},
		},
		unmarshalRecord: func() bebop.Record { return &generated.TestInt32Array{} },
	}, {
		name:            "empty BasicTypes",
		tsName:          "BasicTypes",
		record:          &generated.BasicTypes{},
		unmarshalRecord: func() bebop.Record { return &generated.BasicTypes{} },
	}, {
		name: "BasicTypes",
		// hangs
		tsName: "BasicTypes",
		record: &generated.BasicTypes{
			A_bool:    true,
			A_byte:    4,
			A_int16:   330,
			A_date:    time.Unix(444444, 0).UTC(),
			A_float64: 3.3333,
		},
		unmarshalRecord: func() bebop.Record { return &generated.BasicTypes{} },
	}, {
		name:            "empty DocS",
		tsName:          "DocS",
		record:          &generated.DocS{},
		unmarshalRecord: func() bebop.Record { return &generated.DocS{} },
	}, {
		name:   "DocS",
		tsName: "DocS",
		// fails js side before during decode
		record: &generated.DocS{
			X: 203202003,
		},
		unmarshalRecord: func() bebop.Record { return &generated.DocS{} },
	}, {
		name:            "empty DepM",
		tsName:          "DepM",
		record:          &generated.DepM{},
		unmarshalRecord: func() bebop.Record { return &generated.DepM{} },
	}, {
		name: "DepM",
		// bebop ts fails to write anything for this type, it will always output
		// [1,0,0,0,0]
		tsName: "DepM",
		record: &generated.DepM{
			X: int32p(444),
		},
		unmarshalRecord: func() bebop.Record { return &generated.DepM{} },
	}, {
		name:            "empty DocM",
		record:          &generated.DocM{},
		unmarshalRecord: func() bebop.Record { return &generated.DocM{} },
	}, {
		name: "DocM",
		record: &generated.DocM{
			X: int32p(14123),
			Y: int32p(12314502),
			Z: int32p(-2),
		},
		unmarshalRecord: func() bebop.Record { return &generated.DocM{} },
	}, {
		name:            "empty Foo",
		record:          &generated.Foo{},
		unmarshalRecord: func() bebop.Record { return &generated.Foo{} },
	}, {
		name: "Foo",
		record: &generated.Foo{
			Bar: generated.Bar{
				X: float64p(3.21312),
				Y: float64p(3.21333312),
				Z: float64p(3.21312421),
			},
		},
		unmarshalRecord: func() bebop.Record { return &generated.Foo{} },
	}, {
		name:            "empty Bar",
		record:          &generated.Bar{},
		unmarshalRecord: func() bebop.Record { return &generated.Bar{} },
	}, {
		name: "Bar",
		record: &generated.Bar{
			Y: float64p(19999999999999999.2),
		},
		unmarshalRecord: func() bebop.Record { return &generated.Bar{} },
	}, {
		name:            "empty Musician",
		record:          &generated.Musician{},
		unmarshalRecord: func() bebop.Record { return &generated.Musician{} },
		skipEquality:    true,
	}, {
		name: "empty Library",
		record: &generated.Library{
			Songs: map[[16]uint8]generated.Song{},
		},
		unmarshalRecord: func() bebop.Record { return &generated.Library{} },
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
		unmarshalRecord: func() bebop.Record { return &generated.Library{} },
	}, {
		name:            "empty Song",
		record:          &generated.Song{},
		unmarshalRecord: func() bebop.Record { return &generated.Song{} },
	}, {name: "Song",
		record: &generated.Song{
			Title: stringp("song-title2"),
			Year:  uint16p(20342),
		},
		unmarshalRecord: func() bebop.Record { return &generated.Song{} },
	}, {
		name: "empty VideoData",
		record: &generated.VideoData{
			Fragment: []byte{},
		},
		unmarshalRecord: func() bebop.Record { return &generated.VideoData{} },
	}, {
		name: "VideoData",
		record: &generated.VideoData{
			Time:     -2042.122,
			Width:    9333,
			Height:   123,
			Fragment: []byte{0, 123, 5, 1, 3, 50, 123, 3, 3, 3, 3, 3},
		},
		unmarshalRecord: func() bebop.Record { return &generated.VideoData{} },
	}, {
		name:            "empty MediaMessage",
		record:          &generated.MediaMessage{},
		unmarshalRecord: func() bebop.Record { return &generated.MediaMessage{} },
	}, {
		name:            "empty SkipTestOld",
		record:          &generated.SkipTestOld{},
		unmarshalRecord: func() bebop.Record { return &generated.SkipTestOld{} },
	}, {
		name: "SkipTestOld",
		record: &generated.SkipTestOld{
			X: int32p(2222),
			Y: int32p(12315),
		},
		unmarshalRecord: func() bebop.Record { return &generated.SkipTestOld{} },
	}, {
		name:            "empty SkipTestNew",
		record:          &generated.SkipTestNew{},
		unmarshalRecord: func() bebop.Record { return &generated.SkipTestNew{} },
	}, {
		name: "SkipTestNew",
		record: &generated.SkipTestNew{
			X: int32p(222322),
			Y: int32p(123125),
			Z: int32p(-12344444),
		},
		unmarshalRecord: func() bebop.Record { return &generated.SkipTestNew{} },
	}, {
		name:            "empty SkipTestOldContainer",
		record:          &generated.SkipTestOldContainer{},
		unmarshalRecord: func() bebop.Record { return &generated.SkipTestOldContainer{} },
	}, {
		name:            "empty SkipTestNewContainer",
		record:          &generated.SkipTestNewContainer{},
		unmarshalRecord: func() bebop.Record { return &generated.SkipTestNewContainer{} },
	}, {
		name:            "empty S",
		record:          &generated.S{},
		unmarshalRecord: func() bebop.Record { return &generated.S{} },
		skipEquality:    true,
	}, {
		name: "empty SomeMaps",
		record: &generated.SomeMaps{
			M1: map[bool]bool{},
			M2: map[string]map[string]string{},
			M3: []map[int32][]map[bool]generated.S{},
			M4: []map[string][]float32{},
			M5: map[[16]byte]generated.M{},
		},
		unmarshalRecord: func() bebop.Record { return &generated.SomeMaps{} },
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
		unmarshalRecord: func() bebop.Record { return &generated.SomeMaps{} },
		skipEquality:    true,
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
		// 	unmarshalRecord: func() bebop.Record { return  &generated.SomeMaps{} },
		//  unmarshalTo2:  &generated.SomeMaps{},
		// 	skipEquality: true,
		// },
		{
			name:            "empty M",
			record:          &generated.M{},
			unmarshalRecord: func() bebop.Record { return &generated.M{} },
		}, {
			name: "empty MsgpackComparison",
			record: &generated.MsgpackComparison{
				ARRAY0: []int32{},
				ARRAY1: []string{},
				ARRAY8: []int32{},
			},
			unmarshalRecord: func() bebop.Record { return &generated.MsgpackComparison{} },
		}, {
			name:            "empty Furniture",
			record:          &generated.Furniture{},
			unmarshalRecord: func() bebop.Record { return &generated.Furniture{} },
			skipEquality:    true,
		}, {
			name:            "empty RequestResponse",
			record:          &generated.RequestResponse{},
			unmarshalRecord: func() bebop.Record { return &generated.RequestResponse{} },
			skipEquality:    true,
		}, {
			name:            "empty RequestCatalog",
			record:          &generated.RequestCatalog{},
			unmarshalRecord: func() bebop.Record { return &generated.RequestCatalog{} },
		}, {
			name:            "empty ReadOnlyMap",
			record:          &generated.ReadOnlyMap{},
			unmarshalRecord: func() bebop.Record { return &generated.ReadOnlyMap{} },
			skipEquality:    true,
			// Empty unions are not supported
			// }, {
			// 	name:         "empty Union U",
			// 	record:       &generated.U{},
			// 	unmarshalRecord: func() bebop.Record { return  &generated.U{} },
			// 	unmarshalTo2: &generated.U{},
		}, {
			name: "Union U: A",
			record: &generated.U{
				A: &generated.A{
					B: uint32p(2),
				},
			},
			unmarshalRecord: func() bebop.Record { return &generated.U{} },
		}, {
			name: "Union U: B",
			record: &generated.U{
				B: &generated.B{
					C: true,
				},
			},
			unmarshalRecord: func() bebop.Record { return &generated.U{} },
		}, {
			name:   "Union U: C",
			tsName: "U",
			record: &generated.U{
				C: &generated.C{},
			},
			unmarshalRecord: func() bebop.Record { return &generated.U{} },
		}, {
			name:   "Union List",
			tsName: "List",
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
			unmarshalRecord: func() bebop.Record { return &generated.List{} },
			// }, {
			// 	name: "Date MyObj",
			// 	// compilation + empty value tests-- times are rounded down
			// 	// so comparing them exactly is tricky
			// 	record: &generated.MyObj{
			// 		Start: nowP(),
			// 		End:   nowP(),
			// 	},
			// 	unmarshalRecord: func() bebop.Record { return  &generated.MyObj{} },
			// 	unmarshalTo2: &generated.MyObj{},
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
			umt := tc.unmarshalRecord()
			err = umt.DecodeBebop(bytes.NewBuffer(marshalData))
			if err != nil {
				t.Fatalf("initial unmarshal failed: %v", err)
			}
			buf = &bytes.Buffer{}
			err = umt.EncodeBebop(buf)
			marshalData2 := buf.Bytes()
			if err != nil {
				t.Fatalf("second marshal failed: %v", err)
			}
			// casting to string for easy equality
			if string(marshalData) != string(marshalData2) {
				fmt.Println(umt)
				fmt.Println(marshalData)
				fmt.Println(marshalData2)
				t.Fatal("second marshal did not have same bytes as first")
			}
			if !tc.skipEquality {
				if diff := cmp.Diff(tc.record, umt); diff != "" {
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
			umt = tc.unmarshalRecord()
			err = umt.UnmarshalBebop(noWriterMarshal)
			if err != nil {
				t.Fatalf("second unmarshal failed: %v", err)
			}
			if !tc.skipEquality {
				if diff := cmp.Diff(tc.record, umt); diff != "" {
					fmt.Println(diff)
					t.Fatal("original did not match unmarshaled")
				}
			}
			marshalData3 := umt.MarshalBebop()
			if string(marshalData) != string(marshalData3) {
				fmt.Println(marshalData)
				fmt.Println(marshalData3)
				t.Fatal("no-writer unmarshal did not have same bytes as prior unmarshals")
			}

			type MustUnmarshaler interface {
				MustUnmarshalBebop([]byte)
			}

			if mu, ok := umt.(MustUnmarshaler); ok {
				mu.MustUnmarshalBebop(marshalData3)
				marshalData4 := umt.MarshalBebop()
				if string(marshalData) != string(marshalData4) {
					fmt.Println(marshalData)
					fmt.Println(marshalData4)
					t.Fatal("must unmarshal did not have same bytes as prior unmarshals")
				}
			}

			if tc.tsName == "" {
				return
			}
			fmt.Println("execing js")

			marshalData5 := umt.MarshalBebop()
			inputB64 := base64.URLEncoding.EncodeToString(marshalData5)
			fmt.Println(marshalData5, inputB64)
			jsQuery := fmt.Sprintf(`

			var Buffer = require('buffer').Buffer

			var ToBase64 = function (u8) {
				return Buffer.from(String.fromCharCode.apply(null, u8)).toString('base64')
			}

			var FromBase64 = function (str) {
				console.log(str)
				let buf = Buffer.from(str, 'base64')
				console.log(buf)
				return buf;
			}

			let bbp = exports.%[1]s.decode(FromBase64(%[2]q))
			console.log(bbp)
			let binary = exports.%[1]s.encode(bbp)
			console.log(binary)
			console.log(ToBase64(binary))
						`, tc.tsName, inputB64)
			cmd := exec.Command("node", `out.js`)
			cmd.Stdin = bytes.NewReader([]byte(jsQuery))
			wd, _ := os.Getwd()
			cmd.Dir = filepath.Join(filepath.Dir(wd), "ts")
			stdout := new(bytes.Buffer)
			stderr := new(bytes.Buffer)
			cmd.Stdout = stdout
			cmd.Stderr = stderr
			err = cmd.Run()

			outputB64, _ := io.ReadAll(stdout)
			if err != nil {
				allErr, _ := io.ReadAll(stderr)
				fmt.Println(string(allErr))
				t.Fatalf("node exec failed: %v", err)
			}
			fmt.Println(err)
			fmt.Println("out:", string(outputB64))
			fmt.Println("in:", len([]byte(inputB64)))

			outBinary, _ := base64.StdEncoding.DecodeString(string(outputB64))
			umt = tc.unmarshalRecord()
			err = umt.DecodeBebop(bytes.NewBuffer(outBinary))
			if err != nil {
				t.Fatalf("js unmarshal failed: %v", err)
			}
			buf = &bytes.Buffer{}
			err = umt.EncodeBebop(buf)
			marshalData6 := buf.Bytes()
			if err != nil {
				t.Fatalf("js marshal failed: %v", err)
			}
			// casting to string for easy equality
			if string(marshalData5) != string(marshalData6) {
				fmt.Println(marshalData5)
				fmt.Println(marshalData6)
				t.Fatal("js marshal did not have same bytes as first")
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
