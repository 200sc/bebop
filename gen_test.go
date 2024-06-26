package bebop

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateIncompatibleError(t *testing.T) {
	t.Parallel()
	type testCase struct {
		file string
		err  string
	}
	tcs := []testCase{{
		file: "recursive_struct",
		err:  "recursively includes itself as a required field",
	}, {
		file: "invalid_enum_primitive",
		err:  "enum shares primitive type name uint8",
	}, {
		file: "invalid_struct_primitive",
		err:  "struct shares primitive type name string",
	}, {
		file: "invalid_message_primitive",
		err:  "message shares primitive type name guid",
	}}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.file, func(t *testing.T) {
			t.Parallel()
			f, err := os.Open(filepath.Join("testdata", "incompatible", tc.file+fileExt))
			if err != nil {
				t.Fatalf("failed to open test file %s: %v", tc.file+fileExt, err)
			}
			defer f.Close()
			bopf, _, err := ReadFile(f)
			if err != nil {
				t.Fatalf("failed to read file %s: %v", tc.file+fileExt, err)
			}
			err = bopf.Validate()
			if err == nil {
				t.Fatalf("validation did not fail")
			}
			if !strings.Contains(err.Error(), tc.err) {
				t.Fatalf("validation did not have expected error: got %q expected %q", err.Error(), tc.err)
			}
		})
	}
}

func TestValidateError(t *testing.T) {
	t.Parallel()
	type testCase struct {
		file string
		err  string
	}
	tcs := []testCase{{
		file: "invalid_enum_duplicate",
		err:  "enum has duplicated name myEnum",
	}, {
		file: "invalid_struct_duplicate",
		err:  "struct has duplicated name mystruct",
	}, {
		file: "invalid_message_duplicate",
		err:  "message has duplicated name mymessage",
	}, {
		file: "invalid_struct_unknown",
		err:  "type whereisthistype undefined",
	}, {
		file: "invalid_message_unknown",
		err:  "type whereisthistype undefined",
	}, {
		file: "invalid_enum_duplicate_index",
		err:  "enum MyEnum has duplicate option value 1",
	}, {
		file: "invalid_enum_duplicate_name",
		err:  "enum MyEnum has duplicate option name A",
	}, {
		file: "invalid_message_duplicate_name",
		err:  "message Test has duplicate field name foo",
	}, {
		file: "invalid_struct_duplicate_name",
		err:  "struct Test has duplicate field name foo",
	}, {
		file: "invalid_const_duplicate_name",
		err:  "const has duplicated name hello",
	}, {
		file: "invalid_union_duplicate_name",
		err:  "union Test has duplicate field name foo",
	}, {
		file: "invalid_union_primitive_name",
		err:  "union shares primitive type name uint8",
	}, {
		file: "invalid_op_code_10",
		err:  "struct InvalidOpCode10B has duplicate opcode 34333231 (duplicated in InvalidOpCode10A)",
	}, {
		file: "invalid_op_code_11",
		err:  "struct InvalidOpCode11B has duplicate opcode 34333231 (duplicated in InvalidOpCode11A)",
	}, {
		file: "invalid_op_code_12",
		err:  "message InvalidOpCode12B has duplicate opcode 37363534 (duplicated in InvalidOpCode12A)",
	}, {
		file: "invalid_op_code_13",
		err:  "union InvalidOpCode13B has duplicate opcode 30313938 (duplicated in InvalidOpCode13A)",
	}}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.file, func(t *testing.T) {
			t.Parallel()
			f, err := os.Open(filepath.Join("testdata", "invalid", tc.file+fileExt))
			if err != nil {
				t.Fatalf("failed to open test file %s: %v", tc.file+fileExt, err)
			}
			defer f.Close()
			bopf, _, err := ReadFile(f)
			if err != nil {
				t.Fatalf("failed to read file %s: %v", tc.file+fileExt, err)
			}
			err = bopf.Validate()
			if err == nil {
				t.Fatalf("validation did not fail")
			}
			if !strings.Contains(err.Error(), tc.err) {
				t.Fatalf("validation did not have expected error: got %q expected %q", err.Error(), tc.err)
			}
		})
	}
}

func TestGenerate_Error(t *testing.T) {
	t.Parallel()
	type testCase struct {
		file string
		err  string
	}
	tcs := []testCase{{
		file: "invalid_import_file_not_found",
		err:  "failed to open imported file ../../hello_world.bop",
	}, {
		file: "invalid_import_file_not_parsable",
		err:  "failed to parse imported file ./invalid_array_no_close_square.bop",
	}}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.file, func(t *testing.T) {
			t.Parallel()
			f, err := os.Open(filepath.Join("testdata", "invalid", tc.file+fileExt))
			if err != nil {
				t.Fatalf("failed to open test file %s: %v", tc.file+fileExt, err)
			}
			defer f.Close()
			bopf, _, err := ReadFile(f)
			if err != nil {
				t.Fatalf("failed to read file %s: %v", tc.file+fileExt, err)
			}
			err = bopf.Generate(bytes.NewBuffer([]byte{}), GenerateSettings{})
			if err == nil {
				t.Fatalf("validation did not fail")
			}
			if !strings.Contains(err.Error(), tc.err) {
				t.Fatalf("validation did not have expected error: got %q expected %q", err.Error(), tc.err)
			}
		})
	}
}

func TestGenerateToFile_SeperateImports(t *testing.T) {
	t.Parallel()
	type file struct {
		filename string
		outfile  string
	}
	files := []file{
		{
			filename: "import_separate_a",
			outfile:  filepath.Join("generated", "import_separate_a.go"),
		}, {
			filename: "import_separate_b",
			outfile:  filepath.Join("generatedtwo", "import_separate_b.go"),
		}, {
			filename: "import_separate_c",
			outfile:  filepath.Join("generatedthree", "import_separate_c.go"),
		},
	}
	for _, filedef := range files {
		filedef := filedef
		t.Run(filedef.filename, func(t *testing.T) {
			t.Parallel()
			f, err := os.Open(filepath.Join("testdata", "incompatible", filedef.filename+fileExt))
			if err != nil {
				t.Fatalf("failed to open test file %s: %v", filedef.filename+fileExt, err)
			}
			defer f.Close()
			bopf, _, err := ReadFile(f)
			if err != nil {
				t.Fatalf("failed to read file %s: %v", filedef.filename+fileExt, err)
			}
			// use a separate directory to ensure duplicate definitions in combined mode
			// do not cause compilation failures
			err = os.MkdirAll(filepath.Join("testdata", "incompatible", filepath.Dir(filedef.outfile)), 0777)
			if err != nil {
				t.Fatalf("failed to mkdir: %v", err)
			}
			outFile := filepath.Join("testdata", "incompatible", filedef.outfile)
			out, err := os.Create(outFile)
			if err != nil {
				t.Fatalf("failed to open out file %s: %v", outFile, err)
			}
			defer out.Close()
			err = bopf.Generate(out, GenerateSettings{
				GenerateUnsafeMethods: true,
				SharedMemoryStrings:   false,
				ImportGenerationMode:  ImportGenerationModeSeparate,
			})
			if err != nil {
				t.Fatalf("generation failed: %v", err)
			}
		})
	}
}

func TestGenerateToFile_SeperateImports_ImportCycle(t *testing.T) {
	t.Parallel()
	type file struct {
		filename   string
		errMessage string
	}
	files := []file{
		{
			filename:   "import_loop_a",
			errMessage: "import cycle found:", // the cycle itself is not deterministic, depending on which node is scanned first
		},
	}
	for _, filedef := range files {
		filedef := filedef
		t.Run(filedef.filename, func(t *testing.T) {
			t.Parallel()
			f, err := os.Open(filepath.Join("testdata", "incompatible", filedef.filename+fileExt))
			if err != nil {
				t.Fatalf("failed to open test file %s: %v", filedef.filename+fileExt, err)
			}
			defer f.Close()
			bopf, _, err := ReadFile(f)
			if err != nil {
				t.Fatalf("failed to read file %s: %v", filedef.filename+fileExt, err)
			}
			err = bopf.Generate(io.Discard, GenerateSettings{
				GenerateUnsafeMethods: true,
				SharedMemoryStrings:   false,
				ImportGenerationMode:  ImportGenerationModeSeparate,
			})
			if err == nil {
				t.Fatalf("generate had no error: expected %q", filedef.errMessage)
			}
			if !strings.HasPrefix(err.Error(), filedef.errMessage) {
				t.Fatalf("generate had wrong error: got %q, expected %q", err.Error(), filedef.errMessage)
			}
		})
	}
}

func TestGenerateToFile_Error(t *testing.T) {
	t.Parallel()
	type testCase struct {
		name string
		// provide one of the following two:
		File
		filename string
		GenerateSettings
		errMessage string
	}
	tcs := []testCase{
		{
			name:             "no package definition",
			File:             File{},
			GenerateSettings: GenerateSettings{},
			errMessage:       "no package name is defined, provide a go_package const or an explicit package name setting",
		}, {
			name: "bad import strategy",
			File: File{},
			GenerateSettings: GenerateSettings{
				ImportGenerationMode: 5,
			},
			errMessage: "invalid generation settings: unknown import mode: 5",
		}, {
			name:     "no go package const in import",
			filename: "invalid_import_no_const",
			GenerateSettings: GenerateSettings{
				GenerateUnsafeMethods: true,
				SharedMemoryStrings:   true,
				ImportGenerationMode:  ImportGenerationModeSeparate,
			},
			errMessage: "cannot import quoted_string.bop: file has no go_package const",
		}, {
			name: "undefined map key type",
			File: File{
				Structs: []Struct{
					{
						Fields: []Field{
							{
								Name: "F",
								FieldType: FieldType{
									Map: &MapType{
										Key: "boofus",
										Value: FieldType{
											Simple: "bofus",
										},
									},
								},
							},
						},
					},
				},
			},
			GenerateSettings: GenerateSettings{
				GenerateUnsafeMethods: true,
				SharedMemoryStrings:   true,
				ImportGenerationMode:  ImportGenerationModeSeparate,
			},
			errMessage: "cannot generate file: map key type boofus undefined",
		},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if tc.filename != "" {
				f, err := os.Open(filepath.Join("testdata", "incompatible", tc.filename+fileExt))
				if err != nil {
					t.Fatalf("failed to open test file %s: %v", tc.filename+fileExt, err)
				}
				defer f.Close()
				tc.File, _, err = ReadFile(f)
				if err != nil {
					t.Fatalf("failed to read file %s: %v", tc.filename+fileExt, err)
				}
			}
			err := tc.File.Generate(io.Discard, tc.GenerateSettings)
			if err == nil {
				t.Fatalf("generate had no error: expected %q", tc.errMessage)
			}
			if err.Error() != tc.errMessage {
				t.Fatalf("generate had wrong error: got %q, expected %q", err.Error(), tc.errMessage)
			}
		})
	}
}

var importFiles = []string{
	"import_b",
}

func TestGenerateToFile_Imports(t *testing.T) {
	t.Parallel()
	for _, filename := range importFiles {
		filename := filename
		t.Run(filename, func(t *testing.T) {
			t.Parallel()
			f, err := os.Open(filepath.Join("testdata", "base", filename+fileExt))
			if err != nil {
				t.Fatalf("failed to open test file %s: %v", filename+fileExt, err)
			}
			defer f.Close()
			bopf, _, err := ReadFile(f)
			if err != nil {
				t.Fatalf("failed to read file %s: %v", filename+fileExt, err)
			}
			// use a separate directory to ensure duplicate definitions in combined mode
			// do not cause compilation failures
			err = os.MkdirAll(filepath.Join("testdata", "generated", filename), 0777)
			if err != nil {
				t.Fatalf("failed to mkdir: %v", err)
			}
			outFile := filepath.Join("testdata", "generated", filename, filename+".go")
			out, err := os.Create(outFile)
			if err != nil {
				t.Fatalf("failed to open out file %s: %v", outFile, err)
			}
			defer out.Close()
			err = bopf.Generate(out, GenerateSettings{
				PackageName:           "filename",
				GenerateUnsafeMethods: true,
				SharedMemoryStrings:   false,
				ImportGenerationMode:  ImportGenerationModeCombined,
			})
			if err != nil {
				t.Fatalf("generation failed: %v", err)
			}
		})
	}
}

var genTestFiles = []string{
	"all_consts",
	"array_of_strings",
	"arrays",
	"basic_arrays",
	"basic_types",
	"bitflags",
	"documentation",
	"enums",
	"foo",
	"fruit",
	"jazz",
	"lab",
	"map_types",
	"message_inline",
	"message_map",
	"msgpack_comparison",
	"request",
	"server",
	"union",
	"union_field",
	"date",
	"message_1",
	"tags",
	"typed_enums",
}

func TestGenerateToFile(t *testing.T) {
	t.Parallel()
	for _, filename := range genTestFiles {
		filename := filename
		t.Run(filename, func(t *testing.T) {
			t.Parallel()
			f, err := os.Open(filepath.Join("testdata", "base", filename+fileExt))
			if err != nil {
				t.Fatalf("failed to open test file %s: %v", filename+fileExt, err)
			}
			defer f.Close()
			bopf, _, err := ReadFile(f)
			if err != nil {
				t.Fatalf("failed to read file %s: %v", filename+fileExt, err)
			}
			outFile := filepath.Join("testdata", "generated", filename+".go")
			out, err := os.Create(outFile)
			if err != nil {
				t.Fatalf("failed to open out file %s: %v", outFile, err)
			}
			defer out.Close()
			err = bopf.Generate(out, GenerateSettings{
				PackageName:           "generated",
				GenerateUnsafeMethods: true,
				SharedMemoryStrings:   false,
				GenerateFieldTags:     true,
			})
			if err != nil {
				t.Fatalf("generation failed: %v", err)
			}
		})
	}
}

func TestGenerateToFile_Private(t *testing.T) {
	t.Parallel()
	err := os.MkdirAll(filepath.Join("testdata", "generated-private"), 0666)
	if err != nil {
		t.Fatal(err)
	}
	for _, filename := range genTestFiles {
		filename := filename
		t.Run(filename, func(t *testing.T) {
			t.Parallel()
			f, err := os.Open(filepath.Join("testdata", "base", filename+fileExt))
			if err != nil {
				t.Fatalf("failed to open test file %s: %v", filename+fileExt, err)
			}
			defer f.Close()
			bopf, _, err := ReadFile(f)
			if err != nil {
				t.Fatalf("failed to read file %s: %v", filename+fileExt, err)
			}
			outFile := filepath.Join("testdata", "generated-private", filename+".go")
			out, err := os.Create(outFile)
			if err != nil {
				t.Fatalf("failed to open out file %s: %v", outFile, err)
			}
			defer out.Close()
			err = bopf.Generate(out, GenerateSettings{
				PackageName:           "generated",
				GenerateUnsafeMethods: true,
				SharedMemoryStrings:   false,
				GenerateFieldTags:     true,
				PrivateDefinitions:    true,
			})
			if err != nil {
				t.Fatalf("generation failed: %v", err)
			}
		})
	}
}

func TestGenerateToFile_AlwaysPointers(t *testing.T) {
	t.Parallel()
	err := os.MkdirAll(filepath.Join("testdata", "generated-always-pointers"), 0666)
	if err != nil {
		t.Fatal(err)
	}
	for _, filename := range genTestFiles {
		filename := filename
		t.Run(filename, func(t *testing.T) {
			t.Parallel()
			f, err := os.Open(filepath.Join("testdata", "base", filename+fileExt))
			if err != nil {
				t.Fatalf("failed to open test file %s: %v", filename+fileExt, err)
			}
			defer f.Close()
			bopf, _, err := ReadFile(f)
			if err != nil {
				t.Fatalf("failed to read file %s: %v", filename+fileExt, err)
			}
			outFile := filepath.Join("testdata", "generated-always-pointers", filename+".go")
			out, err := os.Create(outFile)
			if err != nil {
				t.Fatalf("failed to open out file %s: %v", outFile, err)
			}
			defer out.Close()
			err = bopf.Generate(out, GenerateSettings{
				PackageName:               "generated",
				GenerateUnsafeMethods:     true,
				SharedMemoryStrings:       false,
				GenerateFieldTags:         true,
				AlwaysUsePointerReceivers: true,
			})
			if err != nil {
				t.Fatalf("generation failed: %v", err)
			}
		})
	}
}

func TestGenerateToFileIncompatible(t *testing.T) {
	t.Parallel()
	for _, filename := range testIncompatibleFiles {
		filename := filename
		t.Run(filename, func(t *testing.T) {
			t.Parallel()
			f, err := os.Open(filepath.Join("testdata", "incompatible", filename+fileExt))
			if err != nil {
				t.Fatalf("failed to open test file %s: %v", filename+fileExt, err)
			}
			defer f.Close()
			bopf, _, err := ReadFile(f)
			if err != nil {
				t.Fatalf("failed to read file %s: %v", filename+fileExt, err)
			}
			outFile := filepath.Join("testdata", "generated", filename+".go")
			out, err := os.Create(outFile)
			if err != nil {
				t.Fatalf("failed to open out file %s: %v", outFile, err)
			}
			defer out.Close()
			err = bopf.Generate(out, GenerateSettings{
				PackageName:           "generated",
				GenerateUnsafeMethods: true,
				SharedMemoryStrings:   false,
			})
			if err != nil {
				t.Fatalf("generation failed: %v", err)
			}
		})
	}
}
