package bebop

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadFile(t *testing.T) {
	t.Parallel()
	type testCase struct {
		file     string
		expected File
	}
	tcs := []testCase{
		{
			file: "typed_enums",
			expected: File{
				Structs: []Struct{
					{
						Name: "UsesAllEnums",
						Fields: []Field{
							{FieldType: FieldType{Simple: "EnumU8"}, Name: "one"},
							{FieldType: FieldType{Simple: "EnumU16"}, Name: "two"},
							{FieldType: FieldType{Simple: "EnumU32"}, Name: "three"},
							{FieldType: FieldType{Simple: "EnumU64"}, Name: "four"},
							{FieldType: FieldType{Simple: "Enum16"}, Name: "five"},
							{FieldType: FieldType{Simple: "Enum32"}, Name: "six"},
							{FieldType: FieldType{Simple: "Enum64"}, Name: "seven"},
						},
					},
				},
				Enums: []Enum{
					{
						Name:       "EnumU8",
						SimpleType: "uint8",
						Unsigned:   true,
						Options: []EnumOption{
							{
								Name:      "EnumU81",
								UintValue: 1,
							},
							{
								Name:      "EnumU82",
								UintValue: 2,
							},
						},
					},
					{
						Name:       "EnumByte",
						SimpleType: "byte",
						Unsigned:   true,
						Options: []EnumOption{
							{
								Name:      "EnumByte1",
								UintValue: 1,
							},
							{
								Name:      "EnumByte2",
								UintValue: 2,
							},
						},
					},
					{
						Name:       "EnumU16",
						SimpleType: "uint16",
						Unsigned:   true,
						Options: []EnumOption{
							{
								Name:      "EnumU161",
								UintValue: 1,
							},
							{
								Name:      "EnumU162",
								UintValue: 2,
							},
						},
					},
					{
						Name:       "EnumU32",
						SimpleType: "uint32",
						Unsigned:   true,
						Options: []EnumOption{
							{
								Name:      "EnumU321",
								UintValue: 1,
							},
							{
								Name:      "EnumU322",
								UintValue: 2,
							},
						},
					},
					{
						Name:       "EnumU64",
						SimpleType: "uint64",
						Unsigned:   true,
						Options: []EnumOption{
							{
								Name:      "EnumU641",
								UintValue: 1,
							},
							{
								Name:      "EnumU642",
								UintValue: 2,
							},
						},
					},
					{
						Name:       "Enum16",
						SimpleType: "int16",
						Options: []EnumOption{
							{
								Name:  "Enum161",
								Value: 1,
							},
							{
								Name:  "Enum162",
								Value: 2,
							},
						},
					},
					{
						Name:       "Enum32",
						SimpleType: "int32",
						Options: []EnumOption{
							{
								Name:  "Enum321",
								Value: 1,
							},
							{
								Name:  "Enum322",
								Value: 2,
							},
						},
					},
					{
						Name:       "Enum64",
						SimpleType: "int64",
						Options: []EnumOption{
							{
								Name:  "Enum641",
								Value: 1,
							},
							{
								Name:  "Enum642",
								Value: 2,
							},
						},
					},
				},
			},
		},
		{
			file: "arrays",
			expected: File{
				Structs: []Struct{
					{
						Name: "ArraySamples",
						Fields: []Field{
							{
								Name: "bytes",
								FieldType: FieldType{
									Array: &FieldType{
										Array: &FieldType{
											Array: &FieldType{
												Simple: "byte",
											},
										},
									},
								},
							},
							{
								Name: "bytes2",
								FieldType: FieldType{
									Array: &FieldType{
										Array: &FieldType{
											Array: &FieldType{
												Simple: "byte",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			file: "opcodes",
			expected: File{
				Structs: []Struct{
					{
						Name:   "NumericalASCIIOpCode",
						OpCode: 0x34333231,
					},
					{
						Name:     "NumericalASCIIOpCode2",
						ReadOnly: true,
						OpCode:   825373493,
					},
					{
						Name:   "NumericalASCIIOpCode3",
						OpCode: 842150708,
					},
					{
						Name:     "NumericalASCIIOpCode4",
						ReadOnly: true,
						OpCode:   875770418,
					},
				},
			},
		},
		{
			file: "bitflags",
			expected: File{
				Enums: []Enum{
					{
						Name:       "TestFlags",
						SimpleType: "uint32",
						Unsigned:   true,
						Options: []EnumOption{
							{
								Name:      "None",
								UintValue: 0,
							}, {
								Name:      "Read",
								UintValue: 1,
							}, {
								Name:      "Write",
								UintValue: 2,
							}, {
								Name:      "ReadWrite",
								UintValue: 1 | 2,
							}, {
								Name:      "Complex",
								UintValue: (1 | 2) | 0xF0&0x1F,
							},
						},
					},
					{
						Name:       "TestFlags2",
						SimpleType: "int64",
						Options: []EnumOption{
							{
								Name:  "None",
								Value: 0,
							}, {
								Name:  "Read",
								Value: 1,
							}, {
								Name:  "Write",
								Value: 2,
							}, {
								Name:  "ReadWrite",
								Value: 1 | 2,
							}, {
								Name:  "Complex",
								Value: (1 | 2) | 0xF0&0x1F,
							},
						},
					},
				},
			},
		},
		{
			file: "union_2",
			expected: File{
				Unions: []Union{
					{
						Name: "U2",
						Fields: map[uint8]UnionField{
							1: {
								Struct: &Struct{
									Name: "U3",
									Fields: []Field{
										{
											Name: "hello",
											FieldType: FieldType{
												Simple: "uint32",
											},
										},
									},
								},
							},
							2: {
								Message: &Message{
									Name: "U4",
									Fields: map[uint8]Field{
										1: {
											Name: "goodbye",
											FieldType: FieldType{
												Simple: "uint32",
											},
										},
									},
								},
							},
							3: {
								Message: &Message{
									Name: "U5",
									Fields: map[uint8]Field{
										1: {
											Name: "goodbye",
											FieldType: FieldType{
												Simple: "uint32",
											},
										},
									},
								},
							},
							4: {
								Struct: &Struct{
									Name: "U6",
									Fields: []Field{
										{
											Name: "hello",
											FieldType: FieldType{
												Simple: "uint32",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			file: "tags",
			expected: File{
				Structs: []Struct{
					{
						Name: "TaggedStruct",
						Fields: []Field{
							{
								Name: "foo",
								FieldType: FieldType{
									Simple: "string",
								},
								Comment: "[tag(json:\"foo,omitempty\")]",
								Tags: []Tag{
									{
										Key:   "json",
										Value: "foo,omitempty",
									},
								},
							},
						},
					},
				},
				Messages: []Message{
					{
						Name: "TaggedMessage",
						Fields: map[uint8]Field{
							1: {
								Name:    "bar",
								Comment: "[tag(db:\"bar\")]",
								FieldType: FieldType{
									Simple: "uint8",
								},
								Tags: []Tag{
									{
										Key:   "db",
										Value: "bar",
									},
								},
							},
						},
					},
				},
				Unions: []Union{
					{
						Name: "TaggedUnion",
						Fields: map[uint8]UnionField{
							1: {
								Tags: []Tag{
									{
										Key:   "one",
										Value: "one",
									},
									{
										Key:   "two",
										Value: "two",
									},
									{
										Key:     "boolean",
										Boolean: true,
									},
								},
								Struct: &Struct{
									Name:    "TaggedSubStruct",
									Comment: "[tag(one:\"one\")]\n[tag(two:\"two\")]\n[tag(boolean)]",
									Fields: []Field{
										{
											Name: "biz",
											FieldType: FieldType{
												Simple: "guid",
											},
											Comment: "[tag(four:\"four\")]",
											Tags: []Tag{
												{
													Key:   "four",
													Value: "four",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			file: "import",
			expected: File{
				Structs: []Struct{
					{
						Name: "Hello",
						Fields: []Field{
							{
								Name: "Yes",
								FieldType: FieldType{
									Simple: typeInt32,
								},
							},
							{
								Name: "No",
								FieldType: FieldType{
									Simple: typeString,
								},
							},
						},
					},
				},
			},
		},
		{
			file: "import_b",
			expected: File{
				Structs: []Struct{
					{
						Name: "Test22",
						Fields: []Field{
							{
								Name: "noisemaker",
								FieldType: FieldType{
									Simple: "Instrument",
								},
							},
						},
					},
				},
			},
		},
		{
			file: "enum_hex_int",
			expected: File{
				Enums: []Enum{
					{
						Name:       "MyEnum",
						SimpleType: "uint32",
						Unsigned:   true,
						Options: []EnumOption{
							{
								Name:      "One",
								UintValue: 1,
							},
						},
					},
				},
			},
		},
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
									Array: &FieldType{Simple: typeString},
								},
							},
						},
					},
				},
			},
		},
		{
			file: "all_consts",
			expected: File{
				Consts: []Const{
					{
						SimpleType: typeUint8,
						Name:       "uint8const",
						Value:      "1",
					},
					{
						SimpleType: typeUint16,
						Name:       "uint16const",
						Value:      "1",
					},
					{
						SimpleType: typeUint32,
						Name:       "uint32const",
						Value:      "1",
					},
					{
						SimpleType: typeUint64,
						Name:       "uint64const",
						Value:      "1",
					},
					{
						SimpleType: typeByte,
						Name:       "int8const",
						Value:      "1",
					},
					{
						SimpleType: typeInt16,
						Name:       "int16const",
						Value:      "1",
					},
					{
						SimpleType: typeInt32,
						Name:       "int32const",
						Value:      "1",
					},
					{
						SimpleType: typeInt64,
						Name:       "int64const",
						Value:      "1",
					},
					{
						SimpleType: typeFloat32,
						Name:       "float32const",
						Value:      "1",
					},
					{
						SimpleType: typeFloat64,
						Name:       "float64const",
						Value:      "1.5",
					},
					{
						SimpleType: typeFloat64,
						Name:       "float64infconst",
						Value:      "math.Inf(1)",
					},
					{
						SimpleType: typeFloat64,
						Name:       "float64ninfconst",
						Value:      "math.Inf(-1)",
					},
					{
						SimpleType: typeFloat64,
						Name:       "float64nanconst",
						Value:      "math.NaN()",
					},
					{
						SimpleType: typeBool,
						Name:       "boolconst",
						Value:      "true",
					},
					{
						SimpleType: typeString,
						Name:       "stringconst",
						Value:      `"1"`,
					},
					{
						SimpleType: typeGUID,
						Name:       "guidconst",
						Value:      `"e2722bf7-022a-496a-9f01-7029d7d5563d"`,
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
								FieldType: FieldType{Array: &FieldType{Simple: typeBool}},
								Name:      "a_bool",
							}, {
								FieldType: FieldType{Array: &FieldType{Simple: typeByte}},
								Name:      "a_byte",
							}, {
								FieldType: FieldType{Array: &FieldType{Simple: typeInt16}},
								Name:      "a_int16",
							}, {
								FieldType: FieldType{Array: &FieldType{Simple: typeUint16}},
								Name:      "a_uint16",
							}, {
								FieldType: FieldType{Array: &FieldType{Simple: typeInt32}},
								Name:      "a_int32",
							}, {
								FieldType: FieldType{Array: &FieldType{Simple: typeUint32}},
								Name:      "a_uint32",
							}, {
								FieldType: FieldType{Array: &FieldType{Simple: typeInt64}},
								Name:      "a_int64",
							}, {
								FieldType: FieldType{Array: &FieldType{Simple: typeUint64}},
								Name:      "a_uint64",
							}, {
								FieldType: FieldType{Array: &FieldType{Simple: typeFloat32}},
								Name:      "a_float32",
							}, {
								FieldType: FieldType{Array: &FieldType{Simple: typeFloat64}},
								Name:      "a_float64",
							}, {
								FieldType: FieldType{Array: &FieldType{Simple: typeString}},
								Name:      "a_string",
							}, {
								FieldType: FieldType{Array: &FieldType{Simple: typeGUID}},
								Name:      "a_guid",
							},
						},
					}, {
						Name: "TestInt32Array",
						Fields: []Field{
							{
								FieldType: FieldType{Array: &FieldType{Simple: typeInt32}},
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
								FieldType: FieldType{Simple: typeBool},
								Name:      "a_bool",
							}, {
								FieldType: FieldType{Simple: typeByte},
								Name:      "a_byte",
							}, {
								FieldType: FieldType{Simple: typeInt16},
								Name:      "a_int16",
							}, {
								FieldType: FieldType{Simple: typeUint16},
								Name:      "a_uint16",
							}, {
								FieldType: FieldType{Simple: typeInt32},
								Name:      "a_int32",
							}, {
								FieldType: FieldType{Simple: typeUint32},
								Name:      "a_uint32",
							}, {
								FieldType: FieldType{Simple: typeInt64},
								Name:      "a_int64",
							}, {
								FieldType: FieldType{Simple: typeUint64},
								Name:      "a_uint64",
							}, {
								FieldType: FieldType{Simple: typeFloat32},
								Name:      "a_float32",
							}, {
								FieldType: FieldType{Simple: typeFloat64},
								Name:      "a_float64",
							}, {
								FieldType: FieldType{Simple: typeString},
								Name:      "a_string",
							}, {
								FieldType: FieldType{Simple: typeGUID},
								Name:      "a_guid",
							}, {
								FieldType: FieldType{Simple: typeDate},
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
						Name:       "BlockComments",
						Comment:    " block \n line",
						SimpleType: "uint32",
						Unsigned:   true,
						Options: []EnumOption{
							{
								UintValue: 1,
								Name:      "Block",
								Comment:   " block \n line",
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
									Simple: typeInt16,
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
									Simple: typeInt32,
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
									Simple: typeInt32,
								},
							},
							2: {
								Name:              "y",
								Deprecated:        true,
								DeprecatedMessage: "y in DocM",
								FieldType: FieldType{
									Simple: typeInt32,
								},
							},
							3: {
								Name:              "z",
								Comment:           " Deprecated, documented field ",
								Deprecated:        true,
								DeprecatedMessage: "z in DocM",
								FieldType: FieldType{
									Simple: typeInt32,
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
									Simple: typeInt32,
								},
							},
						},
					},
				},
				Enums: []Enum{
					{
						Name:       "DepE",
						SimpleType: "uint32",
						Unsigned:   true,
						Options: []EnumOption{
							{
								Name:              "X",
								UintValue:         1,
								Deprecated:        true,
								DeprecatedMessage: "X in DepE",
							},
						},
					}, {
						Name:       "DocE",
						SimpleType: "uint32",
						Unsigned:   true,
						Comment:    " Documented enum ",
						Options: []EnumOption{
							{
								Name:      "X",
								UintValue: 1,
								Comment:   " Documented constant ",
							}, {
								Name:              "Y",
								UintValue:         2,
								Deprecated:        true,
								DeprecatedMessage: "Y in DocE",
							}, {
								Name:              "Z",
								UintValue:         3,
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
						Name:       "Test",
						SimpleType: "uint32",
						Unsigned:   true,
						Options: []EnumOption{
							{
								Name:      "Start",
								UintValue: 1,
							}, {
								Name:      "End",
								UintValue: 2,
							}, {
								Name:      "Middle",
								UintValue: 3,
							}, {
								Name:              "Beginning",
								UintValue:         4,
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
						Name:       "Test2",
						SimpleType: "uint32",
						Unsigned:   true,
						Comment:    " test 2 has a line comment",
						Options: []EnumOption{
							{
								Name:      "Start",
								UintValue: 1,
							}, {
								Name:      "End",
								Comment:   " end has a line comment too",
								UintValue: 2,
							}, {
								Name:      "Middle",
								UintValue: 3,
							}, {
								Name:              "Beginning",
								UintValue:         4,
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
									Simple: typeFloat64,
								},
							},
							2: {
								Name: "y",
								FieldType: FieldType{
									Simple: typeFloat64,
								},
							},
							3: {
								Name: "z",
								FieldType: FieldType{
									Simple: typeFloat64,
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
						Name:       "Instrument",
						SimpleType: "uint32",
						Unsigned:   true,
						Options: []EnumOption{
							{
								Name:      "Sax",
								UintValue: 0,
							},
							{
								Name:      "Trumpet",
								UintValue: 1,
							},
							{
								Name:      "Clarinet",
								UintValue: 2,
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
									Simple: typeString,
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
										Key: typeGUID,
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
									Simple: typeString,
								},
							},
							2: {
								Name: "year",
								FieldType: FieldType{
									Simple: typeUint16,
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
						Name:       "VideoCodec",
						SimpleType: "uint32",
						Unsigned:   true,
						Options: []EnumOption{
							{
								Name:      "H264",
								UintValue: 0,
							},
							{
								Name:      "H265",
								UintValue: 1,
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
										Simple: typeInt32,
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
										Simple: typeUint32,
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
										Simple: typeFloat32,
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
										Simple: typeInt64,
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
										Simple: typeUint64,
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
										Simple: typeFloat64,
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
									Simple: typeFloat64,
								},
							},
							{
								Name: "width",
								FieldType: FieldType{
									Simple: typeUint32,
								},
							},
							{
								Name: "height",
								FieldType: FieldType{
									Simple: typeUint32,
								},
							},
							{
								Name: "fragment",
								FieldType: FieldType{
									Array: &FieldType{
										Simple: typeByte,
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
									Simple: typeInt32,
								},
							},
							2: {
								Name: "y",
								FieldType: FieldType{
									Simple: typeInt32,
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
									Simple: typeInt32,
								},
							},
							2: {
								Name: "y",
								FieldType: FieldType{
									Simple: typeInt32,
								},
							},
							3: {
								Name: "z",
								FieldType: FieldType{
									Simple: typeInt32,
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
									Simple: typeInt32,
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
									Simple: typeInt32,
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
									Simple: typeInt32,
								},
							},
							{
								Name: "y",
								FieldType: FieldType{
									Simple: typeInt32,
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
										Key: typeBool,
										Value: FieldType{
											Simple: typeBool,
										},
									},
								},
							},
							{
								Name: "m2",
								FieldType: FieldType{
									Map: &MapType{
										Key: typeString,
										Value: FieldType{
											Map: &MapType{
												Key: typeString,
												Value: FieldType{
													Simple: typeString,
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
											Key: typeInt32,
											Value: FieldType{
												Array: &FieldType{
													Map: &MapType{
														Key: typeBool,
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
											Key: typeString,
											Value: FieldType{
												Array: &FieldType{
													Simple: typeFloat32,
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
										Key: typeGUID,
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
									Simple: typeFloat32,
								},
							},
							2: {
								Name: "b",
								FieldType: FieldType{
									Simple: typeFloat64,
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
									Simple: typeUint8,
								},
							},
							{
								Name: "iNT1",
								FieldType: FieldType{
									Simple: typeUint8,
								},
							},
							{
								Name: "iNT1_",
								FieldType: FieldType{
									Simple: typeInt16,
								},
							},
							{
								Name: "iNT8",
								FieldType: FieldType{
									Simple: typeUint8,
								},
							},
							{
								Name: "iNT8_",
								FieldType: FieldType{
									Simple: typeInt16,
								},
							},
							{
								Name: "iNT16",
								FieldType: FieldType{
									Simple: typeInt16,
								},
							},
							{
								Name: "iNT16_",
								FieldType: FieldType{
									Simple: typeInt16,
								},
							},
							{
								Name: "iNT32",
								FieldType: FieldType{
									Simple: typeInt32,
								},
							},
							{
								Name: "iNT32_",
								FieldType: FieldType{
									Simple: typeInt32,
								},
							},
							{
								Name:    "tRUE",
								Comment: " int8 nIL; // \"nil\": null,",
								FieldType: FieldType{
									Simple: typeBool,
								},
							},
							{
								Name: "fALSE",
								FieldType: FieldType{
									Simple: typeBool,
								},
							},
							{
								Name: "fLOAT",
								FieldType: FieldType{
									Simple: typeFloat64,
								},
							},
							{
								Name: "fLOAT_x",
								FieldType: FieldType{
									Simple: typeFloat64,
								},
							},
							{
								Name: "sTRING0",
								FieldType: FieldType{
									Simple: typeString,
								},
							},
							{
								Name: "sTRING1",
								FieldType: FieldType{
									Simple: typeString,
								},
							},
							{
								Name: "sTRING4",
								FieldType: FieldType{
									Simple: typeString,
								},
							},
							{
								Name: "sTRING8",
								FieldType: FieldType{
									Simple: typeString,
								},
							},
							{
								Name: "sTRING16",
								FieldType: FieldType{
									Simple: typeString,
								},
							},
							{
								Name: "aRRAY0",
								FieldType: FieldType{
									Array: &FieldType{
										Simple: typeInt32,
									},
								},
							},
							{
								Name: "aRRAY1",
								FieldType: FieldType{
									Array: &FieldType{
										Simple: typeString,
									},
								},
							},
							{
								Name: "aRRAY8",
								FieldType: FieldType{
									Array: &FieldType{
										Simple: typeInt32,
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
						Name:       "FurnitureFamily",
						SimpleType: "uint32",
						Unsigned:   true,
						Options: []EnumOption{
							{
								Name:      "Bed",
								UintValue: 0,
							},
							{
								Name:      "Table",
								UintValue: 1,
							},
							{
								Name:      "Shoe",
								UintValue: 2,
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
									Simple: typeString,
								},
							},
							{
								Name: "price",
								FieldType: FieldType{
									Simple: typeUint32,
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
						OpCode: 0x41454B49,
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
									Simple: typeString,
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
						Comment: "*\n * This union is so documented!\n ",
						Name:    "U",
						OpCode:  bytesToOpCode([4]byte{'y', 'e', 'a', 'h'}),
						Fields: map[uint8]UnionField{
							1: {
								Message: &Message{
									Name: "A",
									Fields: map[uint8]Field{
										1: {
											Name: "b",
											FieldType: FieldType{
												Simple: typeUint32,
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
											Name: "c",
											FieldType: FieldType{
												Simple: typeBool,
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
												Simple: typeUint32,
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
									Comment: " nil is empty",
									Name:    "Null",
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
			t.Parallel()
			f, err := os.Open(filepath.Join("testdata", "base", tc.file+".bop"))
			if err != nil {
				t.Fatalf("failed to open test file %s: %v", tc.file+".bop", err)
			}
			defer f.Close()
			bf, _, err := ReadFile(f)
			if err != nil {
				t.Fatalf("read file errored: %v", err)
			}
			if err := bf.equals(tc.expected); err != nil {
				t.Fatal("parsed file did not match expected:", err)
			}
		})
	}
}

func TestReadFileErrorWarnings(t *testing.T) {
	t.Parallel()
	type testCase struct {
		file       string
		errMessage string
	}
	tcs := []testCase{
		{file: "invalid_const_unparseable_uint", errMessage: "[0:81] strconv.ParseUint: parsing \"2222222222222222222222222222222222222222222222222222222222222222\": value out of range"},
		{file: "invalid_const_unparseable_int", errMessage: "[0:81] strconv.ParseInt: parsing \"33333333333333333333333333333333333333333333333333333333333333333\": value out of range"},
		{file: "invalid_const_unparseable_float", errMessage: "[0:45] strconv.ParseFloat: parsing \"1.7976931348623159e308\": value out of range"},
		{file: "invalid_const_unparseable_float_2", errMessage: "[0:103] strconv.ParseInt: parsing \"6666666666666666666666666666666666666666666666666666666666666666666666666666666666666\": value out of range"},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.file, func(t *testing.T) {
			t.Parallel()
			f, err := os.Open(filepath.Join("testdata", "warning", tc.file+".bop"))
			if err != nil {
				t.Fatalf("failed to open test file %s: %v", tc.file+".bop", err)
			}
			defer f.Close()
			_, warnings, err := ReadFile(f)
			if err != nil {
				t.Fatalf("read file should not have errored: %v", err)
			}
			if warnings[0] != tc.errMessage {
				t.Fatalf("read file had wrong warning: got %q, expected %q", warnings[0], tc.errMessage)
			}
		})
	}
}

func TestReadFileError(t *testing.T) {
	t.Parallel()
	type testCase struct {
		file       string
		errMessage string
	}
	tcs := []testCase{
		{file: "invalid_import_no_file", errMessage: "[0:6] expected (String Literal), got no token"},
		{file: "invalid_const_no_semi", errMessage: "[0:34] expected (Semicolon), got no token"},
		{file: "invalid_const_float_no_semi", errMessage: "[0:36] expected (Semicolon), got no token"},
		{file: "invalid_enum_with_op_code", errMessage: "[1:4] enums may not have attached op codes"},
		{file: "invalid_op_code_1", errMessage: "[0:2] expected (OpCode, Flags) got Close Square"},
		{file: "invalid_op_code_2", errMessage: "[0:6] expected (OpCode, Flags) got Ident"},
		{file: "invalid_op_code_3", errMessage: "[0:15] opcode string \"12345\" not 4 ascii characters"},
		{file: "invalid_op_code_4", errMessage: "[0:8] expected (Open Paren) got Open Square"},
		{file: "invalid_op_code_5", errMessage: "[0:81] strconv.ParseUint: parsing \"1111111111111111111111111111111111111111111111111111111111111111111111111\": value out of range"},
		{file: "invalid_op_code_6", errMessage: "[0:15] expected (Close Paren) got Close Square"},
		{file: "invalid_op_code_7", errMessage: "[0:16] expected (Close Square) got Equals"},
		{file: "invalid_op_code_8", errMessage: "[0:13] expected (Integer Literal, String Literal) got Ident"},
		{file: "invalid_op_code_9", errMessage: "[0:13] opcode string \"123\" not 4 ascii characters"},
		{file: "invalid_enum_bad_deprecated", errMessage: "[1:17] expected (String Literal) got Equals"},
		{file: "invalid_enum_double_deprecated", errMessage: "[2:5] expected enum option following deprecated annotation"},
		{file: "invalid_enum_no_close", errMessage: "[2:0] enum definition ended early"},
		{file: "invalid_enum_no_curly", errMessage: "[1:0] expected (Colon, Open Curly) got Newline"},
		{file: "invalid_enum_no_eq", errMessage: "[1:9] expected (Equals) got Integer Literal"},
		{file: "invalid_enum_no_int", errMessage: "[1:10] expected (Integer Literal) got Semicolon"},
		{file: "invalid_enum_no_name", errMessage: "[0:6] expected (Ident) got Open Curly"},
		{file: "invalid_enum_no_semi", errMessage: "[2:0] expected (Semicolon) got Newline"},
		{file: "invalid_struct_bad_deprecated", errMessage: "[1:20] expected (String Literal) got Ident"},
		{file: "invalid_struct_bad_type", errMessage: "[1:9] expected (Ident, Array, Map) got Open Square"},
		{file: "invalid_struct_double_deprecated", errMessage: "[2:5] expected field following deprecated annotation"},
		{file: "invalid_struct_no_close", errMessage: "[1:14] struct definition ended early"},
		{file: "invalid_struct_no_curly", errMessage: "[1:0] expected (Open Curly) got Newline"},
		{file: "invalid_struct_no_field_name", errMessage: "[1:10] expected (Ident) got Semicolon"},
		{file: "invalid_struct_no_name", errMessage: "[0:8] expected (Ident) got Open Curly"},
		{file: "invalid_struct_no_semi", errMessage: "[2:0] expected (Semicolon) got Newline"},
		{file: "invalid_message_bad_deprecated", errMessage: "[1:18] expected (String Literal) got Arrow"},
		{file: "invalid_message_bad_type", errMessage: "[1:14] expected (Ident, Array, Map) got Open Square"},
		{file: "invalid_message_double_deprecated", errMessage: "[2:5] expected field following deprecated annotation"},
		{file: "invalid_message_hex_int", errMessage: "[1:7] strconv.ParseUint: parsing \"0x1\": invalid syntax"},
		{file: "invalid_message_no_arrow", errMessage: "[1:11] expected (Arrow) got Ident"},
		{file: "invalid_message_no_close", errMessage: "[1:19] expected (Newline, Integer Literal, Open Square, Block Comment, Line Comment, Close Curly), got no token"},
		{file: "invalid_message_no_curly", errMessage: "[1:0] expected (Open Curly) got Newline"},
		{file: "invalid_message_no_field_name", errMessage: "[1:15] expected (Ident) got Semicolon"},
		{file: "invalid_message_no_name", errMessage: "[0:9] expected (Ident) got Open Curly"},
		{file: "invalid_message_no_semi", errMessage: "[2:0] expected (Semicolon) got Newline"},
		{file: "invalid_enum_reserved", errMessage: "[0:11] expected (Ident) got Struct"},
		{file: "invalid_struct_reserved", errMessage: "[0:12] expected (Ident) got Array"},
		{file: "invalid_message_reserved", errMessage: "[0:11] expected (Ident) got Map"},
		{file: "invalid_message_duplicate_index", errMessage: "[2:2] message has duplicate field index 1"},
		{file: "invalid_readonly_enum", errMessage: "[0:13] expected (Struct) got (Enum)"},
		{file: "invalid_readonly_message", errMessage: "[0:16] expected (Struct) got (Message)"},
		{file: "invalid_readonly_comment", errMessage: "[0:19] expected (Struct) got (Block Comment)"},
		{file: "invalid_nested_union", errMessage: "[1:14] union fields must be messages or structs"},
		{file: "invalid_union_double_deprecated", errMessage: "[2:5] expected field following deprecated annotation"},
		{file: "invalid_union_invalid_deprecated", errMessage: "[1:17] unexpected token '!', expected number, letter, or control sequence"},
		{file: "invalid_union_invalid_message", errMessage: "[2:10] strconv.ParseUint: parsing \"-1\": invalid syntax"},
		{file: "invalid_union_invalid_struct", errMessage: "[2:19] expected (Ident) got Semicolon"},
		{file: "invalid_union_invalid_field_number", errMessage: "[1:6] strconv.ParseUint: parsing \"-1\": invalid syntax"},
		{file: "invalid_union_duplicate_field_number", errMessage: "[2:5] union has duplicate field index 1"},
		{file: "invalid_union_missing_arrow", errMessage: "[1:12] expected (Arrow) got Struct"},
		{file: "invalid_union_eof", errMessage: "[0:9] union definition ended early"},
		{file: "invalid_union_eof_2", errMessage: "[0:5] expected (Ident, Open Curly), got no token"},
		{file: "invalid_readonly_eof", errMessage: "[0:8] expected (Struct) got no token"},
		{file: "invalid_const_eof", errMessage: "[0:14] expected (Ident, Ident, Equals), got no token"},
		{file: "invalid_const_eof_2", errMessage: "[0:17] expected value following const type"},
		{file: "invalid_const_unassignable_uint", errMessage: "[0:30] String Literal unassignable to uint32"},
		{file: "invalid_const_unassignable_int", errMessage: "[0:21] Floating Point Literal unassignable to int64"},
		{file: "invalid_const_unassignable_float", errMessage: "[0:23] String Literal unassignable to float32"},
		{file: "invalid_const_unassignable_guid", errMessage: "[0:16] Integer Literal unassignable to guid"},
		{file: "invalid_const_invalid_guid", errMessage: "[0:31] \"guid-guid-guid\" has wrong length for guid"},
		{file: "invalid_const_unassignable_string", errMessage: "[0:19] Integer Literal unassignable to string"},
		{file: "invalid_const_unassignable_bool", errMessage: "[0:21] String Literal unassignable to bool"},
		{file: "invalid_const_invalid_const_type", errMessage: "[0:24] invalid type \"date\" for const"},
		{file: "invalid_const_opcode", errMessage: "[1:5] consts may not have attached op codes"},
		{file: "invalid_map_no_square", errMessage: "[1:8] expected (Open Square) got Semicolon"},
		{file: "invalid_map_keys", errMessage: "[1:7] map must begin with simple type"},
		{file: "invalid_map_keys_2", errMessage: "[1:7] map must begin with simple type"},
		{file: "invalid_map_keys_3", errMessage: "[1:19] map must begin with simple type"},
		{file: "invalid_map_no_comma", errMessage: "[1:15] expected (Comma) got Close Square"},
		{file: "invalid_map_no_close_square", errMessage: "[1:28] expected (Close Square) got Ident"},
		{file: "invalid_array_no_open_square", errMessage: "[1:18] expected (Open Square) got Ident"},
		{file: "invalid_array_bad_key", errMessage: "[1:15] expected (Ident, Array, Map) got Close Square"},
		{file: "invalid_array_no_close_square", errMessage: "[1:19] expected (Close Square) got Ident"},
		{file: "invalid_array_suffix__no_close_square", errMessage: "[1:13] expected (Close Square) got Ident"},
		{file: "invalid_union_no_message_int", errMessage: "[5:14] expected (Newline, Integer Literal, Open Square, Block Comment, Line Comment, Close Curly) got Ident"},
		{file: "invalid_bitflags_unknown_name_uint", errMessage: "[2:9] enum option B undefined"},
		{file: "invalid_bitflags_unknown_name", errMessage: "[2:9] enum option B undefined"},
		{file: "invalid_bitflags_unparseable_int", errMessage: "strconv.ParseInt: parsing \"1111111111111111111111111111111111111111111111111111111111111111111111\": value out of range"},
		{file: "invalid_bitflags_unparseable_uint", errMessage: "strconv.ParseUint: parsing \"-1\": invalid syntax"},
		{file: "invalid_bitflags_unparseable_rhs", errMessage: "strconv.ParseInt: parsing \"1111111111111111111111111111111111111111111111111111111111111111111111\": value out of range"},
		{file: "invalid_bitflags_unparseable_uint_rhs", errMessage: "strconv.ParseUint: parsing \"-1\": invalid syntax"},
		{file: "invalid_array_no_close", errMessage: "[1:12] expected (Ident, Semicolon), got no token"},
		{file: "invalid_enum_bad_type", errMessage: "[0:15] expected an integer enum type"},
		{file: "invalid_enum_no_type", errMessage: "[0:10] expected (Ident) got Open Curly"},
		{file: "invalid_enum_unparseable", errMessage: "strconv.ParseInt: parsing \"77777777\": value out of range"},
		{file: "invalid_enum_unparseable_uint", errMessage: "strconv.ParseUint: parsing \"77777777\": value out of range"},
		{file: "invalid_bitflags_no_semi", errMessage: "[3:0] eof reading until Semicolon"},
		{file: "invalid_bitflags_on_struct", errMessage: "[1:6] structs may not use bitflags"},
		{file: "invalid_bitflags_on_message", errMessage: "[1:7] messages may not use bitflags"},
		{file: "invalid_bitflags_on_union", errMessage: "[1:5] unions may not use bitflags"},
		{file: "invalid_bitflags_on_const", errMessage: "[1:5] consts may not use bitflags"},
		{file: "invalid_bitflags_no_close_bracket", errMessage: "[0:6] expected (Close Square), got no token"},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.file, func(t *testing.T) {
			t.Parallel()
			f, err := os.Open(filepath.Join("testdata", "invalid", tc.file+".bop"))
			if err != nil {
				t.Fatalf("failed to open test file %s: %v", tc.file+".bop", err)
			}
			defer f.Close()
			_, _, err = ReadFile(f)
			if err == nil {
				t.Fatalf("read file should have errored")
			}
			if err.Error() != tc.errMessage {
				t.Fatalf("read file had wrong error: got %q, expected %q", err.Error(), tc.errMessage)
			}
		})
	}
}

func Test_parseCommentTag(t *testing.T) {
	t.Parallel()
	t.Run("un-unquoteable", func(t *testing.T) {
		t.Parallel()
		s := "[tag(k:\"foo)]"
		_, ok := parseCommentTag(s)
		if ok {
			t.Fatalf("parseCommentTag should have failed")
		}
	})
}
