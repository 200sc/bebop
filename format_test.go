package bebop

import (
	"os"
	"path/filepath"
	"testing"
)

var testFormatFiles = []string{
	"all_consts",
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
	"request",
	"server",
	"union",
}

const fileExt = ".bop"

func TestTokenizeFormat(t *testing.T) {
	t.Parallel()
	for _, filename := range testFormatFiles {
		filename := filename
		t.Run(filename, func(t *testing.T) {
			t.Parallel()
			formatFile(t, "base", filename)
		})
	}
}

var testIncompatibleFiles = []string{
	"quoted_string",
}

func TestTokenizeFormatIncompatible(t *testing.T) {
	t.Parallel()
	for _, filename := range testIncompatibleFiles {
		filename := filename
		t.Run(filename, func(t *testing.T) {
			t.Parallel()
			formatFile(t, "incompatible", filename)
		})
	}
}

func formatFile(t *testing.T, subdir, filename string) {
	t.Helper()
	const formattedExt = "_formatted.bop"

	f, err := os.Open(filepath.Join("testdata", subdir, filename+fileExt))
	if err != nil {
		t.Fatalf("failed to open test file %s: %v", filename+fileExt, err)
	}

	if _, _, err := ReadFile(f); err != nil {
		f.Close()
		t.Fatalf("can not format unparsable files")
	}
	f.Close()
	f, err = os.Open(filepath.Join("testdata", subdir, filename+fileExt))
	if err != nil {
		t.Fatalf("failed to open test file (for format) %s: %v", filename+fileExt, err)
	}

	tr := newTokenReader(f)
	out, err := os.Create(filepath.Join("testdata", "formatted", filename+formattedExt))
	if err != nil {
		t.Fatalf("failed to open out file %s: %v", filename+formattedExt, err)
	}
	defer out.Close()
	err = format(tr, out)
	if err != nil {
		t.Fatalf("failed to format %s: %v", filename+formattedExt, err)
	}
}
