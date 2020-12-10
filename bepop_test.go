package bebop

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

var readFileTestFiles = []string{
	"enums",
}

func TestReadFile(t *testing.T) {
	for _, filename := range readFileTestFiles {
		t.Run(filename, func(t *testing.T) {
			f, err := os.Open(filepath.Join("testdata", filename+".bop"))
			if err != nil {
				t.Fatalf("failed to open test file %s: %v", filename+".bop", err)
			}
			defer f.Close()
			bf, err := ReadFile(f)
			if err != nil {
				t.Fatalf("read file errored: %v", err)
			}
			fmt.Println(bf)
		})
	}
}
