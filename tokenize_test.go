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
	"enums",
	"enums_doc",
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
		filename := filename
		t.Run(filename, func(t *testing.T) {
			filename += ".bop"
			f, err := os.Open(filepath.Join("testdata", filename))
			if err != nil {
				t.Fatalf("failed to open test file %s: %v", filename, err)
			}
			defer f.Close()
			tr := newTokenReader(f)
			for tr.Next() {
			}
			if tr.Err() != nil {
				t.Fatalf("token reader errored: %v", tr.Err())
			}
		})
	}
}

var testFormatFiles = []string{
	"array_of_strings",
	"basic_arrays",
	"basic_types",
	"documentation",
	"enums",
	"enums_doc",
	"foo",
	"jazz",
	"lab",
	"map_types",
	"msgpack_comparison",
	"quoted_string",
	"request",
}

func TestTokenizeFormat(t *testing.T) {
	for _, filename := range testFormatFiles {
		filename := filename
		t.Run(filename, func(t *testing.T) {
			f, err := os.Open(filepath.Join("testdata", filename+".bop"))
			if err != nil {
				t.Fatalf("failed to open test file %s: %v", filename+".bop", err)
			}
			if _, err := ReadFile(f); err != nil {
				f.Close()
				t.Fatalf("can not format unparseable files")
			}
			f.Close()
			f, err = os.Open(filepath.Join("testdata", filename+".bop"))
			if err != nil {
				t.Fatalf("failed to open test file (for format) %s: %v", filename+".bop", err)
			}

			tr := newTokenReader(f)
			out, err := os.Create(filepath.Join("testdata", "formatted", filename+"_formatted.bop"))
			if err != nil {
				t.Fatalf("failed to open out file %s: %v", filename+"_formatted.bop", err)
			}
			defer out.Close()
			format(tr, out)
		})
	}
}
