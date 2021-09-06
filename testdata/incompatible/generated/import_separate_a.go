// Code generated by bebopc-go; DO NOT EDIT.

package generated

import (
	"github.com/200sc/bebop/testdata/incompatible/generatedtwo"
	"io"
	"github.com/200sc/bebop"
	"github.com/200sc/bebop/iohelp"
)

const (
	Go_package = "github.com/200sc/bebop/testdata/incompatible/generated"
)

var _ bebop.Record = &UsesImport{}

type UsesImport struct {
	Imported generatedtwo.ImportedType
}

func (bbp UsesImport) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func (bbp UsesImport) MarshalBebopTo(buf []byte) int {
	at := 0
	(bbp.Imported).MarshalBebopTo(buf[at:])
	at += (bbp.Imported).Size()
	return at
}

func (bbp *UsesImport) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	bbp.Imported, err = generatedtwo.MakeImportedTypeFromBytes(buf[at:])
	if err != nil{
		return err
	}
	at += (bbp.Imported).Size()
	return nil
}

func (bbp *UsesImport) MustUnmarshalBebop(buf []byte) {
	at := 0
	bbp.Imported = generatedtwo.MustMakeImportedTypeFromBytes(buf[at:])
	at += (bbp.Imported).Size()
}

func (bbp UsesImport) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	err = (bbp.Imported).EncodeBebop(w)
	if err != nil{
		return err
	}
	return w.Err
}

func (bbp *UsesImport) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	(bbp.Imported), err = generatedtwo.MakeImportedType(r)
	if err != nil{
		return err
	}
	return r.Err
}

func (bbp UsesImport) Size() int {
	bodyLen := 0
	bodyLen += (bbp.Imported).Size()
	return bodyLen
}

func MakeUsesImport(r iohelp.ErrorReader) (UsesImport, error) {
	v := UsesImport{}
	err := v.DecodeBebop(r)
	return v, err
}

func MakeUsesImportFromBytes(buf []byte) (UsesImport, error) {
	v := UsesImport{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func MustMakeUsesImportFromBytes(buf []byte) UsesImport {
	v := UsesImport{}
	v.MustUnmarshalBebop(buf)
	return v
}
