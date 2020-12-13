package bebop

import (
	"os"
	"path/filepath"
	"testing"
)

var testFiles = []string{
	"array_of_strings",
	"basic_arrays",
	"basic_types",
	"documentation",
	"foo",
	"invalid_map_keys",
	"invalid_syntax",
	"jazz",
	"lab",
	"map_types",
	"msgpack_comparison",
	"quoted_string",
	"request",
}

func TestTokenize(t *testing.T) {
	for _, filename := range testFiles {
		t.Run(filename, func(t *testing.T) {
			filename += ".bop"
			f, err := os.Open(filepath.Join("testdata", filename))
			if err != nil {
				t.Fatalf("failed to open test file %s: %v", filename, err)
			}
			defer f.Close()
			tr := newTokenReader(f)
			tokens := []token{}
			for tr.Next() {
				tokens = append(tokens, tr.Token())
			}
			if tr.Err() != nil {
				t.Fatalf("token reader errored: %v", tr.Err())
			}
		})
	}
}

func TestTokenizeFormat(t *testing.T) {
	for _, filename := range testFiles {
		t.Run(filename, func(t *testing.T) {
			f, err := os.Open(filepath.Join("testdata", filename+".bop"))
			if err != nil {
				t.Fatalf("failed to open test file %s: %v", filename+".bop", err)
			}
			defer f.Close()
			tr := newTokenReader(f)
			tokens := []token{}
			for tr.Next() {
				tokens = append(tokens, tr.Token())
			}
			if tr.Err() != nil {
				t.Fatalf("token reader errored: %v", tr.Err())
			}
			out, err := os.Create(filepath.Join("testdata", filename+"_formatted.bop"))
			if err != nil {
				t.Fatalf("failed to open out file %s: %v", filename+"_formatted.bop", err)
			}
			defer out.Close()
			format(tokens, out)
		})
	}
}
