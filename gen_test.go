package bebop

import (
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestValidateError(t *testing.T) {
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
	}, {
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
	}}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.file, func(t *testing.T) {
			f, err := os.Open(filepath.Join("testdata", "base", tc.file+".bop"))
			if err != nil {
				t.Fatalf("failed to open test file %s: %v", tc.file+".bop", err)
			}
			defer f.Close()
			bopf, err := ReadFile(f)
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

var genTestFiles = []string{
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
	"quoted_string",
	"request",
	"server",
	"union",
	"date",
	"message_1",
}

func TestGenerateToFile(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	for _, filename := range genTestFiles {
		filename := filename
		t.Run(filename, func(t *testing.T) {
			f, err := os.Open(filepath.Join("testdata", "base", filename+".bop"))
			if err != nil {
				t.Fatalf("failed to open test file %s: %v", filename+".bop", err)
			}
			defer f.Close()
			bopf, err := ReadFile(f)
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
				SharedMemoryStrings:   rand.Float64() < .5,
			})
			if err != nil {
				t.Fatalf("generation failed: %v", err)
			}
		})
	}
}
