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
