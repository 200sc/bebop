package bebop

import (
	"bytes"
	"io"
)

type Record interface {
	EncodeBebop(io.Writer) error
	DecodeBebop(io.Reader) error
}

func Marshal(r Record) ([]byte, error) {
	var buf bytes.Buffer
	err := r.EncodeBebop(&buf)
	return buf.Bytes(), err
}

func Unmarshal(data []byte, r Record) error {
	buf := bytes.NewReader(data)
	return r.DecodeBebop(buf)
}
