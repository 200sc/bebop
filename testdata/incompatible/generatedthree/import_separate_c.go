// Code generated by bebopc-go; DO NOT EDIT.

package generatedthree

import (
	"github.com/200sc/bebop"
	"github.com/200sc/bebop/iohelp"
	"io"
)

const (
	Go_package = "github.com/200sc/bebop/testdata/incompatible/generatedthree"
)

var _ bebop.Record = &NotImported{}

type NotImported struct {
}

func (bbp NotImported) MarshalBebopTo(buf []byte) int {
	return 0
}

func (bbp *NotImported) UnmarshalBebop(buf []byte) (err error) {
	return nil
}

func (bbp *NotImported) MustUnmarshalBebop(buf []byte) {
}

func (bbp NotImported) EncodeBebop(iow io.Writer) (err error) {
	return nil
}

func (bbp *NotImported) DecodeBebop(ior io.Reader) (err error) {
	return nil
}

func (bbp NotImported) Size() int {
	return 0
}

func (bbp NotImported) MarshalBebop() []byte {
	return []byte{}
}

func MakeNotImported(r *iohelp.ErrorReader) (NotImported, error) {
	return NotImported{}, nil
}

func MakeNotImportedFromBytes(buf []byte) (NotImported, error) {
	return NotImported{}, nil
}

func MustMakeNotImportedFromBytes(buf []byte) NotImported {
	return NotImported{}
}

