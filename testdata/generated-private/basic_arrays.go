// Code generated by bebopc-go; DO NOT EDIT.

package generated

import (
	"github.com/200sc/bebop"
	"github.com/200sc/bebop/iohelp"
	"io"
)

var _ bebop.Record = &basicArrays{}

type basicArrays struct {
	a_bool []bool
	a_byte []byte
	a_int16 []int16
	a_uint16 []uint16
	a_int32 []int32
	a_uint32 []uint32
	a_int64 []int64
	a_uint64 []uint64
	a_float32 []float32
	a_float64 []float64
	a_string []string
	a_guid [][16]byte
}

func (bbp basicArrays) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.a_bool)))
	at += 4
	for _, v1 := range bbp.a_bool {
		iohelp.WriteBoolBytes(buf[at:], v1)
		at += 1
	}
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.a_byte)))
	at += 4
	copy(buf[at:at+len(bbp.a_byte)], bbp.a_byte)
	at += len(bbp.a_byte)
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.a_int16)))
	at += 4
	for _, v1 := range bbp.a_int16 {
		iohelp.WriteInt16Bytes(buf[at:], v1)
		at += 2
	}
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.a_uint16)))
	at += 4
	for _, v1 := range bbp.a_uint16 {
		iohelp.WriteUint16Bytes(buf[at:], v1)
		at += 2
	}
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.a_int32)))
	at += 4
	for _, v1 := range bbp.a_int32 {
		iohelp.WriteInt32Bytes(buf[at:], v1)
		at += 4
	}
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.a_uint32)))
	at += 4
	for _, v1 := range bbp.a_uint32 {
		iohelp.WriteUint32Bytes(buf[at:], v1)
		at += 4
	}
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.a_int64)))
	at += 4
	for _, v1 := range bbp.a_int64 {
		iohelp.WriteInt64Bytes(buf[at:], v1)
		at += 8
	}
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.a_uint64)))
	at += 4
	for _, v1 := range bbp.a_uint64 {
		iohelp.WriteUint64Bytes(buf[at:], v1)
		at += 8
	}
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.a_float32)))
	at += 4
	for _, v1 := range bbp.a_float32 {
		iohelp.WriteFloat32Bytes(buf[at:], v1)
		at += 4
	}
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.a_float64)))
	at += 4
	for _, v1 := range bbp.a_float64 {
		iohelp.WriteFloat64Bytes(buf[at:], v1)
		at += 8
	}
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.a_string)))
	at += 4
	for _, v1 := range bbp.a_string {
		iohelp.WriteUint32Bytes(buf[at:], uint32(len(v1)))
		copy(buf[at+4:at+4+len(v1)], []byte(v1))
		at += 4 + len(v1)
	}
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.a_guid)))
	at += 4
	for _, v1 := range bbp.a_guid {
		iohelp.WriteGUIDBytes(buf[at:], v1)
		at += 16
	}
	return at
}

func (bbp *basicArrays) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	if len(buf[at:]) < 4 {
		return io.ErrUnexpectedEOF
	}
	bbp.a_bool = make([]bool, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	if len(buf[at:]) < len(bbp.a_bool)*1 {
		return io.ErrUnexpectedEOF
	}
	for i1 := range bbp.a_bool {
		(bbp.a_bool)[i1] = iohelp.ReadBoolBytes(buf[at:])
		at += 1
	}
	if len(buf[at:]) < 4 {
		return io.ErrUnexpectedEOF
	}
	bbp.a_byte = make([]byte, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	if len(buf[at:]) < len(bbp.a_byte)*1 {
		return io.ErrUnexpectedEOF
	}
	copy(bbp.a_byte, buf[at:at+len(bbp.a_byte)])
	at += len(bbp.a_byte)
	if len(buf[at:]) < 4 {
		return io.ErrUnexpectedEOF
	}
	bbp.a_int16 = make([]int16, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	if len(buf[at:]) < len(bbp.a_int16)*2 {
		return io.ErrUnexpectedEOF
	}
	for i1 := range bbp.a_int16 {
		(bbp.a_int16)[i1] = iohelp.ReadInt16Bytes(buf[at:])
		at += 2
	}
	if len(buf[at:]) < 4 {
		return io.ErrUnexpectedEOF
	}
	bbp.a_uint16 = make([]uint16, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	if len(buf[at:]) < len(bbp.a_uint16)*2 {
		return io.ErrUnexpectedEOF
	}
	for i1 := range bbp.a_uint16 {
		(bbp.a_uint16)[i1] = iohelp.ReadUint16Bytes(buf[at:])
		at += 2
	}
	if len(buf[at:]) < 4 {
		return io.ErrUnexpectedEOF
	}
	bbp.a_int32 = make([]int32, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	if len(buf[at:]) < len(bbp.a_int32)*4 {
		return io.ErrUnexpectedEOF
	}
	for i1 := range bbp.a_int32 {
		(bbp.a_int32)[i1] = iohelp.ReadInt32Bytes(buf[at:])
		at += 4
	}
	if len(buf[at:]) < 4 {
		return io.ErrUnexpectedEOF
	}
	bbp.a_uint32 = make([]uint32, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	if len(buf[at:]) < len(bbp.a_uint32)*4 {
		return io.ErrUnexpectedEOF
	}
	for i1 := range bbp.a_uint32 {
		(bbp.a_uint32)[i1] = iohelp.ReadUint32Bytes(buf[at:])
		at += 4
	}
	if len(buf[at:]) < 4 {
		return io.ErrUnexpectedEOF
	}
	bbp.a_int64 = make([]int64, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	if len(buf[at:]) < len(bbp.a_int64)*8 {
		return io.ErrUnexpectedEOF
	}
	for i1 := range bbp.a_int64 {
		(bbp.a_int64)[i1] = iohelp.ReadInt64Bytes(buf[at:])
		at += 8
	}
	if len(buf[at:]) < 4 {
		return io.ErrUnexpectedEOF
	}
	bbp.a_uint64 = make([]uint64, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	if len(buf[at:]) < len(bbp.a_uint64)*8 {
		return io.ErrUnexpectedEOF
	}
	for i1 := range bbp.a_uint64 {
		(bbp.a_uint64)[i1] = iohelp.ReadUint64Bytes(buf[at:])
		at += 8
	}
	if len(buf[at:]) < 4 {
		return io.ErrUnexpectedEOF
	}
	bbp.a_float32 = make([]float32, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	if len(buf[at:]) < len(bbp.a_float32)*4 {
		return io.ErrUnexpectedEOF
	}
	for i1 := range bbp.a_float32 {
		(bbp.a_float32)[i1] = iohelp.ReadFloat32Bytes(buf[at:])
		at += 4
	}
	if len(buf[at:]) < 4 {
		return io.ErrUnexpectedEOF
	}
	bbp.a_float64 = make([]float64, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	if len(buf[at:]) < len(bbp.a_float64)*8 {
		return io.ErrUnexpectedEOF
	}
	for i1 := range bbp.a_float64 {
		(bbp.a_float64)[i1] = iohelp.ReadFloat64Bytes(buf[at:])
		at += 8
	}
	if len(buf[at:]) < 4 {
		return io.ErrUnexpectedEOF
	}
	bbp.a_string = make([]string, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	for i1 := range bbp.a_string {
		(bbp.a_string)[i1], err = iohelp.ReadStringBytes(buf[at:])
		if err != nil {
			return err
		}
		at += 4 + len((bbp.a_string)[i1])
	}
	if len(buf[at:]) < 4 {
		return io.ErrUnexpectedEOF
	}
	bbp.a_guid = make([][16]byte, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	if len(buf[at:]) < len(bbp.a_guid)*16 {
		return io.ErrUnexpectedEOF
	}
	for i1 := range bbp.a_guid {
		(bbp.a_guid)[i1] = iohelp.ReadGUIDBytes(buf[at:])
		at += 16
	}
	return nil
}

func (bbp *basicArrays) MustUnmarshalBebop(buf []byte) {
	at := 0
	bbp.a_bool = make([]bool, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	for i1 := range bbp.a_bool {
		(bbp.a_bool)[i1] = iohelp.ReadBoolBytes(buf[at:])
		at += 1
	}
	bbp.a_byte = make([]byte, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	copy(bbp.a_byte, buf[at:at+len(bbp.a_byte)])
	at += len(bbp.a_byte)
	bbp.a_int16 = make([]int16, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	for i1 := range bbp.a_int16 {
		(bbp.a_int16)[i1] = iohelp.ReadInt16Bytes(buf[at:])
		at += 2
	}
	bbp.a_uint16 = make([]uint16, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	for i1 := range bbp.a_uint16 {
		(bbp.a_uint16)[i1] = iohelp.ReadUint16Bytes(buf[at:])
		at += 2
	}
	bbp.a_int32 = make([]int32, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	for i1 := range bbp.a_int32 {
		(bbp.a_int32)[i1] = iohelp.ReadInt32Bytes(buf[at:])
		at += 4
	}
	bbp.a_uint32 = make([]uint32, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	for i1 := range bbp.a_uint32 {
		(bbp.a_uint32)[i1] = iohelp.ReadUint32Bytes(buf[at:])
		at += 4
	}
	bbp.a_int64 = make([]int64, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	for i1 := range bbp.a_int64 {
		(bbp.a_int64)[i1] = iohelp.ReadInt64Bytes(buf[at:])
		at += 8
	}
	bbp.a_uint64 = make([]uint64, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	for i1 := range bbp.a_uint64 {
		(bbp.a_uint64)[i1] = iohelp.ReadUint64Bytes(buf[at:])
		at += 8
	}
	bbp.a_float32 = make([]float32, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	for i1 := range bbp.a_float32 {
		(bbp.a_float32)[i1] = iohelp.ReadFloat32Bytes(buf[at:])
		at += 4
	}
	bbp.a_float64 = make([]float64, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	for i1 := range bbp.a_float64 {
		(bbp.a_float64)[i1] = iohelp.ReadFloat64Bytes(buf[at:])
		at += 8
	}
	bbp.a_string = make([]string, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	for i1 := range bbp.a_string {
		(bbp.a_string)[i1] = iohelp.MustReadStringBytes(buf[at:])
		at += 4 + len((bbp.a_string)[i1])
	}
	bbp.a_guid = make([][16]byte, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	for i1 := range bbp.a_guid {
		(bbp.a_guid)[i1] = iohelp.ReadGUIDBytes(buf[at:])
		at += 16
	}
}

func (bbp basicArrays) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(len(bbp.a_bool)))
	for _, elem := range bbp.a_bool {
		iohelp.WriteBool(w, elem)
	}
	iohelp.WriteUint32(w, uint32(len(bbp.a_byte)))
	w.Write(bbp.a_byte)
	iohelp.WriteUint32(w, uint32(len(bbp.a_int16)))
	for _, elem := range bbp.a_int16 {
		iohelp.WriteInt16(w, elem)
	}
	iohelp.WriteUint32(w, uint32(len(bbp.a_uint16)))
	for _, elem := range bbp.a_uint16 {
		iohelp.WriteUint16(w, elem)
	}
	iohelp.WriteUint32(w, uint32(len(bbp.a_int32)))
	for _, elem := range bbp.a_int32 {
		iohelp.WriteInt32(w, elem)
	}
	iohelp.WriteUint32(w, uint32(len(bbp.a_uint32)))
	for _, elem := range bbp.a_uint32 {
		iohelp.WriteUint32(w, elem)
	}
	iohelp.WriteUint32(w, uint32(len(bbp.a_int64)))
	for _, elem := range bbp.a_int64 {
		iohelp.WriteInt64(w, elem)
	}
	iohelp.WriteUint32(w, uint32(len(bbp.a_uint64)))
	for _, elem := range bbp.a_uint64 {
		iohelp.WriteUint64(w, elem)
	}
	iohelp.WriteUint32(w, uint32(len(bbp.a_float32)))
	for _, elem := range bbp.a_float32 {
		iohelp.WriteFloat32(w, elem)
	}
	iohelp.WriteUint32(w, uint32(len(bbp.a_float64)))
	for _, elem := range bbp.a_float64 {
		iohelp.WriteFloat64(w, elem)
	}
	iohelp.WriteUint32(w, uint32(len(bbp.a_string)))
	for _, elem := range bbp.a_string {
		iohelp.WriteUint32(w, uint32(len(elem)))
		w.Write([]byte(elem))
	}
	iohelp.WriteUint32(w, uint32(len(bbp.a_guid)))
	for _, elem := range bbp.a_guid {
		iohelp.WriteGUID(w, elem)
	}
	return w.Err
}

func (bbp *basicArrays) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bbp.a_bool = make([]bool, iohelp.ReadUint32(r))
	for i1 := range bbp.a_bool {
		(bbp.a_bool[i1]) = iohelp.ReadBool(r)
	}
	bbp.a_byte = make([]byte, iohelp.ReadUint32(r))
	r.Read(bbp.a_byte)
	bbp.a_int16 = make([]int16, iohelp.ReadUint32(r))
	for i1 := range bbp.a_int16 {
		(bbp.a_int16[i1]) = iohelp.ReadInt16(r)
	}
	bbp.a_uint16 = make([]uint16, iohelp.ReadUint32(r))
	for i1 := range bbp.a_uint16 {
		(bbp.a_uint16[i1]) = iohelp.ReadUint16(r)
	}
	bbp.a_int32 = make([]int32, iohelp.ReadUint32(r))
	for i1 := range bbp.a_int32 {
		(bbp.a_int32[i1]) = iohelp.ReadInt32(r)
	}
	bbp.a_uint32 = make([]uint32, iohelp.ReadUint32(r))
	for i1 := range bbp.a_uint32 {
		(bbp.a_uint32[i1]) = iohelp.ReadUint32(r)
	}
	bbp.a_int64 = make([]int64, iohelp.ReadUint32(r))
	for i1 := range bbp.a_int64 {
		(bbp.a_int64[i1]) = iohelp.ReadInt64(r)
	}
	bbp.a_uint64 = make([]uint64, iohelp.ReadUint32(r))
	for i1 := range bbp.a_uint64 {
		(bbp.a_uint64[i1]) = iohelp.ReadUint64(r)
	}
	bbp.a_float32 = make([]float32, iohelp.ReadUint32(r))
	for i1 := range bbp.a_float32 {
		(bbp.a_float32[i1]) = iohelp.ReadFloat32(r)
	}
	bbp.a_float64 = make([]float64, iohelp.ReadUint32(r))
	for i1 := range bbp.a_float64 {
		(bbp.a_float64[i1]) = iohelp.ReadFloat64(r)
	}
	bbp.a_string = make([]string, iohelp.ReadUint32(r))
	for i1 := range bbp.a_string {
		(bbp.a_string[i1]) = iohelp.ReadString(r)
	}
	bbp.a_guid = make([][16]byte, iohelp.ReadUint32(r))
	for i1 := range bbp.a_guid {
		(bbp.a_guid[i1]) = iohelp.ReadGUID(r)
	}
	return r.Err
}

func (bbp basicArrays) Size() int {
	bodyLen := 0
	bodyLen += 4
	bodyLen += len(bbp.a_bool) * 1
	bodyLen += 4
	bodyLen += len(bbp.a_byte) * 1
	bodyLen += 4
	bodyLen += len(bbp.a_int16) * 2
	bodyLen += 4
	bodyLen += len(bbp.a_uint16) * 2
	bodyLen += 4
	bodyLen += len(bbp.a_int32) * 4
	bodyLen += 4
	bodyLen += len(bbp.a_uint32) * 4
	bodyLen += 4
	bodyLen += len(bbp.a_int64) * 8
	bodyLen += 4
	bodyLen += len(bbp.a_uint64) * 8
	bodyLen += 4
	bodyLen += len(bbp.a_float32) * 4
	bodyLen += 4
	bodyLen += len(bbp.a_float64) * 8
	bodyLen += 4
	for _, elem := range bbp.a_string {
		bodyLen += 4 + len(elem)
	}
	bodyLen += 4
	bodyLen += len(bbp.a_guid) * 16
	return bodyLen
}

func (bbp basicArrays) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func makebasicArrays(r *iohelp.ErrorReader) (basicArrays, error) {
	v := basicArrays{}
	err := v.DecodeBebop(r)
	return v, err
}

func makebasicArraysFromBytes(buf []byte) (basicArrays, error) {
	v := basicArrays{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func mustMakebasicArraysFromBytes(buf []byte) basicArrays {
	v := basicArrays{}
	v.MustUnmarshalBebop(buf)
	return v
}

var _ bebop.Record = &testInt32Array{}

type testInt32Array struct {
	a []int32
}

func (bbp testInt32Array) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.a)))
	at += 4
	for _, v1 := range bbp.a {
		iohelp.WriteInt32Bytes(buf[at:], v1)
		at += 4
	}
	return at
}

func (bbp *testInt32Array) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	if len(buf[at:]) < 4 {
		return io.ErrUnexpectedEOF
	}
	bbp.a = make([]int32, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	if len(buf[at:]) < len(bbp.a)*4 {
		return io.ErrUnexpectedEOF
	}
	for i1 := range bbp.a {
		(bbp.a)[i1] = iohelp.ReadInt32Bytes(buf[at:])
		at += 4
	}
	return nil
}

func (bbp *testInt32Array) MustUnmarshalBebop(buf []byte) {
	at := 0
	bbp.a = make([]int32, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	for i1 := range bbp.a {
		(bbp.a)[i1] = iohelp.ReadInt32Bytes(buf[at:])
		at += 4
	}
}

func (bbp testInt32Array) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(len(bbp.a)))
	for _, elem := range bbp.a {
		iohelp.WriteInt32(w, elem)
	}
	return w.Err
}

func (bbp *testInt32Array) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bbp.a = make([]int32, iohelp.ReadUint32(r))
	for i1 := range bbp.a {
		(bbp.a[i1]) = iohelp.ReadInt32(r)
	}
	return r.Err
}

func (bbp testInt32Array) Size() int {
	bodyLen := 0
	bodyLen += 4
	bodyLen += len(bbp.a) * 4
	return bodyLen
}

func (bbp testInt32Array) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func maketestInt32Array(r *iohelp.ErrorReader) (testInt32Array, error) {
	v := testInt32Array{}
	err := v.DecodeBebop(r)
	return v, err
}

func maketestInt32ArrayFromBytes(buf []byte) (testInt32Array, error) {
	v := testInt32Array{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func mustMaketestInt32ArrayFromBytes(buf []byte) testInt32Array {
	v := testInt32Array{}
	v.MustUnmarshalBebop(buf)
	return v
}

