// Code generated by bebopc-go; DO NOT EDIT.

package generated

import (
	"github.com/200sc/bebop"
	"github.com/200sc/bebop/iohelp"
	"io"
)

var _ bebop.Record = &withUnionField{}

type withUnionField struct {
	test *list2
}

func (bbp withUnionField) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(bbp.Size()-4))
	at += 4
	if bbp.test != nil {
		buf[at] = 1
		at++
		(*bbp.test).MarshalBebopTo(buf[at:])
		at += (*bbp.test).Size()
	}
	return at
}

func (bbp *withUnionField) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.test = new(list2)
			(*bbp.test), err = makelist2FromBytes(buf[at:])
			if err != nil {
				return err
			}
			at += ((*bbp.test)).Size()
		default:
			return nil
		}
	}
}

func (bbp *withUnionField) MustUnmarshalBebop(buf []byte) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.test = new(list2)
			(*bbp.test) = mustMakelist2FromBytes(buf[at:])
			at += ((*bbp.test)).Size()
		default:
			return
		}
	}
}

func (bbp withUnionField) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(bbp.Size()-4))
	if bbp.test != nil {
		w.Write([]byte{1})
		err = (*bbp.test).EncodeBebop(w)
		if err != nil {
			return err
		}
	}
	w.Write([]byte{0})
	return w.Err
}

func (bbp *withUnionField) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bodyLen := iohelp.ReadUint32(r)
	r.Reader = &io.LimitedReader{R:r.Reader, N:int64(bodyLen)}
	for {
		switch iohelp.ReadByte(r) {
		case 1:
			bbp.test = new(list2)
			(*bbp.test), err = makelist2(r)
			if err != nil {
				return err
			}
		default:
			io.ReadAll(r)
			return r.Err
		}
	}
}

func (bbp withUnionField) Size() int {
	bodyLen := 5
	if bbp.test != nil {
		bodyLen += 1
		bodyLen += (*bbp.test).Size()
	}
	return bodyLen
}

func (bbp *withUnionField) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func makewithUnionField(r iohelp.ErrorReader) (withUnionField, error) {
	v := withUnionField{}
	err := v.DecodeBebop(r)
	return v, err
}

func makewithUnionFieldFromBytes(buf []byte) (withUnionField, error) {
	v := withUnionField{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func mustMakewithUnionFieldFromBytes(buf []byte) withUnionField {
	v := withUnionField{}
	v.MustUnmarshalBebop(buf)
	return v
}

var _ bebop.Record = &cons2{}

type cons2 struct {
	head uint32
	tail list
}

func (bbp *cons2) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], bbp.head)
	at += 4
	
	return at
}

func (bbp *cons2) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	if len(buf[at:]) < 4 {
		return io.ErrUnexpectedEOF
	}
	bbp.head = iohelp.ReadUint32Bytes(buf[at:])
	at += 4
	
	return nil
}

func (bbp *cons2) MustUnmarshalBebop(buf []byte) {
	at := 0
	bbp.head = iohelp.ReadUint32Bytes(buf[at:])
	at += 4
	
}

func (bbp *cons2) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, bbp.head)
	
	return w.Err
}

func (bbp *cons2) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bbp.head = iohelp.ReadUint32(r)
	
	return r.Err
}

func (bbp *cons2) Size() int {
	bodyLen := 0
	bodyLen += 4
	
	return bodyLen
}

func (bbp *cons2) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func makecons2(r iohelp.ErrorReader) (cons2, error) {
	v := cons2{}
	err := v.DecodeBebop(r)
	return v, err
}

func makecons2FromBytes(buf []byte) (cons2, error) {
	v := cons2{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func mustMakecons2FromBytes(buf []byte) cons2 {
	v := cons2{}
	v.MustUnmarshalBebop(buf)
	return v
}

var _ bebop.Record = &nil2{}

// nil is empty
type nil2 struct {
}

func (bbp *nil2) MarshalBebopTo(buf []byte) int {
	return 0
}

func (bbp *nil2) UnmarshalBebop(buf []byte) (err error) {
	return nil
}

func (bbp *nil2) MustUnmarshalBebop(buf []byte) {
}

func (bbp *nil2) EncodeBebop(iow io.Writer) (err error) {
	return nil
}

func (bbp *nil2) DecodeBebop(ior io.Reader) (err error) {
	return nil
}

func (bbp *nil2) Size() int {
	return 0
}

func (bbp *nil2) MarshalBebop() []byte {
	return []byte{}
}

func makenil2(r iohelp.ErrorReader) (nil2, error) {
	return nil2{}, nil
}

func makenil2FromBytes(buf []byte) (nil2, error) {
	return nil2{}, nil
}

func mustMakenil2FromBytes(buf []byte) nil2 {
	return nil2{}
}

var _ bebop.Record = &list2{}

type list2 struct {
	cons2 *cons2
	nil2 *nil2
}

func (bbp list2) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(bbp.Size()-5))
	at += 4
	if bbp.cons2 != nil {
		buf[at] = 1
		at++
		(*bbp.cons2).MarshalBebopTo(buf[at:])
		at += (*bbp.cons2).Size()
		return at
	}
	if bbp.nil2 != nil {
		buf[at] = 2
		at++
		(*bbp.nil2).MarshalBebopTo(buf[at:])
		at += (*bbp.nil2).Size()
		return at
	}
	return at
}

func (bbp *list2) UnmarshalBebop(buf []byte) (err error) {
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
			bbp.cons2 = new(cons2)
			(*bbp.cons2), err = makecons2FromBytes(buf[at:])
			if err != nil {
				return err
			}
			at += ((*bbp.cons2)).Size()
			return nil
		case 2:
			at += 1
			bbp.nil2 = new(nil2)
			(*bbp.nil2), err = makenil2FromBytes(buf[at:])
			if err != nil {
				return err
			}
			at += ((*bbp.nil2)).Size()
			return nil
		default:
			return nil
		}
	}
}

func (bbp *list2) MustUnmarshalBebop(buf []byte) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.cons2 = new(cons2)
			(*bbp.cons2) = mustMakecons2FromBytes(buf[at:])
			at += ((*bbp.cons2)).Size()
			return
		case 2:
			at += 1
			bbp.nil2 = new(nil2)
			(*bbp.nil2) = mustMakenil2FromBytes(buf[at:])
			at += ((*bbp.nil2)).Size()
			return
		default:
			return
		}
	}
}

func (bbp list2) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(bbp.Size()-5))
	if bbp.cons2 != nil {
		w.Write([]byte{1})
		err = (*bbp.cons2).EncodeBebop(w)
		if err != nil {
			return err
		}
		return w.Err
	}
	if bbp.nil2 != nil {
		w.Write([]byte{2})
		err = (*bbp.nil2).EncodeBebop(w)
		if err != nil {
			return err
		}
		return w.Err
	}
	return w.Err
}

func (bbp *list2) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bodyLen := iohelp.ReadUint32(r)
	r.Reader = &io.LimitedReader{R: r.Reader, N: int64(bodyLen) + 1}
	for {
		switch iohelp.ReadByte(r) {
		case 1:
			bbp.cons2 = new(cons2)
			(*bbp.cons2), err = makecons2(r)
			if err != nil {
				return err
			}
			io.ReadAll(r)
			return r.Err
		case 2:
			bbp.nil2 = new(nil2)
			(*bbp.nil2), err = makenil2(r)
			if err != nil {
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

func (bbp list2) Size() int {
	bodyLen := 4
	if bbp.cons2 != nil {
		bodyLen += 1
		bodyLen += (*bbp.cons2).Size()
		return bodyLen
	}
	if bbp.nil2 != nil {
		bodyLen += 1
		bodyLen += (*bbp.nil2).Size()
		return bodyLen
	}
	return bodyLen
}

func (bbp *list2) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func makelist2(r iohelp.ErrorReader) (list2, error) {
	v := list2{}
	err := v.DecodeBebop(r)
	return v, err
}

func makelist2FromBytes(buf []byte) (list2, error) {
	v := list2{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func mustMakelist2FromBytes(buf []byte) list2 {
	v := list2{}
	v.MustUnmarshalBebop(buf)
	return v
}

