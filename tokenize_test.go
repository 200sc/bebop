package bebop

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
)

var testTokenizeFiles = []string{
	"array_of_strings",
	"all_consts",
	"basic_arrays",
	"basic_types",
	"documentation",
	"enums",
	"enums_doc",
	"foo",
	"import",
	"import_b",
	"jazz",
	"lab",
	"map_types",
	"msgpack_comparison",
	"request",
	"union",
	"bitflags",
	"typed_enums",
	"decorations",
}

func TestTokenize(t *testing.T) {
	t.Parallel()
	for _, filename := range testTokenizeFiles {
		filename := filename
		t.Run(filename, func(t *testing.T) {
			t.Parallel()
			filename += fileExt
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

func TestTokenizeIncompatible(t *testing.T) {
	for _, filename := range testIncompatibleFiles {
		filename := filename
		t.Run(filename, func(t *testing.T) {
			t.Parallel()
			filename += fileExt
			f, err := os.Open(filepath.Join("testdata", "incompatible", filename))
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

var testFilesInvalidNoError = []string{
	"invalid_map_keys",
	"invalid_syntax",
}

func TestTokenizeInvalidNoError(t *testing.T) {
	t.Parallel()
	for _, filename := range testFilesInvalidNoError {
		filename := filename
		t.Run(filename, func(t *testing.T) {
			t.Parallel()
			filename += fileExt
			f, err := os.Open(filepath.Join("testdata", "invalid", filename))
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

func TestTokenizeNoSemis(t *testing.T) {
	t.Parallel()
	failingCases := map[string]string{
		"all_consts":         "we do not inject a semi at EOF",
		"foo":                "inline semis between fields must stay",
		"lab":                "inline semis between fields must stay",
		"map_types":          "inline semis between fields must stay",
		"union":              "inline semis between fields must stay",
		"import_b":           "import semis cannot be added",
		"import":             "import semis cannot be added",
		"msgpack_comparison": "we are naively removing semis from within comments",
		"decorations":        "TODO",
	}
	for _, filename := range testTokenizeFiles {
		filename := filename
		t.Run(filename, func(t *testing.T) {
			t.Parallel()
			if reason, ok := failingCases[filename]; ok {
				t.Skip(reason)
			}
			origfilename := filename + fileExt
			f, err := os.Open(filepath.Join("testdata", "base", origfilename))
			if err != nil {
				t.Fatalf("failed to open test file %s: %v", origfilename, err)
			}

			semiBytes, _ := io.ReadAll(f)

			f.Close()

			noSemiBytes := bytes.ReplaceAll(semiBytes, []byte{';'}, []byte{})

			tokens := []token{}
			tr := newTokenReader(bytes.NewBuffer(semiBytes))
			tr.optionalSemicolons = false
			for tr.Next() {
				tokens = append(tokens, tr.Token())
			}
			if tr.Err() != nil {
				t.Fatalf("token reader errored: %v", tr.Err())
			}

			tokens2 := []token{}
			tr2 := newTokenReader(bytes.NewBuffer(noSemiBytes))
			tr2.optionalSemicolons = true
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
					t.Fatalf("tokens at pos %d differed: concrete %v vs %v", i, string(tk.concrete), string(tk2.concrete))
				}
			}
		})
	}
}

func TestTokenizeError(t *testing.T) {
	t.Parallel()
	type testCase struct {
		file string
		err  string
	}
	tcs := []testCase{
		{file: "invalid_token_arrow_eof", err: "[0:1] unexpected EOF: waiting for [>, i, number] after '-'"},
		{file: "invalid_token_arrow", err: "[0:2] unexpected token '-' waiting for [>, i, number] after '-'"},
		{file: "invalid_token_comment_eof", err: "[0:1] unexpected EOF: waiting for [*, /] after '/'"},
		{file: "invalid_token_comment", err: `[0:2] unexpected token 'p' waiting for [*, /] after '/'
[0:2] unexpected EOF: block comment missing end token`},
		{file: "invalid_token_string", err: "[0:13] unexpected EOF: waiting for string end quote"},
		{file: "invalid_token_unknown", err: "[0:1] unexpected token '*', expected number, letter, or control sequence"},
		{file: "invalid_token_block_comment_eof", err: "[0:28] unexpected EOF: block comment missing end token"},
		{file: "invalid_token_float_literal", err: "[0:4] unexpected token ' ', expected number following \"-1.\""},
		{file: "invalid_token_float_literal_eof", err: "[0:3] unexpected EOF: expected number following \"-1.\""},
		{file: "invalid_token_ninf_1", err: "[0:3] unexpected EOF: waiting for [f] after '-in'"},
		{file: "invalid_token_ninf_2", err: "[0:2] unexpected EOF: waiting for [n] after '-i'"},
		{file: "invalid_token_ninf_3", err: `[0:3] unexpected token 'i' waiting for [n] after '-i'
[0:4] unexpected token 'i' waiting for [f] after '-in'`},
		{file: "invalid_token_float_literal_two_periods", err: "[0:5] unexpected second period in float following \"-1.0\""},
		{file: "invalid_token_multi_err_float", err: `[0:3] unexpected token ' ', expected number following "1."
[0:6] unexpected token ' ', expected number following "1."
[0:9] unexpected token ' ', expected number following "1."`},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.file, func(t *testing.T) {
			t.Parallel()
			origfilename := tc.file + fileExt
			f, err := os.Open(filepath.Join("testdata", "invalid", origfilename))
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
