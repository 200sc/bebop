// Code generated by bebopc-go; DO NOT EDIT.

package generated

import (
	"github.com/200sc/bebop"
	"github.com/200sc/bebop/iohelp"
	"io"
)

var _ bebop.Record = &QuotedString{}

type QuotedString struct {
	// Deprecated: "deprecated"
	X int32
	// Deprecated: escaped slash: \
	Y int32
	// Deprecated: escaped" "slashes:\\"" \
	Z int32
}

func (bbp QuotedString) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteInt32Bytes(buf[at:], bbp.X)
	at += 4
	iohelp.WriteInt32Bytes(buf[at:], bbp.Y)
	at += 4
	iohelp.WriteInt32Bytes(buf[at:], bbp.Z)
	at += 4
	return at
}

func (bbp *QuotedString) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	if len(buf[at:]) < 4 {
		return io.ErrUnexpectedEOF
	}
	bbp.X = iohelp.ReadInt32Bytes(buf[at:])
	at += 4
	if len(buf[at:]) < 4 {
		return io.ErrUnexpectedEOF
	}
	bbp.Y = iohelp.ReadInt32Bytes(buf[at:])
	at += 4
	if len(buf[at:]) < 4 {
		return io.ErrUnexpectedEOF
	}
	bbp.Z = iohelp.ReadInt32Bytes(buf[at:])
	at += 4
	return nil
}

func (bbp *QuotedString) MustUnmarshalBebop(buf []byte) {
	at := 0
	bbp.X = iohelp.ReadInt32Bytes(buf[at:])
	at += 4
	bbp.Y = iohelp.ReadInt32Bytes(buf[at:])
	at += 4
	bbp.Z = iohelp.ReadInt32Bytes(buf[at:])
	at += 4
}

func (bbp QuotedString) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteInt32(w, bbp.X)
	iohelp.WriteInt32(w, bbp.Y)
	iohelp.WriteInt32(w, bbp.Z)
	return w.Err
}

func (bbp *QuotedString) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bbp.X = iohelp.ReadInt32(r)
	bbp.Y = iohelp.ReadInt32(r)
	bbp.Z = iohelp.ReadInt32(r)
	return r.Err
}

func (bbp QuotedString) Size() int {
	bodyLen := 0
	bodyLen += 4
	bodyLen += 4
	bodyLen += 4
	return bodyLen
}

func (bbp QuotedString) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func MakeQuotedString(r *iohelp.ErrorReader) (QuotedString, error) {
	v := QuotedString{}
	err := v.DecodeBebop(r)
	return v, err
}

func MakeQuotedStringFromBytes(buf []byte) (QuotedString, error) {
	v := QuotedString{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func MustMakeQuotedStringFromBytes(buf []byte) QuotedString {
	v := QuotedString{}
	v.MustUnmarshalBebop(buf)
	return v
}

