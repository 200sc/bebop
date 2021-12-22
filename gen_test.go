package bebop

import (
	"bytes"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
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
			f, err := os.Open(filepath.Join("testdata", "incompatible", tc.file+".bop"))
			if err != nil {
				t.Fatalf("failed to open test file %s: %v", tc.file+".bop", err)
			}
			defer f.Close()
			bopf, _, err := ReadFile(f)
			if err != nil {
				t.Fatalf("failed to read file %s: %v", tc.file+".bop", err)
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
	}}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.file, func(t *testing.T) {
			t.Parallel()
			f, err := os.Open(filepath.Join("testdata", "invalid", tc.file+".bop"))
			if err != nil {
				t.Fatalf("failed to open test file %s: %v", tc.file+".bop", err)
			}
			defer f.Close()
			bopf, _, err := ReadFile(f)
			if err != nil {
				t.Fatalf("failed to read file %s: %v", tc.file+".bop", err)
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
			f, err := os.Open(filepath.Join("testdata", "invalid", tc.file+".bop"))
			if err != nil {
				t.Fatalf("failed to open test file %s: %v", tc.file+".bop", err)
			}
			defer f.Close()
			bopf, _, err := ReadFile(f)
			if err != nil {
				t.Fatalf("failed to read file %s: %v", tc.file+".bop", err)
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
			f, err := os.Open(filepath.Join("testdata", "incompatible", filedef.filename+".bop"))
			if err != nil {
				t.Fatalf("failed to open test file %s: %v", filedef.filename+".bop", err)
			}
			defer f.Close()
			bopf, _, err := ReadFile(f)
			if err != nil {
				t.Fatalf("failed to read file %s: %v", filedef.filename+".bop", err)
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
			f, err := os.Open(filepath.Join("testdata", "incompatible", filedef.filename+".bop"))
			if err != nil {
				t.Fatalf("failed to open test file %s: %v", filedef.filename+".bop", err)
			}
			defer f.Close()
			bopf, _, err := ReadFile(f)
			if err != nil {
				t.Fatalf("failed to read file %s: %v", filedef.filename+".bop", err)
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
				f, err := os.Open(filepath.Join("testdata", "incompatible", tc.filename+".bop"))
				if err != nil {
					t.Fatalf("failed to open test file %s: %v", tc.filename+".bop", err)
				}
				defer f.Close()
				tc.File, _, err = ReadFile(f)
				if err != nil {
					t.Fatalf("failed to read file %s: %v", tc.filename+".bop", err)
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
	rand.Seed(time.Now().UnixNano())
	for _, filename := range importFiles {
		filename := filename
		t.Run(filename, func(t *testing.T) {
			t.Parallel()
			f, err := os.Open(filepath.Join("testdata", "base", filename+".bop"))
			if err != nil {
				t.Fatalf("failed to open test file %s: %v", filename+".bop", err)
			}
			defer f.Close()
			bopf, _, err := ReadFile(f)
			if err != nil {
				t.Fatalf("failed to read file %s: %v", filename+".bop", err)
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
	"basic_arrays",
	"basic_types",
	"documentation",
	"enums",
	"foo",
	"jazz",
	"lab",
	"map_types",
	"message_map",
	"msgpack_comparison",
	"request",
	"server",
	"union",
	"union_field",
	"date",
	"message_1",
	"tags",
}

func TestGenerateToFile(t *testing.T) {
	t.Parallel()
	rand.Seed(time.Now().UnixNano())
	for _, filename := range genTestFiles {
		filename := filename
		t.Run(filename, func(t *testing.T) {
			t.Parallel()
			f, err := os.Open(filepath.Join("testdata", "base", filename+".bop"))
			if err != nil {
				t.Fatalf("failed to open test file %s: %v", filename+".bop", err)
			}
			defer f.Close()
			bopf, _, err := ReadFile(f)
			if err != nil {
				t.Fatalf("failed to read file %s: %v", filename+".bop", err)
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

func TestGenerateToFileIncompatible(t *testing.T) {
	t.Parallel()
	rand.Seed(time.Now().UnixNano())
	for _, filename := range testIncompatibleFiles {
		filename := filename
		t.Run(filename, func(t *testing.T) {
			t.Parallel()
			f, err := os.Open(filepath.Join("testdata", "incompatible", filename+".bop"))
			if err != nil {
				t.Fatalf("failed to open test file %s: %v", filename+".bop", err)
			}
			defer f.Close()
			bopf, _, err := ReadFile(f)
			if err != nil {
				t.Fatalf("failed to read file %s: %v", filename+".bop", err)
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
