package bebop

import (
	"bytes"
	"io"
)

// A Record can be serialized to and from a bebop structure.
type Record interface {
	EncodeBebop(io.Writer) error
	DecodeBebop(io.Reader) error
}

// Marshal writes a record out into in memory bytes, in bebop format.
func Marshal(r Record) ([]byte, error) {
	var buf bytes.Buffer
	err := r.EncodeBebop(&buf)
	return buf.Bytes(), err
}

// Unmarshal populates a record with an in memory bebop-formatted byte slice.
func Unmarshal(data []byte, r Record) error {
	buf := bytes.NewReader(data)
	return r.DecodeBebop(buf)
}
