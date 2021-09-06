// Code generated by bebopc-go; DO NOT EDIT.

package generatedtwo

import (
	"io"
	"github.com/200sc/bebop"
	"github.com/200sc/bebop/iohelp"
)

const (
	Go_package = "github.com/200sc/bebop/testdata/incompatible/generatedtwo"
)

var _ bebop.Record = &ImportedType{}

type ImportedType struct {
	Foobar string
}

func (bbp ImportedType) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func (bbp ImportedType) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.Foobar)))
	copy(buf[at+4:at+4+len(bbp.Foobar)], []byte(bbp.Foobar))
	at += 4 + len(bbp.Foobar)
	return at
}

func (bbp *ImportedType) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	bbp.Foobar, err = iohelp.ReadStringBytes(buf[at:])
	if err != nil{
		return err
	}
	at += 4 + len(bbp.Foobar)
	return nil
}

func (bbp *ImportedType) MustUnmarshalBebop(buf []byte) {
	at := 0
	bbp.Foobar =  iohelp.MustReadStringBytes(buf[at:])
	at += 4 + len(bbp.Foobar)
}

func (bbp ImportedType) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(len(bbp.Foobar)))
	w.Write([]byte(bbp.Foobar))
	return w.Err
}

func (bbp *ImportedType) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bbp.Foobar = iohelp.ReadString(r)
	return r.Err
}

func (bbp ImportedType) Size() int {
	bodyLen := 0
	bodyLen += 4 + len(bbp.Foobar)
	return bodyLen
}

func MakeImportedType(r iohelp.ErrorReader) (ImportedType, error) {
	v := ImportedType{}
	err := v.DecodeBebop(r)
	return v, err
}

func MakeImportedTypeFromBytes(buf []byte) (ImportedType, error) {
	v := ImportedType{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func MustMakeImportedTypeFromBytes(buf []byte) ImportedType {
	v := ImportedType{}
	v.MustUnmarshalBebop(buf)
	return v
}

