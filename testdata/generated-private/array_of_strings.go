// Code generated by bebopc-go; DO NOT EDIT.

package generated

import (
	"github.com/200sc/bebop"
	"github.com/200sc/bebop/iohelp"
	"io"
)

var _ bebop.Record = &arrayOfStrings{}

type arrayOfStrings struct {
	strings []string
}

func (bbp arrayOfStrings) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.strings)))
	at += 4
	for _, v1 := range bbp.strings {
		iohelp.WriteUint32Bytes(buf[at:], uint32(len(v1)))
		copy(buf[at+4:at+4+len(v1)], []byte(v1))
		at += 4 + len(v1)
	}
	return at
}

func (bbp *arrayOfStrings) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	if len(buf[at:]) < 4 {
		return io.ErrUnexpectedEOF
	}
	bbp.strings = make([]string, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	for i1 := range bbp.strings {
		(bbp.strings)[i1], err = iohelp.ReadStringBytes(buf[at:])
		if err != nil {
			return err
		}
		at += 4 + len((bbp.strings)[i1])
	}
	return nil
}

func (bbp *arrayOfStrings) MustUnmarshalBebop(buf []byte) {
	at := 0
	bbp.strings = make([]string, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	for i1 := range bbp.strings {
		(bbp.strings)[i1] = iohelp.MustReadStringBytes(buf[at:])
		at += 4 + len((bbp.strings)[i1])
	}
}

func (bbp arrayOfStrings) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(len(bbp.strings)))
	for _, elem := range bbp.strings {
		iohelp.WriteUint32(w, uint32(len(elem)))
		w.Write([]byte(elem))
	}
	return w.Err
}

func (bbp *arrayOfStrings) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bbp.strings = make([]string, iohelp.ReadUint32(r))
	for i1 := range bbp.strings {
		(bbp.strings[i1]) = iohelp.ReadString(r)
	}
	return r.Err
}

func (bbp arrayOfStrings) Size() int {
	bodyLen := 0
	bodyLen += 4
	for _, elem := range bbp.strings {
		bodyLen += 4 + len(elem)
	}
	return bodyLen
}

func (bbp arrayOfStrings) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func makearrayOfStrings(r *iohelp.ErrorReader) (arrayOfStrings, error) {
	v := arrayOfStrings{}
	err := v.DecodeBebop(r)
	return v, err
}

func makearrayOfStringsFromBytes(buf []byte) (arrayOfStrings, error) {
	v := arrayOfStrings{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func mustMakearrayOfStringsFromBytes(buf []byte) arrayOfStrings {
	v := arrayOfStrings{}
	v.MustUnmarshalBebop(buf)
	return v
}

