package bebop

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateError(t *testing.T) {
	type testCase struct {
		file string
		err  string
	}
	tcs := []testCase{{
		file: "recursive_struct",
		err:  "recursively includes itself as a required field",
	}}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.file, func(t *testing.T) {
			f, err := os.Open(filepath.Join("testdata", tc.file+".bop"))
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
	"msgpack_comparison",
	"quoted_string",
	"request",
}

func TestGenerateToFile(t *testing.T) {
	for _, filename := range genTestFiles {
		t.Run(filename, func(t *testing.T) {
			f, err := os.Open(filepath.Join("testdata", filename+".bop"))
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
				PackageName: "generated",
			})
			if err != nil {
				t.Fatalf("generation failed: %v", err)
			}
		})
	}
}
