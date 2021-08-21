package bebop

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	skipIfUpstreamMissing(t)

	cmd := exec.Command(upsteamCompilerName, "--ts", "./out.ts", "--dir", filepath.Join(".", "testdata", "base"))
	printed, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(printed))
		t.Fatalf("%s failed: %v", upsteamCompilerName, err)
	}
}

func TestUpstreamCompatiblityFailures(t *testing.T) {
	skipIfUpstreamMissing(t)

	files, err := os.ReadDir(filepath.Join(".", "testdata", "invalid"))
	if err != nil {
		t.Fatalf("failed to list invalid files: %v", err)
	}

	var exceptions = map[string]string{
		"invalid_nested_union.bop": "merged into upstream but unreleased",
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		filename := f.Name()
		t.Run(filename, func(t *testing.T) {
			if reason := exceptions[filename]; reason != "" {
				t.Skip(reason)
			}
			cmd := exec.Command(upsteamCompilerName, "--ts", "./out.ts", "--files", filepath.Join(".", "testdata", "invalid", filename))
			_, err := cmd.CombinedOutput()
			if err == nil {
				t.Fatalf("%s should have errored", upsteamCompilerName)
			}
			//fmt.Println(string(printed))
		})
	}
}