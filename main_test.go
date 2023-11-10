package bebop

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestMain(m *testing.M) {
	code := m.Run()
	if code != 0 {
		os.Exit(code)
	}
	fmt.Println("Bebop tests passed; running all generated code tests")

	dirs := []string{"./testdata/generated", "./testdata/generated-always-pointers", "./internal/importgraph", "./testdata/incompatible/...", "./testdata/generated-private"}

	for _, dir := range dirs {
		var cmd *exec.Cmd
		if testing.Verbose() {
			cmd = exec.Command("go", "test", "-v", dir)
		} else {
			cmd = exec.Command("go", "test", dir)
		}
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			os.Exit(1)
		}
	}
	os.Exit(0)
}
