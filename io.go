package bebop

import (
	"io"
)

// A Record can be serialized to and from a bebop structure.
type Record interface {
	// MarshalBebop converts a bebop record to wire format. It is recommended over
	// EncodeBebop for performance.
	MarshalBebop() []byte
	// MarshalBebopTo writes a bebop record to an existing byte slice. It is primarily
	// used internally, and performs no checks to ensure the given byte slice is large
	// enough to contain the record.
	MarshalBebopTo([]byte)
	// EncodeBebop writes a bebop record in wire format to a writer. It is slower (~6x)
	// than MarshalBebop, and is only recommended for uses where the record size is both
	// larger than a network packet and able to be acted upon as writer receives the byte
	// stream, not only after the entire message has been received.
	EncodeBebop(io.Writer) error
	DecodeBebop(io.Reader) error
}
