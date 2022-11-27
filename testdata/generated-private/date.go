// Code generated by bebopc-go; DO NOT EDIT.

package generated

import (
	"github.com/200sc/bebop"
	"github.com/200sc/bebop/iohelp"
	"io"
	"time"
)

var _ bebop.Record = &myObj{}

type myObj struct {
	start *time.Time
	end *time.Time
}

func (bbp myObj) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(bbp.Size()-4))
	at += 4
	if bbp.start != nil {
		buf[at] = 1
		at++
		if *bbp.start != (time.Time{}) {
			iohelp.WriteInt64Bytes(buf[at:], ((*bbp.start).UnixNano() / 100))
		} else {
			iohelp.WriteInt64Bytes(buf[at:], 0)
		}
		at += 8
	}
	if bbp.end != nil {
		buf[at] = 2
		at++
		if *bbp.end != (time.Time{}) {
			iohelp.WriteInt64Bytes(buf[at:], ((*bbp.end).UnixNano() / 100))
		} else {
			iohelp.WriteInt64Bytes(buf[at:], 0)
		}
		at += 8
	}
	return at
}

func (bbp *myObj) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.start = new(time.Time)
			if len(buf[at:]) < 8 {
				return io.ErrUnexpectedEOF
			}
			(*bbp.start) = iohelp.ReadDateBytes(buf[at:])
			at += 8
		case 2:
			at += 1
			bbp.end = new(time.Time)
			if len(buf[at:]) < 8 {
				return io.ErrUnexpectedEOF
			}
			(*bbp.end) = iohelp.ReadDateBytes(buf[at:])
			at += 8
		default:
			return nil
		}
	}
}

func (bbp *myObj) MustUnmarshalBebop(buf []byte) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.start = new(time.Time)
			(*bbp.start) = iohelp.ReadDateBytes(buf[at:])
			at += 8
		case 2:
			at += 1
			bbp.end = new(time.Time)
			(*bbp.end) = iohelp.ReadDateBytes(buf[at:])
			at += 8
		default:
			return
		}
	}
}

func (bbp myObj) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(bbp.Size()-4))
	if bbp.start != nil {
		w.Write([]byte{1})
		if *bbp.start != (time.Time{}) {
			iohelp.WriteInt64(w, ((*bbp.start).UnixNano() / 100))
		} else {
			iohelp.WriteInt64(w, 0)
		}
	}
	if bbp.end != nil {
		w.Write([]byte{2})
		if *bbp.end != (time.Time{}) {
			iohelp.WriteInt64(w, ((*bbp.end).UnixNano() / 100))
		} else {
			iohelp.WriteInt64(w, 0)
		}
	}
	w.Write([]byte{0})
	return w.Err
}

func (bbp *myObj) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bodyLen := iohelp.ReadUint32(r)
	r.Reader = &io.LimitedReader{R:r.Reader, N:int64(bodyLen)}
	for {
		switch iohelp.ReadByte(r) {
		case 1:
			bbp.start = new(time.Time)
			*bbp.start = iohelp.ReadDate(r)
		case 2:
			bbp.end = new(time.Time)
			*bbp.end = iohelp.ReadDate(r)
		default:
			io.ReadAll(r)
			return r.Err
		}
	}
}

func (bbp myObj) Size() int {
	bodyLen := 5
	if bbp.start != nil {
		bodyLen += 1
		bodyLen += 8
	}
	if bbp.end != nil {
		bodyLen += 1
		bodyLen += 8
	}
	return bodyLen
}

func (bbp *myObj) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func makemyObj(r iohelp.ErrorReader) (myObj, error) {
	v := myObj{}
	err := v.DecodeBebop(r)
	return v, err
}

func makemyObjFromBytes(buf []byte) (myObj, error) {
	v := myObj{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func mustMakemyObjFromBytes(buf []byte) myObj {
	v := myObj{}
	v.MustUnmarshalBebop(buf)
	return v
}

