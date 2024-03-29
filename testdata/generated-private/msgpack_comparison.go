// Code generated by bebopc-go; DO NOT EDIT.

package generated

import (
	"github.com/200sc/bebop"
	"github.com/200sc/bebop/iohelp"
	"io"
)

var _ bebop.Record = &msgpackComparison{}

// These field names are extremely weirdly capitalized, because I wanted the
// key names in JSON to be the same length while not coinciding with Bebop keywords.
type msgpackComparison struct {
	iNT0 uint8
	iNT1 uint8
	iNT1_ int16
	iNT8 uint8
	iNT8_ int16
	iNT16 int16
	iNT16_ int16
	iNT32 int32
	iNT32_ int32
	// int8 nIL; // "nil": null,
	tRUE bool
	fALSE bool
	fLOAT float64
	fLOAT_x float64
	sTRING0 string
	sTRING1 string
	sTRING4 string
	sTRING8 string
	sTRING16 string
	aRRAY0 []int32
	aRRAY1 []string
	aRRAY8 []int32
}

func (bbp msgpackComparison) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint8Bytes(buf[at:], bbp.iNT0)
	at += 1
	iohelp.WriteUint8Bytes(buf[at:], bbp.iNT1)
	at += 1
	iohelp.WriteInt16Bytes(buf[at:], bbp.iNT1_)
	at += 2
	iohelp.WriteUint8Bytes(buf[at:], bbp.iNT8)
	at += 1
	iohelp.WriteInt16Bytes(buf[at:], bbp.iNT8_)
	at += 2
	iohelp.WriteInt16Bytes(buf[at:], bbp.iNT16)
	at += 2
	iohelp.WriteInt16Bytes(buf[at:], bbp.iNT16_)
	at += 2
	iohelp.WriteInt32Bytes(buf[at:], bbp.iNT32)
	at += 4
	iohelp.WriteInt32Bytes(buf[at:], bbp.iNT32_)
	at += 4
	iohelp.WriteBoolBytes(buf[at:], bbp.tRUE)
	at += 1
	iohelp.WriteBoolBytes(buf[at:], bbp.fALSE)
	at += 1
	iohelp.WriteFloat64Bytes(buf[at:], bbp.fLOAT)
	at += 8
	iohelp.WriteFloat64Bytes(buf[at:], bbp.fLOAT_x)
	at += 8
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.sTRING0)))
	copy(buf[at+4:at+4+len(bbp.sTRING0)], []byte(bbp.sTRING0))
	at += 4 + len(bbp.sTRING0)
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.sTRING1)))
	copy(buf[at+4:at+4+len(bbp.sTRING1)], []byte(bbp.sTRING1))
	at += 4 + len(bbp.sTRING1)
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.sTRING4)))
	copy(buf[at+4:at+4+len(bbp.sTRING4)], []byte(bbp.sTRING4))
	at += 4 + len(bbp.sTRING4)
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.sTRING8)))
	copy(buf[at+4:at+4+len(bbp.sTRING8)], []byte(bbp.sTRING8))
	at += 4 + len(bbp.sTRING8)
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.sTRING16)))
	copy(buf[at+4:at+4+len(bbp.sTRING16)], []byte(bbp.sTRING16))
	at += 4 + len(bbp.sTRING16)
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.aRRAY0)))
	at += 4
	for _, v1 := range bbp.aRRAY0 {
		iohelp.WriteInt32Bytes(buf[at:], v1)
		at += 4
	}
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.aRRAY1)))
	at += 4
	for _, v1 := range bbp.aRRAY1 {
		iohelp.WriteUint32Bytes(buf[at:], uint32(len(v1)))
		copy(buf[at+4:at+4+len(v1)], []byte(v1))
		at += 4 + len(v1)
	}
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.aRRAY8)))
	at += 4
	for _, v1 := range bbp.aRRAY8 {
		iohelp.WriteInt32Bytes(buf[at:], v1)
		at += 4
	}
	return at
}

func (bbp *msgpackComparison) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	if len(buf[at:]) < 1 {
		return io.ErrUnexpectedEOF
	}
	bbp.iNT0 = iohelp.ReadUint8Bytes(buf[at:])
	at += 1
	if len(buf[at:]) < 1 {
		return io.ErrUnexpectedEOF
	}
	bbp.iNT1 = iohelp.ReadUint8Bytes(buf[at:])
	at += 1
	if len(buf[at:]) < 2 {
		return io.ErrUnexpectedEOF
	}
	bbp.iNT1_ = iohelp.ReadInt16Bytes(buf[at:])
	at += 2
	if len(buf[at:]) < 1 {
		return io.ErrUnexpectedEOF
	}
	bbp.iNT8 = iohelp.ReadUint8Bytes(buf[at:])
	at += 1
	if len(buf[at:]) < 2 {
		return io.ErrUnexpectedEOF
	}
	bbp.iNT8_ = iohelp.ReadInt16Bytes(buf[at:])
	at += 2
	if len(buf[at:]) < 2 {
		return io.ErrUnexpectedEOF
	}
	bbp.iNT16 = iohelp.ReadInt16Bytes(buf[at:])
	at += 2
	if len(buf[at:]) < 2 {
		return io.ErrUnexpectedEOF
	}
	bbp.iNT16_ = iohelp.ReadInt16Bytes(buf[at:])
	at += 2
	if len(buf[at:]) < 4 {
		return io.ErrUnexpectedEOF
	}
	bbp.iNT32 = iohelp.ReadInt32Bytes(buf[at:])
	at += 4
	if len(buf[at:]) < 4 {
		return io.ErrUnexpectedEOF
	}
	bbp.iNT32_ = iohelp.ReadInt32Bytes(buf[at:])
	at += 4
	if len(buf[at:]) < 1 {
		return io.ErrUnexpectedEOF
	}
	bbp.tRUE = iohelp.ReadBoolBytes(buf[at:])
	at += 1
	if len(buf[at:]) < 1 {
		return io.ErrUnexpectedEOF
	}
	bbp.fALSE = iohelp.ReadBoolBytes(buf[at:])
	at += 1
	if len(buf[at:]) < 8 {
		return io.ErrUnexpectedEOF
	}
	bbp.fLOAT = iohelp.ReadFloat64Bytes(buf[at:])
	at += 8
	if len(buf[at:]) < 8 {
		return io.ErrUnexpectedEOF
	}
	bbp.fLOAT_x = iohelp.ReadFloat64Bytes(buf[at:])
	at += 8
	bbp.sTRING0, err = iohelp.ReadStringBytes(buf[at:])
	if err != nil {
		return err
	}
	at += 4 + len(bbp.sTRING0)
	bbp.sTRING1, err = iohelp.ReadStringBytes(buf[at:])
	if err != nil {
		return err
	}
	at += 4 + len(bbp.sTRING1)
	bbp.sTRING4, err = iohelp.ReadStringBytes(buf[at:])
	if err != nil {
		return err
	}
	at += 4 + len(bbp.sTRING4)
	bbp.sTRING8, err = iohelp.ReadStringBytes(buf[at:])
	if err != nil {
		return err
	}
	at += 4 + len(bbp.sTRING8)
	bbp.sTRING16, err = iohelp.ReadStringBytes(buf[at:])
	if err != nil {
		return err
	}
	at += 4 + len(bbp.sTRING16)
	if len(buf[at:]) < 4 {
		return io.ErrUnexpectedEOF
	}
	bbp.aRRAY0 = make([]int32, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	if len(buf[at:]) < len(bbp.aRRAY0)*4 {
		return io.ErrUnexpectedEOF
	}
	for i1 := range bbp.aRRAY0 {
		(bbp.aRRAY0)[i1] = iohelp.ReadInt32Bytes(buf[at:])
		at += 4
	}
	if len(buf[at:]) < 4 {
		return io.ErrUnexpectedEOF
	}
	bbp.aRRAY1 = make([]string, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	for i1 := range bbp.aRRAY1 {
		(bbp.aRRAY1)[i1], err = iohelp.ReadStringBytes(buf[at:])
		if err != nil {
			return err
		}
		at += 4 + len((bbp.aRRAY1)[i1])
	}
	if len(buf[at:]) < 4 {
		return io.ErrUnexpectedEOF
	}
	bbp.aRRAY8 = make([]int32, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	if len(buf[at:]) < len(bbp.aRRAY8)*4 {
		return io.ErrUnexpectedEOF
	}
	for i1 := range bbp.aRRAY8 {
		(bbp.aRRAY8)[i1] = iohelp.ReadInt32Bytes(buf[at:])
		at += 4
	}
	return nil
}

func (bbp *msgpackComparison) MustUnmarshalBebop(buf []byte) {
	at := 0
	bbp.iNT0 = iohelp.ReadUint8Bytes(buf[at:])
	at += 1
	bbp.iNT1 = iohelp.ReadUint8Bytes(buf[at:])
	at += 1
	bbp.iNT1_ = iohelp.ReadInt16Bytes(buf[at:])
	at += 2
	bbp.iNT8 = iohelp.ReadUint8Bytes(buf[at:])
	at += 1
	bbp.iNT8_ = iohelp.ReadInt16Bytes(buf[at:])
	at += 2
	bbp.iNT16 = iohelp.ReadInt16Bytes(buf[at:])
	at += 2
	bbp.iNT16_ = iohelp.ReadInt16Bytes(buf[at:])
	at += 2
	bbp.iNT32 = iohelp.ReadInt32Bytes(buf[at:])
	at += 4
	bbp.iNT32_ = iohelp.ReadInt32Bytes(buf[at:])
	at += 4
	bbp.tRUE = iohelp.ReadBoolBytes(buf[at:])
	at += 1
	bbp.fALSE = iohelp.ReadBoolBytes(buf[at:])
	at += 1
	bbp.fLOAT = iohelp.ReadFloat64Bytes(buf[at:])
	at += 8
	bbp.fLOAT_x = iohelp.ReadFloat64Bytes(buf[at:])
	at += 8
	bbp.sTRING0 = iohelp.MustReadStringBytes(buf[at:])
	at += 4 + len(bbp.sTRING0)
	bbp.sTRING1 = iohelp.MustReadStringBytes(buf[at:])
	at += 4 + len(bbp.sTRING1)
	bbp.sTRING4 = iohelp.MustReadStringBytes(buf[at:])
	at += 4 + len(bbp.sTRING4)
	bbp.sTRING8 = iohelp.MustReadStringBytes(buf[at:])
	at += 4 + len(bbp.sTRING8)
	bbp.sTRING16 = iohelp.MustReadStringBytes(buf[at:])
	at += 4 + len(bbp.sTRING16)
	bbp.aRRAY0 = make([]int32, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	for i1 := range bbp.aRRAY0 {
		(bbp.aRRAY0)[i1] = iohelp.ReadInt32Bytes(buf[at:])
		at += 4
	}
	bbp.aRRAY1 = make([]string, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	for i1 := range bbp.aRRAY1 {
		(bbp.aRRAY1)[i1] = iohelp.MustReadStringBytes(buf[at:])
		at += 4 + len((bbp.aRRAY1)[i1])
	}
	bbp.aRRAY8 = make([]int32, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	for i1 := range bbp.aRRAY8 {
		(bbp.aRRAY8)[i1] = iohelp.ReadInt32Bytes(buf[at:])
		at += 4
	}
}

func (bbp msgpackComparison) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint8(w, bbp.iNT0)
	iohelp.WriteUint8(w, bbp.iNT1)
	iohelp.WriteInt16(w, bbp.iNT1_)
	iohelp.WriteUint8(w, bbp.iNT8)
	iohelp.WriteInt16(w, bbp.iNT8_)
	iohelp.WriteInt16(w, bbp.iNT16)
	iohelp.WriteInt16(w, bbp.iNT16_)
	iohelp.WriteInt32(w, bbp.iNT32)
	iohelp.WriteInt32(w, bbp.iNT32_)
	iohelp.WriteBool(w, bbp.tRUE)
	iohelp.WriteBool(w, bbp.fALSE)
	iohelp.WriteFloat64(w, bbp.fLOAT)
	iohelp.WriteFloat64(w, bbp.fLOAT_x)
	iohelp.WriteUint32(w, uint32(len(bbp.sTRING0)))
	w.Write([]byte(bbp.sTRING0))
	iohelp.WriteUint32(w, uint32(len(bbp.sTRING1)))
	w.Write([]byte(bbp.sTRING1))
	iohelp.WriteUint32(w, uint32(len(bbp.sTRING4)))
	w.Write([]byte(bbp.sTRING4))
	iohelp.WriteUint32(w, uint32(len(bbp.sTRING8)))
	w.Write([]byte(bbp.sTRING8))
	iohelp.WriteUint32(w, uint32(len(bbp.sTRING16)))
	w.Write([]byte(bbp.sTRING16))
	iohelp.WriteUint32(w, uint32(len(bbp.aRRAY0)))
	for _, elem := range bbp.aRRAY0 {
		iohelp.WriteInt32(w, elem)
	}
	iohelp.WriteUint32(w, uint32(len(bbp.aRRAY1)))
	for _, elem := range bbp.aRRAY1 {
		iohelp.WriteUint32(w, uint32(len(elem)))
		w.Write([]byte(elem))
	}
	iohelp.WriteUint32(w, uint32(len(bbp.aRRAY8)))
	for _, elem := range bbp.aRRAY8 {
		iohelp.WriteInt32(w, elem)
	}
	return w.Err
}

func (bbp *msgpackComparison) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bbp.iNT0 = iohelp.ReadUint8(r)
	bbp.iNT1 = iohelp.ReadUint8(r)
	bbp.iNT1_ = iohelp.ReadInt16(r)
	bbp.iNT8 = iohelp.ReadUint8(r)
	bbp.iNT8_ = iohelp.ReadInt16(r)
	bbp.iNT16 = iohelp.ReadInt16(r)
	bbp.iNT16_ = iohelp.ReadInt16(r)
	bbp.iNT32 = iohelp.ReadInt32(r)
	bbp.iNT32_ = iohelp.ReadInt32(r)
	bbp.tRUE = iohelp.ReadBool(r)
	bbp.fALSE = iohelp.ReadBool(r)
	bbp.fLOAT = iohelp.ReadFloat64(r)
	bbp.fLOAT_x = iohelp.ReadFloat64(r)
	bbp.sTRING0 = iohelp.ReadString(r)
	bbp.sTRING1 = iohelp.ReadString(r)
	bbp.sTRING4 = iohelp.ReadString(r)
	bbp.sTRING8 = iohelp.ReadString(r)
	bbp.sTRING16 = iohelp.ReadString(r)
	bbp.aRRAY0 = make([]int32, iohelp.ReadUint32(r))
	for i1 := range bbp.aRRAY0 {
		(bbp.aRRAY0[i1]) = iohelp.ReadInt32(r)
	}
	bbp.aRRAY1 = make([]string, iohelp.ReadUint32(r))
	for i1 := range bbp.aRRAY1 {
		(bbp.aRRAY1[i1]) = iohelp.ReadString(r)
	}
	bbp.aRRAY8 = make([]int32, iohelp.ReadUint32(r))
	for i1 := range bbp.aRRAY8 {
		(bbp.aRRAY8[i1]) = iohelp.ReadInt32(r)
	}
	return r.Err
}

func (bbp msgpackComparison) Size() int {
	bodyLen := 0
	bodyLen += 1
	bodyLen += 1
	bodyLen += 2
	bodyLen += 1
	bodyLen += 2
	bodyLen += 2
	bodyLen += 2
	bodyLen += 4
	bodyLen += 4
	bodyLen += 1
	bodyLen += 1
	bodyLen += 8
	bodyLen += 8
	bodyLen += 4 + len(bbp.sTRING0)
	bodyLen += 4 + len(bbp.sTRING1)
	bodyLen += 4 + len(bbp.sTRING4)
	bodyLen += 4 + len(bbp.sTRING8)
	bodyLen += 4 + len(bbp.sTRING16)
	bodyLen += 4
	bodyLen += len(bbp.aRRAY0) * 4
	bodyLen += 4
	for _, elem := range bbp.aRRAY1 {
		bodyLen += 4 + len(elem)
	}
	bodyLen += 4
	bodyLen += len(bbp.aRRAY8) * 4
	return bodyLen
}

func (bbp msgpackComparison) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func makemsgpackComparison(r *iohelp.ErrorReader) (msgpackComparison, error) {
	v := msgpackComparison{}
	err := v.DecodeBebop(r)
	return v, err
}

func makemsgpackComparisonFromBytes(buf []byte) (msgpackComparison, error) {
	v := msgpackComparison{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func mustMakemsgpackComparisonFromBytes(buf []byte) msgpackComparison {
	v := msgpackComparison{}
	v.MustUnmarshalBebop(buf)
	return v
}

