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

type ImportedEnum uint32

const (
	ImportedEnum_One ImportedEnum = 1
)

var _ bebop.Record = &ImportedType{}

type ImportedType struct {
	Foobar string
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
	bbp.Foobar = iohelp.MustReadStringBytes(buf[at:])
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

func (bbp ImportedType) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
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

var _ bebop.Record = &ImportedMessage{}

type ImportedMessage struct {
	Foo *ImportedEnum
	Bar *ImportedType
	Unin *ImportedUnion
}

func (bbp ImportedMessage) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(bbp.Size()-4))
	at += 4
	if bbp.Foo != nil {
		buf[at] = 1
		at++
		iohelp.WriteUint32Bytes(buf[at:], uint32(*bbp.Foo))
		at += 4
	}
	if bbp.Bar != nil {
		buf[at] = 2
		at++
		(*bbp.Bar).MarshalBebopTo(buf[at:])
		at += (*bbp.Bar).Size()
	}
	if bbp.Unin != nil {
		buf[at] = 3
		at++
		(*bbp.Unin).MarshalBebopTo(buf[at:])
		at += (*bbp.Unin).Size()
	}
	return at
}

func (bbp *ImportedMessage) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.Foo = new(ImportedEnum)
			(*bbp.Foo) = ImportedEnum(iohelp.ReadUint32Bytes(buf[at:]))
			at += 4
		case 2:
			at += 1
			bbp.Bar = new(ImportedType)
			(*bbp.Bar) = MustMakeImportedTypeFromBytes(buf[at:])
			if err != nil{
				return err
			}
			at += ((*bbp.Bar)).Size()
		case 3:
			at += 1
			bbp.Unin = new(ImportedUnion)
			(*bbp.Unin) = MustMakeImportedUnionFromBytes(buf[at:])
			if err != nil{
				return err
			}
			at += ((*bbp.Unin)).Size()
		default:
			return nil
		}
	}
}

func (bbp *ImportedMessage) MustUnmarshalBebop(buf []byte) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.Foo = new(ImportedEnum)
			(*bbp.Foo) = ImportedEnum(iohelp.ReadUint32Bytes(buf[at:]))
			at += 4
		case 2:
			at += 1
			bbp.Bar = new(ImportedType)
			(*bbp.Bar) = MustMakeImportedTypeFromBytes(buf[at:])
			at += ((*bbp.Bar)).Size()
		case 3:
			at += 1
			bbp.Unin = new(ImportedUnion)
			(*bbp.Unin) = MustMakeImportedUnionFromBytes(buf[at:])
			at += ((*bbp.Unin)).Size()
		default:
			return
		}
	}
}

func (bbp ImportedMessage) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(bbp.Size()-4))
	if bbp.Foo != nil {
		w.Write([]byte{1})
		iohelp.WriteUint32(w, uint32(*bbp.Foo))
	}
	if bbp.Bar != nil {
		w.Write([]byte{2})
		err = (*bbp.Bar).EncodeBebop(w)
		if err != nil{
			return err
		}
	}
	if bbp.Unin != nil {
		w.Write([]byte{3})
		err = (*bbp.Unin).EncodeBebop(w)
		if err != nil{
			return err
		}
	}
	w.Write([]byte{0})
	return w.Err
}

func (bbp *ImportedMessage) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bodyLen := iohelp.ReadUint32(r)
	r.Reader = &io.LimitedReader{R:r.Reader, N:int64(bodyLen)}
	for {
		switch iohelp.ReadByte(r) {
		case 1:
			bbp.Foo = new(ImportedEnum)
			*bbp.Foo = ImportedEnum(iohelp.ReadUint32(r))
		case 2:
			bbp.Bar = new(ImportedType)
			(*bbp.Bar), err = MakeImportedType(r)
			if err != nil{
				return err
			}
		case 3:
			bbp.Unin = new(ImportedUnion)
			(*bbp.Unin), err = MakeImportedUnion(r)
			if err != nil{
				return err
			}
		default:
			io.ReadAll(r)
			return r.Err
		}
	}
}

func (bbp ImportedMessage) Size() int {
	bodyLen := 5
	if bbp.Foo != nil {
		bodyLen += 1
		bodyLen += 4
	}
	if bbp.Bar != nil {
		bodyLen += 1
		bodyLen += (*bbp.Bar).Size()
	}
	if bbp.Unin != nil {
		bodyLen += 1
		bodyLen += (*bbp.Unin).Size()
	}
	return bodyLen
}

func (bbp ImportedMessage) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func MakeImportedMessage(r iohelp.ErrorReader) (ImportedMessage, error) {
	v := ImportedMessage{}
	err := v.DecodeBebop(r)
	return v, err
}

func MakeImportedMessageFromBytes(buf []byte) (ImportedMessage, error) {
	v := ImportedMessage{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func MustMakeImportedMessageFromBytes(buf []byte) ImportedMessage {
	v := ImportedMessage{}
	v.MustUnmarshalBebop(buf)
	return v
}

var _ bebop.Record = &WhyAreTheseInline{}

type WhyAreTheseInline struct {
}

func (bbp WhyAreTheseInline) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(bbp.Size()-4))
	at += 4
	return at
}

func (bbp *WhyAreTheseInline) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		default:
			return nil
		}
	}
}

func (bbp *WhyAreTheseInline) MustUnmarshalBebop(buf []byte) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		default:
			return
		}
	}
}

func (bbp WhyAreTheseInline) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(bbp.Size()-4))
	w.Write([]byte{0})
	return w.Err
}

func (bbp *WhyAreTheseInline) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bodyLen := iohelp.ReadUint32(r)
	r.Reader = &io.LimitedReader{R:r.Reader, N:int64(bodyLen)}
	for {
		switch iohelp.ReadByte(r) {
		default:
			io.ReadAll(r)
			return r.Err
		}
	}
}

func (bbp WhyAreTheseInline) Size() int {
	bodyLen := 5
	return bodyLen
}

func (bbp WhyAreTheseInline) MarshalBebop() []byte {
	return []byte{}
}

func MakeWhyAreTheseInline(r iohelp.ErrorReader) (WhyAreTheseInline, error) {
	return WhyAreTheseInline{}, nil
}

func MakeWhyAreTheseInlineFromBytes(buf []byte) (WhyAreTheseInline, error) {
	return WhyAreTheseInline{}, nil
}

func MustMakeWhyAreTheseInlineFromBytes(buf []byte) WhyAreTheseInline {
	return WhyAreTheseInline{}
}

var _ bebop.Record = &Really{}

type Really struct {
}

func (bbp Really) MarshalBebopTo(buf []byte) int {
	return 0
}

func (bbp *Really) UnmarshalBebop(buf []byte) (err error) {
	return nil
}

func (bbp *Really) MustUnmarshalBebop(buf []byte) {
}

func (bbp Really) EncodeBebop(iow io.Writer) (err error) {
	return nil
}

func (bbp *Really) DecodeBebop(ior io.Reader) (err error) {
	return nil
}

func (bbp Really) Size() int {
	return 0
}

func (bbp Really) MarshalBebop() []byte {
	return []byte{}
}

func MakeReally(r iohelp.ErrorReader) (Really, error) {
	return Really{}, nil
}

func MakeReallyFromBytes(buf []byte) (Really, error) {
	return Really{}, nil
}

func MustMakeReallyFromBytes(buf []byte) Really {
	return Really{}
}

var _ bebop.Record = &ImportedUnion{}

type ImportedUnion struct {
	WhyAreTheseInline *WhyAreTheseInline
	Really *Really
}

func (bbp ImportedUnion) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(bbp.Size()-5))
	at += 4
	if bbp.WhyAreTheseInline != nil {
		buf[at] = 1
		at++
		(*bbp.WhyAreTheseInline).MarshalBebopTo(buf[at:])
		at += (*bbp.WhyAreTheseInline).Size()
		return at
	}
	if bbp.Really != nil {
		buf[at] = 2
		at++
		(*bbp.Really).MarshalBebopTo(buf[at:])
		at += (*bbp.Really).Size()
		return at
	}
	return at
}

func (bbp *ImportedUnion) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	if len(buf) == 0 {
		return iohelp.ErrUnpopulatedUnion
	}
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.WhyAreTheseInline = new(WhyAreTheseInline)
			(*bbp.WhyAreTheseInline) = MustMakeWhyAreTheseInlineFromBytes(buf[at:])
			if err != nil{
				return err
			}
			at += ((*bbp.WhyAreTheseInline)).Size()
			return nil
		case 2:
			at += 1
			bbp.Really = new(Really)
			(*bbp.Really) = MustMakeReallyFromBytes(buf[at:])
			if err != nil{
				return err
			}
			at += ((*bbp.Really)).Size()
			return nil
		default:
			return nil
		}
	}
}

func (bbp *ImportedUnion) MustUnmarshalBebop(buf []byte) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.WhyAreTheseInline = new(WhyAreTheseInline)
			(*bbp.WhyAreTheseInline) = MustMakeWhyAreTheseInlineFromBytes(buf[at:])
			at += ((*bbp.WhyAreTheseInline)).Size()
			return
		case 2:
			at += 1
			bbp.Really = new(Really)
			(*bbp.Really) = MustMakeReallyFromBytes(buf[at:])
			at += ((*bbp.Really)).Size()
			return
		default:
			return
		}
	}
}

func (bbp ImportedUnion) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(bbp.Size()-5))
	if bbp.WhyAreTheseInline != nil {
		w.Write([]byte{1})
		err = (*bbp.WhyAreTheseInline).EncodeBebop(w)
		if err != nil{
			return err
		}
		return w.Err
	}
	if bbp.Really != nil {
		w.Write([]byte{2})
		err = (*bbp.Really).EncodeBebop(w)
		if err != nil{
			return err
		}
		return w.Err
	}
	return w.Err
}

func (bbp *ImportedUnion) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bodyLen := iohelp.ReadUint32(r)
	r.Reader = &io.LimitedReader{R:r.Reader, N:int64(bodyLen)+1}
	for {
		switch iohelp.ReadByte(r) {
		case 1:
			bbp.WhyAreTheseInline = new(WhyAreTheseInline)
			(*bbp.WhyAreTheseInline), err = MakeWhyAreTheseInline(r)
			if err != nil{
				return err
			}
			io.ReadAll(r)
			return r.Err
		case 2:
			bbp.Really = new(Really)
			(*bbp.Really), err = MakeReally(r)
			if err != nil{
				return err
			}
			io.ReadAll(r)
			return r.Err
		default:
			io.ReadAll(r)
			return r.Err
		}
	}
}

func (bbp ImportedUnion) Size() int {
	bodyLen := 4
	if bbp.WhyAreTheseInline != nil {
		bodyLen += 1
		bodyLen += (*bbp.WhyAreTheseInline).Size()
		return bodyLen
	}
	if bbp.Really != nil {
		bodyLen += 1
		bodyLen += (*bbp.Really).Size()
		return bodyLen
	}
	return bodyLen
}

func (bbp ImportedUnion) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func MakeImportedUnion(r iohelp.ErrorReader) (ImportedUnion, error) {
	v := ImportedUnion{}
	err := v.DecodeBebop(r)
	return v, err
}

func MakeImportedUnionFromBytes(buf []byte) (ImportedUnion, error) {
	v := ImportedUnion{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func MustMakeImportedUnionFromBytes(buf []byte) ImportedUnion {
	v := ImportedUnion{}
	v.MustUnmarshalBebop(buf)
	return v
}

