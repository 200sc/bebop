package bebop

import "testing"

func TestTokenKindStrings(t *testing.T) {
	// Verify every token kind has a string encoding
	for tok := tokenKindInvalid; tok < tokenKindFinal; tok++ {
		s := tok.String()
		if s == "" {
			t.Errorf("token %v has no string encoding", int(tok))
		}
	}
}
