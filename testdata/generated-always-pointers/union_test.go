package generated_test

import (
	"testing"

	generated "github.com/200sc/bebop/testdata/generated-always-pointers"
)

func TestListUnpopulated(t *testing.T) {
	l := generated.List{}
	err := l.UnmarshalBebop([]byte{0, 0, 0, 0})
	if err == nil {
		t.Fatalf("unmarshal of empty union should fail")
	}
}
