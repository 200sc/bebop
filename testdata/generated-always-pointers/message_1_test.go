package generated_test

import (
	"testing"

	generated "github.com/200sc/bebop/testdata/generated-always-pointers"
)

// c.f. https://github.com/RainwayApp/bebop/wiki/Wire-format#messages
func TestMessage1ExampleEncoding(t *testing.T) {
	x := uint8(15)
	z := int32(5)
	msg := generated.ExampleMessage{
		X: &x,
		Z: &z,
	}
	out := msg.MarshalBebop()
	expected := []byte{8, 0, 0, 0, 1, 0x0f, 3, 5, 0, 0, 0, 0}
	if len(out) != len(expected) {
		t.Fatalf("length mismatch: %v vs %v", len(out), len(expected))
	}
	for i, v := range out {
		if expected[i] != v {
			t.Fatalf("index %d mismatch: %v vs %v", i, v, expected[i])
		}
	}
}
