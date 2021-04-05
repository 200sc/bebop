package bebop

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadFile(t *testing.T) {
	type testCase struct {
		file     string
		expected File
	}
	tcs := []testCase{
		{
			file: "array_of_strings",
			expected: File{
				Structs: []Struct{
					{
						Name: "ArrayOfStrings",
						Fields: []Field{
							{
								Name: "strings",
								FieldType: FieldType{
									Array: &FieldType{Simple: "string"},
								},
							},
						},
					},
				},
			},
		},
		{
			file: "basic_arrays",
			expected: File{
				Structs: []Struct{
					{
						Name: "BasicArrays",
						Fields: []Field{
							{
								FieldType: FieldType{Array: &FieldType{Simple: "bool"}},
								Name:      "a_bool",
							}, {
								FieldType: FieldType{Array: &FieldType{Simple: "byte"}},
								Name:      "a_byte",
							}, {
								FieldType: FieldType{Array: &FieldType{Simple: "int16"}},
								Name:      "a_int16",
							}, {
								FieldType: FieldType{Array: &FieldType{Simple: "uint16"}},
								Name:      "a_uint16",
							}, {
								FieldType: FieldType{Array: &FieldType{Simple: "int32"}},
								Name:      "a_int32",
							}, {
								FieldType: FieldType{Array: &FieldType{Simple: "uint32"}},
								Name:      "a_uint32",
							}, {
								FieldType: FieldType{Array: &FieldType{Simple: "int64"}},
								Name:      "a_int64",
							}, {
								FieldType: FieldType{Array: &FieldType{Simple: "uint64"}},
								Name:      "a_uint64",
							}, {
								FieldType: FieldType{Array: &FieldType{Simple: "float32"}},
								Name:      "a_float32",
							}, {
								FieldType: FieldType{Array: &FieldType{Simple: "float64"}},
								Name:      "a_float64",
							}, {
								FieldType: FieldType{Array: &FieldType{Simple: "string"}},
								Name:      "a_string",
							}, {
								FieldType: FieldType{Array: &FieldType{Simple: "guid"}},
								Name:      "a_guid",
							},
						},
					}, {
						Name: "TestInt32Array",
						Fields: []Field{
							{
								FieldType: FieldType{Array: &FieldType{Simple: "int32"}},
								Name:      "a",
							},
						},
					},
				},
			},
		},
		{
			file: "basic_types",
			expected: File{
				Structs: []Struct{
					{
						Name: "BasicTypes",
						Fields: []Field{
							{
								FieldType: FieldType{Simple: "bool"},
								Name:      "a_bool",
							}, {
								FieldType: FieldType{Simple: "byte"},
								Name:      "a_byte",
							}, {
								FieldType: FieldType{Simple: "int16"},
								Name:      "a_int16",
							}, {
								FieldType: FieldType{Simple: "uint16"},
								Name:      "a_uint16",
							}, {
								FieldType: FieldType{Simple: "int32"},
								Name:      "a_int32",
							}, {
								FieldType: FieldType{Simple: "uint32"},
								Name:      "a_uint32",
							}, {
								FieldType: FieldType{Simple: "int64"},
								Name:      "a_int64",
							}, {
								FieldType: FieldType{Simple: "uint64"},
								Name:      "a_uint64",
							}, {
								FieldType: FieldType{Simple: "float32"},
								Name:      "a_float32",
							}, {
								FieldType: FieldType{Simple: "float64"},
								Name:      "a_float64",
							}, {
								FieldType: FieldType{Simple: "string"},
								Name:      "a_string",
							}, {
								FieldType: FieldType{Simple: "guid"},
								Name:      "a_guid",
							}, {
								FieldType: FieldType{Simple: "date"},
								Name:      "a_date",
							},
						},
					},
				},
			},
		},
		{
			file: "block_comments",
			expected: File{
				Enums: []Enum{
					{
						Name:    "BlockComments",
						Comment: " block \n line",
						Options: []EnumOption{
							{
								Value:   1,
								Name:    "Block",
								Comment: " block \n line",
							},
						},
					},
				},
				Structs: []Struct{
					{
						Name: "BlockComments2",
						Fields: []Field{
							{
								FieldType: FieldType{
									Simple: "int16",
								},
								Comment: " block \n line",
								Name:    "f",
							},
						},
					},
				},
				Messages: []Message{
					{
						Name:   "BlockComments3",
						Fields: map[uint8]Field{},
					},
				},
			},
		},
		{
			file: "documentation",
			expected: File{
				Messages: []Message{
					{
						Name: "DepM",
						Fields: map[uint8]Field{
							1: {
								Name:              "x",
								Deprecated:        true,
								DeprecatedMessage: "x in DepM",
								FieldType: FieldType{
									Simple: "int32",
								},
							},
						},
					}, {
						Name:    "DocM",
						Comment: " Documented message ",
						Fields: map[uint8]Field{
							1: {
								Name:    "x",
								Comment: " Documented field ",
								FieldType: FieldType{
									Simple: "int32",
								},
							},
							2: {
								Name:              "y",
								Deprecated:        true,
								DeprecatedMessage: "y in DocM",
								FieldType: FieldType{
									Simple: "int32",
								},
							},
							3: {
								Name:              "z",
								Comment:           " Deprecated, documented field ",
								Deprecated:        true,
								DeprecatedMessage: "z in DocM",
								FieldType: FieldType{
									Simple: "int32",
								},
							},
						},
					},
				},
				Structs: []Struct{
					{
						Name:    "DocS",
						Comment: " Documented struct ",
						Fields: []Field{
							{
								Comment: " Documented field ",
								Name:    "x",
								FieldType: FieldType{
									Simple: "int32",
								},
							},
						},
					},
				},
				Enums: []Enum{
					{
						Name: "DepE",
						Options: []EnumOption{
							{
								Name:              "X",
								Value:             1,
								Deprecated:        true,
								DeprecatedMessage: "X in DepE",
							},
						},
					}, {
						Name:    "DocE",
						Comment: " Documented enum ",
						Options: []EnumOption{
							{
								Name:    "X",
								Value:   1,
								Comment: " Documented constant ",
							}, {
								Name:              "Y",
								Value:             2,
								Deprecated:        true,
								DeprecatedMessage: "Y in DocE",
							}, {
								Name:              "Z",
								Value:             3,
								Comment:           " Deprecated, documented constant ",
								Deprecated:        true,
								DeprecatedMessage: "Z in DocE",
							},
						},
					},
				},
			},
		},
		{
			file: "enums",
			expected: File{
				Enums: []Enum{
					{
						Name: "Test",
						Options: []EnumOption{
							{
								Name:  "Start",
								Value: 1,
							}, {
								Name:  "End",
								Value: 2,
							}, {
								Name:  "Middle",
								Value: 3,
							}, {
								Name:              "Beginning",
								Value:             4,
								DeprecatedMessage: "who knows",
								Deprecated:        true,
							},
						},
					},
				},
			},
		},
		{
			file: "enums_doc",
			expected: File{
				Enums: []Enum{
					{
						Name:    "Test2",
						Comment: " test 2 has a line comment",
						Options: []EnumOption{
							{
								Name:  "Start",
								Value: 1,
							}, {
								Name:    "End",
								Comment: " end has a line comment too",
								Value:   2,
							}, {
								Name:  "Middle",
								Value: 3,
							}, {
								Name:              "Beginning",
								Value:             4,
								DeprecatedMessage: "who knows",
								Deprecated:        true,
							},
						},
					},
				},
			},
		},
		{
			file: "foo",
			expected: File{
				Structs: []Struct{
					{
						Name: "Foo",
						Fields: []Field{
							{
								Name: "bar",
								FieldType: FieldType{
									Simple: "Bar",
								},
							},
						},
					},
				},
				Messages: []Message{
					{
						Name: "Bar",
						Fields: map[uint8]Field{
							1: {
								Name: "x",
								FieldType: FieldType{
									Simple: "float64",
								},
							},
							2: {
								Name: "y",
								FieldType: FieldType{
									Simple: "float64",
								},
							},
							3: {
								Name: "z",
								FieldType: FieldType{
									Simple: "float64",
								},
							},
						},
					},
				},
			},
		},
		{
			file: "jazz",
			expected: File{
				Enums: []Enum{
					{
						Name: "Instrument",
						Options: []EnumOption{
							{
								Name:  "Sax",
								Value: 0,
							},
							{
								Name:  "Trumpet",
								Value: 1,
							},
							{
								Name:  "Clarinet",
								Value: 2,
							},
						},
					},
				},
				Structs: []Struct{
					{
						Name:     "Musician",
						ReadOnly: true,
						Fields: []Field{
							{
								Name: "name",
								FieldType: FieldType{
									Simple: "string",
								},
							},
							{
								Name: "plays",
								FieldType: FieldType{
									Simple: "Instrument",
								},
							},
						},
					},
					{
						Name: "Library",
						Fields: []Field{
							{
								Name: "songs",
								FieldType: FieldType{
									Map: &MapType{
										Key: "guid",
										Value: FieldType{
											Simple: "Song",
										},
									},
								},
							},
						},
					},
				},
				Messages: []Message{
					{
						Name: "Song",
						Fields: map[uint8]Field{
							1: {
								Name: "title",
								FieldType: FieldType{
									Simple: "string",
								},
							},
							2: {
								Name: "year",
								FieldType: FieldType{
									Simple: "uint16",
								},
							},
							3: {
								Name: "performers",
								FieldType: FieldType{
									Array: &FieldType{
										Simple: "Musician",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			file: "lab",
			expected: File{
				Enums: []Enum{
					{
						Name: "VideoCodec",
						Options: []EnumOption{
							{
								Name:  "H264",
								Value: 0,
							},
							{
								Name:  "H265",
								Value: 1,
							},
						},
					},
				},
				Structs: []Struct{
					{
						Name: "Int32s",
						Fields: []Field{
							{
								Name: "a",
								FieldType: FieldType{
									Array: &FieldType{
										Simple: "int32",
									},
								},
							},
						},
					},
					{
						Name: "Uint32s",
						Fields: []Field{
							{
								Name: "a",
								FieldType: FieldType{
									Array: &FieldType{
										Simple: "uint32",
									},
								},
							},
						},
					},
					{
						Name: "Float32s",
						Fields: []Field{
							{
								Name: "a",
								FieldType: FieldType{
									Array: &FieldType{
										Simple: "float32",
									},
								},
							},
						},
					},
					{
						Name: "Int64s",
						Fields: []Field{
							{
								Name: "a",
								FieldType: FieldType{
									Array: &FieldType{
										Simple: "int64",
									},
								},
							},
						},
					},
					{
						Name: "Uint64s",
						Fields: []Field{
							{
								Name: "a",
								FieldType: FieldType{
									Array: &FieldType{
										Simple: "uint64",
									},
								},
							},
						},
					},
					{
						Name: "Float64s",
						Fields: []Field{
							{
								Name: "a",
								FieldType: FieldType{
									Array: &FieldType{
										Simple: "float64",
									},
								},
							},
						},
					},
					{
						Name: "VideoData",
						Fields: []Field{
							{
								Name: "time",
								FieldType: FieldType{
									Simple: "float64",
								},
							},
							{
								Name: "width",
								FieldType: FieldType{
									Simple: "uint32",
								},
							},
							{
								Name: "height",
								FieldType: FieldType{
									Simple: "uint32",
								},
							},
							{
								Name: "fragment",
								FieldType: FieldType{
									Array: &FieldType{
										Simple: "byte",
									},
								},
							},
						},
					},
				},
				Messages: []Message{
					{
						Name: "MediaMessage",
						Fields: map[uint8]Field{
							1: {
								Name: "codec",
								FieldType: FieldType{
									Simple: "VideoCodec",
								},
							},
							2: {
								Name: "data",
								FieldType: FieldType{
									Simple: "VideoData",
								},
							},
						},
					},
					{
						Name:    "SkipTestOld",
						Comment: " Should be able to decode a \"SkipTestNewContainer\" as a \"SkipTestOldContainer\".",
						Fields: map[uint8]Field{
							1: {
								Name: "x",
								FieldType: FieldType{
									Simple: "int32",
								},
							},
							2: {
								Name: "y",
								FieldType: FieldType{
									Simple: "int32",
								},
							},
						},
					},
					{
						Name: "SkipTestNew",
						Fields: map[uint8]Field{
							1: {
								Name: "x",
								FieldType: FieldType{
									Simple: "int32",
								},
							},
							2: {
								Name: "y",
								FieldType: FieldType{
									Simple: "int32",
								},
							},
							3: {
								Name: "z",
								FieldType: FieldType{
									Simple: "int32",
								},
							},
						},
					},
					{
						Name: "SkipTestOldContainer",
						Fields: map[uint8]Field{
							1: {
								Name: "s",
								FieldType: FieldType{
									Simple: "SkipTestOld",
								},
							},
							2: {
								Name: "after",
								FieldType: FieldType{
									Simple: "int32",
								},
							},
						},
					},
					{
						Name: "SkipTestNewContainer",
						Fields: map[uint8]Field{
							1: {
								Name: "s",
								FieldType: FieldType{
									Simple: "SkipTestNew",
								},
							},
							2: {
								Name: "after",
								FieldType: FieldType{
									Simple: "int32",
								},
							},
						},
					},
				},
			},
		},
		{
			file: "map_types",
			expected: File{
				Structs: []Struct{
					{
						Name:     "S",
						ReadOnly: true,
						Fields: []Field{
							{
								Name: "x",
								FieldType: FieldType{
									Simple: "int32",
								},
							},
							{
								Name: "y",
								FieldType: FieldType{
									Simple: "int32",
								},
							},
						},
					},
					{
						Name: "SomeMaps",
						Fields: []Field{
							{
								Name: "m1",
								FieldType: FieldType{
									Map: &MapType{
										Key: "bool",
										Value: FieldType{
											Simple: "bool",
										},
									},
								},
							},
							{
								Name: "m2",
								FieldType: FieldType{
									Map: &MapType{
										Key: "string",
										Value: FieldType{
											Map: &MapType{
												Key: "string",
												Value: FieldType{
													Simple: "string",
												},
											},
										},
									},
								},
							},
							{
								Name: "m3",
								FieldType: FieldType{
									Array: &FieldType{
										Map: &MapType{
											Key: "int32",
											Value: FieldType{
												Array: &FieldType{
													Map: &MapType{
														Key: "bool",
														Value: FieldType{
															Simple: "S",
														},
													},
												},
											},
										},
									},
								},
							},
							{
								Name: "m4",
								FieldType: FieldType{
									Array: &FieldType{
										Map: &MapType{
											Key: "string",
											Value: FieldType{
												Array: &FieldType{
													Simple: "float32",
												},
											},
										},
									},
								},
							},
							{
								Name: "m5",
								FieldType: FieldType{
									Map: &MapType{
										Key: "guid",
										Value: FieldType{
											Simple: "M",
										},
									},
								},
							},
						},
					},
				},
				Messages: []Message{
					{
						Name: "M",
						Fields: map[uint8]Field{
							1: {
								Name: "a",
								FieldType: FieldType{
									Simple: "float32",
								},
							},
							2: {
								Name: "b",
								FieldType: FieldType{
									Simple: "float64",
								},
							},
						},
					},
				},
			},
		},
		{
			file: "msgpack_comparison",
			expected: File{
				Structs: []Struct{
					{
						Name: "MsgpackComparison",
						Comment: " These field names are extremely weirdly capitalized, because I wanted the\n" +
							" key names in JSON to be the same length while not coinciding with Bebop keywords.",
						Fields: []Field{
							{
								Name: "iNT0",
								FieldType: FieldType{
									Simple: "uint8",
								},
							},
							{
								Name: "iNT1",
								FieldType: FieldType{
									Simple: "uint8",
								},
							},
							{
								Name: "iNT1_",
								FieldType: FieldType{
									Simple: "int16",
								},
							},
							{
								Name: "iNT8",
								FieldType: FieldType{
									Simple: "uint8",
								},
							},
							{
								Name: "iNT8_",
								FieldType: FieldType{
									Simple: "int16",
								},
							},
							{
								Name: "iNT16",
								FieldType: FieldType{
									Simple: "int16",
								},
							},
							{
								Name: "iNT16_",
								FieldType: FieldType{
									Simple: "int16",
								},
							},
							{
								Name: "iNT32",
								FieldType: FieldType{
									Simple: "int32",
								},
							},
							{
								Name: "iNT32_",
								FieldType: FieldType{
									Simple: "int32",
								},
							},
							{
								Name:    "tRUE",
								Comment: " int8 nIL; // \"nil\": null,",
								FieldType: FieldType{
									Simple: "bool",
								},
							},
							{
								Name: "fALSE",
								FieldType: FieldType{
									Simple: "bool",
								},
							},
							{
								Name: "fLOAT",
								FieldType: FieldType{
									Simple: "float64",
								},
							},
							{
								Name: "fLOAT_",
								FieldType: FieldType{
									Simple: "float64",
								},
							},
							{
								Name: "sTRING0",
								FieldType: FieldType{
									Simple: "string",
								},
							},
							{
								Name: "sTRING1",
								FieldType: FieldType{
									Simple: "string",
								},
							},
							{
								Name: "sTRING4",
								FieldType: FieldType{
									Simple: "string",
								},
							},
							{
								Name: "sTRING8",
								FieldType: FieldType{
									Simple: "string",
								},
							},
							{
								Name: "sTRING16",
								FieldType: FieldType{
									Simple: "string",
								},
							},
							{
								Name: "aRRAY0",
								FieldType: FieldType{
									Array: &FieldType{
										Simple: "int32",
									},
								},
							},
							{
								Name: "aRRAY1",
								FieldType: FieldType{
									Array: &FieldType{
										Simple: "string",
									},
								},
							},
							{
								Name: "aRRAY8",
								FieldType: FieldType{
									Array: &FieldType{
										Simple: "int32",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			file: "request",
			expected: File{
				Enums: []Enum{
					{
						Name: "FurnitureFamily",
						Options: []EnumOption{
							{
								Name:  "Bed",
								Value: 0,
							},
							{
								Name:  "Table",
								Value: 1,
							},
							{
								Name:  "Shoe",
								Value: 2,
							},
						},
					},
				},
				Structs: []Struct{
					{
						Name:     "Furniture",
						ReadOnly: true,
						Fields: []Field{
							{
								Name: "name",
								FieldType: FieldType{
									Simple: "string",
								},
							},
							{
								Name: "price",
								FieldType: FieldType{
									Simple: "uint32",
								},
							},
							{
								Name: "family",
								FieldType: FieldType{
									Simple: "FurnitureFamily",
								},
							},
						},
					},
					{
						Name:     "RequestResponse",
						ReadOnly: true,
						OpCode:   0x31323334,
						Fields: []Field{
							{
								Name: "availableFurniture",
								FieldType: FieldType{
									Array: &FieldType{
										Simple: "Furniture",
									},
								},
							},
						},
					},
				},
				Messages: []Message{
					{
						Name:   "RequestCatalog",
						OpCode: bytesToOpCode([]byte("IKEA")),
						Fields: map[uint8]Field{
							1: {
								Name: "family",
								FieldType: FieldType{
									Simple: "FurnitureFamily",
								},
							},
							2: {
								Name:              "secretTunnel",
								Deprecated:        true,
								DeprecatedMessage: "Nobody react to what I'm about to say...",
								FieldType: FieldType{
									Simple: "string",
								},
							},
						},
					},
				},
			},
		}, {
			file: "union",
			expected: File{
				Unions: []Union{
					{
						Name:   "U",
						OpCode: bytesToOpCode([]byte("yeah")),
						Fields: map[uint8]UnionField{
							1: {
								Message: &Message{
									Name: "A",
									Fields: map[uint8]Field{
										1: {
											Name: "a",
											FieldType: FieldType{
												Simple: "uint32",
											},
										},
									},
								},
							},
							2: {
								Struct: &Struct{
									Comment: "*\n     * This branch is, too!\n     ",
									Name:    "B",
									Fields: []Field{
										{
											Name: "b",
											FieldType: FieldType{
												Simple: "bool",
											},
										},
									},
								},
							},
							3: {
								Struct: &Struct{
									Name: "C",
								},
							},
							4: {
								Union: &Union{
									Name: "W",
									Fields: map[uint8]UnionField{
										1: {
											Struct: &Struct{
												Name: "D",
												Fields: []Field{{
													Name: "s",
													FieldType: FieldType{
														Simple: "string",
													},
												}},
											},
										},
										2: {
											Struct: &Struct{
												Name: "X",
												Fields: []Field{{
													Name: "x",
													FieldType: FieldType{
														Simple: "bool",
													},
												}},
											},
										},
									},
								},
							},
						},
					},
					{
						Name: "List",
						Fields: map[uint8]UnionField{
							1: {
								Struct: &Struct{
									Name: "Cons",
									Fields: []Field{
										{
											Name: "head",
											FieldType: FieldType{
												Simple: "uint32",
											},
										}, {
											Name: "tail",
											FieldType: FieldType{
												Simple: "List",
											},
										},
									},
								},
							},
							2: {
								Struct: &Struct{
									Name: "Nil",
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.file, func(t *testing.T) {
			f, err := os.Open(filepath.Join("testdata", "base", tc.file+".bop"))
			if err != nil {
				t.Fatalf("failed to open test file %s: %v", tc.file+".bop", err)
			}
			defer f.Close()
			bf, err := ReadFile(f)
			if err != nil {
				t.Fatalf("read file errored: %v", err)
			}
			if err := bf.equals(tc.expected); err != nil {
				t.Fatal("parsed file did not match expected:", err)
			}
		})
	}
}

func TestReadFileError(t *testing.T) {
	type testCase struct {
		file       string
		errMessage string
	}
	tcs := []testCase{
		{file: "invalid_enum_with_op_code", errMessage: "[1:5] enums may not have attached op codes"},
		{file: "invalid_op_code_1", errMessage: "[0:2] expected (OpCode) got Close Square"},
		{file: "invalid_op_code_2", errMessage: "[0:7] expected (OpCode) got Ident"},
		{file: "invalid_op_code_3", errMessage: "[0:16] opcode string 12345 exceeds 4 ascii characters"},
		{file: "invalid_op_code_4", errMessage: "[0:9] expected (Open Paren) got Open Square"},
		{file: "invalid_op_code_5", errMessage: "[0:82] strconv.ParseInt: parsing \"1111111111111111111111111111111111111111111111111111111111111111111111111\": value out of range"},
		{file: "invalid_op_code_6", errMessage: "[0:16] expected (Close Paren) got Close Square"},
		{file: "invalid_op_code_7", errMessage: "[0:17] expected (Close Square) got Equals"},
		{file: "invalid_op_code_8", errMessage: "[0:15] expected ([Integer String Literal]) got Ident"},
		{file: "invalid_enum_bad_deprecated", errMessage: "[1:18] expected (String Literal) got Equals"},
		{file: "invalid_enum_double_deprecated", errMessage: "[2:5] expected enum option following deprecated annotation"},
		{file: "invalid_enum_hex_int", errMessage: "[1:14] strconv.ParseInt: parsing \"0x1\": invalid syntax"},
		{file: "invalid_enum_no_close", errMessage: "[2:0] enum definition ended early"},
		{file: "invalid_enum_no_curly", errMessage: "[1:0] expected (Open Curly) got Newline"},
		{file: "invalid_enum_no_eq", errMessage: "[1:10] expected (Equals) got Integer"},
		{file: "invalid_enum_no_int", errMessage: "[1:11] expected (Integer) got Semicolon"},
		{file: "invalid_enum_no_name", errMessage: "[0:7] expected (Ident) got Open Curly"},
		{file: "invalid_enum_no_semi", errMessage: "[2:0] expected (Semicolon) got Newline"},
		{file: "invalid_struct_bad_deprecated", errMessage: "[1:22] expected (String Literal) got Ident"},
		{file: "invalid_struct_bad_type", errMessage: "[1:10] expected ([Ident Array Map]) got Open Square"},
		{file: "invalid_struct_double_deprecated", errMessage: "[2:5] expected field following deprecated annotation"},
		{file: "invalid_struct_no_close", errMessage: "[1:16] struct definition ended early"},
		{file: "invalid_struct_no_curly", errMessage: "[1:0] expected (Open Curly) got Newline"},
		{file: "invalid_struct_no_field_name", errMessage: "[1:11] expected (Ident) got Semicolon"},
		{file: "invalid_struct_no_name", errMessage: "[0:9] expected (Ident) got Open Curly"},
		{file: "invalid_struct_no_semi", errMessage: "[2:0] expected (Semicolon) got Newline"},
		{file: "invalid_message_bad_deprecated", errMessage: "[1:19] expected (String Literal) got Arrow"},
		{file: "invalid_message_bad_type", errMessage: "[1:15] expected ([Ident Array Map]) got Open Square"},
		{file: "invalid_message_double_deprecated", errMessage: "[2:5] expected field following deprecated annotation"},
		{file: "invalid_message_hex_int", errMessage: "[1:7] strconv.ParseInt: parsing \"0x1\": invalid syntax"},
		{file: "invalid_message_no_arrow", errMessage: "[1:12] expected (Arrow) got Ident"},
		{file: "invalid_message_no_close", errMessage: "[1:21] message definition ended early"},
		{file: "invalid_message_no_curly", errMessage: "[1:0] expected (Open Curly) got Newline"},
		{file: "invalid_message_no_field_name", errMessage: "[1:16] expected (Ident) got Semicolon"},
		{file: "invalid_message_no_name", errMessage: "[0:10] expected (Ident) got Open Curly"},
		{file: "invalid_message_no_semi", errMessage: "[2:0] expected (Semicolon) got Newline"},
		{file: "invalid_enum_reserved", errMessage: "[0:13] expected (Ident) got Struct"},
		{file: "invalid_struct_reserved", errMessage: "[0:14] expected (Ident) got Array"},
		{file: "invalid_message_reserved", errMessage: "[0:13] expected (Ident) got Map"},
		{file: "invalid_message_duplicate_index", errMessage: "[2:2] message has duplicate field index 1"},
		{file: "invalid_readonly_enum", errMessage: "[0:15] expected (Struct) got (Enum)"},
		{file: "invalid_readonly_message", errMessage: "[0:18] expected (Struct) got (Message)"},
		{file: "invalid_readonly_comment", errMessage: "[0:20] expected (Struct) got (Block Comment)"},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.file, func(t *testing.T) {
			f, err := os.Open(filepath.Join("testdata", "base", tc.file+".bop"))
			if err != nil {
				t.Fatalf("failed to open test file %s: %v", tc.file+".bop", err)
			}
			defer f.Close()
			_, err = ReadFile(f)
			if err == nil {
				t.Fatalf("read file should have errored")
			}
			if err.Error() != tc.errMessage {
				t.Fatalf("read file had wrong error: got %q, expected %q", err.Error(), tc.errMessage)
			}
		})
	}
}
