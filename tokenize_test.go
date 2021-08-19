package bebop

import (
	"os"
	"path/filepath"
	"testing"
)

var testFiles = []string{
	"array_of_strings",
	"all_consts",
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
	"union",
}

func TestTokenize(t *testing.T) {
	for _, filename := range testFiles {
		filename := filename
		t.Run(filename, func(t *testing.T) {
			filename += ".bop"
			f, err := os.Open(filepath.Join("testdata", "base", filename))
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

var testNoSemisFiles = []string{
	"jazz",
}

func TestTokenizeNoSemis(t *testing.T) {
	optionalSemicolons = true
	for _, filename := range testNoSemisFiles {
		filename := filename
		t.Run(filename, func(t *testing.T) {
			origfilename := filename + ".bop"
			f, err := os.Open(filepath.Join("testdata", "base", origfilename))
			if err != nil {
				t.Fatalf("failed to open test file %s: %v", origfilename, err)
			}
			defer f.Close()

			tokens := []token{}
			tr := newTokenReader(f)
			for tr.Next() {
				tokens = append(tokens, tr.Token())
			}
			if tr.Err() != nil {
				t.Fatalf("token reader errored: %v", tr.Err())
			}

			noSemiFilepath := filename + "_nosemis.bop"
			f2, err := os.Open(filepath.Join("testdata", "base", noSemiFilepath))
			if err != nil {
				t.Fatalf("failed to open test file %s: %v", noSemiFilepath, err)
			}
			defer f2.Close()

			tokens2 := []token{}
			tr2 := newTokenReader(f2)
			for tr2.Next() {
				tokens2 = append(tokens2, tr2.Token())
			}
			if tr2.Err() != nil {
				t.Fatalf("token reader errored (nosemi): %v", tr.Err())
			}

			if len(tokens) != len(tokens2) {
				t.Fatalf("tokens had different lengths: %v vs %v", len(tokens), len(tokens2))
			}

			for i, tk := range tokens {
				tk2 := tokens2[i]
				if tk2.kind != tk.kind {
					t.Fatalf("tokens at pos %d differed: kind %v vs %v", i, tk.kind, tk2.kind)
				}
				if string(tk2.concrete) != string(tk.concrete) {
					t.Fatalf("tokens at pos %d differed: concrete %v vs %v", i, tk.concrete, tk2.concrete)
				}
			}
		})
	}
	optionalSemicolons = false
}

func TestTokenizeError(t *testing.T) {
	type testCase struct {
		file string
		err  string
	}
	tcs := []testCase{
		{file: "invalid_token_arrow_eof", err: "[0:2] eof waiting for (['>', number]) after '-'"},
		{file: "invalid_token_arrow", err: "[0:2] unexpected token '-' waiting for (['>', number]) after '-'"},
		{file: "invalid_token_comment_eof", err: "[0:2] eof waiting for '[/, *]' after '/'"},
		{file: "invalid_token_comment", err: "[0:2] unexpected token 'p' waiting for '[/, *]' after '/'"},
		{file: "invalid_token_string", err: "[0:14] eof waiting for string end quote"},
		{file: "invalid_token_unknown", err: "[0:1] unexpected token '*', expected number, letter, or control sequence"},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.file, func(t *testing.T) {
			origfilename := tc.file + ".bop"
			f, err := os.Open(filepath.Join("testdata", "base", origfilename))
			if err != nil {
				t.Fatalf("failed to open test file %s: %v", origfilename, err)
			}
			defer f.Close()

			tr := newTokenReader(f)
			for tr.Next() {
			}
			err = tr.Err()
			if err == nil {
				t.Fatalf("token reader did not error")
			}
			if err.Error() != tc.err {
				t.Fatalf("tokenization had unexpected error: got %q, wanted %q", err.Error(), tc.err)
			}
		})
	}
}
