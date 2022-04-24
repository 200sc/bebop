// Code generated by bebopc-go; DO NOT EDIT.

package generated

import (
	"io"
	"github.com/200sc/bebop"
	"github.com/200sc/bebop/iohelp"
)

var _ bebop.Record = &print{}

type print struct {
	printout string
}

func (bbp print) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.printout)))
	copy(buf[at+4:at+4+len(bbp.printout)], []byte(bbp.printout))
	at += 4 + len(bbp.printout)
	return at
}

func (bbp *print) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	bbp.printout, err = iohelp.ReadStringBytes(buf[at:])
	if err != nil{
		return err
	}
	at += 4 + len(bbp.printout)
	return nil
}

func (bbp *print) MustUnmarshalBebop(buf []byte) {
	at := 0
	bbp.printout = iohelp.MustReadStringBytes(buf[at:])
	at += 4 + len(bbp.printout)
}

func (bbp print) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(len(bbp.printout)))
	w.Write([]byte(bbp.printout))
	return w.Err
}

func (bbp *print) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bbp.printout = iohelp.ReadString(r)
	return r.Err
}

func (bbp print) Size() int {
	bodyLen := 0
	bodyLen += 4 + len(bbp.printout)
	return bodyLen
}

func (bbp print) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func makeprint(r iohelp.ErrorReader) (print, error) {
	v := print{}
	err := v.DecodeBebop(r)
	return v, err
}

func makeprintFromBytes(buf []byte) (print, error) {
	v := print{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func mustMakeprintFromBytes(buf []byte) print {
	v := print{}
	v.MustUnmarshalBebop(buf)
	return v
}

var _ bebop.Record = &add{}

type add struct {
	a int32
	b int32
}

func (bbp add) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteInt32Bytes(buf[at:], bbp.a)
	at += 4
	iohelp.WriteInt32Bytes(buf[at:], bbp.b)
	at += 4
	return at
}

func (bbp *add) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	if len(buf[at:]) < 4 {
		 return io.ErrUnexpectedEOF
	}
	bbp.a = iohelp.ReadInt32Bytes(buf[at:])
	at += 4
	if len(buf[at:]) < 4 {
		 return io.ErrUnexpectedEOF
	}
	bbp.b = iohelp.ReadInt32Bytes(buf[at:])
	at += 4
	return nil
}

func (bbp *add) MustUnmarshalBebop(buf []byte) {
	at := 0
	bbp.a = iohelp.ReadInt32Bytes(buf[at:])
	at += 4
	bbp.b = iohelp.ReadInt32Bytes(buf[at:])
	at += 4
}

func (bbp add) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteInt32(w, bbp.a)
	iohelp.WriteInt32(w, bbp.b)
	return w.Err
}

func (bbp *add) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bbp.a = iohelp.ReadInt32(r)
	bbp.b = iohelp.ReadInt32(r)
	return r.Err
}

func (bbp add) Size() int {
	bodyLen := 0
	bodyLen += 4
	bodyLen += 4
	return bodyLen
}

func (bbp add) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func makeadd(r iohelp.ErrorReader) (add, error) {
	v := add{}
	err := v.DecodeBebop(r)
	return v, err
}

func makeaddFromBytes(buf []byte) (add, error) {
	v := add{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func mustMakeaddFromBytes(buf []byte) add {
	v := add{}
	v.MustUnmarshalBebop(buf)
	return v
}

var _ bebop.Record = &addResponse{}

type addResponse struct {
	c int32
}

func (bbp addResponse) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteInt32Bytes(buf[at:], bbp.c)
	at += 4
	return at
}

func (bbp *addResponse) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	if len(buf[at:]) < 4 {
		 return io.ErrUnexpectedEOF
	}
	bbp.c = iohelp.ReadInt32Bytes(buf[at:])
	at += 4
	return nil
}

func (bbp *addResponse) MustUnmarshalBebop(buf []byte) {
	at := 0
	bbp.c = iohelp.ReadInt32Bytes(buf[at:])
	at += 4
}

func (bbp addResponse) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteInt32(w, bbp.c)
	return w.Err
}

func (bbp *addResponse) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bbp.c = iohelp.ReadInt32(r)
	return r.Err
}

func (bbp addResponse) Size() int {
	bodyLen := 0
	bodyLen += 4
	return bodyLen
}

func (bbp addResponse) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func makeaddResponse(r iohelp.ErrorReader) (addResponse, error) {
	v := addResponse{}
	err := v.DecodeBebop(r)
	return v, err
}

func makeaddResponseFromBytes(buf []byte) (addResponse, error) {
	v := addResponse{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func mustMakeaddResponseFromBytes(buf []byte) addResponse {
	v := addResponse{}
	v.MustUnmarshalBebop(buf)
	return v
}

const printRequestOpCode = 0x2

var _ bebop.Record = &printRequest{}

type printRequest struct {
	print *print
}

func (bbp printRequest) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(bbp.Size()-4))
	at += 4
	if bbp.print != nil {
		buf[at] = 1
		at++
		(*bbp.print).MarshalBebopTo(buf[at:])
		at += (*bbp.print).Size()
	}
	return at
}

func (bbp *printRequest) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.print = new(print)
			(*bbp.print) = mustMakeprintFromBytes(buf[at:])
			if err != nil{
				return err
			}
			at += ((*bbp.print)).Size()
		default:
			return nil
		}
	}
}

func (bbp *printRequest) MustUnmarshalBebop(buf []byte) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.print = new(print)
			(*bbp.print) = mustMakeprintFromBytes(buf[at:])
			at += ((*bbp.print)).Size()
		default:
			return
		}
	}
}

func (bbp printRequest) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(bbp.Size()-4))
	if bbp.print != nil {
		w.Write([]byte{1})
		err = (*bbp.print).EncodeBebop(w)
		if err != nil{
			return err
		}
	}
	w.Write([]byte{0})
	return w.Err
}

func (bbp *printRequest) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bodyLen := iohelp.ReadUint32(r)
	r.Reader = &io.LimitedReader{R:r.Reader, N:int64(bodyLen)}
	for {
		switch iohelp.ReadByte(r) {
		case 1:
			bbp.print = new(print)
			(*bbp.print), err = makeprint(r)
			if err != nil{
				return err
			}
		default:
			io.ReadAll(r)
			return r.Err
		}
	}
}

func (bbp printRequest) Size() int {
	bodyLen := 5
	if bbp.print != nil {
		bodyLen += 1
		bodyLen += (*bbp.print).Size()
	}
	return bodyLen
}

func (bbp printRequest) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func makeprintRequest(r iohelp.ErrorReader) (printRequest, error) {
	v := printRequest{}
	err := v.DecodeBebop(r)
	return v, err
}

func makeprintRequestFromBytes(buf []byte) (printRequest, error) {
	v := printRequest{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func mustMakeprintRequestFromBytes(buf []byte) printRequest {
	v := printRequest{}
	v.MustUnmarshalBebop(buf)
	return v
}

const addRequestOpCode = 0x1

var _ bebop.Record = &addRequest{}

type addRequest struct {
	add *add
}

func (bbp addRequest) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(bbp.Size()-4))
	at += 4
	if bbp.add != nil {
		buf[at] = 1
		at++
		(*bbp.add).MarshalBebopTo(buf[at:])
		at += (*bbp.add).Size()
	}
	return at
}

func (bbp *addRequest) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.add = new(add)
			(*bbp.add) = mustMakeaddFromBytes(buf[at:])
			if err != nil{
				return err
			}
			at += ((*bbp.add)).Size()
		default:
			return nil
		}
	}
}

func (bbp *addRequest) MustUnmarshalBebop(buf []byte) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.add = new(add)
			(*bbp.add) = mustMakeaddFromBytes(buf[at:])
			at += ((*bbp.add)).Size()
		default:
			return
		}
	}
}

func (bbp addRequest) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(bbp.Size()-4))
	if bbp.add != nil {
		w.Write([]byte{1})
		err = (*bbp.add).EncodeBebop(w)
		if err != nil{
			return err
		}
	}
	w.Write([]byte{0})
	return w.Err
}

func (bbp *addRequest) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bodyLen := iohelp.ReadUint32(r)
	r.Reader = &io.LimitedReader{R:r.Reader, N:int64(bodyLen)}
	for {
		switch iohelp.ReadByte(r) {
		case 1:
			bbp.add = new(add)
			(*bbp.add), err = makeadd(r)
			if err != nil{
				return err
			}
		default:
			io.ReadAll(r)
			return r.Err
		}
	}
}

func (bbp addRequest) Size() int {
	bodyLen := 5
	if bbp.add != nil {
		bodyLen += 1
		bodyLen += (*bbp.add).Size()
	}
	return bodyLen
}

func (bbp addRequest) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func makeaddRequest(r iohelp.ErrorReader) (addRequest, error) {
	v := addRequest{}
	err := v.DecodeBebop(r)
	return v, err
}

func makeaddRequestFromBytes(buf []byte) (addRequest, error) {
	v := addRequest{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func mustMakeaddRequestFromBytes(buf []byte) addRequest {
	v := addRequest{}
	v.MustUnmarshalBebop(buf)
	return v
}

