package bebop

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const upsteamCompilerName = "bebopc"

func skipIfUpstreamMissing(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath(upsteamCompilerName); err != nil {
		t.Skipf("missing upstream %s compiler", upsteamCompilerName)
	}
}

func TestUpstreamCompatiblitySuccess(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("upstream tests skipped by --short")
	}
	skipIfUpstreamMissing(t)

	outfile := "./compsuccess-out.ts"
	defer os.Remove(outfile)
	cmd := exec.Command(upsteamCompilerName, "--ts", outfile, "--dir", filepath.Join(".", "testdata", "base"))
	printed, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(printed))
		t.Fatalf("%s failed: %v", upsteamCompilerName, err)
	}
}

func TestUpstreamCompatiblityFailures(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("upstream tests skipped by --short")
	}
	skipIfUpstreamMissing(t)

	files, err := os.ReadDir(filepath.Join(".", "testdata", "invalid"))
	if err != nil {
		t.Fatalf("failed to list invalid files: %v", err)
	}

	var exceptions = map[string]string{
		"invalid_readonly_comment.bop": "bebopc 2.2.4 errors where 2.3.0 does not, without a changelog note",
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		filename := f.Name()
		t.Run(filename, func(t *testing.T) {
			t.Parallel()
			if reason := exceptions[filename]; reason != "" {
				t.Skip(reason)
			}
			outfile := "./compfail-" + filename + "-out.ts"
			defer os.Remove(outfile)
			cmd := exec.Command(upsteamCompilerName, "--ts", outfile, "--files", filepath.Join(".", "testdata", "invalid", filename))
			err := cmd.Run()
			if err == nil {
				t.Fatalf("%s should have errored", upsteamCompilerName)
			}
		})
	}
}

func TestIncompatibilityExpectations_200sc(t *testing.T) {
	t.Parallel()
	files, err := os.ReadDir(filepath.Join(".", "testdata", "incompatible"))
	if err != nil {
		t.Fatalf("failed to list incompatible files: %v", err)
	}

	failures := map[string]struct{}{
		"import_loop_a.bop":             {},
		"import_loop_b.bop":             {},
		"import_loop_c.bop":             {},
		"invalid_enum_primitive.bop":    {},
		"invalid_import_no_const.bop":   {},
		"invalid_message_primitive.bop": {},
		"invalid_struct_primitive.bop":  {},
		"recursive_struct.bop":          {},
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		filename := f.Name()
		if !strings.HasSuffix(filename, ".bop") {
			continue
		}
		t.Run(filename, func(t *testing.T) {
			t.Parallel()
			f, err := os.Open(filepath.Join("testdata", "incompatible", filename))
			if err != nil {
				t.Fatalf("failed to open test file %s: %v", filename, err)
			}
			defer f.Close()
			bopf, _, err := ReadFile(f)
			if err != nil {
				t.Fatalf("failed to read file %s: %v", filename, err)
			}
			err = bopf.Generate(bytes.NewBuffer([]byte{}), GenerateSettings{
				PackageName: "generated",
			})
			_, shouldFail := failures[filename]
			if shouldFail && err == nil {
				t.Fatal("expected generation failure")
			}
			if !shouldFail && err != nil {
				t.Fatalf("expected generation success: %v", err)
			}
		})
	}
}

func TestIncompatibilityExpectations_Rainway(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("upstream tests skipped by --short")
	}
	skipIfUpstreamMissing(t)

	files, err := os.ReadDir(filepath.Join(".", "testdata", "incompatible"))
	if err != nil {
		t.Fatalf("failed to list incompatible files: %v", err)
	}

	failures := map[string]struct{}{
		"import_loop_a.bop":           {},
		"import_loop_b.bop":           {},
		"import_loop_c.bop":           {},
		"import_separate_a.bop":       {},
		"invalid_import_no_const.bop": {},
		"quoted_string.bop":           {},
		"union.bop":                   {},
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		filename := f.Name()
		if !strings.HasSuffix(filename, ".bop") {
			continue
		}
		t.Run(filename, func(t *testing.T) {
			t.Parallel()
			outfile := "./compexpect-" + filename + "-out.ts"
			defer os.Remove(outfile)
			cmd := exec.Command(upsteamCompilerName, "--ts", outfile, "--files", filepath.Join(".", "testdata", "incompatible", filename))
			out, err := cmd.CombinedOutput()
			bytes.TrimSuffix(out, []byte("\n"))

			_, shouldFail := failures[filename]
			if shouldFail && err == nil {
				t.Fatal("expected generation failure")
			}
			if !shouldFail && err != nil {
				t.Fatalf("expected generation success: %v", err)
			}

			fmt.Printf("filename: %v err: %v out:%v\n", filename, err, string(out))
		})
	}
}
