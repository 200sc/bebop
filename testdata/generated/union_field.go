// Code generated by bebopc-go; DO NOT EDIT.

package generated

import (
	"github.com/200sc/bebop"
	"github.com/200sc/bebop/iohelp"
	"io"
)

var _ bebop.Record = &WithUnionField{}

type WithUnionField struct {
	Test *List2
}

func (bbp WithUnionField) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(bbp.Size()-4))
	at += 4
	if bbp.Test != nil {
		buf[at] = 1
		at++
		(*bbp.Test).MarshalBebopTo(buf[at:])
		{
			tmp := (*bbp.Test)
			at += tmp.Size()
		}
		
	}
	return at
}

func (bbp *WithUnionField) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.Test = new(List2)
			(*bbp.Test), err = MakeList2FromBytes(buf[at:])
			if err != nil {
				return err
			}
			{
				tmp := ((*bbp.Test))
				at += tmp.Size()
			}
			
		default:
			return nil
		}
	}
}

func (bbp *WithUnionField) MustUnmarshalBebop(buf []byte) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.Test = new(List2)
			(*bbp.Test) = MustMakeList2FromBytes(buf[at:])
			{
				tmp := ((*bbp.Test))
				at += tmp.Size()
			}
			
		default:
			return
		}
	}
}

func (bbp WithUnionField) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(bbp.Size()-4))
	if bbp.Test != nil {
		w.Write([]byte{1})
		err = (*bbp.Test).EncodeBebop(w)
		if err != nil {
			return err
		}
	}
	w.Write([]byte{0})
	return w.Err
}

func (bbp *WithUnionField) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bodyLen := iohelp.ReadUint32(r)
	baseReader := r.Reader
	r.Reader = &io.LimitedReader{R: baseReader, N: int64(bodyLen)}
	for {
		switch iohelp.ReadByte(r) {
		case 1:
			bbp.Test = new(List2)
			(*bbp.Test), err = MakeList2(r)
			if err != nil {
				return err
			}
		default:
			r.Drain()
			r.Reader = baseReader
			return r.Err
		}
	}
}

func (bbp WithUnionField) Size() int {
	bodyLen := 5
	if bbp.Test != nil {
		bodyLen += 1
		{
			tmp := (*bbp.Test)
			bodyLen += tmp.Size()
		}
		
	}
	return bodyLen
}

func (bbp WithUnionField) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func MakeWithUnionField(r *iohelp.ErrorReader) (WithUnionField, error) {
	v := WithUnionField{}
	err := v.DecodeBebop(r)
	return v, err
}

func MakeWithUnionFieldFromBytes(buf []byte) (WithUnionField, error) {
	v := WithUnionField{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func MustMakeWithUnionFieldFromBytes(buf []byte) WithUnionField {
	v := WithUnionField{}
	v.MustUnmarshalBebop(buf)
	return v
}

var _ bebop.Record = &Cons2{}

type Cons2 struct {
	Head uint32
	Tail List
}

func (bbp Cons2) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], bbp.Head)
	at += 4
	
	return at
}

func (bbp *Cons2) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	if len(buf[at:]) < 4 {
		return io.ErrUnexpectedEOF
	}
	bbp.Head = iohelp.ReadUint32Bytes(buf[at:])
	at += 4
	
	return nil
}

func (bbp *Cons2) MustUnmarshalBebop(buf []byte) {
	at := 0
	bbp.Head = iohelp.ReadUint32Bytes(buf[at:])
	at += 4
	
}

func (bbp Cons2) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, bbp.Head)
	
	return w.Err
}

func (bbp *Cons2) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bbp.Head = iohelp.ReadUint32(r)
	
	return r.Err
}

func (bbp Cons2) Size() int {
	bodyLen := 0
	bodyLen += 4
	
	return bodyLen
}

func (bbp Cons2) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func MakeCons2(r *iohelp.ErrorReader) (Cons2, error) {
	v := Cons2{}
	err := v.DecodeBebop(r)
	return v, err
}

func MakeCons2FromBytes(buf []byte) (Cons2, error) {
	v := Cons2{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func MustMakeCons2FromBytes(buf []byte) Cons2 {
	v := Cons2{}
	v.MustUnmarshalBebop(buf)
	return v
}

var _ bebop.Record = &Nil2{}

// nil is empty
type Nil2 struct {
}

func (bbp Nil2) MarshalBebopTo(buf []byte) int {
	return 0
}

func (bbp *Nil2) UnmarshalBebop(buf []byte) (err error) {
	return nil
}

func (bbp *Nil2) MustUnmarshalBebop(buf []byte) {
}

func (bbp Nil2) EncodeBebop(iow io.Writer) (err error) {
	return nil
}

func (bbp *Nil2) DecodeBebop(ior io.Reader) (err error) {
	return nil
}

func (bbp Nil2) Size() int {
	return 0
}

func (bbp Nil2) MarshalBebop() []byte {
	return []byte{}
}

func MakeNil2(r *iohelp.ErrorReader) (Nil2, error) {
	return Nil2{}, nil
}

func MakeNil2FromBytes(buf []byte) (Nil2, error) {
	return Nil2{}, nil
}

func MustMakeNil2FromBytes(buf []byte) Nil2 {
	return Nil2{}
}

var _ bebop.Record = &List2{}

type List2 struct {
	Cons2 *Cons2
	Nil2 *Nil2
}

func (bbp List2) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(bbp.Size()-5))
	at += 4
	if bbp.Cons2 != nil {
		buf[at] = 1
		at++
		(*bbp.Cons2).MarshalBebopTo(buf[at:])
		{
			tmp := (*bbp.Cons2)
			at += tmp.Size()
		}
		
		return at
	}
	if bbp.Nil2 != nil {
		buf[at] = 2
		at++
		(*bbp.Nil2).MarshalBebopTo(buf[at:])
		{
			tmp := (*bbp.Nil2)
			at += tmp.Size()
		}
		
		return at
	}
	return at
}

func (bbp *List2) UnmarshalBebop(buf []byte) (err error) {
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
			bbp.Cons2 = new(Cons2)
			(*bbp.Cons2), err = MakeCons2FromBytes(buf[at:])
			if err != nil {
				return err
			}
			{
				tmp := ((*bbp.Cons2))
				at += tmp.Size()
			}
			
			return nil
		case 2:
			at += 1
			bbp.Nil2 = new(Nil2)
			(*bbp.Nil2), err = MakeNil2FromBytes(buf[at:])
			if err != nil {
				return err
			}
			{
				tmp := ((*bbp.Nil2))
				at += tmp.Size()
			}
			
			return nil
		default:
			return nil
		}
	}
}

func (bbp *List2) MustUnmarshalBebop(buf []byte) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.Cons2 = new(Cons2)
			(*bbp.Cons2) = MustMakeCons2FromBytes(buf[at:])
			{
				tmp := ((*bbp.Cons2))
				at += tmp.Size()
			}
			
			return
		case 2:
			at += 1
			bbp.Nil2 = new(Nil2)
			(*bbp.Nil2) = MustMakeNil2FromBytes(buf[at:])
			{
				tmp := ((*bbp.Nil2))
				at += tmp.Size()
			}
			
			return
		default:
			return
		}
	}
}

func (bbp List2) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(bbp.Size()-5))
	if bbp.Cons2 != nil {
		w.Write([]byte{1})
		err = (*bbp.Cons2).EncodeBebop(w)
		if err != nil {
			return err
		}
		return w.Err
	}
	if bbp.Nil2 != nil {
		w.Write([]byte{2})
		err = (*bbp.Nil2).EncodeBebop(w)
		if err != nil {
			return err
		}
		return w.Err
	}
	return w.Err
}

func (bbp *List2) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bodyLen := iohelp.ReadUint32(r)
	baseReader := r.Reader
	limitReader := &io.LimitedReader{R: baseReader, N: int64(bodyLen)+1}
	r.Reader = limitReader
	for {
		switch iohelp.ReadByte(r) {
		case 1:
			bbp.Cons2 = new(Cons2)
			(*bbp.Cons2), err = MakeCons2(r)
			if err != nil {
				return err
			}
			r.Drain()
			r.Reader = baseReader
			return r.Err
		case 2:
			bbp.Nil2 = new(Nil2)
			(*bbp.Nil2), err = MakeNil2(r)
			if err != nil {
				return err
			}
			r.Drain()
			r.Reader = baseReader
			return r.Err
		default:
			r.Drain()
			r.Reader = baseReader
			return r.Err
		}
	}
}

func (bbp List2) Size() int {
	bodyLen := 4
	if bbp.Cons2 != nil {
		bodyLen += 1
		{
			tmp := (*bbp.Cons2)
			bodyLen += tmp.Size()
		}
		
		return bodyLen
	}
	if bbp.Nil2 != nil {
		bodyLen += 1
		{
			tmp := (*bbp.Nil2)
			bodyLen += tmp.Size()
		}
		
		return bodyLen
	}
	return bodyLen
}

func (bbp List2) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func MakeList2(r *iohelp.ErrorReader) (List2, error) {
	v := List2{}
	err := v.DecodeBebop(r)
	return v, err
}

func MakeList2FromBytes(buf []byte) (List2, error) {
	v := List2{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func MustMakeList2FromBytes(buf []byte) List2 {
	v := List2{}
	v.MustUnmarshalBebop(buf)
	return v
}

