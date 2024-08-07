// Code generated by bebopc-go; DO NOT EDIT.

package generated

import (
	"github.com/200sc/bebop"
	"github.com/200sc/bebop/iohelp"
	"io"
)

var _ bebop.Record = &Print{}

type Print struct {
	Printout string
}

func (bbp Print) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.Printout)))
	copy(buf[at+4:at+4+len(bbp.Printout)], []byte(bbp.Printout))
	at += 4 + len(bbp.Printout)
	return at
}

func (bbp *Print) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	bbp.Printout, err = iohelp.ReadStringBytes(buf[at:])
	if err != nil {
		return err
	}
	at += 4 + len(bbp.Printout)
	return nil
}

func (bbp *Print) MustUnmarshalBebop(buf []byte) {
	at := 0
	bbp.Printout = iohelp.MustReadStringBytes(buf[at:])
	at += 4 + len(bbp.Printout)
}

func (bbp Print) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(len(bbp.Printout)))
	w.Write([]byte(bbp.Printout))
	return w.Err
}

func (bbp *Print) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bbp.Printout = iohelp.ReadString(r)
	return r.Err
}

func (bbp Print) Size() int {
	bodyLen := 0
	bodyLen += 4 + len(bbp.Printout)
	return bodyLen
}

func (bbp Print) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func MakePrint(r *iohelp.ErrorReader) (Print, error) {
	v := Print{}
	err := v.DecodeBebop(r)
	return v, err
}

func MakePrintFromBytes(buf []byte) (Print, error) {
	v := Print{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func MustMakePrintFromBytes(buf []byte) Print {
	v := Print{}
	v.MustUnmarshalBebop(buf)
	return v
}

var _ bebop.Record = &Add{}

type Add struct {
	A int32
	B int32
}

func (bbp Add) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteInt32Bytes(buf[at:], bbp.A)
	at += 4
	iohelp.WriteInt32Bytes(buf[at:], bbp.B)
	at += 4
	return at
}

func (bbp *Add) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	if len(buf[at:]) < 4 {
		return io.ErrUnexpectedEOF
	}
	bbp.A = iohelp.ReadInt32Bytes(buf[at:])
	at += 4
	if len(buf[at:]) < 4 {
		return io.ErrUnexpectedEOF
	}
	bbp.B = iohelp.ReadInt32Bytes(buf[at:])
	at += 4
	return nil
}

func (bbp *Add) MustUnmarshalBebop(buf []byte) {
	at := 0
	bbp.A = iohelp.ReadInt32Bytes(buf[at:])
	at += 4
	bbp.B = iohelp.ReadInt32Bytes(buf[at:])
	at += 4
}

func (bbp Add) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteInt32(w, bbp.A)
	iohelp.WriteInt32(w, bbp.B)
	return w.Err
}

func (bbp *Add) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bbp.A = iohelp.ReadInt32(r)
	bbp.B = iohelp.ReadInt32(r)
	return r.Err
}

func (bbp Add) Size() int {
	bodyLen := 0
	bodyLen += 4
	bodyLen += 4
	return bodyLen
}

func (bbp Add) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func MakeAdd(r *iohelp.ErrorReader) (Add, error) {
	v := Add{}
	err := v.DecodeBebop(r)
	return v, err
}

func MakeAddFromBytes(buf []byte) (Add, error) {
	v := Add{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func MustMakeAddFromBytes(buf []byte) Add {
	v := Add{}
	v.MustUnmarshalBebop(buf)
	return v
}

var _ bebop.Record = &AddResponse{}

type AddResponse struct {
	C int32
}

func (bbp AddResponse) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteInt32Bytes(buf[at:], bbp.C)
	at += 4
	return at
}

func (bbp *AddResponse) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	if len(buf[at:]) < 4 {
		return io.ErrUnexpectedEOF
	}
	bbp.C = iohelp.ReadInt32Bytes(buf[at:])
	at += 4
	return nil
}

func (bbp *AddResponse) MustUnmarshalBebop(buf []byte) {
	at := 0
	bbp.C = iohelp.ReadInt32Bytes(buf[at:])
	at += 4
}

func (bbp AddResponse) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteInt32(w, bbp.C)
	return w.Err
}

func (bbp *AddResponse) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bbp.C = iohelp.ReadInt32(r)
	return r.Err
}

func (bbp AddResponse) Size() int {
	bodyLen := 0
	bodyLen += 4
	return bodyLen
}

func (bbp AddResponse) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func MakeAddResponse(r *iohelp.ErrorReader) (AddResponse, error) {
	v := AddResponse{}
	err := v.DecodeBebop(r)
	return v, err
}

func MakeAddResponseFromBytes(buf []byte) (AddResponse, error) {
	v := AddResponse{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func MustMakeAddResponseFromBytes(buf []byte) AddResponse {
	v := AddResponse{}
	v.MustUnmarshalBebop(buf)
	return v
}

const PrintRequestOpCode = 0x2

var _ bebop.Record = &PrintRequest{}

type PrintRequest struct {
	Print *Print
}

func (bbp PrintRequest) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(bbp.Size()-4))
	at += 4
	if bbp.Print != nil {
		buf[at] = 1
		at++
		(*bbp.Print).MarshalBebopTo(buf[at:])
		{
			tmp := (*bbp.Print)
			at += tmp.Size()
		}
		
	}
	return at
}

func (bbp *PrintRequest) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.Print = new(Print)
			(*bbp.Print), err = MakePrintFromBytes(buf[at:])
			if err != nil {
				return err
			}
			{
				tmp := ((*bbp.Print))
				at += tmp.Size()
			}
			
		default:
			return nil
		}
	}
}

func (bbp *PrintRequest) MustUnmarshalBebop(buf []byte) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.Print = new(Print)
			(*bbp.Print) = MustMakePrintFromBytes(buf[at:])
			{
				tmp := ((*bbp.Print))
				at += tmp.Size()
			}
			
		default:
			return
		}
	}
}

func (bbp PrintRequest) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(bbp.Size()-4))
	if bbp.Print != nil {
		w.Write([]byte{1})
		err = (*bbp.Print).EncodeBebop(w)
		if err != nil {
			return err
		}
	}
	w.Write([]byte{0})
	return w.Err
}

func (bbp *PrintRequest) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bodyLen := iohelp.ReadUint32(r)
	baseReader := r.Reader
	r.Reader = &io.LimitedReader{R: baseReader, N: int64(bodyLen)}
	for {
		switch iohelp.ReadByte(r) {
		case 1:
			bbp.Print = new(Print)
			(*bbp.Print), err = MakePrint(r)
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

func (bbp PrintRequest) Size() int {
	bodyLen := 5
	if bbp.Print != nil {
		bodyLen += 1
		{
			tmp := (*bbp.Print)
			bodyLen += tmp.Size()
		}
		
	}
	return bodyLen
}

func (bbp PrintRequest) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func MakePrintRequest(r *iohelp.ErrorReader) (PrintRequest, error) {
	v := PrintRequest{}
	err := v.DecodeBebop(r)
	return v, err
}

func MakePrintRequestFromBytes(buf []byte) (PrintRequest, error) {
	v := PrintRequest{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func MustMakePrintRequestFromBytes(buf []byte) PrintRequest {
	v := PrintRequest{}
	v.MustUnmarshalBebop(buf)
	return v
}

const AddRequestOpCode = 0x1

var _ bebop.Record = &AddRequest{}

type AddRequest struct {
	Add *Add
}

func (bbp AddRequest) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(bbp.Size()-4))
	at += 4
	if bbp.Add != nil {
		buf[at] = 1
		at++
		(*bbp.Add).MarshalBebopTo(buf[at:])
		{
			tmp := (*bbp.Add)
			at += tmp.Size()
		}
		
	}
	return at
}

func (bbp *AddRequest) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.Add = new(Add)
			(*bbp.Add), err = MakeAddFromBytes(buf[at:])
			if err != nil {
				return err
			}
			{
				tmp := ((*bbp.Add))
				at += tmp.Size()
			}
			
		default:
			return nil
		}
	}
}

func (bbp *AddRequest) MustUnmarshalBebop(buf []byte) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.Add = new(Add)
			(*bbp.Add) = MustMakeAddFromBytes(buf[at:])
			{
				tmp := ((*bbp.Add))
				at += tmp.Size()
			}
			
		default:
			return
		}
	}
}

func (bbp AddRequest) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(bbp.Size()-4))
	if bbp.Add != nil {
		w.Write([]byte{1})
		err = (*bbp.Add).EncodeBebop(w)
		if err != nil {
			return err
		}
	}
	w.Write([]byte{0})
	return w.Err
}

func (bbp *AddRequest) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bodyLen := iohelp.ReadUint32(r)
	baseReader := r.Reader
	r.Reader = &io.LimitedReader{R: baseReader, N: int64(bodyLen)}
	for {
		switch iohelp.ReadByte(r) {
		case 1:
			bbp.Add = new(Add)
			(*bbp.Add), err = MakeAdd(r)
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

func (bbp AddRequest) Size() int {
	bodyLen := 5
	if bbp.Add != nil {
		bodyLen += 1
		{
			tmp := (*bbp.Add)
			bodyLen += tmp.Size()
		}
		
	}
	return bodyLen
}

func (bbp AddRequest) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func MakeAddRequest(r *iohelp.ErrorReader) (AddRequest, error) {
	v := AddRequest{}
	err := v.DecodeBebop(r)
	return v, err
}

func MakeAddRequestFromBytes(buf []byte) (AddRequest, error) {
	v := AddRequest{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func MustMakeAddRequestFromBytes(buf []byte) AddRequest {
	v := AddRequest{}
	v.MustUnmarshalBebop(buf)
	return v
}

