package generated_test

import (
	"testing"

	"github.com/200sc/bebop/testdata/generated"
)

func TestListUnpopulated(t *testing.T) {
	l := generated.List{}
	err := l.UnmarshalBebop([]byte{0, 0, 0, 0})
	if err == nil {
		t.Fatalf("unmarshal of empty union should fail")
	}
}
