package bebop

import (
	"os"
	"path/filepath"
	"testing"
)

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
			out, err := os.Create(filepath.Join("testdata", filename+".go"))
			if err != nil {
				t.Fatalf("failed to open out file %s: %v", filename+"_formatted.bop", err)
			}
			defer out.Close()
			bopf.Generate(out, GenerateSettings{
				PackageName: "testdata",
			})
		})
	}
}
