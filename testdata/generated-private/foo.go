// Code generated by bebopc-go; DO NOT EDIT.

package generated

import (
	"github.com/200sc/bebop"
	"github.com/200sc/bebop/iohelp"
	"io"
)

var _ bebop.Record = &foo{}

type foo struct {
	bar bar
}

func (bbp foo) MarshalBebopTo(buf []byte) int {
	at := 0
	(bbp.bar).MarshalBebopTo(buf[at:])
	tmp7570 := (bbp.bar); at += tmp7570.Size()
	return at
}

func (bbp *foo) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	bbp.bar, err = makebarFromBytes(buf[at:])
	if err != nil {
		return err
	}
	tmp7574 := (bbp.bar); at += tmp7574.Size()
	return nil
}

func (bbp *foo) MustUnmarshalBebop(buf []byte) {
	at := 0
	bbp.bar = mustMakebarFromBytes(buf[at:])
	tmp7581 := (bbp.bar); at += tmp7581.Size()
}

func (bbp foo) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	err = (bbp.bar).EncodeBebop(w)
	if err != nil {
		return err
	}
	return w.Err
}

func (bbp *foo) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	(bbp.bar), err = makebar(r)
	if err != nil {
		return err
	}
	return r.Err
}

func (bbp foo) Size() int {
	bodyLen := 0
	tmp7596 := (bbp.bar); bodyLen += tmp7596.Size()
	return bodyLen
}

func (bbp foo) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func makefoo(r iohelp.ErrorReader) (foo, error) {
	v := foo{}
	err := v.DecodeBebop(r)
	return v, err
}

func makefooFromBytes(buf []byte) (foo, error) {
	v := foo{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func mustMakefooFromBytes(buf []byte) foo {
	v := foo{}
	v.MustUnmarshalBebop(buf)
	return v
}

var _ bebop.Record = &bar{}

type bar struct {
	x *float64
	y *float64
	z *float64
}

func (bbp bar) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(bbp.Size()-4))
	at += 4
	if bbp.x != nil {
		buf[at] = 1
		at++
		iohelp.WriteFloat64Bytes(buf[at:], *bbp.x)
		at += 8
	}
	if bbp.y != nil {
		buf[at] = 2
		at++
		iohelp.WriteFloat64Bytes(buf[at:], *bbp.y)
		at += 8
	}
	if bbp.z != nil {
		buf[at] = 3
		at++
		iohelp.WriteFloat64Bytes(buf[at:], *bbp.z)
		at += 8
	}
	return at
}

func (bbp *bar) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.x = new(float64)
			if len(buf[at:]) < 8 {
				return io.ErrUnexpectedEOF
			}
			(*bbp.x) = iohelp.ReadFloat64Bytes(buf[at:])
			at += 8
		case 2:
			at += 1
			bbp.y = new(float64)
			if len(buf[at:]) < 8 {
				return io.ErrUnexpectedEOF
			}
			(*bbp.y) = iohelp.ReadFloat64Bytes(buf[at:])
			at += 8
		case 3:
			at += 1
			bbp.z = new(float64)
			if len(buf[at:]) < 8 {
				return io.ErrUnexpectedEOF
			}
			(*bbp.z) = iohelp.ReadFloat64Bytes(buf[at:])
			at += 8
		default:
			return nil
		}
	}
}

func (bbp *bar) MustUnmarshalBebop(buf []byte) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.x = new(float64)
			(*bbp.x) = iohelp.ReadFloat64Bytes(buf[at:])
			at += 8
		case 2:
			at += 1
			bbp.y = new(float64)
			(*bbp.y) = iohelp.ReadFloat64Bytes(buf[at:])
			at += 8
		case 3:
			at += 1
			bbp.z = new(float64)
			(*bbp.z) = iohelp.ReadFloat64Bytes(buf[at:])
			at += 8
		default:
			return
		}
	}
}

func (bbp bar) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(bbp.Size()-4))
	if bbp.x != nil {
		w.Write([]byte{1})
		iohelp.WriteFloat64(w, *bbp.x)
	}
	if bbp.y != nil {
		w.Write([]byte{2})
		iohelp.WriteFloat64(w, *bbp.y)
	}
	if bbp.z != nil {
		w.Write([]byte{3})
		iohelp.WriteFloat64(w, *bbp.z)
	}
	w.Write([]byte{0})
	return w.Err
}

func (bbp *bar) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bodyLen := iohelp.ReadUint32(r)
	r.Reader = &io.LimitedReader{R:r.Reader, N:int64(bodyLen)}
	for {
		switch iohelp.ReadByte(r) {
		case 1:
			bbp.x = new(float64)
			*bbp.x = iohelp.ReadFloat64(r)
		case 2:
			bbp.y = new(float64)
			*bbp.y = iohelp.ReadFloat64(r)
		case 3:
			bbp.z = new(float64)
			*bbp.z = iohelp.ReadFloat64(r)
		default:
			io.ReadAll(r)
			return r.Err
		}
	}
}

func (bbp bar) Size() int {
	bodyLen := 5
	if bbp.x != nil {
		bodyLen += 1
		bodyLen += 8
	}
	if bbp.y != nil {
		bodyLen += 1
		bodyLen += 8
	}
	if bbp.z != nil {
		bodyLen += 1
		bodyLen += 8
	}
	return bodyLen
}

func (bbp bar) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func makebar(r iohelp.ErrorReader) (bar, error) {
	v := bar{}
	err := v.DecodeBebop(r)
	return v, err
}

func makebarFromBytes(buf []byte) (bar, error) {
	v := bar{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func mustMakebarFromBytes(buf []byte) bar {
	v := bar{}
	v.MustUnmarshalBebop(buf)
	return v
}

